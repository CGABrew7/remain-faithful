import AVFoundation
import Foundation

enum FutureSelfCopy {
    static let title = "Record why you started"
    static let body = "In a weak moment, Remain Faithful can play this back to you. It stays on this phone."
    static let recordButton = "Start recording"
    static let stopButton = "Stop recording"
    static let playTitle = "Hear why you started"
    static let playAgain = "Play again"
    static let skip = "Skip"
    static let rememberTitle = "What you want to remember"
    static let micReason = "Remain Faithful uses the microphone so you can record a short message to yourself."
    static let later = "You can record this later in Settings."
    static let shortTake = "Hold a little longer, then save."
    static let preview = "Preview"
    static let delete = "Delete recording"
    static let saved = "Saved on this phone."

    static var userFacing: [String] {
        [title, body, recordButton, stopButton, playTitle, playAgain, skip, rememberTitle, micReason, later, shortTake, preview, delete, saved]
    }
}

enum FutureSelfRules {
    static let minSeconds: TimeInterval = 1
    static let maxSeconds: TimeInterval = 45
    static let countKey = "futureSelfClipCount"
    static let lastPlayedKey = "futureSelfLastPlayedAt"

    /// A take under a second is dropped. A take past the cap is dropped.
    static func keep(duration: TimeInterval) -> Bool {
        duration >= minSeconds && duration < maxSeconds + 1
    }

    static var skipSendsPartnerNotice: Bool { false }
}

enum FutureSelfPlayback {
    static func shouldAutoPlay(hasClip: Bool, alreadyPlayedThisOpen: Bool) -> Bool {
        hasClip && !alreadyPlayedThisOpen
    }
}

enum FutureSelfClip: String, CaseIterable {
    case why = "why.m4a"
    case remember = "remember.m4a"
}

enum FutureSelfStore {
    static func directory(base: URL = defaultBase()) -> URL {
        base.appendingPathComponent("FutureSelf", isDirectory: true)
    }

    static func url(for clip: FutureSelfClip, base: URL = defaultBase()) -> URL {
        directory(base: base).appendingPathComponent(clip.rawValue)
    }

    static func exists(_ clip: FutureSelfClip, base: URL = defaultBase()) -> Bool {
        FileManager.default.fileExists(atPath: url(for: clip, base: base).path)
    }

    static func count(base: URL = defaultBase()) -> Int {
        FutureSelfClip.allCases.reduce(0) { total, clip in
            total + (exists(clip, base: base) ? 1 : 0)
        }
    }

    static func prepare(base: URL = defaultBase()) throws -> URL {
        let folder = directory(base: base)
        try FileManager.default.createDirectory(at: folder, withIntermediateDirectories: true)
        var values = URLResourceValues()
        values.isExcludedFromBackup = true
        var mutable = folder
        try? mutable.setResourceValues(values)
        return folder
    }

    /// Replaces one take. The previous file at that name is removed first.
    static func replace(_ clip: FutureSelfClip, with source: URL, duration: TimeInterval, base: URL = defaultBase()) throws -> Bool {
        guard FutureSelfRules.keep(duration: duration) else {
            try? FileManager.default.removeItem(at: source)
            return false
        }
        _ = try prepare(base: base)
        let destination = url(for: clip, base: base)
        if FileManager.default.fileExists(atPath: destination.path) {
            try FileManager.default.removeItem(at: destination)
        }
        try FileManager.default.moveItem(at: source, to: destination)
        var values = URLResourceValues()
        values.isExcludedFromBackup = true
        var stored = destination
        try? stored.setResourceValues(values)
        noteCount(base: base)
        return true
    }

    static func delete(_ clip: FutureSelfClip, base: URL = defaultBase()) {
        let file = url(for: clip, base: base)
        if FileManager.default.fileExists(atPath: file.path) {
            try? FileManager.default.removeItem(at: file)
        }
        noteCount(base: base)
    }

    static func wipe(base: URL = defaultBase()) {
        let folder = directory(base: base)
        if FileManager.default.fileExists(atPath: folder.path) {
            try? FileManager.default.removeItem(at: folder)
        }
        let defaults = UserDefaults.standard
        defaults.removeObject(forKey: FutureSelfRules.countKey)
        defaults.removeObject(forKey: FutureSelfRules.lastPlayedKey)
    }

    static func notePlayed(at date: Date = Date()) {
        UserDefaults.standard.set(date.timeIntervalSince1970, forKey: FutureSelfRules.lastPlayedKey)
    }

    private static func noteCount(base: URL) {
        UserDefaults.standard.set(count(base: base), forKey: FutureSelfRules.countKey)
    }

    private static func defaultBase() -> URL {
        FileManager.default.urls(for: .applicationSupportDirectory, in: .userDomainMask).first
            ?? FileManager.default.temporaryDirectory
    }
}

final class FutureSelfRecorder: NSObject, ObservableObject, AVAudioRecorderDelegate {
    @Published var recording = false
    @Published var message: String?
    @Published var savedClip: FutureSelfClip?
    @Published var active: FutureSelfClip = .why

    private var recorder: AVAudioRecorder?
    private var started: Date?
    private var clip: FutureSelfClip = .why
    private var file: URL?
    private var stopping = false

    func toggle(_ clip: FutureSelfClip) {
        if recording {
            stop()
        } else {
            begin(clip)
        }
    }

    func begin(_ clip: FutureSelfClip) {
        self.clip = clip
        active = clip
        message = nil
        AVAudioSession.sharedInstance().requestRecordPermission { [weak self] granted in
            DispatchQueue.main.async {
                guard let self else { return }
                guard granted else {
                    self.message = FutureSelfCopy.later
                    return
                }
                self.startRecorder()
            }
        }
    }

    func stop() {
        guard !stopping else { return }
        stopping = true
        defer { stopping = false }
        let elapsed = started.map { Date().timeIntervalSince($0) } ?? 0
        recording = false
        started = nil
        recorder?.stop()
        recorder = nil
        try? AVAudioSession.sharedInstance().setActive(false, options: .notifyOthersOnDeactivation)
        guard let file else { return }
        if FutureSelfRules.keep(duration: elapsed) {
            let kept = (try? FutureSelfStore.replace(clip, with: file, duration: elapsed)) ?? false
            message = kept ? FutureSelfCopy.saved : FutureSelfCopy.shortTake
            if kept { savedClip = clip }
        } else {
            try? FileManager.default.removeItem(at: file)
            message = FutureSelfCopy.shortTake
        }
        self.file = nil
    }

    /// A recording that moves to the background is saved when it is long enough.
    func stopForBackground() {
        guard recording else { return }
        stop()
    }

    private func startRecorder() {
        do {
            let session = AVAudioSession.sharedInstance()
            try session.setCategory(.playAndRecord, mode: .default, options: [.defaultToSpeaker])
            try session.setActive(true)
            let folder = try FutureSelfStore.prepare()
            let target = folder.appendingPathComponent("take-\(UUID().uuidString).m4a")
            let settings: [String: Any] = [
                AVFormatIDKey: Int(kAudioFormatMPEG4AAC),
                AVSampleRateKey: 44100,
                AVNumberOfChannelsKey: 1,
                AVEncoderAudioQualityKey: AVAudioQuality.medium.rawValue,
            ]
            let recorder = try AVAudioRecorder(url: target, settings: settings)
            recorder.delegate = self
            recorder.record(forDuration: FutureSelfRules.maxSeconds)
            self.recorder = recorder
            self.file = target
            self.started = Date()
            self.recording = true
        } catch {
            message = FutureSelfCopy.later
        }
    }

    func audioRecorderDidFinishRecording(_ recorder: AVAudioRecorder, successfully flag: Bool) {
        DispatchQueue.main.async {
            guard self.recording else { return }
            self.stop()
        }
    }
}

final class FutureSelfPlayer: ObservableObject {
    private var player: AVAudioPlayer?

    func play(_ clip: FutureSelfClip = .why, duckOthers: Bool = false) {
        let url = FutureSelfStore.url(for: clip)
        guard FileManager.default.fileExists(atPath: url.path) else { return }
        stop()
        do {
            let session = AVAudioSession.sharedInstance()
            try session.setCategory(.playback, mode: .default)
            try session.setActive(true)
            let player = try AVAudioPlayer(contentsOf: url)
            player.volume = 0.8
            player.prepareToPlay()
            player.play()
            self.player = player
            FutureSelfStore.notePlayed()
        } catch {
            self.player = nil
        }
        _ = duckOthers
    }

    func stop() {
        player?.stop()
        player = nil
    }
}
