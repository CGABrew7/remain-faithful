# Rebuild the NSIS installer so it contains the already signed app.
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

$exe = Join-Path $root 'build\bin\RemainFaithful.exe'
$nsiDir = Join-Path $root 'build\windows\installer'
$nsi = Join-Path $nsiDir 'project.nsi'
if (-not (Test-Path -LiteralPath $exe)) {
    throw 'RemainFaithful.exe was not built.'
}
if (-not (Test-Path -LiteralPath $nsi)) {
    throw 'NSIS project was not generated.'
}

function Find-MakeNsis {
    $command = Get-Command makensis -ErrorAction SilentlyContinue
    if ($command) {
        return $command.Source
    }
    $candidates = @(
        "${env:ProgramFiles(x86)}\NSIS\makensis.exe",
        "${env:ProgramFiles(x86)}\NSIS\Bin\makensis.exe",
        'C:\ProgramData\chocolatey\bin\makensis.exe'
    )
    foreach ($candidate in $candidates) {
        if (Test-Path -LiteralPath $candidate) {
            return $candidate
        }
    }
    throw 'makensis was not found.'
}

$makensis = Find-MakeNsis
Push-Location $nsiDir
try {
    & $makensis "-DARG_WAILS_AMD64_BINARY=$exe" 'project.nsi'
    if ($LASTEXITCODE -ne 0) {
        throw 'makensis failed while rebuilding the installer.'
    }
}
finally {
    Pop-Location
}
Write-Host 'Rebuilt the installer around the signed app.'
