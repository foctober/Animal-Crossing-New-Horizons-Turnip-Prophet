# Animal Crossing: New Horizons - Turnip Prophet

**Version 1.3.**

A Windows tracker for Daisy Mae's turnip prices. Enter each price through the week. It saves on this PC and tells you whether to sell or wait.

The prediction uses the same engine as [turnipprophet.io](https://turnipprophet.io): Mike Bryant's [ac-nh-turnip-prices](https://github.com/mikebryant/ac-nh-turnip-prices), based on Ninji's reverse engineering of the game. This project adds a window, a desktop shortcut, and a saved history. It is not the website, and it is not affiliated with Nintendo.

## What it looks like

The left side is your saved weeks. The right side is this week's prices.

![Entering the week's prices. The open shop window has the rainbow border.](docs/prices.png)

![The sell-or-wait banner, the chart, and the credit line](docs/chart.png)

![Every pattern line, with the same numbers as the website](docs/patterns.png)

## Install

Pick one. Do not run both at once. Both store weeks in `%AppData%\TurnipProphet\turnips.json`. Quit Turnip Prophet first if it is already open. More detail is in [INSTALL.md](INSTALL.md).

**Windows installer.** Installs to Program Files, adds a Start menu entry, and uninstalls from Settings. Paste this whole block into PowerShell:

```powershell
$d = Join-Path $env:TEMP 'TurnipProphet'
if (Test-Path $d) { Remove-Item $d -Recurse -Force }
New-Item -ItemType Directory -Path $d | Out-Null
Invoke-WebRequest 'https://github.com/foctober/Animal-Crossing-New-Horizons-Turnip-Prophet/archive/refs/heads/main.zip' -OutFile "$d\src.zip"
Expand-Archive "$d\src.zip" $d -Force
Start-Process -FilePath "$d\Animal-Crossing-New-Horizons-Turnip-Prophet-main\TurnipProphet-Setup.exe" -Verb RunAs -Wait
Remove-Item $d -Recurse -Force
```

Windows asks for permission. SmartScreen may say the installer is unrecognized. Choose More info, then Run anyway. The app is installed to `C:\Program Files\Turnip Prophet`.

**PowerShell.** Same app, without the installer. It is installed to `C:\Tools\Turnips`. Paste this whole block into PowerShell:

```powershell
$d = Join-Path $env:TEMP 'TurnipProphet'
if (Test-Path $d) { Remove-Item $d -Recurse -Force }
New-Item -ItemType Directory -Path $d | Out-Null
Invoke-WebRequest 'https://github.com/foctober/Animal-Crossing-New-Horizons-Turnip-Prophet/archive/refs/heads/main.zip' -OutFile "$d\src.zip"
Expand-Archive "$d\src.zip" $d -Force
powershell -ExecutionPolicy Bypass -File "$d\Animal-Crossing-New-Horizons-Turnip-Prophet-main\Install.ps1"
Remove-Item $d -Recurse -Force
```

The shortcut and window icon are the teal leaf. There is no system tray icon.

## What you see

The left side lists your weeks. Each card shows the date, the buy price, and the highest sell price so far. After you record a sale, that price is shown too. **New Week** starts the next Sunday.

The right side is where you enter this week's prices. The box for the shop window that is open now is marked. Below the prices, the app says whether to sell or wait, shows the chart, and can open the full pattern list.

A line at the bottom credits Mike Bryant and Ninji, and states that this project is not affiliated with Nintendo.

## What you do each week

1. On Sunday morning, enter Daisy Mae's price. It is always 90–110.
2. When you check Nook's Cranny, enter that price in the morning or afternoon box. Morning is 8:00 AM – 11:59 AM. Afternoon is 12:00 PM – 10:00 PM. Do not clear the morning price.
3. Read the banner. It says whether to sell or wait.
4. Next Sunday, click **New Week**. If this week's pattern was certain, the new week starts with it filled in as last week's pattern.

Weeks are stored in `%AppData%\TurnipProphet\turnips.json` (usually `C:\Users\<you>\AppData\Roaming\TurnipProphet\turnips.json`). Neither install stores weeks in the program folder. Closing the window saves, then quits. The Quit button does the same.

## Saving

Changes save on their own, about a second after you edit a price, the note, or the sell time. **Save** writes the file immediately. Quitting does not delete it. Close the window, open the shortcut again, and the week is still there. The first launch of this version copies an old `C:\Tools\Turnips\turnips.json` into AppData, if that file is still there, and renames it `turnips.json.migrated`.

## When you sold

Set **When did you sell?** Each day has two windows: 8:00 AM – 11:59 AM, then 12:00 PM – 10:00 PM. That is when Nook's Cranny is open. The price changes at noon. The banner then reports:

- **Good job** if that was the best price of the week, or within 10 bells of it, and nothing later can beat it.
- **Too soon** if a later price was higher.
- **Too late** if an earlier price was the high, or the turnips spoiled.

## How to read the result

The banner is the main result.

- **Sell now** means every remaining price this week is worse than one you already saw.
- **Wait** means a later price can still beat the current offer.
- The green band on the chart is the range still possible. Gold dots are prices you entered. The faint dashed line is your buy price. The gold dashed vertical line is now. Everything to the right has not happened yet.
- **Guaranteed minimum** and **Potential maximum** mean the same thing as on the website. The minimum is the best price still possible in what remains. The maximum is the highest price still possible.
- **Every pattern line** is the full website output: one row for each way the week can still play out.

Past days you skipped stay as ranges. Those prices are gone. Only slots after your last entered price are still ahead.

## Notes on the patterns

The game picks one pattern per week, per island. The buy price is the base. Sell prices are half-days, Monday morning through Saturday afternoon.

| Pattern | What happens | About how high |
| --- | --- | --- |
| Fluctuating | Up and down all week, no real spike | about 1.4× the buy price |
| Large spike | Drifts down, one large peak, then low prices | up to about 6× |
| Decreasing | Falls all week and does not recover | stays under the buy price |
| Small spike | Drifts down, a modest peak, then down again | about 2× |

A first purchase is always a small spike. **First time buying** forces that, the same way the website does.

Last week's pattern changes the odds for this week:

| Last week | Fluctuating | Large spike | Decreasing | Small spike |
| --- | --- | --- | --- | --- |
| Fluctuating | 20% | 30% | 15% | 35% |
| Large spike | 50% | 5% | 20% | 25% |
| Decreasing | 25% | 45% | 5% | 25% |
| Small spike | 45% | 25% | 15% | 15% |
| Unknown | about 35% | about 25% | about 15% | about 26% |

If you do not know last week, leave the dropdown on Unknown. That row is the same mix the website uses.

Turnips spoil when Sunday starts. Saturday afternoon is the last chance to sell.

## Worked example (week of Oct 18, 2026)

A new install opens a blank week. It does not fill in anyone else's prices. This example matches a link from the website, so you can compare the two. Paste the link into the import box.

| | Mon | Tue | Wed | Thu | Fri | Sat |
| --- | --- | --- | --- | --- | --- | --- |
| AM | | 77 | | | 172 | 149 |
| PM | | | 64 | 132 | 191 | |

`https://turnipprophet.io/?prices=96...77...64..132.172.191.149.&first=true`

With **First time buying** on, the result is a small spike at 100%. Friday, 12:00 PM – 10:00 PM at 191 was the high, just under 2× the 96 buy price. Saturday morning is already down to 149. Saturday afternoon can only land between 39 and 97.

Sell on Saturday morning if you are still holding. Waiting for the afternoon loses bells. The website's guaranteed minimum for what is left is 39, and the potential maximum is 97. If that sale is recorded as Friday afternoon, the banner says it was a good sale.

You can paste that link into the import box on any week. The app fills the grid.

## Run it

Open **Turnip Prophet** from the desktop or the Start menu. Opening it again does not start a second copy. It brings the window back. The PowerShell copy can also be started with `C:\Tools\Turnips\Run.cmd`.

## Files

| File | What it is |
| --- | --- |
| `TurnipProphet-Setup.exe` | Windows installer. Installs `TurnipProphet.exe` to Program Files. |
| `Install.ps1` | PowerShell install. Copies the script app to `C:\Tools\Turnips`. |
| `Turnips.ps1` | The PowerShell app. |
| `Launch.vbs` / `Run.cmd` | Start the PowerShell app without a console window. |
| `TurnipProphet.ico` | The icon for the shortcut and the installer. |
| `www\` | The page: layout, pictures, and font. |
| `predictions.js` | The Turnip Prophet prediction engine, unmodified. |
| `turnips.json` | Your weeks. Stored in `%AppData%\TurnipProphet\`, not beside the program and not in git. |

## Updating

Quit first. Closing the window, or using Quit, writes the weeks file again.

- Installer: run `TurnipProphet-Setup.exe` again over the same folder.
- PowerShell: `git pull`, then `powershell -ExecutionPolicy Bypass -File .\Install.ps1`.

Neither update changes the weeks file in AppData.

## Uninstall

Quit first. Weeks stay in `%AppData%\TurnipProphet\turnips.json`. Delete that folder only if you want the history gone.

- Installer: Windows Settings, Apps, Turnip Prophet, Uninstall.
- PowerShell: delete the desktop shortcut, then delete `C:\Tools\Turnips`.

## License

Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).

The prediction engine is Mike Bryant's, thanks to Ninji. This project is not affiliated with Nintendo.
