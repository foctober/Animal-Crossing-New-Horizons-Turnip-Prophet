# Install

You need Windows 10 or 11. Windows PowerShell 5.1 is already on the machine. You do not need to install Python, Node, or anything else.

Edge or Chrome is used as the window. If neither is installed, your default browser opens instead.

## 1. Get the files

Clone:

```powershell
git clone https://github.com/foctober/Animal-Crossing-New-Horizons-Turnip-Prophet.git
cd Animal-Crossing-New-Horizons-Turnip-Prophet
```

Or download the ZIP from GitHub, extract it, and open PowerShell in that folder.

If Windows marked the download as blocked:

```powershell
Get-ChildItem -Recurse -File | Unblock-File
```

## 2. Install to C:\Tools\Turnips

```powershell
powershell -ExecutionPolicy Bypass -File .\Install.ps1
```

What that does:

- Creates `C:\Tools\Turnips`
- Copies the app into it
- Does **not** delete `turnips.json` if you already have weeks saved
- Puts a **Turnip Prophet** shortcut on your desktop
- Opens the app

The shortcut runs `Launch.vbs`, which starts PowerShell hidden and opens the tracker.

### If Windows will not create C:\Tools

Some PCs block new folders on `C:\`. In that case Install.ps1 stops and says so. Open PowerShell **as Administrator** once and run:

```powershell
New-Item -ItemType Directory -Force -Path C:\Tools\Turnips
icacls C:\Tools /grant "$($env:USERNAME):(OI)(CI)M"
```

Then run `Install.ps1` again as your normal user. The app itself does not need to run as Administrator. It only listens on `127.0.0.1` (this PC), so it should not ask the firewall for a public port.

## 3. Use it

- Desktop shortcut **Turnip Prophet**, or
- `C:\Tools\Turnips\Run.cmd`, or
- from PowerShell:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -STA -File C:\Tools\Turnips\Turnips.ps1
```

The first open loads the week of October 18, 2026 (buy 96, the prices from that week's link) so you can see a finished prediction immediately. Delete that week if you do not want it.

Your history is `C:\Tools\Turnips\turnips.json`.

## 4. Quit

Use **Quit** in the window, or right-click the Turnip Prophet tray icon and choose Quit.

Closing the tracker window does not quit the helper. Double-click the shortcut again and the same window comes back. That is on purpose, so you do not lose the process and start two copies.

## 5. Update later

Quit first. Then:

```powershell
cd \path\to\Animal-Crossing-New-Horizons-Turnip-Prophet
git pull
powershell -ExecutionPolicy Bypass -File .\Install.ps1
```

Open the shortcut again. `turnips.json` stays where it is.

## 6. Remove it

Quit, delete the desktop shortcut, delete `C:\Tools\Turnips`. Take a copy of `turnips.json` first if you care about the old weeks.

## Where the shortcut points

| | |
| --- | --- |
| Target | `wscript.exe` |
| Arguments | `C:\Tools\Turnips\Launch.vbs` |
| Working directory | `C:\Tools\Turnips` |

`Launch.vbs` runs:

```text
powershell.exe -NoProfile -ExecutionPolicy Bypass -STA -File C:\Tools\Turnips\Turnips.ps1
```

`-STA` is required. The tray icon is a Windows Forms control and it misbehaves otherwise. The script relaunches itself in STA if you forget.

## If it does not open

Look at `C:\Tools\Turnips\turnips.log`.

- "No free local port" — something is holding ports 47821–47840. Quit old copies in Task Manager (look for Windows PowerShell) and try again.
- The window is blank — Edge is still starting. Give it a second. If it stays blank, quit from the tray and reopen.
- SmartScreen or an antivirus warning on the `.ps1` — this is a local script, not a signed store app. Unblock the file or allow it if you trust this repo.

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
      "note": ""
    }
  ],
  "activeId": "2026-10-18"
}
```

`id` is the Sunday that starts the week. `prices` is 13 numbers: Sunday buy, then Monday morning, Monday afternoon, Tuesday morning, and so on through Saturday afternoon. `null` means you did not type that one. `previousPattern` is -1 unknown, 0 fluctuating, 1 large spike, 2 decreasing, 3 small spike.
