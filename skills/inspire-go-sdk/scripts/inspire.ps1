# Run the versioned Windows CLI without changing PATH or installing a compiler.
$cliArguments = @($args)
$ErrorActionPreference = 'Stop'
$cliStaging = $null
try {
    if ($env:OS -ne 'Windows_NT') { throw 'Use inspire.sh on Linux or macOS.' }
    $cliVersion = [IO.File]::ReadAllText((Join-Path $PSScriptRoot 'VERSION')).Trim()
    if ($cliVersion -notmatch '^\d+\.\d+\.\d+$') { throw 'Invalid CLI version.' }
    $cliArchitecture = $env:PROCESSOR_ARCHITEW6432
    if (-not $cliArchitecture) { $cliArchitecture = $env:PROCESSOR_ARCHITECTURE }
    switch ($cliArchitecture) {
        'AMD64' { $cliArch = 'amd64' }
        'ARM64' { $cliArch = 'arm64' }
        default { throw 'Unsupported Windows architecture.' }
    }
    $cliCacheRoot = $env:INSPIRE_CLI_CACHE_DIR
    if (-not $cliCacheRoot) { $cliCacheRoot = Join-Path ([Environment]::GetFolderPath('LocalApplicationData')) 'Inspire/cli' }
    $cliCacheRoot = [IO.Path]::GetFullPath($cliCacheRoot)
    $cliCacheDirectory = Join-Path $cliCacheRoot "$cliVersion/windows_$cliArch"
    $cliExecutable = Join-Path $cliCacheDirectory 'inspire.exe'
    if (-not (Test-Path -LiteralPath $cliExecutable -PathType Leaf)) {
        New-Item -ItemType Directory -Path $cliCacheDirectory -Force | Out-Null
        $cliStaging = [IO.Path]::GetFullPath((Join-Path $cliCacheDirectory ('.download-' + [guid]::NewGuid().ToString('N'))))
        if (-not $cliStaging.StartsWith($cliCacheRoot + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) { throw 'Invalid cache staging path.' }
        New-Item -ItemType Directory -Path $cliStaging | Out-Null
        $cliAsset = "inspire_${cliVersion}_windows_${cliArch}.zip"
        $cliRelease = "https://github.com/Lixsp11/inspire-go-sdk/releases/download/v$cliVersion"
        Invoke-WebRequest -UseBasicParsing -Uri "$cliRelease/SHA256SUMS" -OutFile (Join-Path $cliStaging 'SHA256SUMS')
        Invoke-WebRequest -UseBasicParsing -Uri "$cliRelease/$cliAsset" -OutFile (Join-Path $cliStaging $cliAsset)
        $cliMatches = @(Get-Content -LiteralPath (Join-Path $cliStaging 'SHA256SUMS') | Where-Object { $_ -match ('^[a-f0-9]{64}  ' + [regex]::Escape($cliAsset) + '$') })
        if ($cliMatches.Count -ne 1) { throw 'Missing release checksum.' }
        $cliExpected = $cliMatches[0].Substring(0, 64)
        $cliHasher = [Security.Cryptography.SHA256]::Create()
        $cliHashStream = $null
        try {
            $cliHashStream = [IO.File]::OpenRead((Join-Path $cliStaging $cliAsset))
            $cliActual = [BitConverter]::ToString($cliHasher.ComputeHash($cliHashStream)).Replace('-', '').ToLowerInvariant()
        } finally {
            if ($cliHashStream) { $cliHashStream.Dispose() }
            $cliHasher.Dispose()
        }
        if ($cliActual -ne $cliExpected) { throw 'Release checksum mismatch.' }
        Add-Type -AssemblyName System.IO.Compression.FileSystem
        $cliPayload = Join-Path $cliStaging 'payload'
        [IO.Compression.ZipFile]::ExtractToDirectory((Join-Path $cliStaging $cliAsset), $cliPayload)
        if (-not (Test-Path -LiteralPath (Join-Path $cliPayload 'inspire.exe') -PathType Leaf)) { throw 'Missing CLI executable.' }
        Copy-Item -LiteralPath (Join-Path $cliPayload 'LICENSE'),(Join-Path $cliPayload 'NOTICE') -Destination $cliCacheDirectory -Force
        Copy-Item -LiteralPath (Join-Path $cliPayload 'third_party') -Destination $cliCacheDirectory -Recurse -Force
        try { [IO.File]::Move((Join-Path $cliPayload 'inspire.exe'), $cliExecutable) } catch { if (-not (Test-Path -LiteralPath $cliExecutable -PathType Leaf)) { throw } }
    }
} catch {
    [Console]::Error.WriteLine('{"schemaVersion":"v1","error_code":"BootstrapError","message":"Could not prepare the prebuilt Inspire CLI.","recovery":{"action":"inspect","hint":"Check GitHub access, archive tools and the matching release, or install inspire from its release archive."}}')
    exit 1
} finally {
    if ($cliStaging -and (Test-Path -LiteralPath $cliStaging)) {
        if ($cliStaging.StartsWith($cliCacheRoot + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) { Remove-Item -LiteralPath $cliStaging -Recurse -Force }
    }
}
function ConvertTo-CliArgument {
    param([string]$Value)
    if ($Value -ne '' -and $Value -notmatch '[\s"]') { return $Value }
    $cliBuilder = New-Object Text.StringBuilder
    [void]$cliBuilder.Append('"')
    $cliSlashes = 0
    foreach ($cliCharacter in $Value.ToCharArray()) {
        if ($cliCharacter -eq '\') { $cliSlashes++; continue }
        if ($cliCharacter -eq '"') {
            [void]$cliBuilder.Append(('\' * (2 * $cliSlashes + 1)))
        } else {
            [void]$cliBuilder.Append(('\' * $cliSlashes))
        }
        [void]$cliBuilder.Append($cliCharacter)
        $cliSlashes = 0
    }
    [void]$cliBuilder.Append(('\' * (2 * $cliSlashes)))
    [void]$cliBuilder.Append('"')
    return $cliBuilder.ToString()
}
$cliProcess = New-Object Diagnostics.Process
$cliProcess.StartInfo.FileName = $cliExecutable
$cliProcess.StartInfo.UseShellExecute = $false
$cliProcess.StartInfo.CreateNoWindow = $true
$cliProcess.StartInfo.RedirectStandardOutput = $true
$cliProcess.StartInfo.RedirectStandardError = $true
$cliProcess.StartInfo.RedirectStandardInput = $true
$cliProcess.StartInfo.Arguments = (@($cliArguments | ForEach-Object { ConvertTo-CliArgument $_ }) -join ' ')
$cliStarted = $false
try {
    $cliInput = $null
    for ($cliIndex = 0; $cliIndex -lt $cliArguments.Count; $cliIndex++) {
        if ($cliArguments[$cliIndex] -eq '--body-file=-' -or ($cliArguments[$cliIndex] -eq '--body-file' -and $cliIndex + 1 -lt $cliArguments.Count -and $cliArguments[$cliIndex + 1] -eq '-')) {
            $cliPipelineInput = @($input)
            if ($cliPipelineInput.Count -gt 0) {
                $cliInput = [Text.Encoding]::UTF8.GetBytes((@($cliPipelineInput | ForEach-Object { [string]$_ }) -join "`n"))
                $cliInputLength = [Math]::Min($cliInput.Length, (2MB + 1))
            } else {
                $cliInput = New-Object byte[] (2MB + 1)
                $cliInputLength = 0
                $cliInputStream = [Console]::OpenStandardInput()
                while ($cliInputLength -lt $cliInput.Length) {
                    $cliRead = $cliInputStream.Read($cliInput, $cliInputLength, $cliInput.Length - $cliInputLength)
                    if ($cliRead -eq 0) { break }
                    $cliInputLength += $cliRead
                }
            }
            break
        }
    }
    $cliOriginalInputEncoding = [Console]::InputEncoding
    try {
        # Windows PowerShell's redirected writer otherwise emits a UTF-8 BOM.
        [Console]::InputEncoding = New-Object Text.UTF8Encoding($false)
        [void]$cliProcess.Start()
        $cliStarted = $true
    } finally {
        [Console]::InputEncoding = $cliOriginalInputEncoding
    }
    $cliStdoutTask = $cliProcess.StandardOutput.BaseStream.CopyToAsync([Console]::OpenStandardOutput())
    $cliStderrTask = $cliProcess.StandardError.BaseStream.CopyToAsync([Console]::OpenStandardError())
    try {
        if ($cliInput) { $cliProcess.StandardInput.BaseStream.Write($cliInput, 0, $cliInputLength) }
        $cliProcess.StandardInput.BaseStream.Close()
    } catch [IO.IOException] {
        # The CLI may reject its arguments before consuming stdin.
    }
    $cliProcess.WaitForExit()
    [void]$cliStdoutTask.GetAwaiter().GetResult()
    [void]$cliStderrTask.GetAwaiter().GetResult()
    $cliExitCode = $cliProcess.ExitCode
} catch {
    if ($cliStarted) {
        [Console]::Error.WriteLine('{"schemaVersion":"v1","error_code":"LauncherError","message":"CLI execution or output forwarding failed.","recovery":{"action":"reconcile","hint":"Check platform state before repeating any write."}}')
    } else {
        [Console]::Error.WriteLine('{"schemaVersion":"v1","error_code":"BootstrapError","message":"Could not start the cached Inspire CLI.","recovery":{"action":"inspect","hint":"Check the cached executable and execution permissions, or install the matching release."}}')
    }
    exit 1
} finally { $cliProcess.Dispose() }
exit $cliExitCode
