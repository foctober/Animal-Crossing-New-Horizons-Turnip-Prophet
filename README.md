# Animal Crossing: New Horizons - Turnip Prophet

**Version 1.2.** No system tray icon. Closing the window saves, then quits. The installer always adds a Start menu shortcut. A desktop icon is a checkbox, on unless you turn it off. Pinning the window to the taskbar still uses the Turnip Prophet icon. On open, the cursor is already in the sky-blue price box. The title bar and the scrollbars follow the Windows app color, dark or light. Text boxes, the chart, and the rest of the page do not. Zoom stays at 100%. The wood fills the window. Save has a black border. The note is 19px and stays three lines. The First time buying check is the same sky blue as the open price box. Sunday dates in the calendar get a green border and a check on hover. Other days get a red border and an X. Pinning the window to the taskbar uses the Turnip Prophet icon. Unpin the old one first. Closing the window leaves the pin, and opening that pin uses the same icon.

A small Windows tracker for Daisy Mae's turnip prices. Type each price as the week goes, it saves on this PC, and it tells you whether to sell or wait.

The price math is the same engine as [turnipprophet.io](https://turnipprophet.io) (Mike Bryant's [ac-nh-turnip-prices](https://github.com/mikebryant/ac-nh-turnip-prices), from Ninji's reverse engineering of the game). This repo adds a window, a desktop shortcut, and a saved history. It is not the website, and it is not made by Nintendo.

## What it looks like

The left side is the week list, on close-up grass with a white turnip in the middle. The right side is where you type the prices, on orange wood grain tiled small so it repeats.

![Entering the week's prices. Saturday afternoon is the open shop window, so that box has the turning rainbow border.](docs/prices.png)

![The sell-or-wait banner, the chart, and the credit bar](docs/chart.png)

![Every pattern line, the same numbers as the website](docs/patterns.png)

## Install

Pick one. Do not run both at once. Weeks go to `%AppData%\TurnipProphet\turnips.json` either way. Quit Turnip Prophet first if it is already open. More detail is in [INSTALL.md](INSTALL.md).

**Windows installer.** Program Files, Start menu, uninstall from Settings. Paste this whole block into PowerShell:

```powershell
$d = Join-Path $env:TEMP 'TurnipProphet'
if (Test-Path $d) { Remove-Item $d -Recurse -Force }
New-Item -ItemType Directory -Path $d | Out-Null
Invoke-WebRequest 'https://github.com/foctober/Animal-Crossing-New-Horizons-Turnip-Prophet/archive/refs/heads/main.zip' -OutFile "$d\src.zip"
Expand-Archive "$d\src.zip" $d -Force
Start-Process -FilePath "$d\Animal-Crossing-New-Horizons-Turnip-Prophet-main\TurnipProphet-Setup.exe" -Verb RunAs -Wait
Remove-Item $d -Recurse -Force
```

Windows asks for permission. SmartScreen may say the installer is unrecognized. More info, then Run anyway, if you trust this repo. The app goes in `C:\Program Files\Turnip Prophet`.

**PowerShell.** No installer. The app goes in `C:\Tools\Turnips`. Paste this whole block into PowerShell:

```powershell
$d = Join-Path $env:TEMP 'TurnipProphet'
if (Test-Path $d) { Remove-Item $d -Recurse -Force }
New-Item -ItemType Directory -Path $d | Out-Null
Invoke-WebRequest 'https://github.com/foctober/Animal-Crossing-New-Horizons-Turnip-Prophet/archive/refs/heads/main.zip' -OutFile "$d\src.zip"
Expand-Archive "$d\src.zip" $d -Force
powershell -ExecutionPolicy Bypass -File "$d\Animal-Crossing-New-Horizons-Turnip-Prophet-main\Install.ps1"
Remove-Item $d -Recurse -Force
```

The shortcut and the window icon are the teal leaf. There is no tray icon.

## What you see

The left side is the week list, on a close-up field that fades into sky. The grass is the bottom quarter, with the white turnip in the middle. The top of it sits just over the grass and the leaves reach into the sky. The right side is where you type the prices, on orange wood that matches the title. The grain is one seamless tile, repeated at a small size so it does not show up as one big picture. Across the top of the list is the Animal Crossing: New Horizons logo, then the title on two lines in the Animal Crossing font, both the same size: **Animal Crossing: New Horizons** and **Turnip Prophet**, in the same yellow as the logo with the same thick black outline as the yellow text on the price side. Those two lines sit as far below the logo as they sit above the week cards. New Week uses that font too, the same size as Quit, in plain black. New Week, the turnipprophet.io link, and Quit are the same height as the week cards. The week cards use that font too. The date is plain white. Buy is the same sky blue as the open price box, and High is that same yellow. Both prices still pulse. After you pick when you sold, the card adds / Sold and that price. Sold is the New Week green when it was the high, and the Delete Week red when it was lower. It pulses with the other two. Clear when you sold and Sold comes off the card.

The price side uses the Animal Crossing font, FinkHeavy. Yellow text and the green Saved note have a black outline. Shop hours, Daisy Mae, 90, 110, and the saved-file line are white with a black outline. That saved-file line grows to fill the space beside the link and stays on one line. Black wood text has a white outline. Delete Week is white text with a black outline on the red button. A thin red bar under the footer credits Mike Bryant and Ninji and says this project is not affiliated with Nintendo. That bar, the shop-hours lines, the Daisy Mae line, the save note, the banner detail, the pattern odds, the chart key, and the pattern-table note are one size. The top three card titles, First time buying, Sunday Buy Price, the weekdays, the banner headline, and Note are the larger matching size. Cards and input boxes have a thin black border. The price box for the open shop window is a bright sky blue, the same blue as the sky on the left, with a thick rainbow border that turns and a white glow on the same 1.5 second pulse as Buy and High. Text inside the white cards, inputs, Sunday button, and calendar stays plain black.

The top three fields are three equal cards. Inside each card the label is left-aligned with its box. The date and both dropdowns stay left aligned. Their text, the price boxes, and the twelve cards under the chart sit in the vertical middle. The twelve cards are 17px and one and a half times the height of the price boxes. The notes are 18px, and the buttons and inputs are 19px. Price boxes are centered left to right as well. The date, the two dropdowns, every price box, the paste box, and Import, Delete Week, and Save are the same height as the week cards. Those labels are the same size as New Week, Quit, and the turnipprophet.io link. Last week's pattern stays grey while First time buying is checked. Uncheck that and the pattern box is bold like the other two. The price box for the shop window that is open right now is bright sky blue, with a thick rainbow border that turns and a white glow on the same 1.5 second pulse as Buy and High. Every other price box keeps the thin black border. On a week that has not started yet, Sunday Buy Price is the one marked, because that is the first price to enter. On this week, Sunday morning is the buy box, and Monday through Saturday follow 8:00 AM and noon. A finished week is not marked. Outside shop hours nothing is marked.

Under the prices is the prediction: the banner, the odds, the chart, the twelve sell slots, and the full pattern table if you open it. The table text is 14px, the same as the prices on the chart.

## What you do each week

1. Sunday morning, type what Daisy Mae charges. It is always 90–110.
2. When you check Nook's Cranny, type that price in the morning or afternoon box. Morning is 8:00 AM – 11:59 AM. Afternoon is 12:00 PM – 10:00 PM. Do not wipe the morning number.
3. Read the banner. It says sell or wait, in plain words.
4. Next Sunday, click **New Week**. If this week was a sure pattern, the new week starts with that filled in as "last week's pattern."

Weeks are stored in `%AppData%\TurnipProphet\turnips.json` (usually `C:\Users\<you>\AppData\Roaming\TurnipProphet\turnips.json`). Neither install keeps the weeks in the program folder. Closing the window saves, then quits. The Quit button does the same.

## Saving

It saves by itself, about a second after you change a price, the note, or the sell time. The **Save** button writes `%AppData%\TurnipProphet\turnips.json` immediately. Quit does not wipe that file. Close the window, open the shortcut again, and the week is still there. The first time this version runs, it copies an old `C:\Tools\Turnips\turnips.json` into AppData and renames the old file `turnips.json.migrated`.

## When you sold

Pick **When did you sell?** Every day has the same two windows: 8:00 AM – 11:59 AM, then 12:00 PM – 10:00 PM. That is when Nook's Cranny is open, and the price changes at noon. The banner then says:

- **Good job** if that was the best price of the week, or within 10 bells of it, and nothing later can beat it.
- **Too soon** if a later price was higher.
- **Too late** if an earlier price was the high, or the turnips spoiled.

## How to read the result

The banner is the part that matters.

- **Sell now** means every price still left this week is worse than a price you already saw. Sell whatever you are still holding.
- **Wait** means a later slot can still beat the price Nook is offering right now.
- The green band on the chart is the range the math still allows. Gold dots are numbers you typed. The faint dashed line is your buy price. The gold dashed vertical line is "now" — everything to the right has not happened yet.
- **Guaranteed minimum** and **Potential maximum** use the same meaning as the website, so you can compare the two. The minimum is the best floor still hiding in the remaining pattern. The maximum is the highest price the math still allows.
- **Every pattern line** (the collapsed table) is the full website output: one row per way the week can still happen.

Blank days you already passed stay as ranges. You cannot go back and sell those. Only slots after your last typed price are still ahead.

## Notes on the patterns

The game picks one pattern per week, per island. Your buy price is the base. Sell prices are half-days, Monday morning through Saturday afternoon.

| Pattern | What it feels like | About how high |
| --- | --- | --- |
| Fluctuating | Up and down all week, no real spike | about 1.4× the buy price |
| Large spike | Drifts down, then one huge peak, then junk prices | up to about 6× |
| Decreasing | Down all week and it does not come back | stays under the buy price |
| Small spike | Drifts down, a modest peak, then down again | about 2× |

A first-time buyer always gets a small spike. The "First time buying" box forces that, same as the website.

What happened last week changes the odds for this week:

| Last week | Fluctuating | Large spike | Decreasing | Small spike |
| --- | --- | --- | --- | --- |
| Fluctuating | 20% | 30% | 15% | 35% |
| Large spike | 50% | 5% | 20% | 25% |
| Decreasing | 25% | 45% | 5% | 25% |
| Small spike | 45% | 25% | 15% | 15% |
| Unknown | about 35% | about 25% | about 15% | about 26% |

If you do not know last week, leave the dropdown on Unknown. The unknown row is the steady-state mix the website uses.

Turnips spoil when Sunday starts. Saturday afternoon is the last chance to sell.

## Worked example (week of Oct 18, 2026)

A new install opens a blank week for this Sunday. It does not fill in anybody else's prices. This example is only here so you can try the same numbers as the website. Paste the link into the import box.

| | Mon | Tue | Wed | Thu | Fri | Sat |
| --- | --- | --- | --- | --- | --- | --- |
| AM | | 77 | | | 172 | 149 |
| PM | | | 64 | 132 | 191 | |

That link is:

`https://turnipprophet.io/?prices=96...77...64..132.172.191.149.&first=true`

With "First time buying" on, the engine says **small spike, 100%**. Friday, 12:00 PM – 10:00 PM at 191 was the peak (just under 2× 96). Saturday, 8:00 AM – 11:59 AM is already down at 149. Saturday, 12:00 PM – 10:00 PM only comes out between **39 and 97**.

Sell during Saturday morning if you still have any. Waiting for the afternoon loses bells. The website's guaranteed minimum for what is left is 39, and the potential maximum is 97. If that sale is recorded as Friday afternoon, the banner says it was a good job.

You can paste that link into the import box on any week and the app fills the grid.

## Run it

Open **Turnip Prophet** from the desktop or the Start menu. A second launch does not start a second copy. It reopens the window. The PowerShell copy can also be started with `C:\Tools\Turnips\Run.cmd`.

## Files

| File | What it is |
| --- | --- |
| `TurnipProphet-Setup.exe` | Windows installer. Puts `TurnipProphet.exe` in Program Files. |
| `Install.ps1` | PowerShell install. Copies the script app to `C:\Tools\Turnips`. |
| `Turnips.ps1` | The PowerShell app. Local page, no tray icon. |
| `Launch.vbs` / `Run.cmd` | Start the PowerShell app with no PowerShell window. |
| `TurnipProphet.ico` | The leaf icon for the shortcut and the installer. |
| `www\` | The window you actually look at. `logo.png` is the wide title art. `field.jpg` is the left side. `wood.png` is the right side, tiled at 360 by 180. `FinkHeavy.ttf` is the font on both sides. |
| `predictions.js` | Unmodified Turnip Prophet engine. |
| `turnips.json` | Your weeks. Lives in `%AppData%\TurnipProphet\`, not next to the program and not in git. |

## Updating

Quit the app first (close the window, or the Quit button). That writes the weeks file again.

- Installer: run `TurnipProphet-Setup.exe` again over the same folder.
- PowerShell: `git pull`, then `powershell -ExecutionPolicy Bypass -File .\Install.ps1`.

Neither one touches the weeks file in AppData.

## Uninstall

Quit first. The weeks stay in `%AppData%\TurnipProphet\turnips.json`. Delete that folder only if you want the history gone.

- Installer: Windows Settings, Apps, Turnip Prophet, Uninstall.
- PowerShell: delete the desktop shortcut and delete `C:\Tools\Turnips`.

## License

Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).

The prediction engine is Mike Bryant's, thanks to Ninji. This project is not affiliated with Nintendo.
