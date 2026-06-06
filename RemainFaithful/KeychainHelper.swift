import Security
import Foundation

final class KeychainHelper {
    static let shared = KeychainHelper(service: "com.remainfaithful.app")
    let service: String
    init(service: String) { self.service = service }

    func set(_ value: String, for key: String) {
        guard let data = value.data(using: .utf8) else { return }
        setData(data, for: key)
    }

    func get(_ key: String) -> String? {
        guard let data = getData(key) else { return nil }
        return String(data: data, encoding: .utf8)
    }

    func setData(_ data: Data, for key: String) {
        let query: [CFString: Any] = [
            kSecClass: kSecClassGenericPassword,
            kSecAttrService: service,
            kSecAttrAccount: key,
            kSecAttrAccessible: kSecAttrAccessibleAfterFirstUnlock,
            kSecAttrAccessGroup: "group.com.remainfaithful.app",
        ]
        let status = SecItemCopyMatching(query as CFDictionary, nil)
        if status == errSecSuccess {
            SecItemUpdate(query as CFDictionary,
                          [kSecValueData: data] as CFDictionary)
        } else {
            var attrs = query
            attrs[kSecValueData] = data
            SecItemAdd(attrs as CFDictionary, nil)
        }
    }

    func getData(_ key: String) -> Data? {
        let query: [CFString: Any] = [
            kSecClass: kSecClassGenericPassword,
            kSecAttrService: service,
            kSecAttrAccount: key,
            kSecAttrAccessible: kSecAttrAccessibleAfterFirstUnlock,
            kSecAttrAccessGroup: "group.com.remainfaithful.app",
            kSecReturnData: true,
            kSecMatchLimit: kSecMatchLimitOne,
        ]
        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)
        guard status == errSecSuccess else { return nil }
        return result as? Data
    }

    func delete(_ key: String) {
        let query: [CFString: Any] = [
            kSecClass: kSecClassGenericPassword,
            kSecAttrService: service,
            kSecAttrAccount: key,
            kSecAttrAccessible: kSecAttrAccessibleAfterFirstUnlock,
            kSecAttrAccessGroup: "group.com.remainfaithful.app",
        ]
        SecItemDelete(query as CFDictionary)
    }
}
