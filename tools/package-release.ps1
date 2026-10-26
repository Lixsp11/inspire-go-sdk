param(
    [Parameter(Mandatory=$true)][string]$Version,
    [string]$OutputDirectory = 'dist'
)
$ErrorActionPreference = 'Stop'
if ($Version -notmatch '^\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$') { throw 'Invalid release version.' }
$sdkRoot = Split-Path -Parent $PSScriptRoot
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
    & git archive --format=tar.gz ("--prefix=inspire-go-sdk-${Version}/") ("--output=" + (Join-Path $sdkOutput "inspire-go-sdk_${Version}_source.tar.gz")) HEAD
    if ($LASTEXITCODE -ne 0) { throw 'SDK source archive failed.' }
    $sdkSkillRoot = Join-Path $sdkStage 'skill'
    New-Item -ItemType Directory -Path $sdkSkillRoot | Out-Null
    Copy-Item -LiteralPath (Join-Path $sdkRoot 'skills/inspire-go-sdk') -Destination $sdkSkillRoot -Recurse
    Copy-Item -LiteralPath (Join-Path $sdkRoot 'LICENSE') -Destination (Join-Path $sdkSkillRoot 'inspire-go-sdk/LICENSE')
    New-SdkZip $sdkSkillRoot (Join-Path $sdkOutput "inspire-go-sdk-skill_${Version}.zip")
} finally { Pop-Location }
