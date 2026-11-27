import XCTest
@testable import RemainFaithful

final class FutureSelfStoreTests: XCTestCase {

    private var folder: URL!

    override func setUp() {
        super.setUp()
        folder = FileManager.default.temporaryDirectory.appendingPathComponent("rf-voice-\(UUID().uuidString)", isDirectory: true)
        try? FileManager.default.createDirectory(at: folder, withIntermediateDirectories: true)
    }

    override func tearDown() {
        try? FileManager.default.removeItem(at: folder)
        folder = nil
        super.tearDown()
    }

    func testDurationCap() {
        XCTAssertFalse(FutureSelfRules.keep(duration: 0))
        XCTAssertFalse(FutureSelfRules.keep(duration: 0.4))
        XCTAssertTrue(FutureSelfRules.keep(duration: 1))
        XCTAssertTrue(FutureSelfRules.keep(duration: 45))
        XCTAssertFalse(FutureSelfRules.keep(duration: 46))
        XCTAssertFalse(FutureSelfRules.skipSendsPartnerNotice)
        XCTAssertFalse(FutureSelfPlayback.shouldAutoPlay(hasClip: false, alreadyPlayedThisOpen: false))
        XCTAssertTrue(FutureSelfPlayback.shouldAutoPlay(hasClip: true, alreadyPlayedThisOpen: false))
        XCTAssertFalse(FutureSelfPlayback.shouldAutoPlay(hasClip: true, alreadyPlayedThisOpen: true))
    }

    func testReplaceRemovesThePreviousFileAndWipeClearsTheFolder() throws {
        let first = folder.appendingPathComponent("one.m4a")
        let second = folder.appendingPathComponent("two.m4a")
        try Data("first".utf8).write(to: first)
        XCTAssertTrue(try FutureSelfStore.replace(.why, with: first, duration: 12, base: folder))
        XCTAssertTrue(FutureSelfStore.exists(.why, base: folder))
        XCTAssertFalse(FileManager.default.fileExists(atPath: first.path))

        try Data("second".utf8).write(to: second)
        XCTAssertTrue(try FutureSelfStore.replace(.why, with: second, duration: 20, base: folder))
        let stored = try Data(contentsOf: FutureSelfStore.url(for: .why, base: folder))
        XCTAssertEqual(stored, Data("second".utf8))
        XCTAssertEqual(FutureSelfStore.count(base: folder), 1)

        let tooLong = folder.appendingPathComponent("long.m4a")
        try Data("nope".utf8).write(to: tooLong)
        XCTAssertFalse(try FutureSelfStore.replace(.remember, with: tooLong, duration: 46, base: folder))
        XCTAssertFalse(FutureSelfStore.exists(.remember, base: folder))

        FutureSelfStore.wipe(base: folder)
        XCTAssertFalse(FileManager.default.fileExists(atPath: FutureSelfStore.directory(base: folder).path))
        XCTAssertEqual(FutureSelfStore.count(base: folder), 0)
    }

    func testEmptyShieldDoesNotRequireAClip() {
        XCTAssertFalse(FutureSelfStore.exists(.why, base: folder))
        XCTAssertFalse(FutureSelfPlayback.shouldAutoPlay(hasClip: FutureSelfStore.exists(.why, base: folder), alreadyPlayedThisOpen: false))
    }
}
