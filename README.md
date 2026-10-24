# Animal Crossing: New Horizons - Turnip Prophet

A small Windows tracker for Daisy Mae's turnip prices. Type each price as the week goes, it saves on this PC, and it tells you whether to sell or wait.

The price math is the same engine as [turnipprophet.io](https://turnipprophet.io) (Mike Bryant's [ac-nh-turnip-prices](https://github.com/mikebryant/ac-nh-turnip-prices), from Ninji's reverse engineering of the game). This repo adds a window, a desktop shortcut, and a saved history. It is not the website, and it is not made by Nintendo.

## Install

See [INSTALL.md](INSTALL.md). Short version, from this folder:

```powershell
powershell -ExecutionPolicy Bypass -File .\Install.ps1
```

That copies the app to `C:\Tools\Turnips` and puts **Turnip Prophet** on your desktop. The shortcut, the tray icon, and the window icon are the teal leaf.

## What you see

The left side is the week list, on a field that fades into sky. The grass is the bottom quarter. The right side is where you type the prices, on wood. Across the top of the list is the Animal Crossing: New Horizons logo, then the title on two lines: **Animal Crossing: New Horizons** and **Turnip Prophet**.

Week starting Sunday stays on the left. Last week's pattern sits in the middle. When did you sell sits on the right, lined up with the cards.

Under the prices is the prediction: the banner, the odds, the chart, the twelve sell slots, and the full pattern table if you open it.

## What you do each week

1. Sunday morning, type what Daisy Mae charges. It is always 90–110.
2. When you check Nook's Cranny, type that price in the morning or afternoon box. Morning is 8:00 AM – 11:59 AM. Afternoon is 12:00 PM – 10:00 PM. Do not wipe the morning number.
3. Read the banner. It says sell or wait, in plain words.
4. Next Sunday, click **New week**. If this week was a sure pattern, the new week starts with that filled in as "last week's pattern."

Weeks are stored in `%AppData%\TurnipProphet\turnips.json` (usually `C:\Users\<you>\AppData\Roaming\TurnipProphet\turnips.json`). That is your Windows profile, not the program folder, so deleting `C:\Tools\Turnips` does not delete the weeks. Quit from the button in the window, or right-click the tray icon. Closing the window alone leaves the little helper running so the next double-click just reopens it.

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

This is the week filled in the first time you open the app. Buy 96, first time buying, and these sell prices:

| | Mon | Tue | Wed | Thu | Fri | Sat |
| --- | --- | --- | --- | --- | --- | --- |
| AM | | 77 | | | 172 | 149 |
| PM | | | 64 | 132 | 191 | |

That link is:

`https://turnipprophet.io/?prices=96...77...64..132.172.191.149.&first=true`

With "First time buying" on, the engine says **small spike, 100%**. Friday, 12:00 PM – 10:00 PM at 191 was the peak (just under 2× 96). Saturday, 8:00 AM – 11:59 AM is already down at 149. Saturday, 12:00 PM – 10:00 PM only comes out between **39 and 97**.

Sell during Saturday morning if you still have any. Waiting for the afternoon loses bells. The website's guaranteed minimum for what is left is 39, and the potential maximum is 97. The seed week records that sale as Friday afternoon, so the banner says it was a good job.

You can paste that link into the import box on any week and the app fills the grid.

## Run it from PowerShell

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -STA -File C:\Tools\Turnips\Turnips.ps1
```

Or double-click `Run.cmd` in that folder. A second launch does not start a second copy. It reopens the window.

## Files

| File | What it is |
| --- | --- |
| `Turnips.ps1` | The Windows app. Local page plus the tray icon. |
| `TurnipProphet.ico` | The leaf icon for the desktop shortcut, the tray, and the window. |
| `Install.ps1` | Copies everything to `C:\Tools\Turnips` and makes the shortcut. |
| `Launch.vbs` / `Run.cmd` | Starts the app with no PowerShell window. |
| `www\` | The window you actually look at. `logo.png` is the wide title art. `field.jpg` is the left side. `wood.png` is the right side. |
| `predictions.js` | Unmodified Turnip Prophet engine. |
| `turnips.json` | Your weeks. Lives in `%AppData%\TurnipProphet\`, not in the install folder and not in git. |

## Updating

Quit the app (tray icon or the Quit button), then from a clone of this repo:

```powershell
git pull
powershell -ExecutionPolicy Bypass -File .\Install.ps1
```

Install copies the program files again. It does not touch the weeks file in AppData, or the browser profile folder.

## Uninstall

1. Quit the app.
2. Delete the desktop shortcut.
3. Delete `C:\Tools\Turnips`. The weeks stay in `%AppData%\TurnipProphet\turnips.json`. Delete that folder only if you want the history gone.

## License

Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).

The prediction engine is Mike Bryant's, thanks to Ninji. This project is not affiliated with Nintendo.
