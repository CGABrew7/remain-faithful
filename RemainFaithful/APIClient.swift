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
    var emailVerified: Bool? = nil
    var pendingEmail: String? = nil
    enum CodingKeys: String, CodingKey {
        case id, name, email
        case createdAt = "created_at"
        case emailVerified = "email_verified"
        case pendingEmail = "pending_email"
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

struct UserStats: Decodable {
    let streakDays: Int
    let bestStreak: Int
    let week: [StreakDayInput]
    /// `before_account`, `clean`, or `flagged`. Empty on older servers.
    let weekStates: [StreakDayInput]
    enum CodingKeys: String, CodingKey {
        case streakDays = "streak_days"
        case bestStreak = "best_streak"
        case week
        case weekStates = "week_states"
    }

    init(from decoder: Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        streakDays = try c.decode(Int.self, forKey: .streakDays)
        bestStreak = try c.decode(Int.self, forKey: .bestStreak)
        week = try c.decode([StreakDayInput].self, forKey: .week)
        weekStates = (try? c.decode([StreakDayInput].self, forKey: .weekStates)) ?? []
    }
}

struct RemoteAlert: Decodable, Identifiable {
    let id: Int
    let eventId: Int?
    let seen: Bool
    let createdAt: String
    let event: RemoteEvent?
    /// "content" or "protection" when the server sends a kind. Optional today.
    let kind: String?
    let protectionType: String?
    let summary: String?

    enum CodingKeys: String, CodingKey {
        case id, seen, event, kind, summary
        case eventId = "event_id"
        case createdAt = "created_at"
        case protectionType = "protection_type"
        case alertType = "type"
        case contractAlertType = "alert_type"
    }

    init(from decoder: Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        id = try c.decode(Int.self, forKey: .id)
        eventId = try c.decodeIfPresent(Int.self, forKey: .eventId)
        seen = (try? c.decode(Bool.self, forKey: .seen)) ?? false
        createdAt = (try? c.decode(String.self, forKey: .createdAt)) ?? ""
        event = try c.decodeIfPresent(RemoteEvent.self, forKey: .event)
        kind = try c.decodeIfPresent(String.self, forKey: .kind)
        let named = [
            try c.decodeIfPresent(String.self, forKey: .contractAlertType),
            try c.decodeIfPresent(String.self, forKey: .protectionType),
            try c.decodeIfPresent(String.self, forKey: .alertType),
        ]
        protectionType = named.compactMap { $0 }.first { !$0.isEmpty }
        summary = try c.decodeIfPresent(String.self, forKey: .summary)
    }
}

struct PendingPartnerInvite: Decodable, Identifiable {
    let id: Int
    let inviterId: Int
    let inviterName: String
    let createdAt: String
    let type: String
    enum CodingKeys: String, CodingKey {
        case id, type
        case inviterId = "inviter_id"
        case inviterName = "inviter_name"
        case createdAt = "created_at"
    }
}

struct PendingGroupInvite: Decodable, Identifiable {
    let id: Int
    let groupId: Int
    let groupName: String
    let inviterId: Int
    let inviterName: String
    let createdAt: String
    let type: String
    enum CodingKeys: String, CodingKey {
        case id, type
        case groupId = "group_id"
        case groupName = "group_name"
        case inviterId = "inviter_id"
        case inviterName = "inviter_name"
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
    let pinSet: Bool      // true if a protection PIN exists for this relationship
    let partner: RemoteUser

    enum CodingKeys: String, CodingKey {
        case id, type, status, partner
        case userId    = "user_id"
        case partnerId = "partner_id"
        case createdAt = "created_at"
        case isPrimary = "is_primary"
        case pinSet    = "pin_set"
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
        pinSet    = (try? c.decode(Bool.self,     forKey: .pinSet))    ?? false
    }
}

struct RemoteGroup: Decodable {
    let id: Int
    let name: String
    let createdAt: String
    let covenant: String?
    let members: [RemoteGroupMember]?
    enum CodingKeys: String, CodingKey {
        case id, name, members, covenant
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
    case server(String, code: String?)
    case decoding(Error)
    case network(Error)

    var code: String? {
        if case .server(_, let code) = self { return code }
        return nil
    }

    var errorDescription: String? {
        switch self {
        case .badURL:          return "Invalid URL"
        case .unauthenticated: return "Not signed in. Please log in again."
        case .server(let m, _): return m
        case .decoding(let e): return "Unexpected response format: \(e.localizedDescription)"
        case .network(let e):  return "Network error: \(e.localizedDescription)"
        }
    }

    static func fromServer(data: Data, status: Int) -> APIError {
        struct Body: Decodable {
            let error: String?
            let code: String?
        }
        let body = try? JSONDecoder().decode(Body.self, from: data)
        return .server(body?.error ?? "HTTP \(status)", code: body?.code)
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
            AuthState.shared.updateToken(newToken)
        }
    }

    /// Writes the API host into the app group so extensions can upload and heartbeat.
    func publishSharedConfig() {
        UserDefaults(suiteName: "group.com.remainfaithful.app")?
            .set(apiBase, forKey: "apiBaseURL")
    }

    // MARK: - Users

    func me() async throws -> RemoteUser {
        try await get("/users/me")
    }

    func updateMe(name: String, email: String) async throws -> RemoteUser {
        try await put("/users/me", body: ["name": name, "email": email])
    }

    // MARK: - Events

    func getUserStats() async throws -> UserStats {
        try await get("/users/me/stats")
    }

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

    /// POST /alerts/{id}/seen — mark one alert when the person opens it.
    /// A 404 means the route is not deployed yet; the caller ignores that.
    func markAlertSeen(alertID: Int) async throws {
        try await postVoid("/alerts/\(alertID)/seen", body: [String: String]())
    }

    func markAlertDiscussed(alertID: Int) async throws {
        try await patchVoid("/alerts/\(alertID)/discussed")
    }

    // MARK: - Protection (Feature 1 & 2)
    //
    // These endpoints are on the backend at /protection/* and /relationships/{id}/pin.
    // Expected request/response shapes are documented per method below.

    /// Send a protection-change alert to the accountability partner via APNs.
    /// POST /protection/alerts  body: {"type": "...", "detail": "..."}
    /// The server pushes this to the primary partner with notification_type: "PROTECTION_ALERT".
    func sendProtectionAlert(type: String, detail: String) async throws {
        guard isAuthenticated else { return }
        try await postVoid("/protection/alerts", body: ["type": type, "detail": detail])
    }

    /// Returns true if any accountability partner has set a protection PIN for the current user.
    /// GET /protection/pin  → {"is_set": bool}
    func getPartnerPINStatus() async throws -> Bool {
        let resp: [String: Bool] = try await get("/protection/pin")
        return resp["is_set"] ?? false
    }

    /// Verify a protection PIN entered by the monitored user.
    /// POST /protection/pin/verify  body: {"pin": "1234"}
    /// → success, or reason no_pin / wrong_pin. 429 reason is rate_limited.
    /// A 429 still carries the reason, so this does not go through post(), which
    /// drops the body of any non-2xx response.
    func verifyPartnerPIN(pin: String) async throws -> PINCheck {
        struct Body: Decodable {
            let success: Bool
            let reason: String?
        }
        do { try await refreshTokenIfNeeded() } catch { print("[APIClient] token refresh: \(error)") }
        var req = try makeRequest("/protection/pin/verify", method: "POST")
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpBody = try JSONEncoder().encode(["pin": pin])
        try attachToken(&req)
        do {
            let (data, resp) = try await session.data(for: req)
            if let body = try? JSONDecoder().decode(Body.self, from: data) {
                return PINCheck.from(success: body.success, reason: body.reason)
            }
            if let http = resp as? HTTPURLResponse, !(200...299).contains(http.statusCode) {
                throw APIError.fromServer(data: data, status: http.statusCode)
            }
            throw APIError.decoding(DecodingError.dataCorrupted(
                .init(codingPath: [], debugDescription: "PIN response")
            ))
        } catch let e as APIError {
            throw e
        } catch {
            throw APIError.network(error)
        }
    }

    /// Set (or overwrite) the protection PIN for a monitored user — called from the partner's device.
    /// POST /relationships/{id}/pin  body: {"pin": "1234"}
    func setRelationshipPIN(relationshipID: Int, pin: String) async throws {
        try await postVoid("/relationships/\(relationshipID)/pin", body: ["pin": pin])
    }

    /// Remove the protection PIN for a relationship — called from the partner's device.
    /// DELETE /relationships/{id}/pin
    func removeRelationshipPIN(relationshipID: Int) async throws {
        try await deleteVoid("/relationships/\(relationshipID)/pin")
    }

    // MARK: - Relationships

    func deleteRelationship(id: Int) async throws {
        try await deleteVoid("/relationships/\(id)")
    }

    // MARK: - Heartbeat

    func heartbeat(screen: String) async throws {
        try await postVoid("/heartbeat", body: ["screen": screen])
    }

    /// Counts only: Screen Time yes or no, how many apps, the app version, and support taps.
    func sendPhoneSignals(screenTimeGranted: Bool, monitoredAppCount: Int, supportTaps: Int) async {
        guard isAuthenticated else { return }
        struct Body: Encodable {
            let screenTimeGranted: Bool
            let monitoredAppCount: Int
            let appVersion: String
            let supportTaps: Int
            enum CodingKeys: String, CodingKey {
                case screenTimeGranted = "screen_time_granted"
                case monitoredAppCount = "monitored_app_count"
                case appVersion = "app_version"
                case supportTaps = "support_taps"
            }
        }
        let version = Bundle.main.object(forInfoDictionaryKey: "CFBundleShortVersionString") as? String ?? ""
        do {
            try await postVoid("/users/me/signals", body: Body(
                screenTimeGranted: screenTimeGranted,
                monitoredAppCount: monitoredAppCount,
                appVersion: version,
                supportTaps: supportTaps
            ))
        } catch {
            // A server that does not know this yet leaves the counts on the phone.
        }
    }

    /// POST /auth/resend-verification. Returns a short line for the screen.
    func resendVerification() async throws -> String {
        struct Body: Decodable { let status: String? }
        do {
            let body: Body = try await post("/auth/resend-verification", body: [String: String]())
            if body.status == "already_verified" {
                AuthState.shared.markEmailVerified()
            }
            return MailboxNotice.resendStatus(body.status)
        } catch {
            if let api = error as? APIError, api.code == "email_required" {
                AuthState.shared.noteEmailRequired()
            }
            throw error
        }
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

    func deleteAccount() async throws {
        try await deleteVoid("/users/me")
    }

    func getGroup(id: Int) async throws -> RemoteGroup {
        try await get("/groups/\(id)")
    }

    func createGroup(name: String, covenant: String = "") async throws -> RemoteGroup {
        try await post("/groups", body: ["name": name, "covenant": covenant])
    }

    /// PATCH /groups/{id}. Sends only the fields you pass.
    func updateGroup(id: Int, name: String? = nil, covenant: String? = nil) async throws {
        do { try await refreshTokenIfNeeded() } catch { print("[APIClient] token refresh: \(error)") }
        var fields: [String: String] = [:]
        if let name { fields["name"] = name }
        if let covenant { fields["covenant"] = covenant }
        try await patchJSON("/groups/\(id)", body: JSONStringMap(values: fields))
    }

    /// POST /groups/{id}/encouragement. The server writes the message.
    func sendEncouragement(groupID: Int, userID: Int) async throws {
        do { try await refreshTokenIfNeeded() } catch { print("[APIClient] token refresh: \(error)") }
        try await postVoid("/groups/\(groupID)/encouragement", body: ["user_id": userID])
    }

    /// POST /feedback. Name and email come from the signed-in account.
    func sendFeedback(_ fields: [String: String]) async throws {
        do { try await refreshTokenIfNeeded() } catch { print("[APIClient] token refresh: \(error)") }
        try await postVoid("/feedback", body: JSONStringMap(values: fields))
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
        #if DEBUG
        let env = "sandbox"
        #else
        let env = "production"
        #endif
        try await postVoid("/users/device-token", body: ["token": token, "platform": "ios", "environment": env])
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

    /// Invite a partner by email.
    /// Response status: "pending" / "connected" (waiting to accept),
    /// "invited" (no account yet), or "accepted".
    func invitePartner(email: String) async throws -> PartnerInviteResult {
        struct Body: Decodable { let status: String? }
        let body: Body = try await post("/relationships/invite", body: ["email": email])
        return PartnerInviteResult.from(status: body.status ?? "")
    }

    /// GET /relationships/invites/pending. `id` is the invite row, not relationships.id.
    func listPendingPartnerInvites() async throws -> [PendingPartnerInvite] {
        try await get("/relationships/invites/pending")
    }

    /// POST /relationships/invites/{id}/accept
    func acceptPartnerInvite(id: Int) async throws {
        try await postVoid("/relationships/invites/\(id)/accept", body: [String: String]())
    }

    /// POST /relationships/invites/{id}/decline
    func declinePartnerInvite(id: Int) async throws {
        try await postVoid("/relationships/invites/\(id)/decline", body: [String: String]())
    }

    /// GET /groups/invites/pending
    func listPendingGroupInvites() async throws -> [PendingGroupInvite] {
        try await get("/groups/invites/pending")
    }

    /// POST /groups/invites/{id}/accept
    func acceptGroupInvite(id: Int) async throws {
        try await postVoid("/groups/invites/\(id)/accept", body: [String: String]())
    }

    /// POST /groups/invites/{id}/decline
    func declineGroupInvite(id: Int) async throws {
        try await postVoid("/groups/invites/\(id)/decline", body: [String: String]())
    }

    func listRelationships() async throws -> [RemoteRelationship] {
        try await get("/relationships")
    }

    /// Send an email-based group invite.
    /// status is "added", "already_member", or "invited".
    func groupEmailInvite(groupID: Int, email: String) async throws -> GroupInviteOutcome {
        struct Body: Decodable { let status: String? }
        let body: Body = try await post(
            "/groups/\(groupID)/email-invite",
            body: ["email": email]
        )
        if let outcome = GroupInviteOutcome.from(status: body.status ?? "") {
            return outcome
        }
        return .invited
    }

    // MARK: - Primitives

    private func get<T: Decodable>(_ path: String) async throws -> T {
        do { try await refreshTokenIfNeeded() } catch { print("[APIClient] token refresh: \(error)") }
        var req = try makeRequest(path)
        try attachToken(&req)
        return try await send(req)
    }

    private func post<Body: Encodable, T: Decodable>(_ path: String,
                                                      body: Body,
                                                      auth: Bool = true) async throws -> T {
        do { try await refreshTokenIfNeeded() } catch { print("[APIClient] token refresh: \(error)") }
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
            throw APIError.fromServer(data: data, status: http.statusCode)
        }
    }

    /// PUT that returns a decoded response body.
    private func put<Body: Encodable, T: Decodable>(_ path: String,
                                                     body: Body) async throws -> T {
        do { try await refreshTokenIfNeeded() } catch { print("[APIClient] token refresh: \(error)") }
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
            throw APIError.fromServer(data: data, status: http.statusCode)
        }
    }

    /// PATCH with a JSON object body. Throws the server's error string on failure.
    private func patchJSON(_ path: String, body: JSONStringMap) async throws {
        var req = try makeRequest(path, method: "PATCH")
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpBody = try JSONEncoder().encode(body)
        try attachToken(&req)
        let (data, resp) = try await session.data(for: req)
        if let http = resp as? HTTPURLResponse, !(200...299).contains(http.statusCode) {
            throw APIError.fromServer(data: data, status: http.statusCode)
        }
    }

    /// Fire-and-forget PATCH for endpoints with no request body.
    private func patchVoid(_ path: String) async throws {
        var req = try makeRequest(path, method: "PATCH")
        try attachToken(&req)
        let (data, resp) = try await session.data(for: req)
        if let http = resp as? HTTPURLResponse, !(200...299).contains(http.statusCode) {
            throw APIError.fromServer(data: data, status: http.statusCode)
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
            throw APIError.fromServer(data: data, status: http.statusCode)
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
                throw APIError.fromServer(data: data, status: http.statusCode)
            }
            return try JSONDecoder().decode(T.self, from: data)
        } catch let e as APIError { throw e }
        catch let e as DecodingError { throw APIError.decoding(e) }
        catch { throw APIError.network(error) }
    }
}
