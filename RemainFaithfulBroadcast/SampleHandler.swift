import ReplayKit

class SampleHandler: RPBroadcastSampleHandler {

    private let appGroupID = "group.com.remainfaithful.app"

    private var sharedDefaults: UserDefaults? {
        UserDefaults(suiteName: appGroupID)
    }

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

    override func processSampleBuffer(_ sampleBuffer: CMSampleBuffer, with sampleBufferType: RPSampleBufferType) {
        guard sampleBufferType == .video else { return }

        guard let imageBuffer = CMSampleBufferGetImageBuffer(sampleBuffer) else { return }

        let width = CVPixelBufferGetWidth(imageBuffer)
        let height = CVPixelBufferGetHeight(imageBuffer)
        let pts = CMSampleBufferGetPresentationTimeStamp(sampleBuffer)

        print("[RemainFaithfulBroadcast] frame \(width)x\(height) @ \(CMTimeGetSeconds(pts))s")

        sharedDefaults?.set("\(width)x\(height)", forKey: "lastFrameDimensions")
        sharedDefaults?.set(CMTimeGetSeconds(pts), forKey: "lastFrameTimestamp")
    }
}
