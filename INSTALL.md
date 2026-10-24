# Install

You need 64-bit Windows 10 or 11. You do not need Python, Node, or PowerShell extras. Edge or Chrome is used as the window. If neither is installed, your default browser opens instead.

## 1. Run the installer

Quit Turnip Prophet first if it is open.

Download the repo ZIP, extract it, and double-click `TurnipProphet-Setup.exe`. Windows asks for permission because the app installs into Program Files. That is normal.

Or from PowerShell in the extracted folder:

```powershell
powershell -ExecutionPolicy Bypass -File .\Install.ps1
```

That opens the same installer.

Default folder: `C:\Program Files\Turnip Prophet`.

What it does:

- Installs the program there
- Puts **Turnip Prophet** on the desktop and in the Start menu
- Adds an entry in Windows Settings, Apps, so you can uninstall it
- Does **not** delete your weeks. Those live in `%AppData%\TurnipProphet\turnips.json`
- If an older install left `C:\Tools\Turnips\turnips.json`, the first launch copies it into AppData and renames the old file `turnips.json.migrated`

The bottom of the window shows the real path to the weeks file.

## 2. Use it

Open the desktop shortcut or the Start menu item. A second open does not start a second copy. It brings the window back.

Quit from the button in the window, or right-click the tray icon and choose Quit. Closing the window alone leaves the app running on purpose.

## 3. Update later

Quit first. Run `TurnipProphet-Setup.exe` again and install over the same folder. The weeks file stays where it is.

## 4. Remove it

Quit, then Windows Settings, Apps, Turnip Prophet, Uninstall. That removes the program and the shortcuts. It does not remove `%AppData%\TurnipProphet`. Delete that folder only if you want the history gone.

An old `C:\Tools\Turnips` folder is not part of this installer. Delete it yourself once you have seen your weeks in the new app.

## If it does not open

The log is `%LocalAppData%\TurnipProphet\turnips.log`. On a normal account that is `C:\Users\<you>\AppData\Local\TurnipProphet\turnips.log`.

- "No free local port" — something is holding ports 47821–47840. Quit old copies in Task Manager. Look for Turnip Prophet and for Windows PowerShell from the old version.
- The window is blank — Edge is still starting. Give it a second. If it stays blank, quit from the tray and reopen.
- SmartScreen — this installer is not a signed store app. More info, then Run anyway, if you trust this repo.

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
