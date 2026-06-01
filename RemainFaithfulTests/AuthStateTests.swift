import XCTest
import Combine
@testable import RemainFaithful

final class AuthStateTests: XCTestCase {

    var sut: AuthState!
    var testKeychain: KeychainHelper!
    var cancellables = Set<AnyCancellable>()

    override func setUp() {
        super.setUp()
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
