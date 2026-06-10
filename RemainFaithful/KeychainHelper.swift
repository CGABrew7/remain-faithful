import Security
import Foundation
import os

// Root cause of prior silent failures: kSecAttrAccessGroup requires the
// keychain-access-groups entitlement, which was absent from the provisioning
// profile (only com.apple.security.application-groups was present). On a real
// device this returned errSecMissingEntitlement (-34018) for every write.
// The broadcast extension reads authToken from shared UserDefaults, not
// Keychain, so removing the access group does not break IPC.
private let keychainLogger = Logger(subsystem: "com.remainfaithful.app", category: "Keychain")

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
        // Search query uses only identity attrs — no data or accessibility attrs.
        let searchQuery: [CFString: Any] = [
            kSecClass:       kSecClassGenericPassword,
            kSecAttrService: service,
            kSecAttrAccount: key,
        ]
        // Try update first; add only when the item doesn't exist yet.
        let updateStatus = SecItemUpdate(
            searchQuery as CFDictionary,
            [kSecValueData: data] as CFDictionary
        )
        if updateStatus == errSecItemNotFound {
            var addAttrs = searchQuery
            addAttrs[kSecValueData]      = data
            addAttrs[kSecAttrAccessible] = kSecAttrAccessibleAfterFirstUnlock
            let addStatus = SecItemAdd(addAttrs as CFDictionary, nil)
            if addStatus != errSecSuccess {
                keychainLogger.error("SecItemAdd('\(key)'): OSStatus \(addStatus)")
            }
        } else if updateStatus != errSecSuccess {
            keychainLogger.error("SecItemUpdate('\(key)'): OSStatus \(updateStatus)")
        }
    }

    func getData(_ key: String) -> Data? {
        let query: [CFString: Any] = [
            kSecClass:       kSecClassGenericPassword,
            kSecAttrService: service,
            kSecAttrAccount: key,
            kSecReturnData:  true,
            kSecMatchLimit:  kSecMatchLimitOne,
        ]
        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)
        if status != errSecSuccess && status != errSecItemNotFound {
            keychainLogger.error("SecItemCopyMatching('\(key)'): OSStatus \(status)")
        }
        guard status == errSecSuccess else { return nil }
        return result as? Data
    }

    func delete(_ key: String) {
        let query: [CFString: Any] = [
            kSecClass:       kSecClassGenericPassword,
            kSecAttrService: service,
            kSecAttrAccount: key,
        ]
        let status = SecItemDelete(query as CFDictionary)
        if status != errSecSuccess && status != errSecItemNotFound {
            keychainLogger.error("SecItemDelete('\(key)'): OSStatus \(status)")
        }
    }
}
