# PowerShell install. Copies the app to C:\Tools\Turnips and adds a desktop shortcut.
# The other choice is TurnipProphet-Setup.exe, which installs into Program Files.
# Pick one. Do not run both at once. Both save weeks in %AppData%\TurnipProphet\turnips.json.
$ErrorActionPreference = 'Stop'

$here = $PSScriptRoot
if (-not $here) { $here = Split-Path -Parent $MyInvocation.MyCommand.Path }
$dest = 'C:\Tools\Turnips'

function Fail([string]$Message) {
  Write-Host ''
  Write-Host $Message
  Write-Host ''
  Write-Host 'Install stopped.'
  exit 1
}

try {
  if (-not (Test-Path -LiteralPath 'C:\Tools')) {
    New-Item -ItemType Directory -Path 'C:\Tools' -ErrorAction Stop | Out-Null
  }
  if (-not (Test-Path -LiteralPath $dest)) {
    New-Item -ItemType Directory -Path $dest -ErrorAction Stop | Out-Null
  }
  $probe = Join-Path $dest '.write-test'
  [IO.File]::WriteAllText($probe, 'ok')
  Remove-Item -LiteralPath $probe -Force
} catch {
  Fail "Could not write to C:\Tools\Turnips.`r`n`r`nCreate that folder as Administrator, then run Install.ps1 again as your normal user:`r`n  New-Item -ItemType Directory -Force -Path C:\Tools\Turnips"
}

$www = Join-Path $dest 'www'
if (-not (Test-Path -LiteralPath $www)) {
  New-Item -ItemType Directory -Path $www | Out-Null
}

$names = @('Turnips.ps1', 'Launch.vbs', 'Run.cmd', 'predictions.js', 'LICENSE', 'NOTICE', 'README.md', 'INSTALL.md', 'Install.ps1', 'TurnipProphet.ico')
foreach ($name in $names) {
  $src = Join-Path $here $name
  if (-not (Test-Path -LiteralPath $src)) { Fail "Missing $name next to Install.ps1." }
  Copy-Item -LiteralPath $src -Destination (Join-Path $dest $name) -Force
}

foreach ($name in @('index.html', 'app.js', 'styles.css', 'logo.png', 'field.jpg', 'turnip.png', 'stamp.png', 'wood.png', 'FinkHeavy.ttf')) {
  $src = Join-Path (Join-Path $here 'www') $name
  if (-not (Test-Path -LiteralPath $src)) { Fail "Missing www\$name." }
  Copy-Item -LiteralPath $src -Destination (Join-Path $www $name) -Force
}

Get-ChildItem -LiteralPath $dest -Recurse -File | Unblock-File -ErrorAction SilentlyContinue

$desktop = [Environment]::GetFolderPath('Desktop')
$lnk = Join-Path $desktop 'Turnip Prophet.lnk'
$wsh = New-Object -ComObject WScript.Shell
$sc = $wsh.CreateShortcut($lnk)
$sc.TargetPath = Join-Path $env:SystemRoot 'System32\wscript.exe'
$sc.Arguments = '"' + (Join-Path $dest 'Launch.vbs') + '"'
$sc.WorkingDirectory = $dest
$sc.WindowStyle = 7
$sc.Description = 'Animal Crossing: New Horizons - Turnip Prophet'
$leaf = Join-Path $dest 'TurnipProphet.ico'
if (Test-Path -LiteralPath $leaf) {
  $sc.IconLocation = $leaf + ',0'
} else {
  $edge = Join-Path $env:ProgramFiles 'Microsoft\Edge\Application\msedge.exe'
  $psIcon = Join-Path $env:SystemRoot 'System32\WindowsPowerShell\v1.0\powershell.exe'
  if (Test-Path -LiteralPath $edge) {
    $sc.IconLocation = $edge + ',0'
  } else {
    $sc.IconLocation = $psIcon + ',0'
  }
}
$sc.Save()

Write-Host "Installed to $dest"
Write-Host "Shortcut: $lnk"
Write-Host 'Opening Turnip Prophet...'
Write-Host 'Quit from the Quit button, or close the window. Closing saves first. There is no tray icon.'

$launcher = Join-Path $dest 'Launch.vbs'
Start-Process -FilePath (Join-Path $env:SystemRoot 'System32\wscript.exe') -ArgumentList ('"' + $launcher + '"')
