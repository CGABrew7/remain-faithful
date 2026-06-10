import Foundation

// MARK: - Base URL
//
// Simulator + docker-compose on the same Mac: http://localhost:8080 (default)
// Real device + local network:                http://<YOUR_MAC_IP>:8080
// Production:                                 https://remain-faithful-api.fly.dev
//
// Override at runtime by setting the API_BASE_URL environment variable in the
// scheme's Run arguments, or change the constant below.

private let apiBase = ProcessInfo.processInfo.environment["API_BASE_URL"] ?? "https://remain-faithful-api.fly.dev"

// MARK: - Remote response types (mirror the Go API)

struct AuthResponse: Decodable {
    let token: String
    let user: RemoteUser
}

struct RemoteUser: Decodable {
    let id: Int
    let name: String
    let email: String
    let createdAt: String?
    enum CodingKeys: String, CodingKey {
        case id, name, email
        case createdAt = "created_at"
    }
}

struct RemoteEvent: Decodable {
    let id: Int
    let userId: Int
    let category: String
    let severity: String
    let summary: String
    let timestamp: String
    enum CodingKeys: String, CodingKey {
        case id, category, severity, summary, timestamp
        case userId = "user_id"
    }
}

struct RemoteAlert: Decodable, Identifiable {
    let id: Int
    let eventId: Int
    let seen: Bool
    let createdAt: String
    let event: RemoteEvent
    enum CodingKeys: String, CodingKey {
        case id, seen, event
        case eventId = "event_id"
        case createdAt = "created_at"
    }
}

struct RemoteRelationship: Decodable {
    let id: Int
    let userId: Int
    let partnerId: Int
    let type: String
    let status: String
    let createdAt: String
    let isPrimary: Bool
    let partner: RemoteUser

    enum CodingKeys: String, CodingKey {
        case id, type, status, partner
        case userId    = "user_id"
        case partnerId = "partner_id"
        case createdAt = "created_at"
        case isPrimary = "is_primary"
    }

    init(from decoder: Decoder) throws {
        let c  = try decoder.container(keyedBy: CodingKeys.self)
        id        = try c.decode(Int.self,        forKey: .id)
        userId    = try c.decode(Int.self,        forKey: .userId)
        partnerId = try c.decode(Int.self,        forKey: .partnerId)
        type      = try c.decode(String.self,     forKey: .type)
        status    = try c.decode(String.self,     forKey: .status)
        createdAt = try c.decode(String.self,     forKey: .createdAt)
        partner   = try c.decode(RemoteUser.self, forKey: .partner)
        isPrimary = (try? c.decode(Bool.self,     forKey: .isPrimary)) ?? false
    }
}

struct RemoteGroup: Decodable {
    let id: Int
    let name: String
    let createdAt: String
    let members: [RemoteGroupMember]?
    enum CodingKeys: String, CodingKey {
        case id, name, members
        case createdAt = "created_at"
    }
}

struct RemoteGroupMember: Decodable {
    let userId: Int
    let role: String
    let joinedAt: String
    let user: RemoteUser
    let flagsLast30: Int
    let streakDays: Int
    enum CodingKeys: String, CodingKey {
        case role, user
        case userId      = "user_id"
        case joinedAt    = "joined_at"
        case flagsLast30 = "flags_last_30"
        case streakDays  = "streak_days"
    }
}

// MARK: - Errors

enum APIError: LocalizedError {
    case badURL
    case unauthenticated
    case server(String)
    case decoding(Error)
    case network(Error)

    var errorDescription: String? {
        switch self {
        case .badURL:          return "Invalid URL"
        case .unauthenticated: return "Not signed in — please log in again"
        case .server(let m):   return m
        case .decoding(let e): return "Unexpected response format: \(e.localizedDescription)"
        case .network(let e):  return "Network error: \(e.localizedDescription)"
        }
    }
}

// MARK: - APIClient

extension APIClient: EventSending {}

final class APIClient {
    static let shared = APIClient()
    private init() {}

    private let session = URLSession.shared

    /// JWT stored in Keychain. Reads through AuthState so the token is always current.
    var token: String? { AuthState.shared.token }

    var isAuthenticated: Bool { token != nil }

    // MARK: - Auth

    func register(name: String, email: String, password: String) async throws -> RemoteUser {
        try await post("/auth/register",
                       body: ["name": name, "email": email, "password": password],
                       auth: false)
    }

    func login(email: String, password: String) async throws -> AuthResponse {
        let resp: AuthResponse = try await post("/auth/login",
                                                body: ["email": email, "password": password],
                                                auth: false)
        AuthState.shared.setSession(token: resp.token, user: resp.user)
        return resp
    }

    func logout() {
        AuthState.shared.clearSession()
    }

    func refreshTokenIfNeeded() async throws {
        guard let expiry = AuthState.shared.tokenExpiresAt,
              expiry.timeIntervalSinceNow < 3600 else { return }
        let resp: [String: String] = try await post("/auth/refresh", body: [String: String]())
        if let newToken = resp["token"] {
            KeychainHelper.shared.set(newToken, for: "authToken")
        }
    }

    // MARK: - Users

    func me() async throws -> RemoteUser {
        try await get("/users/me")
    }

    func updateMe(name: String, email: String) async throws -> RemoteUser {
        try await put("/users/me", body: ["name": name, "email": email])
    }

    // MARK: - Events

    func listEvents() async throws -> [RemoteEvent] {
        try await get("/events")
    }

    func createEvent(category: String, severity: String, summary: String,
                     timestamp: String? = nil) async throws -> RemoteEvent {
        var body: [String: String] = [
            "category": category,
            "severity": severity,
            "summary":  summary,
        ]
        if let ts = timestamp { body["timestamp"] = ts }
        return try await post("/events", body: body)
    }

    // MARK: - Alerts

    func listAlerts() async throws -> [RemoteAlert] {
        try await get("/alerts")
    }

    func alertUnreadCount() async throws -> Int {
        let resp: [String: Int] = try await get("/alerts/count")
        return resp["unseen"] ?? 0
    }

    func markAlertsSeen() async throws {
        try await postVoid("/alerts/mark-seen", body: [String: String]())
    }

    func markAlertDiscussed(alertID: Int) async throws {
        try await patchVoid("/alerts/\(alertID)/discussed")
    }

    // MARK: - Relationships

    func deleteRelationship(id: Int) async throws {
        try await deleteVoid("/relationships/\(id)")
    }

    // MARK: - Donations

    func createCheckoutSession(amountDollars: Int, monthly: Bool) async throws -> URL {
        let body: [String: Any] = ["amount_dollars": amountDollars, "monthly": monthly]
        let resp: [String: String] = try await postAny(
            "/donations/create-checkout-session", body: body, auth: true)
        guard let urlStr = resp["url"], let url = URL(string: urlStr) else {
            throw APIError.server("Invalid checkout URL")
        }
        return url
    }

    // MARK: - Heartbeat

    func heartbeat(screen: String) async throws {
        try await postVoid("/heartbeat", body: ["screen": screen])
    }

    // MARK: - Groups

    func listMyGroups() async throws -> [RemoteGroup] {
        try await get("/groups")
    }

    func leaveGroup(groupID: Int) async throws {
        try await deleteVoid("/groups/\(groupID)/members/me")
    }

    func leaveAllGroups() async throws {
        try await postVoid("/groups/leave-all", body: [String: String]())
    }

    func getGroup(id: Int) async throws -> RemoteGroup {
        try await get("/groups/\(id)")
    }

    func createGroup(name: String, covenant: String = "") async throws -> RemoteGroup {
        try await post("/groups", body: ["name": name, "covenant": covenant])
    }

    func inviteMember(groupID: Int, email: String) async throws {
        try await postVoid("/groups/\(groupID)/invite", body: ["user_email": email])
    }

    // MARK: - Social auth

    func appleSignIn(identityToken: String, authorizationCode: String,
                     firstName: String, lastName: String) async throws -> AuthResponse {
        let body: [String: Any] = [
            "identity_token":     identityToken,
            "authorization_code": authorizationCode,
            "name": ["firstName": firstName, "lastName": lastName],
        ]
        return try await postAny("/auth/apple", body: body, auth: false)
    }

    func googleSignIn(idToken: String) async throws -> AuthResponse {
        try await post("/auth/google", body: ["id_token": idToken], auth: false)
    }

    func forgotPassword(email: String) async throws {
        try await postVoid("/auth/forgot-password", body: ["email": email], auth: false)
    }

    func resetPassword(token: String, password: String) async throws {
        try await postVoid("/auth/reset-password",
                           body: ["token": token, "password": password], auth: false)
    }

    // MARK: - Push notifications

    func registerDeviceToken(_ token: String) async throws {
        try await postVoid("/users/device-token", body: ["token": token, "platform": "ios"])
    }

    func sendPanicAlert() async throws {
        try await postVoid("/panic", body: [String: String]())
    }

    // MARK: - Relationships

    func setPrimaryPartner(relationshipID: Int) async throws {
        try await putVoid("/relationships/\(relationshipID)/primary", body: [String: String]())
    }

    func createRelationship(partnerEmail: String, type: String = "partner") async throws {
        try await postVoid("/relationships",
                           body: ["partner_email": partnerEmail, "type": type])
    }

    /// Invite a partner by email — creates the relationship if they already have
    /// an account, otherwise sends them an email with a sign-up deep link.
    func invitePartner(email: String) async throws {
        try await postVoid("/relationships/invite", body: ["email": email])
    }

    func listRelationships() async throws -> [RemoteRelationship] {
        try await get("/relationships")
    }

    /// Send an email-based group invite (invitee may not have an account yet).
    func groupEmailInvite(groupID: Int, email: String) async throws {
        try await postVoid("/groups/\(groupID)/email-invite", body: ["email": email])
    }

    // MARK: - Primitives

    private func get<T: Decodable>(_ path: String) async throws -> T {
        try? await refreshTokenIfNeeded()
        var req = try makeRequest(path)
        try attachToken(&req)
        return try await send(req)
    }

    private func post<Body: Encodable, T: Decodable>(_ path: String,
                                                      body: Body,
                                                      auth: Bool = true) async throws -> T {
        try? await refreshTokenIfNeeded()
        var req = try makeRequest(path, method: "POST")
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpBody = try JSONEncoder().encode(body)
        if auth { try attachToken(&req) }
        return try await send(req)
    }

    /// Fire-and-forget POST for endpoints where we don't need the response body.
    private func postVoid<Body: Encodable>(_ path: String, body: Body,
                                           auth: Bool = true) async throws {
        var req = try makeRequest(path, method: "POST")
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpBody = try JSONEncoder().encode(body)
        if auth { try attachToken(&req) }
        let (data, resp) = try await session.data(for: req)
        if let http = resp as? HTTPURLResponse, !(200...299).contains(http.statusCode) {
            let msg = (try? JSONDecoder().decode([String: String].self, from: data))?["error"]
            throw APIError.server(msg ?? "HTTP \(http.statusCode)")
        }
    }

    /// PUT that returns a decoded response body.
    private func put<Body: Encodable, T: Decodable>(_ path: String,
                                                     body: Body) async throws -> T {
        try? await refreshTokenIfNeeded()
        var req = try makeRequest(path, method: "PUT")
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpBody = try JSONEncoder().encode(body)
        try attachToken(&req)
        return try await send(req)
    }

    /// Fire-and-forget DELETE with no request body.
    private func deleteVoid(_ path: String) async throws {
        var req = try makeRequest(path, method: "DELETE")
        try attachToken(&req)
        let (data, resp) = try await session.data(for: req)
        if let http = resp as? HTTPURLResponse, !(200...299).contains(http.statusCode) {
            let msg = (try? JSONDecoder().decode([String: String].self, from: data))?["error"]
            throw APIError.server(msg ?? "HTTP \(http.statusCode)")
        }
    }

    /// Fire-and-forget PATCH for endpoints with no request body.
    private func patchVoid(_ path: String) async throws {
        var req = try makeRequest(path, method: "PATCH")
        try attachToken(&req)
        let (data, resp) = try await session.data(for: req)
        if let http = resp as? HTTPURLResponse, !(200...299).contains(http.statusCode) {
            let msg = (try? JSONDecoder().decode([String: String].self, from: data))?["error"]
            throw APIError.server(msg ?? "HTTP \(http.statusCode)")
        }
    }

    /// Fire-and-forget PUT for endpoints where we don't need the response body.
    private func putVoid<Body: Encodable>(_ path: String, body: Body,
                                          auth: Bool = true) async throws {
        var req = try makeRequest(path, method: "PUT")
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpBody = try JSONEncoder().encode(body)
        if auth { try attachToken(&req) }
        let (data, resp) = try await session.data(for: req)
        if let http = resp as? HTTPURLResponse, !(200...299).contains(http.statusCode) {
            let msg = (try? JSONDecoder().decode([String: String].self, from: data))?["error"]
            throw APIError.server(msg ?? "HTTP \(http.statusCode)")
        }
    }

    /// POST where the body is `[String: Any]` (not Encodable) — needed for nested dicts.
    private func postAny<T: Decodable>(_ path: String, body: [String: Any],
                                        auth: Bool = true) async throws -> T {
        var req = try makeRequest(path, method: "POST")
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpBody = try JSONSerialization.data(withJSONObject: body)
        if auth { try attachToken(&req) }
        return try await send(req)
    }

    private func makeRequest(_ path: String, method: String = "GET") throws -> URLRequest {
        guard let url = URL(string: apiBase + path) else { throw APIError.badURL }
        var req = URLRequest(url: url, timeoutInterval: 15)
        req.httpMethod = method
        return req
    }

    private func attachToken(_ req: inout URLRequest) throws {
        guard let t = token else { throw APIError.unauthenticated }
        req.setValue("Bearer \(t)", forHTTPHeaderField: "Authorization")
    }

    private func send<T: Decodable>(_ req: URLRequest) async throws -> T {
        do {
            let (data, resp) = try await session.data(for: req)
            if let http = resp as? HTTPURLResponse, !(200...299).contains(http.statusCode) {
                let msg = (try? JSONDecoder().decode([String: String].self, from: data))?["error"]
                throw APIError.server(msg ?? "HTTP \(http.statusCode)")
            }
            return try JSONDecoder().decode(T.self, from: data)
        } catch let e as APIError { throw e }
        catch let e as DecodingError { throw APIError.decoding(e) }
        catch { throw APIError.network(error) }
    }
}
