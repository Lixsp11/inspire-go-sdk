param(
    [Parameter(Mandatory=$true)][string]$Version,
    [ValidateSet('windows','linux','darwin')][string]$TargetOS,
    [ValidateSet('amd64','arm64')][string]$TargetArch,
    [string]$OutputDirectory = 'dist',
    [switch]$IncludeSource
)
$ErrorActionPreference = 'Stop'
if ($Version -notmatch '^\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$') { throw 'Invalid release version.' }
if (($TargetOS -and -not $TargetArch) -or ($TargetArch -and -not $TargetOS)) { throw 'TargetOS and TargetArch must be supplied together.' }
if (-not $TargetOS) { $IncludeSource = $true }
$sdkRoot = Split-Path -Parent $PSScriptRoot
$sdkManifest = Get-Content -LiteralPath (Join-Path $sdkRoot 'plugin.json') -Raw | ConvertFrom-Json
if ($sdkManifest.name -ne 'inspire-go-sdk' -or $sdkManifest.version -ne $Version) { throw 'Release version differs from plugin manifest.' }
$sdkVersionText = [IO.File]::ReadAllText((Join-Path $sdkRoot 'src/openapi/client/config.go'))
$sdkVersionMatch = [regex]::Match($sdkVersionText, 'const\s+Version\s*=\s*"([^"]+)"')
if (-not $sdkVersionMatch.Success -or $sdkVersionMatch.Groups[1].Value -ne $Version) { throw 'Release version differs from SDK Version.' }
$sdkOutput = [IO.Path]::GetFullPath((Join-Path $sdkRoot $OutputDirectory))
$sdkComparison = [StringComparison]::Ordinal
if ($env:OS -eq 'Windows_NT') { $sdkComparison = [StringComparison]::OrdinalIgnoreCase }
if (-not $sdkOutput.StartsWith($sdkRoot + [IO.Path]::DirectorySeparatorChar, $sdkComparison)) { throw 'Release output must remain inside the project.' }
New-Item -ItemType Directory -Path $sdkOutput -Force | Out-Null
$sdkStage = Join-Path $sdkRoot ('tmp/package-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $sdkStage -Force | Out-Null
function New-SdkZip {
    param([string]$SourceDirectory, [string]$DestinationFile)
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    Add-Type -AssemblyName System.IO.Compression
    $sourceRoot = [IO.Path]::GetFullPath($SourceDirectory).TrimEnd([char[]]'\/')
    $zip = [IO.Compression.ZipFile]::Open($DestinationFile, [IO.Compression.ZipArchiveMode]::Create)
    try {
        foreach ($file in (Get-ChildItem -LiteralPath $sourceRoot -Recurse -File | Sort-Object FullName)) {
            $entryName = $file.FullName.Substring($sourceRoot.Length + 1).Replace('\', '/')
            [IO.Compression.ZipFileExtensions]::CreateEntryFromFile($zip, $file.FullName, $entryName, [IO.Compression.CompressionLevel]::Optimal) | Out-Null
        }
    } finally { $zip.Dispose() }
}
Push-Location $sdkRoot
try {
    & git diff --quiet HEAD -- .
    if ($LASTEXITCODE -ne 0) { throw 'Commit tracked changes before packaging a release.' }
    if ($TargetOS) {
        $sdkPayload = Join-Path $sdkStage 'cli'
        New-Item -ItemType Directory -Path $sdkPayload | Out-Null
        $sdkBinaryName = 'inspire'
        if ($TargetOS -eq 'windows') { $sdkBinaryName += '.exe' }
        $sdkSaved = @{}
        foreach ($sdkName in @('GOOS','GOARCH','CGO_ENABLED')) { $sdkSaved[$sdkName] = [Environment]::GetEnvironmentVariable($sdkName, 'Process') }
        try {
            $env:GOOS = $TargetOS
            $env:GOARCH = $TargetArch
            $env:CGO_ENABLED = '0'
            & (Join-Path $PSScriptRoot 'go.ps1') build -trimpath '-ldflags=-s -w' -o (Join-Path $sdkPayload $sdkBinaryName) ./src/main
            if ($LASTEXITCODE -ne 0) { throw 'CLI build failed.' }
        } finally {
            foreach ($sdkName in $sdkSaved.Keys) { [Environment]::SetEnvironmentVariable($sdkName, $sdkSaved[$sdkName], 'Process') }
        }
        foreach ($sdkFile in @('LICENSE','NOTICE','README.md')) { Copy-Item -LiteralPath (Join-Path $sdkRoot $sdkFile) -Destination $sdkPayload }
        & (Join-Path $PSScriptRoot 'go.ps1') run ./tools/licenses -os $TargetOS -arch $TargetArch -output (Join-Path $sdkPayload 'third_party')
        if ($LASTEXITCODE -ne 0) { throw 'Dependency license collection failed.' }
        $sdkAssetBase = "inspire_${Version}_${TargetOS}_${TargetArch}"
        if ($TargetOS -eq 'windows') {
            New-SdkZip $sdkPayload (Join-Path $sdkOutput ($sdkAssetBase + '.zip'))
        } else {
            & tar -czf (Join-Path $sdkOutput ($sdkAssetBase + '.tar.gz')) -C $sdkPayload $sdkBinaryName LICENSE NOTICE README.md third_party
            if ($LASTEXITCODE -ne 0) { throw 'CLI archive failed.' }
        }
    }
    if ($IncludeSource) {
    & git archive --format=tar.gz ("--prefix=inspire-go-sdk-${Version}/") ("--output=" + (Join-Path $sdkOutput "inspire-go-sdk_${Version}_source.tar.gz")) HEAD
    if ($LASTEXITCODE -ne 0) { throw 'SDK source archive failed.' }
    $sdkPluginArchive = Join-Path $sdkOutput "inspire-go-sdk-plugin_${Version}.zip"
    & git archive --format=zip --prefix=inspire-go-sdk/ ("--output=" + $sdkPluginArchive) HEAD plugin.json .agents/plugins/marketplace.json skills LICENSE NOTICE README.md
    if ($LASTEXITCODE -ne 0) { throw 'Plugin archive failed.' }
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    [IO.Compression.ZipFile]::ExtractToDirectory($sdkPluginArchive, $sdkStage)
    $sdkPluginRoot = Join-Path $sdkStage 'inspire-go-sdk'
    $sdkSkillRoot = Join-Path $sdkStage 'skill'
    New-Item -ItemType Directory -Path $sdkSkillRoot | Out-Null
    Copy-Item -LiteralPath (Join-Path $sdkPluginRoot 'skills/inspire-go-sdk') -Destination $sdkSkillRoot -Recurse
    Copy-Item -LiteralPath (Join-Path $sdkPluginRoot 'LICENSE') -Destination (Join-Path $sdkSkillRoot 'inspire-go-sdk/LICENSE')
    New-SdkZip $sdkSkillRoot (Join-Path $sdkOutput "inspire-go-sdk-skill_${Version}.zip")
    }
} finally { Pop-Location }
