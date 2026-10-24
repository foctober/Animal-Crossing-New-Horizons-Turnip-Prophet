# Install

You need 64-bit Windows 10 or 11. Edge or Chrome is used as the window. If neither is installed, your default browser opens instead.

Pick one install. Do not run both at once. They share the same weeks file, `%AppData%\TurnipProphet\turnips.json`.

Quit Turnip Prophet first if it is already open.

## Windows installer

Double-click `TurnipProphet-Setup.exe`. Windows asks for permission because the app installs into Program Files.

Default folder: `C:\Program Files\Turnip Prophet`.

- Installs the program there
- Puts **Turnip Prophet** on the public desktop and in the Start menu
- Adds an entry in Windows Settings, Apps
- Does not delete the weeks file

Uninstall later from Windows Settings, Apps, Turnip Prophet. That does not remove `%AppData%\TurnipProphet`.

To update, quit, then run `TurnipProphet-Setup.exe` again over the same folder.

SmartScreen may say the installer is unrecognized. It is not a signed store app. More info, then Run anyway, if you trust this repo.

The log for this copy is `%LocalAppData%\TurnipProphet\turnips.log`.

## PowerShell

From the extracted folder:

```powershell
powershell -ExecutionPolicy Bypass -File .\Install.ps1
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

Then run `Install.ps1` again as your normal user.

To update, quit, then from a fresh copy of this folder run `Install.ps1` again. You can also double-click `C:\Tools\Turnips\Run.cmd`.

To remove it, quit, delete the desktop shortcut, and delete `C:\Tools\Turnips`. That does not remove `%AppData%\TurnipProphet`.

The log for this copy is `C:\Tools\Turnips\turnips.log`.

## Either way

A second open does not start a second copy. It brings the window back.

Quit from the button in the window, or right-click the tray icon and choose Quit. Closing the window alone leaves the app running on purpose.

The bottom of the window shows the real path to the weeks file. On a normal account that is `C:\Users\<you>\AppData\Roaming\TurnipProphet\turnips.json`.

If AppData has no weeks file yet, the first launch copies `C:\Tools\Turnips\turnips.json` when that old file is still there, and renames it `turnips.json.migrated`.

## If it does not open

- "No free local port" means something is holding ports 47821–47840. Quit the other copy. In Task Manager, end **Turnip Prophet** or a Windows PowerShell process whose command line contains `Turnips.ps1`.
- The window is blank — Edge is still starting. Give it a second. If it stays blank, quit from the tray and reopen.

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
