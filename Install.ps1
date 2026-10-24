# Opens the Windows installer. The app goes in Program Files.
$ErrorActionPreference = 'Stop'

$here = $PSScriptRoot
if (-not $here) { $here = Split-Path -Parent $MyInvocation.MyCommand.Path }
$setup = Join-Path $here 'TurnipProphet-Setup.exe'
if (-not (Test-Path -LiteralPath $setup)) {
  Write-Host ''
  Write-Host 'TurnipProphet-Setup.exe is missing next to Install.ps1.'
  Write-Host ''
  exit 1
}

Write-Host 'Windows will ask for permission.'
Write-Host 'Turnip Prophet installs to C:\Program Files\Turnip Prophet.'
Write-Host 'Your weeks stay in %AppData%\TurnipProphet and are not deleted.'
Start-Process -FilePath $setup -Verb RunAs -Wait
