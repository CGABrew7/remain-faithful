# Authenticode-sign files with the PFX certificate from GitHub Actions secrets.
param(
    [Parameter(Mandatory = $true)]
    [string[]]$Files
)
$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrWhiteSpace($env:WINDOWS_PFX_BASE64) -or [string]::IsNullOrWhiteSpace($env:WINDOWS_PFX_PASSWORD)) {
    throw 'WINDOWS_PFX_BASE64 and WINDOWS_PFX_PASSWORD are required for PFX signing.'
}

function Find-SignTool {
    $kits = "${env:ProgramFiles(x86)}\Windows Kits\10\bin"
    if (Test-Path -LiteralPath $kits) {
        $found = Get-ChildItem -Path $kits -Recurse -Filter signtool.exe -ErrorAction SilentlyContinue |
            Where-Object { $_.FullName -match '\\x64\\signtool.exe$' } |
            Sort-Object FullName -Descending |
            Select-Object -First 1
        if ($found) {
            return $found.FullName
        }
    }
    $command = Get-Command signtool -ErrorAction SilentlyContinue
    if ($command) {
        return $command.Source
    }
    throw 'signtool was not found. The Windows SDK is required for PFX signing.'
}

$signtool = Find-SignTool
$pfxPath = Join-Path $env:RUNNER_TEMP 'remain-faithful-codesign.pfx'
try {
    $raw = ($env:WINDOWS_PFX_BASE64 -replace '\s', '')
    [System.IO.File]::WriteAllBytes($pfxPath, [Convert]::FromBase64String($raw))
    foreach ($file in $Files) {
        if (-not (Test-Path -LiteralPath $file)) {
            throw "File to sign was not found: $(Split-Path -Leaf $file)"
        }
        & $signtool sign /fd SHA256 /td SHA256 /tr http://timestamp.digicert.com /f $pfxPath /p $env:WINDOWS_PFX_PASSWORD $file
        if ($LASTEXITCODE -ne 0) {
            throw "signtool could not sign $(Split-Path -Leaf $file)."
        }
        Write-Host "Signed $(Split-Path -Leaf $file) with the PFX certificate."
    }
}
finally {
    if (Test-Path -LiteralPath $pfxPath) {
        Remove-Item -LiteralPath $pfxPath -Force
    }
}
