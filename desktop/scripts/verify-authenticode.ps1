# Fail when a Windows file does not have a valid Authenticode signature.
param(
    [Parameter(Mandatory = $true)]
    [string[]]$Files
)
$ErrorActionPreference = 'Stop'

foreach ($file in $Files) {
    if (-not (Test-Path -LiteralPath $file)) {
        throw "Signed file was not found: $(Split-Path -Leaf $file)"
    }
    $signature = Get-AuthenticodeSignature -LiteralPath $file
    if ($signature.Status -ne 'Valid') {
        throw "Authenticode status for $(Split-Path -Leaf $file) is $($signature.Status)."
    }
    Write-Host "Authenticode valid: $(Split-Path -Leaf $file)"
}
