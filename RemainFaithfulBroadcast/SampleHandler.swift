import ReplayKit
import SensitiveContentAnalysis
import CoreImage
import UIKit

class SampleHandler: RPBroadcastSampleHandler {

    private let appGroupID = "group.com.remainfaithful.app"

    private var sharedDefaults: UserDefaults? {
        UserDefaults(suiteName: appGroupID)
    }

    private let analyzer  = SCSensitivityAnalyzer()
    // GPU-backed context; avoids allocating a new one per frame.
    private let ciContext = CIContext(options: [.useSoftwareRenderer: false])
    private var frameCount = 0

    // MARK: - Broadcast lifecycle

    override func broadcastStarted(withSetupInfo setupInfo: [String: NSObject]?) {
        sharedDefaults?.set(true, forKey: "isBroadcasting")
        sharedDefaults?.set(Date().timeIntervalSince1970, forKey: "broadcastStartTime")
    }

    override func broadcastPaused() {
        sharedDefaults?.set(false, forKey: "isBroadcasting")
    }

    override func broadcastResumed() {
        sharedDefaults?.set(true, forKey: "isBroadcasting")
    }

    override func broadcastFinished() {
        sharedDefaults?.set(false, forKey: "isBroadcasting")
    }

    // MARK: - Sample buffer processing

    override func processSampleBuffer(_ sampleBuffer: CMSampleBuffer,
                                      with sampleBufferType: RPSampleBufferType) {
        guard sampleBufferType == .video else { return }
        guard let imageBuffer = CMSampleBufferGetImageBuffer(sampleBuffer) else { return }

        frameCount += 1

        let width  = CVPixelBufferGetWidth(imageBuffer)
        let height = CVPixelBufferGetHeight(imageBuffer)
        let pts    = CMSampleBufferGetPresentationTimeStamp(sampleBuffer)

        print("[RemainFaithfulBroadcast] frame \(width)x\(height) @ \(CMTimeGetSeconds(pts))s")

        sharedDefaults?.set("\(width)x\(height)", forKey: "lastFrameDimensions")
        sharedDefaults?.set(CMTimeGetSeconds(pts), forKey: "lastFrameTimestamp")

        if frameCount % 30 == 0 {
            scheduleAnalysis(pixelBuffer: imageBuffer,
                             frameNumber: frameCount,
                             timestamp: CMTimeGetSeconds(pts))
        }
    }

    // MARK: - Sensitive content analysis

    // Converts the pixel buffer to a JPEG on disk, passes it to
    // SCSensitivityAnalyzer, then logs the result. Runs in a detached Task
    // so it never stalls the ReplayKit sample-buffer delivery thread.
    //
    // SCSensitivityAnalysis exposes `isSensitive` (Bool) as its public result.
    // The framework internally applies its own threshold; a raw floating-point
    // confidence score is not part of the public API surface.
    private func scheduleAnalysis(pixelBuffer: CVPixelBuffer,
                                  frameNumber: Int,
                                  timestamp: Double) {
        let ciImage = CIImage(cvPixelBuffer: pixelBuffer)
        guard let cgImage = ciContext.createCGImage(ciImage, from: ciImage.extent) else {
            print("[RemainFaithfulBroadcast] frame \(frameNumber): CGImage render failed")
            return
        }
        guard let jpegData = UIImage(cgImage: cgImage).jpegData(compressionQuality: 0.75) else {
            print("[RemainFaithfulBroadcast] frame \(frameNumber): JPEG encode failed")
            return
        }

        // Use a per-frame filename so concurrent Tasks never clobber each other.
        let tmpURL = FileManager.default.temporaryDirectory
            .appendingPathComponent("rf_sca_\(frameNumber).jpg")
        do {
            try jpegData.write(to: tmpURL, options: .atomic)
        } catch {
            print("[RemainFaithfulBroadcast] frame \(frameNumber): temp write failed – \(error)")
            return
        }

        // Capture only what the task needs; SCSensitivityAnalyzer is a class
        // so this is a cheap reference copy.
        let capturedAnalyzer = analyzer
        let capturedDefaults = sharedDefaults

        Task.detached(priority: .utility) {
            defer { try? FileManager.default.removeItem(at: tmpURL) }
            do {
                let analysis = try await capturedAnalyzer.analyzeImage(at: tmpURL)

                // SCSensitivityAnalysis.isSensitive is the framework's binary
                // detection result. Log it alongside the analysis policy so the
                // output is actionable even though no raw confidence float is
                // exposed by the public API.
                print(
                    "[RemainFaithfulBroadcast] SCA frame \(frameNumber)" +
                    " @ \(String(format: "%.2f", timestamp))s" +
                    " — sensitive: \(analysis.isSensitive)" +
                    " confidence: \(analysis.isSensitive ? "above threshold" : "below threshold")" +
                    " (policy: \(capturedAnalyzer.analysisPolicy))"
                )

                if analysis.isSensitive {
                    capturedDefaults?.set(true,      forKey: "sensitiveContentDetected")
                    capturedDefaults?.set(timestamp, forKey: "lastSensitiveContentTimestamp")
                }
            } catch {
                print("[RemainFaithfulBroadcast] SCA error frame \(frameNumber): \(error)")
            }
        }
    }
}
