# Run Go with project-local dependency, build and temporary directories.
$ErrorActionPreference = 'Stop'
$sdkRoot = Split-Path -Parent $PSScriptRoot
$sdkTemporary = Join-Path $sdkRoot 'tmp'
$sdkVariables = @{
    GOPATH = (Join-Path $sdkTemporary 'gopath')
    GOMODCACHE = (Join-Path $sdkTemporary 'gomodcache')
    GOCACHE = (Join-Path $sdkTemporary 'gocache')
    GOTMPDIR = (Join-Path $sdkTemporary 'gotmp')
    TEST_TELEMETRY_DIR = (Join-Path $sdkTemporary 'telemetry')
    TEMP = (Join-Path $sdkTemporary 'gotmp')
    TMP = (Join-Path $sdkTemporary 'gotmp')
    GOTOOLCHAIN = 'local'
    GOENV = 'off'
    GOWORK = 'off'
}
$sdkPrevious = @{}
foreach ($sdkVariable in $sdkVariables.Keys) {
    $sdkPrevious[$sdkVariable] = [Environment]::GetEnvironmentVariable($sdkVariable, 'Process')
}
try {
    foreach ($sdkVariable in $sdkVariables.Keys) {
        if ($sdkVariable -notin @('GOTOOLCHAIN', 'GOENV', 'GOWORK')) {
            New-Item -ItemType Directory -Path $sdkVariables[$sdkVariable] -Force | Out-Null
        }
        [Environment]::SetEnvironmentVariable($sdkVariable, $sdkVariables[$sdkVariable], 'Process')
    }
    [System.IO.File]::WriteAllText((Join-Path $sdkVariables.TEST_TELEMETRY_DIR 'mode'), 'off')
    $sdkTmpModule = Join-Path $sdkTemporary 'go.mod'
    if (-not (Test-Path -LiteralPath $sdkTmpModule)) {
        [System.IO.File]::WriteAllText($sdkTmpModule, "module github.com/Lixsp11/inspire-go-sdk/tmp`n`ngo 1.23.0`n")
    }
    $sdkGo = Join-Path $sdkTemporary 'go/bin/go.exe'
    if (-not (Test-Path -LiteralPath $sdkGo)) { $sdkGo = Join-Path $sdkTemporary 'go/bin/go' }
    if (-not (Test-Path -LiteralPath $sdkGo)) { $sdkGo = (Get-Command go -ErrorAction Stop).Source }
    Push-Location $sdkRoot
    try { & $sdkGo @args; $sdkExitCode = $LASTEXITCODE } finally { Pop-Location }
} finally {
    foreach ($sdkVariable in $sdkVariables.Keys) {
        [Environment]::SetEnvironmentVariable($sdkVariable, $sdkPrevious[$sdkVariable], 'Process')
    }
}
if ($sdkExitCode -ne 0) { exit $sdkExitCode }
