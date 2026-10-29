# Install

Version 1.2. You need 64-bit Windows 10 or 11. The window is the app. If it cannot start, your browser opens instead. There is no system tray icon.

Pick one install. Do not run both at once. They share the same weeks file, `%AppData%\TurnipProphet\turnips.json`.

Quit Turnip Prophet first if it is already open. Paste one whole block into PowerShell. It downloads the repo, installs, and deletes the temp folder.

## Windows installer

```powershell
$d = Join-Path $env:TEMP 'TurnipProphet'
if (Test-Path $d) { Remove-Item $d -Recurse -Force }
New-Item -ItemType Directory -Path $d | Out-Null
Invoke-WebRequest 'https://github.com/foctober/Animal-Crossing-New-Horizons-Turnip-Prophet/archive/refs/heads/main.zip' -OutFile "$d\src.zip"
Expand-Archive "$d\src.zip" $d -Force
Start-Process -FilePath "$d\Animal-Crossing-New-Horizons-Turnip-Prophet-main\TurnipProphet-Setup.exe" -Verb RunAs -Wait
Remove-Item $d -Recurse -Force
```

Windows asks for permission because the app installs into Program Files.

Default folder: `C:\Program Files\Turnip Prophet`.

- Installs the program there
- Puts **Turnip Prophet** in the Start menu
- Desktop icon is a checkbox on the installer. It is checked unless you turn it off
- Adds an entry in Windows Settings, Apps
- Does not delete the weeks file

Uninstall later from Windows Settings, Apps, Turnip Prophet. That does not remove `%AppData%\TurnipProphet`.

To update, quit, then paste the same block again.

SmartScreen may say the installer is unrecognized. It is not a signed store app. More info, then Run anyway, if you trust this repo.

The log for this copy is `%LocalAppData%\TurnipProphet\turnips.log`.

## PowerShell

Same app as the installer. Same page, same weeks file, same 1360 by 960 window. Closing the window saves and quits. There is no tray icon.

```powershell
$d = Join-Path $env:TEMP 'TurnipProphet'
if (Test-Path $d) { Remove-Item $d -Recurse -Force }
New-Item -ItemType Directory -Path $d | Out-Null
Invoke-WebRequest 'https://github.com/foctober/Animal-Crossing-New-Horizons-Turnip-Prophet/archive/refs/heads/main.zip' -OutFile "$d\src.zip"
Expand-Archive "$d\src.zip" $d -Force
powershell -ExecutionPolicy Bypass -File "$d\Animal-Crossing-New-Horizons-Turnip-Prophet-main\Install.ps1"
Remove-Item $d -Recurse -Force
```

- Copies the app to `C:\Tools\Turnips`
- Puts **Turnip Prophet** on your own desktop
- Opens the app
- Does not need an admin prompt unless Windows blocks new folders on `C:\`

The shortcut runs `Launch.vbs`, which starts PowerShell hidden.

If Windows will not create `C:\Tools`, open PowerShell **as Administrator** once:

```powershell
New-Item -ItemType Directory -Force -Path C:\Tools\Turnips
icacls C:\Tools /grant "$($env:USERNAME):(OI)(CI)M"
```

Then paste the PowerShell block again as your normal user.

To update, quit, then paste the same block again. You can also double-click `C:\Tools\Turnips\Run.cmd`.

To remove it, quit, delete the desktop shortcut, and delete `C:\Tools\Turnips`. That does not remove `%AppData%\TurnipProphet`.

The log for this copy is `C:\Tools\Turnips\turnips.log`.

## Either way

A second open does not start a second copy. It brings the window back.

Closing the window saves the weeks file, then quits. The Quit button does the same. There is no tray icon. On open, the cursor is already in the sky-blue box, so you can type the next price. The title bar and the scrollbars match the Windows app color. Dark Windows gets a dark bar and dark scrollbars. Light Windows gets light ones. The prices, boxes, note, and chart stay the same. This is still version 1.2. Week cards show Buy in the open-box sky blue and High in the title yellow. Sold appears only after you pick when you sold. That price is New Week green when it matches the high, and Delete Week red when it is lower. The day next to the open price box is that same sky blue, with the same black outline. The other days stay yellow. Sunday Buy Price: sits above the box, left aligned. It is black text with a white outline, same as First time buying. When First time buying is checked, that label's outline turns the open-box sky blue. Unchecked, the outline stays white. Weekly Sell Prices: sits under that box, left aligned, in the same size and style. The window opens at 1360 by 960 so that save note sits just inside the bottom. Width is unchanged. The word noon in the shop-hours line is white with a black outline.

The bottom of the window shows the real path to the weeks file. On a normal account that is `C:\Users\<you>\AppData\Roaming\TurnipProphet\turnips.json`.

If AppData has no weeks file yet, the first launch copies `C:\Tools\Turnips\turnips.json` when that old file is still there, and renames it `turnips.json.migrated`.

## If it does not open

- "No free local port" means something is holding ports 47821–47840. Quit the other copy. In Task Manager, end **Turnip Prophet** or a Windows PowerShell process whose command line contains `Turnips.ps1`.
- The installer says it cannot open `TurnipProphet.exe` for writing. The app is still running. Click Abort, close the Turnip Prophet window, then run the installer again. If the window is already gone, end **Turnip Prophet** in Task Manager.

## Data

One file, one island:

```json
{
  "weeks": [
    {
      "id": "2026-10-18",
      "firstBuy": true,
      "previousPattern": -1,
      "prices": [96, null, null, 77, null, null, 64, null, 132, 172, 191, 149, null],
      "soldAt": 10,
      "likelyPattern": 3,
      "note": ""
    }
  ],
  "activeId": "2026-10-18"
}
```

`id` is the Sunday that starts the week. `prices` is 13 numbers: Sunday buy, then Monday morning, Monday afternoon, Tuesday morning, and so on through Saturday afternoon. `null` means you did not type that one.

Morning is 8:00 AM – 11:59 AM. Afternoon is 12:00 PM – 10:00 PM. Nook's Cranny is open 8:00 AM – 10:00 PM, and the sell price changes at noon.

`soldAt` is when you sold. `null` means you are still holding. `"spoil"` means they rotted. A number is the sell window: 1 is Monday morning, 2 is Monday afternoon, and so on. 10 is Friday, 12:00 PM – 10:00 PM. 12 is Saturday afternoon.

`previousPattern` is -1 unknown, 0 fluctuating, 1 large spike, 2 decreasing, 3 small spike. `likelyPattern` is the pattern this week is almost certainly on (80% or more), or -1 if it is still mixed. A new week copies that into "last week's pattern."
