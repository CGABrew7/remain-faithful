# Signing the Mac and Windows apps

This page is for the person who owns the Apple and Microsoft accounts. It explains how to create the certificates and how to save them as GitHub Actions secrets. After that, a version tag can publish downloads that macOS and Windows will accept.

The workflow reads those secrets from GitHub. It does not read certificates from the repository. This page has no real certificate, password, key, or team id. Do not put those values in the repository, in a pull request, or in an email.

The iPhone push key already in the project is only for the iPhone app. Leave that key where it is. Do not upload it here, and do not replace it with these new keys.

## Where a secret goes

1. Open the Remain Faithful repository on GitHub.
2. Choose **Settings**, then **Secrets and variables**, then **Actions**.
3. Choose **New repository secret**.
4. Type the name exactly as it appears below. Names are case sensitive.
5. Paste the value into the secret field and save it.

GitHub will not show the value again. You can replace a secret later with **Update**. The Desktop CI log lists a secret name when it is missing. The log does not print the value.

If a required secret is missing, Desktop CI still finishes successfully and uploads an unsigned file. A release tag does the same, and the release notes say the file is unsigned. Add the secrets, then run the workflow again before you publish a version you want people to install.

## Mac

You need an Apple Developer Program membership for the publisher name you want people to see. You also need a Mac for the certificate steps below.

Create two things:

- A **Developer ID Application** certificate. This is the certificate that signs the app. It is not an Apple Development certificate, not an Apple Distribution certificate, and not the iPhone push key.
- An **App Store Connect API key**. The workflow uses this key to send the app to Apple for notarization. Notarization is Apple's check that lets a downloaded app open.

### Create the Developer ID Application certificate

1. On a Mac, open **Keychain Access**.
2. In the menu bar, choose **Keychain Access**, then **Certificate Assistant**, then **Request a Certificate From a Certificate Authority**.
3. Enter the email on the Apple developer account and the publisher name you want on the certificate. Leave the certificate authority email blank. Choose **Saved to disk**. Save the request file. This file is not the secret. It is only used in the next step.
4. Open [Apple Developer account](https://developer.apple.com/account) and sign in.
5. Open **Certificates, Identifiers & Profiles**, then **Certificates**.
6. Choose the add button. Under software, choose **Developer ID Application**. Continue.
7. Upload the request file from step 3. Download the certificate Apple gives you.
8. Double-click the downloaded certificate so it opens in Keychain Access. Confirm that it has a private key under it. If there is no private key, the request and the certificate were made on different Macs. Start again on the Mac that will export the certificate.
9. Select the certificate and its private key together. Choose **File**, then **Export Items**. Choose the Personal Information Exchange format, which is a `.p12` file. Set a password and remember it. You will save the password as its own secret.

Turn the `.p12` file into one line of text. In Terminal, on the Mac that has the file:

```bash
base64 -i /path/to/developer-id.p12 | pbcopy
```

The copied text is the secret. It is not something you commit.

| Secret name | What to paste |
| --- | --- |
| `MACOS_CERT_P12_BASE64` | The base64 text from the command above |
| `MACOS_CERT_PASSWORD` | The password you set when you exported the `.p12` file |

### Create the App Store Connect API key

1. Open [App Store Connect](https://appstoreconnect.apple.com) and sign in as an Account Holder or Admin.
2. Open **Users and Access**, then **Integrations**, then **App Store Connect API**.
3. Copy the **Issuer ID** shown at the top of that page. That value is `APPLE_API_ISSUER_ID`.
4. Choose to generate a team key. Name it something you will recognize, such as desktop notarization. Set the access to **Admin**.
5. Download the key file. Apple lets you download it once. It is a small text file whose name starts with `AuthKey`. Keep the download somewhere private until the secret is saved.
6. Copy the **Key ID** shown next to that key. That value is `APPLE_API_KEY_ID`.

Turn the key file into one line of text:

```bash
base64 -i /path/to/AuthKey.p8 | pbcopy
```

Paste that text into `APPLE_API_KEY_P8_BASE64`. Paste the base64 text, not the original contents of the key file.

### Copy the Team ID

1. On [Apple Developer account](https://developer.apple.com/account), open **Membership details**.
2. Copy the **Team ID**.
3. Save it as `APPLE_TEAM_ID`.

The workflow uses the Team ID only to pick the Developer ID Application certificate out of the `.p12` file. It does not print the Team ID in the log.

### Mac secret checklist

Save all six of these. If any one is missing, the Mac build stays unsigned.

- `MACOS_CERT_P12_BASE64`
- `MACOS_CERT_PASSWORD`
- `APPLE_API_KEY_ID`
- `APPLE_API_ISSUER_ID`
- `APPLE_API_KEY_P8_BASE64`
- `APPLE_TEAM_ID`

When they are all present, Desktop CI signs the app and the monitor helper with the Developer ID Application certificate, turns on the hardened runtime, keeps the entitlements already in the project, sends the app to Apple for notarization, staples the ticket, and puts the result in a signed disk image.

## Windows with Azure Artifact Signing

This is the signing path to use. Azure Artifact Signing used to be called Trusted Signing. In the Azure portal, search for **Artifact Signing**.

A normal Windows PC trusts this kind of signature when the certificate profile is **Public Trust**. A private trust profile is only for inside an organization. Choose Public Trust.

Microsoft reviews the organization before a public certificate profile can sign. The Desktop CI log keeps saying the Windows file is unsigned until that review is finished and the secrets below are saved.

### Create the signing account

1. Sign in to the [Azure portal](https://portal.azure.com) with an account that can create resources.
2. Register the **Artifact Signing** resource provider for the subscription if Azure asks you to.
3. Create an **Artifact Signing account**. The region you pick becomes part of the endpoint. Copy the account name. That name, not the long resource id, is `AZURE_SIGNING_ACCOUNT`.
4. On the account page, copy the endpoint. It is a URL that ends in `codesigning.azure.net`. That URL is `AZURE_SIGNING_ENDPOINT`.
5. Start **identity validation**. Enter the organization's legal details in the Azure form. Those details stay in Azure. Do not copy them into GitHub.
6. Wait until the identity validation status is valid. Then create a **certificate profile**. Choose Public Trust and the validated identity. Copy the certificate profile name. That name is `AZURE_CERT_PROFILE`.

The setup pages Microsoft maintains are [Artifact Signing overview](https://learn.microsoft.com/en-us/azure/artifact-signing/overview) and the quickstart linked from there.

### Let GitHub sign on your behalf

1. In the Azure portal, open **Microsoft Entra ID**, then **App registrations**, then **New registration**.
2. Name it something you will recognize, such as Remain Faithful desktop signing. A single-tenant registration is enough.
3. Open the new registration. Copy the **Application (client) ID**. That is `AZURE_CLIENT_ID`. Copy the **Directory (tenant) ID**. That is `AZURE_TENANT_ID`.
4. Open **Certificates & secrets**, then **New client secret**. Copy the secret value immediately. Azure shows it once. That value is `AZURE_CLIENT_SECRET`. You will need to replace this GitHub secret before the client secret expires.
5. Open the Artifact Signing certificate profile. Open **Access control (IAM)**, then **Add role assignment**.
6. Choose the role **Artifact Signing Certificate Profile Signer**. If the portal still lists **Trusted Signing Certificate Profile Signer**, that is the same role. Assign it to the app registration from step 2.

### Azure secret checklist

- `AZURE_TENANT_ID`
- `AZURE_CLIENT_ID`
- `AZURE_CLIENT_SECRET`
- `AZURE_SIGNING_ENDPOINT`
- `AZURE_SIGNING_ACCOUNT`
- `AZURE_CERT_PROFILE`

When all six are set, the workflow signs `RemainFaithful.exe` and then rebuilds the installer so the installer contains that signed app. It then signs the installer. Azure is used even if a PFX certificate is also saved.

## Windows with a certificate file

Use this only when you are not using Azure Artifact Signing. Buy an **Organization Validation code signing** certificate from a public certificate authority. A domain-only certificate will not do this job. Complete the organization's validation with that authority. Those details stay with the authority.

Export or download a `.pfx` file that includes the private key, and set a password.

On Windows, PowerShell can copy the base64 text:

```powershell
[Convert]::ToBase64String([IO.File]::ReadAllBytes("C:\path\to\codesign.pfx")) | Set-Clipboard
```

| Secret name | What to paste |
| --- | --- |
| `WINDOWS_PFX_BASE64` | The base64 text from the command above |
| `WINDOWS_PFX_PASSWORD` | The `.pfx` password |

The workflow uses the PFX certificate when the Azure secrets are not all present and both PFX secrets are present. If only some of the Azure secrets are set, the log names the missing Azure secrets and uses the PFX certificate instead.

A signed installer can still show a Windows SmartScreen notice while the file is new. That notice is separate from the signature. An unsigned installer does not name a publisher.

## Confirm a build

1. On GitHub, open **Actions**, then **Desktop CI**, then **Run workflow**, and choose the `main` branch.
2. Open the finished Mac job. The summary should say `Mac disk image: signed.` If secrets are missing, it names them and says the disk image is unsigned.
3. Open the finished Windows job. The summary should say the installer was signed with Azure Artifact Signing, or with a PFX certificate.

You can also download `RemainFaithful-macos.dmg` from the workflow artifacts and open it on a Mac. A signed build opens without a warning that the developer cannot be verified. Download `RemainFaithful-windows-setup.exe` and check that Windows shows a publisher on the installer.

## Publish a version

Publishing is a tag push. This repository does not publish a desktop release until somebody pushes a tag.

1. Confirm the Mac and Windows summaries say the files are signed.
2. On the commit you want to ship, create a tag named like `desktop-v1.0.0`. The word `desktop-v` and the version numbers are required. `desktop-v1.0.0-rc1` is treated as a prerelease.
3. Push that tag to GitHub.

The Desktop release workflow builds the apps again, signs them when the secrets are present, and creates a GitHub Release. The release files keep the same names for every version:

- `RemainFaithful-macos.dmg`
- `RemainFaithful-windows-setup.exe`
- `SHA256SUMS.txt`

These are the download links, with the version replaced by `latest` after the first release exists:

```text
https://github.com/CGABrew7/remain-faithful/releases/latest/download/RemainFaithful-macos.dmg
https://github.com/CGABrew7/remain-faithful/releases/latest/download/RemainFaithful-windows-setup.exe
https://github.com/CGABrew7/remain-faithful/releases/latest/download/SHA256SUMS.txt
```

If you push the tag before the secrets are saved, GitHub still creates the release and the notes say the files are unsigned. Delete that release and that tag, save the secrets, and push the tag again.

The website is not changed by the signing workflow. It can point at the links above when you want desktop downloads on the site.
