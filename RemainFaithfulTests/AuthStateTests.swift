import XCTest
import Combine
@testable import RemainFaithful

final class AuthStateTests: XCTestCase {

    var sut: AuthState!
    var testKeychain: KeychainHelper!
    var cancellables = Set<AnyCancellable>()

    private var savedPrimaryGroupID: Any?

    override func setUp() {
        super.setUp()
        savedPrimaryGroupID = UserDefaults.standard.object(forKey: AuthState.primaryGroupIDKey)
        testKeychain = KeychainHelper(service: "com.remainfaithful.test")
        testKeychain.delete("authToken")
        testKeychain.delete("currentUser")
        sut = AuthState(keychain: testKeychain)
    }

    override func tearDown() {
        testKeychain.delete("authToken")
        testKeychain.delete("currentUser")
        cancellables = []
        sut = nil
        if let saved = savedPrimaryGroupID {
            UserDefaults.standard.set(saved, forKey: AuthState.primaryGroupIDKey)
        } else {
            UserDefaults.standard.removeObject(forKey: AuthState.primaryGroupIDKey)
        }
        super.tearDown()
    }

    // MARK: - clearSession

    func testClearSession_setsIsAuthenticatedFalse() {
        testKeychain.set("sometoken.payload.sig", for: "authToken")
        sut = AuthState(keychain: testKeychain)
        XCTAssertTrue(sut.isAuthenticated, "precondition: should be authenticated")

        let exp = expectation(description: "isAuthenticated → false")
        sut.$isAuthenticated
            .dropFirst()
            .first(where: { !$0 })
            .sink { _ in exp.fulfill() }
            .store(in: &cancellables)

        sut.clearSession()
        wait(for: [exp], timeout: 1)
        XCTAssertFalse(sut.isAuthenticated)
    }

    func testClearSession_setsCurrentUserNil() {
        let user = RemoteUser(id: 1, name: "Bob", email: "bob@test.com", createdAt: nil)
        sut.setSession(token: "tok.pay.sig", user: user)

        let exp = expectation(description: "currentUser → nil")
        sut.$currentUser
            .dropFirst()
            .first(where: { $0 == nil })
            .sink { _ in exp.fulfill() }
            .store(in: &cancellables)

        sut.clearSession()
        wait(for: [exp], timeout: 1)
        XCTAssertNil(sut.currentUser)
    }

    func testClearSession_clearsStoredGroupSelection() {
        UserDefaults.standard.set(42, forKey: AuthState.primaryGroupIDKey)
        sut.clearSession()
        XCTAssertEqual(UserDefaults.standard.integer(forKey: AuthState.primaryGroupIDKey), 0,
                       "a signed-out phone must not keep the previous account's group")
    }

    // MARK: - setSession and stored group selection

    func testSetSession_differentAccountClearsStoredGroupSelection() {
        signIn(id: 1)
        UserDefaults.standard.set(42, forKey: AuthState.primaryGroupIDKey)

        sut.setSession(token: "tok.pay.sig",
                       user: RemoteUser(id: 2, name: "Other", email: "o@test.com", createdAt: nil))
        XCTAssertEqual(UserDefaults.standard.integer(forKey: AuthState.primaryGroupIDKey), 0)
    }

    func testSetSession_sameAccountKeepsStoredGroupSelection() {
        signIn(id: 1)
        UserDefaults.standard.set(42, forKey: AuthState.primaryGroupIDKey)

        sut.setSession(token: "tok2.pay.sig",
                       user: RemoteUser(id: 1, name: "Bob", email: "bob@test.com", createdAt: nil))
        XCTAssertEqual(UserDefaults.standard.integer(forKey: AuthState.primaryGroupIDKey), 42)
    }

    // MARK: - tokenExpiresAt

    func testTokenExpiresAt_nilForMalformedJWT() {
        testKeychain.set("notajwt", for: "authToken")
        sut = AuthState(keychain: testKeychain)
        XCTAssertNil(sut.tokenExpiresAt)
    }

    func testTokenExpiresAt_nilForTwoPartJWT() {
        testKeychain.set("header.payload", for: "authToken")
        sut = AuthState(keychain: testKeychain)
        XCTAssertNil(sut.tokenExpiresAt)
    }

    func testTokenExpiresAt_validDateForWellFormedJWT() {
        let futureExp = Int(Date(timeIntervalSinceNow: 3600).timeIntervalSince1970)
        let jwt = makeJWT(exp: futureExp)

        testKeychain.set(jwt, for: "authToken")
        sut = AuthState(keychain: testKeychain)

        XCTAssertNotNil(sut.tokenExpiresAt)
        XCTAssertEqual(
            sut.tokenExpiresAt!.timeIntervalSince1970,
            Double(futureExp),
            accuracy: 1.0
        )
    }

    // MARK: - setSession

    func testSetSession_setsIsAuthenticatedTrue() {
        XCTAssertFalse(sut.isAuthenticated, "precondition: fresh state is unauthenticated")

        let exp = expectation(description: "isAuthenticated → true")
        sut.$isAuthenticated
            .dropFirst()
            .first(where: { $0 })
            .sink { _ in exp.fulfill() }
            .store(in: &cancellables)

        sut.setSession(token: "tok.pay.sig",
                       user: RemoteUser(id: 5, name: "Alice", email: "a@b.com", createdAt: nil))
        wait(for: [exp], timeout: 1)
        XCTAssertTrue(sut.isAuthenticated)
    }

    func testSetSession_populatesCurrentUserCorrectly() {
        let exp = expectation(description: "currentUser populated")
        sut.$currentUser
            .compactMap { $0 }
            .first()
            .sink { _ in exp.fulfill() }
            .store(in: &cancellables)

        sut.setSession(token: "tok.pay.sig",
                       user: RemoteUser(id: 7, name: "Alice", email: "alice@example.com", createdAt: nil))
        wait(for: [exp], timeout: 1)

        XCTAssertEqual(sut.currentUser?.id, 7)
        XCTAssertEqual(sut.currentUser?.name, "Alice")
        XCTAssertEqual(sut.currentUser?.email, "alice@example.com")
    }

    // MARK: - Helpers

    /// Signs in and waits for the async currentUser publish to land.
    private func signIn(id: Int) {
        let exp = expectation(description: "signed in as \(id)")
        sut.$currentUser
            .compactMap { $0 }
            .first(where: { $0.id == id })
            .sink { _ in exp.fulfill() }
            .store(in: &cancellables)
        sut.setSession(token: "tok.pay.sig",
                       user: RemoteUser(id: id, name: "Bob", email: "bob@test.com", createdAt: nil))
        wait(for: [exp], timeout: 1)
    }

    private func makeJWT(exp: Int) -> String {
        let header  = base64url(Data("{\"alg\":\"HS256\",\"typ\":\"JWT\"}".utf8))
        let payload = base64url(Data("{\"user_id\":1,\"email\":\"t@t.com\",\"exp\":\(exp)}".utf8))
        return "\(header).\(payload).fakesig"
    }

    private func base64url(_ data: Data) -> String {
        data.base64EncodedString()
            .replacingOccurrences(of: "+", with: "-")
            .replacingOccurrences(of: "/", with: "_")
            .trimmingCharacters(in: .init(charactersIn: "="))
    }
}

// MARK: - Group tab selection

final class GroupSelectionTests: XCTestCase {

    func testKeepsStoredGroupWhenStillAMember() {
        XCTAssertEqual(GroupSelection.resolvePrimaryGroupID(stored: 7, memberGroupIDs: [3, 7]), 7)
    }

    func testStaleGroupFallsBackToFirstGroup() {
        // Group 99 belonged to the previous account on this phone.
        XCTAssertEqual(GroupSelection.resolvePrimaryGroupID(stored: 99, memberGroupIDs: [3, 7]), 3)
    }

    func testStaleGroupWithNoGroupsShowsCreateState() {
        XCTAssertEqual(GroupSelection.resolvePrimaryGroupID(stored: 99, memberGroupIDs: []), 0)
    }

    func testNoStoredGroupPicksFirst() {
        XCTAssertEqual(GroupSelection.resolvePrimaryGroupID(stored: 0, memberGroupIDs: [5]), 5)
    }

    func testNoStoredGroupAndNoGroups() {
        XCTAssertEqual(GroupSelection.resolvePrimaryGroupID(stored: 0, memberGroupIDs: []), 0)
    }

    func testInviteLinkUsesLiveSite() {
        XCTAssertEqual(GroupSelection.inviteLink(groupID: 12), "www.remainfaithful.com/join/12")
        XCTAssertNotNil(URL(string: "https://" + GroupSelection.inviteLink(groupID: 12)))
    }
}
