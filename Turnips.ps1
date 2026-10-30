# Animal Crossing: New Horizons - Turnip Prophet
# Same page as the Program Files app. Weeks stay in %AppData%\TurnipProphet.
$ErrorActionPreference = 'Stop'

$psExe = Join-Path $env:SystemRoot 'System32\WindowsPowerShell\v1.0\powershell.exe'
if ([Threading.Thread]::CurrentThread.ApartmentState -ne 'STA') {
  if ($env:TURNIP_STA_RELAUNCH -eq '1') {
    Add-Type -AssemblyName System.Windows.Forms
    [System.Windows.Forms.MessageBox]::Show('Could not start Turnip Prophet in STA mode.', 'Turnip Prophet') | Out-Null
    exit 1
  }
  $env:TURNIP_STA_RELAUNCH = '1'
  $relaunch = '-NoProfile -ExecutionPolicy Bypass -STA -File "{0}"' -f $PSCommandPath
  Start-Process -FilePath $psExe -WindowStyle Hidden -ArgumentList $relaunch
  exit 0
}

Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

$Root = $PSScriptRoot
if (-not $Root) { $Root = Split-Path -Parent $MyInvocation.MyCommand.Path }

$script:OwnsServer = $false
$script:Cleaned = $false
$script:StopServer = $false
$script:Listener = $null
$script:Icon = $null
$script:AppIcon = $null
$script:Timer = $null
$script:Port = 0
$script:DataFile = $null

function Get-TurnipDataFile {
  if ($script:DataFile) { return $script:DataFile }
  $base = $env:APPDATA
  if (-not $base) { $base = $Root }
  $dir = Join-Path $base 'TurnipProphet'
  if (-not (Test-Path -LiteralPath $dir)) {
    New-Item -ItemType Directory -Path $dir -Force | Out-Null
  }
  $dest = Join-Path $dir 'turnips.json'
  $old = Join-Path $Root 'turnips.json'
  if (-not (Test-Path -LiteralPath $dest) -and (Test-Path -LiteralPath $old)) {
    Copy-Item -LiteralPath $old -Destination $dest -Force
    $bak = Join-Path $Root 'turnips.json.migrated'
    if (-not (Test-Path -LiteralPath $bak)) {
      Rename-Item -LiteralPath $old -NewName 'turnips.json.migrated' -ErrorAction SilentlyContinue
    }
    Write-Log ("Copied weeks to " + $dest)
  }
  $script:DataFile = $dest
  return $dest
}

function Write-Log([string]$Message) {
  try {
    $line = (Get-Date -Format 'yyyy-MM-dd HH:mm:ss') + '  ' + $Message + "`r`n"
    $enc = New-Object System.Text.UTF8Encoding $false
    [IO.File]::AppendAllText((Join-Path $Root 'turnips.log'), $line, $enc)
  } catch {
  }
}

function Show-Error([string]$Message) {
  Write-Log $Message
  try {
    [System.Windows.Forms.MessageBox]::Show($Message, 'Turnip Prophet') | Out-Null
  } catch {
  }
}

function Find-Browser {
  $candidates = New-Object System.Collections.Generic.List[string]
  $names = @(
    'Microsoft\Edge\Application\msedge.exe',
    'Google\Chrome\Application\chrome.exe'
  )
  $roots = @($env:ProgramFiles, ${env:ProgramFiles(x86)}, $env:LocalAppData)
  foreach ($root in $roots) {
    if (-not $root) { continue }
    foreach ($name in $names) {
      [void]$candidates.Add((Join-Path $root $name))
    }
  }
  foreach ($path in $candidates) {
    if ($path -and (Test-Path -LiteralPath $path)) { return $path }
  }
  return $null
}

function Set-ChromeWindowSize {
  $pref = Join-Path $Root 'profile\Default\Preferences'
  if (-not (Test-Path -LiteralPath $pref)) { return }
  try {
    $text = [IO.File]::ReadAllText($pref)
    $m = [regex]::Match($text, '"window_placement":\{[^{}]*\}')
    if (-not $m.Success) { return }
    $block = $m.Value
    $block2 = [regex]::Replace($block, '"maximized":(true|false)', '"maximized":false')
    $block2 = [regex]::Replace($block2, '"left":-?\d+', '"left":40')
    $block2 = [regex]::Replace($block2, '"top":-?\d+', '"top":40')
    $block2 = [regex]::Replace($block2, '"right":-?\d+', '"right":1400')
    $block2 = [regex]::Replace($block2, '"bottom":-?\d+', '"bottom":1000')
    if ($block2 -eq $block) { return }
    $text2 = $text.Remove($m.Index, $m.Length).Insert($m.Index, $block2)
    $enc = New-Object System.Text.UTF8Encoding $false
    [IO.File]::WriteAllText($pref, $text2, $enc)
  } catch {
  }
}

function Open-AppWindow([int]$Port, [bool]$Fresh) {
  $url = "http://127.0.0.1:$Port/"
  $browser = Find-Browser
  if ($browser) {
    $profile = Join-Path $Root 'profile'
    if (-not (Test-Path -LiteralPath $profile)) {
      New-Item -ItemType Directory -Path $profile | Out-Null
    }
    if ($Fresh) { Set-ChromeWindowSize }
    Start-Process -FilePath $browser -ArgumentList @(
      "--user-data-dir=$profile",
      "--app=$url",
      '--window-size=1360,960',
      '--no-first-run',
      '--disable-sync'
    ) | Out-Null
  } else {
    Start-Process $url | Out-Null
  }
}

function Test-BrowserOpen {
  $needle = '*' + $Root + '\profile*'
  try {
    $hits = @(Get-CimInstance Win32_Process -ErrorAction Stop | Where-Object {
      ($_.Name -eq 'msedge.exe' -or $_.Name -eq 'chrome.exe') -and $_.CommandLine -like $needle
    })
    return $hits.Count -gt 0
  } catch {
    return $true
  }
}

function Test-TurnipPort([int]$Port) {
  $client = $null
  try {
    $client = New-Object System.Net.Sockets.TcpClient
    $client.ReceiveTimeout = 800
    $client.SendTimeout = 800
    $client.Connect('127.0.0.1', $Port)
    $stream = $client.GetStream()
    $req = [Text.Encoding]::ASCII.GetBytes("GET /api/info HTTP/1.1`r`nHost: 127.0.0.1`r`nConnection: close`r`n`r`n")
    $stream.Write($req, 0, $req.Length)
    $buf = New-Object byte[] 2048
    $n = $stream.Read($buf, 0, $buf.Length)
    if ($n -le 0) { return $false }
    $text = [Text.Encoding]::ASCII.GetString($buf, 0, $n)
    return $text.Contains('"app":"turnip-prophet"')
  } catch {
    return $false
  } finally {
    if ($client) { $client.Close() }
  }
}

function Send-Response($Stream, [string]$Status, [string]$ContentType, [byte[]]$Body) {
  if ($null -eq $Body) { $Body = New-Object byte[] 0 }
  $header = "HTTP/1.1 $Status`r`nContent-Type: $ContentType`r`nContent-Length: $($Body.Length)`r`nConnection: close`r`nCache-Control: no-store`r`n`r`n"
  $hb = [Text.Encoding]::ASCII.GetBytes($header)
  $Stream.Write($hb, 0, $hb.Length)
  if ($Body.Length -gt 0) { $Stream.Write($Body, 0, $Body.Length) }
  $Stream.Flush()
}

function Receive-HttpRequest($Client) {
  $stream = $Client.GetStream()
  $stream.ReadTimeout = 8000
  $Client.ReceiveTimeout = 8000
  $ms = New-Object System.IO.MemoryStream
  $tmp = New-Object byte[] 4096
  $headerEnd = -1
  while ($headerEnd -lt 0) {
    $n = $stream.Read($tmp, 0, $tmp.Length)
    if ($n -le 0) { return $null }
    $ms.Write($tmp, 0, $n)
    $buf = $ms.ToArray()
    $limit = $buf.Length - 4
    for ($i = 0; $i -le $limit; $i++) {
      if ($buf[$i] -eq 13 -and $buf[$i + 1] -eq 10 -and $buf[$i + 2] -eq 13 -and $buf[$i + 3] -eq 10) {
        $headerEnd = $i
        break
      }
    }
    if ($headerEnd -lt 0 -and $ms.Length -gt 65536) { return $null }
  }
  $all = $ms.ToArray()
  $headerText = [Text.Encoding]::ASCII.GetString($all, 0, $headerEnd)
  $lines = $headerText -split "`r`n"
  $req = $lines[0].Split(' ')
  if ($req.Length -lt 2) { return $null }
  $method = $req[0]
  $path = $req[1]
  $q = $path.IndexOf('?')
  if ($q -ge 0) { $path = $path.Substring(0, $q) }
  $length = 0
  foreach ($line in $lines) {
    if ($line -match '^(?i)Content-Length:\s*(\d+)\s*$') {
      $length = [int]$Matches[1]
    }
  }
  if ($length -gt 2000000) { return $null }
  $body = New-Object byte[] $length
  $filled = 0
  $have = $all.Length - ($headerEnd + 4)
  if ($have -gt 0 -and $length -gt 0) {
    $take = [Math]::Min($have, $length)
    [Array]::Copy($all, $headerEnd + 4, $body, 0, $take)
    $filled = $take
  }
  while ($filled -lt $length) {
    $n = $stream.Read($tmp, 0, $tmp.Length)
    if ($n -le 0) { break }
    $take = [Math]::Min($n, $length - $filled)
    [Array]::Copy($tmp, 0, $body, $filled, $take)
    $filled += $take
  }
  return @{ Method = $method; Path = $path; Body = $body; Stream = $stream }
}

function Get-StateBytes {
  $path = Get-TurnipDataFile
  if (Test-Path -LiteralPath $path) {
    $bytes = [IO.File]::ReadAllBytes($path)
    return ,$bytes
  }
  $enc = New-Object System.Text.UTF8Encoding $false
  $bytes = $enc.GetBytes('{"weeks":[],"activeId":null}')
  return ,$bytes
}

function Test-StateJson([string]$Text) {
  if ([string]::IsNullOrWhiteSpace($Text)) { return $false }
  if ($Text.Length -gt 2000000) { return $false }
  try {
    $obj = $Text | ConvertFrom-Json
  } catch {
    return $false
  }
  if ($null -eq $obj.PSObject.Properties['weeks']) { return $false }
  return $true
}

function Save-StateBytes([byte[]]$Body) {
  $path = Get-TurnipDataFile
  $dir = Split-Path -Parent $path
  $tmp = Join-Path $dir 'turnips.json.tmp'
  [IO.File]::WriteAllBytes($tmp, $Body)
  [IO.File]::Copy($tmp, $path, $true)
  Remove-Item -LiteralPath $tmp -Force
}

function Serve-Client {
  $client = $script:Listener.AcceptTcpClient()
  try {
    $req = Receive-HttpRequest $client
    if (-not $req) { return }
    $method = $req.Method.ToUpperInvariant()
    $path = $req.Path
    if ($method -eq 'GET' -and $path -eq '/api/info') {
      $jsonRoot = $Root.Replace('\', '\\')
      $dataPath = (Get-TurnipDataFile).Replace('\', '\\')
      $dark = 'false'
      try {
        $theme = Get-ItemProperty -Path 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize' -Name AppsUseLightTheme -ErrorAction Stop
        if ([int]$theme.AppsUseLightTheme -eq 0) { $dark = 'true' }
      } catch {}
      $json = '{"app":"turnip-prophet","version":1,"root":"' + $jsonRoot + '","data":"' + $dataPath + '","dark":' + $dark + '}'
      $payload = [Text.Encoding]::UTF8.GetBytes($json)
      Send-Response $req.Stream '200 OK' 'application/json; charset=utf-8' $payload
      return
    }
    if ($method -eq 'GET' -and $path -eq '/api/state') {
      $payload = Get-StateBytes
      Send-Response $req.Stream '200 OK' 'application/json; charset=utf-8' $payload
      return
    }
    if ($method -eq 'POST' -and $path -eq '/api/state') {
      $text = [Text.Encoding]::UTF8.GetString($req.Body)
      if (-not (Test-StateJson $text)) {
        $payload = [Text.Encoding]::UTF8.GetBytes('{"ok":false}')
        Send-Response $req.Stream '400 Bad Request' 'application/json; charset=utf-8' $payload
        return
      }
      Save-StateBytes $req.Body
      $payload = [Text.Encoding]::UTF8.GetBytes('{"ok":true}')
      Send-Response $req.Stream '200 OK' 'application/json; charset=utf-8' $payload
      return
    }
    if ($method -eq 'POST' -and $path -eq '/api/exit-saved') {
      $payload = [Text.Encoding]::UTF8.GetBytes('{"ok":true}')
      Send-Response $req.Stream '200 OK' 'application/json; charset=utf-8' $payload
      return
    }
    if ($method -eq 'POST' -and $path -eq '/api/shutdown') {
      $payload = [Text.Encoding]::UTF8.GetBytes('{"ok":true}')
      Send-Response $req.Stream '200 OK' 'application/json; charset=utf-8' $payload
      $script:StopServer = $true
      return
    }
    if ($method -eq 'GET' -and ($path -eq '/favicon.ico' -or $path -eq '/TurnipProphet.ico')) {
      $full = Join-Path $Root 'TurnipProphet.ico'
      if (Test-Path -LiteralPath $full) {
        $payload = [IO.File]::ReadAllBytes($full)
        Send-Response $req.Stream '200 OK' 'image/x-icon' $payload
      } else {
        Send-Response $req.Stream '204 No Content' 'text/plain; charset=utf-8' (New-Object byte[] 0)
      }
      return
    }
    if ($method -eq 'GET' -and $script:Files.ContainsKey($path)) {
      $full = Join-Path $Root $script:Files[$path].Rel
      if (-not (Test-Path -LiteralPath $full)) {
        $payload = [Text.Encoding]::ASCII.GetBytes('missing')
        Send-Response $req.Stream '404 Not Found' 'text/plain; charset=utf-8' $payload
        return
      }
      $payload = [IO.File]::ReadAllBytes($full)
      Send-Response $req.Stream '200 OK' $script:Files[$path].Type $payload
      return
    }
    $payload = [Text.Encoding]::ASCII.GetBytes('not found')
    Send-Response $req.Stream '404 Not Found' 'text/plain; charset=utf-8' $payload
  } catch {
    Write-Log $_.Exception.Message
  } finally {
    $client.Close()
  }
}

function Start-TurnipListener {
  for ($p = 47821; $p -le 47840; $p++) {
    $listener = New-Object System.Net.Sockets.TcpListener ([System.Net.IPAddress]::Loopback, $p)
    try {
      $listener.Start()
      return @{ Listener = $listener; Port = $p }
    } catch {
      try { $listener.Stop() } catch {}
    }
  }
  return $null
}

function Stop-Turnip {
  if ($script:Cleaned) { return }
  $script:Cleaned = $true
  if (-not $script:OwnsServer) { return }
  try { if ($script:Timer) { $script:Timer.Stop(); $script:Timer.Dispose() } } catch {}
  try { if ($script:Listener) { $script:Listener.Stop() } } catch {}
  try {
    if ($script:Icon) {
      $script:Icon.Visible = $false
      $script:Icon.Dispose()
    }
    if ($script:AppIcon) { $script:AppIcon.Dispose() }
  } catch {}
  $portFile = Join-Path $Root 'turnips.port'
  if (Test-Path -LiteralPath $portFile) {
    Remove-Item -LiteralPath $portFile -Force -ErrorAction SilentlyContinue
  }
  Write-Log 'Stopped.'
}

function Main {
  Get-TurnipDataFile | Out-Null
  foreach ($needed in @('predictions.js', 'TurnipProphet.ico', 'www\index.html', 'www\app.js', 'www\styles.css', 'www\logo.png', 'www\field.jpg', 'www\turnip.png', 'www\stamp.png', 'www\wood.png')) {
    $full = Join-Path $Root $needed
    if (-not (Test-Path -LiteralPath $full)) {
      Show-Error "Missing $needed. Run Install.ps1 again from the repo."
      return
    }
  }

  $portFile = Join-Path $Root 'turnips.port'
  if (Test-Path -LiteralPath $portFile) {
    $raw = (Get-Content -LiteralPath $portFile -Raw).Trim()
    $old = 0
    if ([int]::TryParse($raw, [ref]$old) -and $old -gt 0 -and (Test-TurnipPort $old)) {
      Open-AppWindow -Port $old -Fresh $false
      return
    }
    Remove-Item -LiteralPath $portFile -Force -ErrorAction SilentlyContinue
  }

  $started = Start-TurnipListener
  if (-not $started) {
    Show-Error 'No free local port between 47821 and 47840. Quit old Turnip Prophet copies and try again.'
    return
  }

  $script:OwnsServer = $true
  $script:Listener = $started.Listener
  $script:Port = [int]$started.Port
  $enc = New-Object System.Text.UTF8Encoding $false
  [IO.File]::WriteAllText($portFile, "$($script:Port)", $enc)
  Write-Log ("Listening on 127.0.0.1:" + $script:Port)

  $script:Files = @{
    '/' = @{ Rel = 'www\index.html'; Type = 'text/html; charset=utf-8' }
    '/index.html' = @{ Rel = 'www\index.html'; Type = 'text/html; charset=utf-8' }
    '/styles.css' = @{ Rel = 'www\styles.css'; Type = 'text/css; charset=utf-8' }
    '/app.js' = @{ Rel = 'www\app.js'; Type = 'text/javascript; charset=utf-8' }
    '/predictions.js' = @{ Rel = 'predictions.js'; Type = 'text/javascript; charset=utf-8' }
    '/logo.png' = @{ Rel = 'www\logo.png'; Type = 'image/png' }
    '/field.jpg' = @{ Rel = 'www\field.jpg'; Type = 'image/jpeg' }
    '/turnip.png' = @{ Rel = 'www\turnip.png'; Type = 'image/png' }
    '/stamp.png' = @{ Rel = 'www\stamp.png'; Type = 'image/png' }
    '/wood.png' = @{ Rel = 'www\wood.png'; Type = 'image/png' }
  }

  $script:AppContext = New-Object System.Windows.Forms.ApplicationContext
  $ctx = $script:AppContext
  $script:WindowSeen = $false
  $script:WindowTick = 0
  $script:Timer = New-Object System.Windows.Forms.Timer
  $script:Timer.Interval = 40
  $script:Timer.Add_Tick({
    try {
      $guard = 0
      while ($script:Listener.Pending() -and $guard -lt 30) {
        Serve-Client
        $guard++
      }
    } catch {
      Write-Log $_.Exception.Message
    }
    $script:WindowTick++
    if (($script:WindowTick % 15) -eq 0) {
      $open = Test-BrowserOpen
      if ($open) { $script:WindowSeen = $true }
      elseif ($script:WindowSeen) { $script:StopServer = $true }
    }
    if ($script:StopServer) {
      $script:Timer.Stop()
      $script:AppContext.ExitThread()
    }
  })
  $script:Timer.Start()

  try { Open-AppWindow -Port $script:Port -Fresh $true } catch { Write-Log $_.Exception.Message }

  [System.Windows.Forms.Application]::Run($ctx)
}

try {
  Main
} catch {
  Show-Error $_.Exception.Message
} finally {
  Stop-Turnip
}
