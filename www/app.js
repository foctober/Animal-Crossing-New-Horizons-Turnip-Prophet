"use strict";

var PATTERN_NAME = ["Fluctuating", "Large spike", "Decreasing", "Small spike", "All patterns"];
var SHORT = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
var SLOT = [
  ["Monday", "morning"],
  ["Monday", "afternoon"],
  ["Tuesday", "morning"],
  ["Tuesday", "afternoon"],
  ["Wednesday", "morning"],
  ["Wednesday", "afternoon"],
  ["Thursday", "morning"],
  ["Thursday", "afternoon"],
  ["Friday", "morning"],
  ["Friday", "afternoon"],
  ["Saturday", "morning"],
  ["Saturday", "afternoon"]
];

var SEED = {
  weeks: [{
    id: "2026-10-18",
    firstBuy: true,
    previousPattern: -1,
    likelyPattern: -1,
    prices: [96, null, null, 77, null, null, 64, null, 132, 172, 191, 149, null],
    soldAt: 10,
    note: "",
    updatedAt: "2026-10-24T12:00:00.000Z"
  }],
  activeId: "2026-10-18"
};

var state = { weeks: [], activeId: null };
var gen = 0;
var saveTimer = null;
var predictTimer = null;
var lastChart = null;

function $(id) { return document.getElementById(id); }

function fmt(d) {
  var y = d.getFullYear();
  var m = String(d.getMonth() + 1).padStart(2, "0");
  var day = String(d.getDate()).padStart(2, "0");
  return y + "-" + m + "-" + day;
}

function pretty(id) {
  var p = id.split("-").map(Number);
  var dt = new Date(p[0], p[1] - 1, p[2]);
  return dt.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" });
}

function thisSunday() {
  var n = new Date();
  n.setHours(12, 0, 0, 0);
  n.setDate(n.getDate() - n.getDay());
  return fmt(n);
}

function snapSunday(value) {
  if (!value) return null;
  var p = value.split("-").map(Number);
  var dt = new Date(p[0], p[1] - 1, p[2]);
  dt.setDate(dt.getDate() - dt.getDay());
  return fmt(dt);
}

function slotName(idx14) {
  var pair = SLOT[idx14 - 2];
  var hours = pair[1] === "morning" ? "8:00 AM – 11:59 AM" : "12:00 PM – 10:00 PM";
  return pair[0] + ", " + hours;
}

function active() {
  for (var i = 0; i < state.weeks.length; i++) {
    if (state.weeks[i].id === state.activeId) return state.weeks[i];
  }
  return null;
}

function byNewest(a, b) {
  if (a.id < b.id) return 1;
  if (a.id > b.id) return -1;
  return 0;
}

function status(text) { $("status").textContent = text || ""; }

function normPrice(v) {
  if (v === null || v === undefined || v === "") return NaN;
  var n = Number(v);
  return Number.isFinite(n) ? n : NaN;
}

function readNum(input) {
  var t = String(input.value || "").trim();
  if (!t) return null;
  var n = Number(t);
  if (!Number.isInteger(n) || n < 1 || n > 999) return null;
  return n;
}

function prices14(week) {
  var buy = normPrice(week.prices[0]);
  var out = [buy, buy];
  for (var i = 1; i <= 12; i++) out.push(normPrice(week.prices[i]));
  return out;
}

function displayPercentage(fraction) {
  if (!Number.isFinite(fraction)) return "\u2014";
  var percent = fraction * 100;
  if (percent >= 1) return percent.toPrecision(3) + "%";
  if (percent >= 0.01) return percent.toFixed(2) + "%";
  return "<0.01%";
}

function priceClass(buy, max) {
  if (!Number.isFinite(max)) return "";
  var base = Number.isFinite(buy) ? buy : 100;
  var diff = max - base;
  if (diff >= 200) return "hot";
  if (diff >= 30) return "good";
  if (diff >= 0) return "ok";
  if (diff >= -30) return "meh";
  return "bad";
}

function permalink(week) {
  var sells = week.prices.slice(1).map(function (v) { return v == null ? "" : String(v); }).join(".");
  var buy = week.prices[0] == null ? "" : String(week.prices[0]);
  var url = "https://turnipprophet.io/?prices=" + buy + "." + sells;
  if (week.firstBuy) url += "&first=true";
  if (week.previousPattern >= 0) url += "&pattern=" + week.previousPattern;
  return url;
}

function readFormInto(week) {
  if (!week) return;
  week.firstBuy = $("firstBuy").checked;
  week.previousPattern = Number($("prevPattern").value);
  week.note = $("note").value;
  var sold = $("soldAt").value;
  if (sold === "" || sold === "hold") week.soldAt = null;
  else if (sold === "spoil") week.soldAt = "spoil";
  else week.soldAt = Number(sold);
  var prices = [readNum($("buy"))];
  for (var i = 0; i < 12; i++) prices.push(readNum($("p" + i)));
  week.prices = prices;
  week.updatedAt = new Date().toISOString();
}

function fillForm(week) {
  $("sunday").value = week.id;
  syncSundayFace();
  $("firstBuy").checked = !!week.firstBuy;
  $("prevPattern").value = String(week.previousPattern == null ? -1 : week.previousPattern);
  $("note").value = week.note || "";
  $("soldAt").value = week.soldAt == null || week.soldAt === "hold" ? "" : String(week.soldAt);
  $("buy").value = week.prices[0] == null ? "" : week.prices[0];
  for (var i = 0; i < 12; i++) {
    var v = week.prices[i + 1];
    $("p" + i).value = v == null ? "" : v;
  }
  syncPrevLock();
  $("siteLink").href = permalink(week);
  document.title = "Turnip Prophet \u00b7 " + pretty(week.id);
  markDue();
}

function dueSlotId() {
  var week = active();
  if (!week || week.id !== thisSunday()) return null;
  var now = new Date();
  var day = now.getDay();
  var mins = now.getHours() * 60 + now.getMinutes();
  if (day === 0) return mins < 12 * 60 ? "buy" : null;
  if (mins < 8 * 60 || mins >= 22 * 60) return null;
  var half = mins < 12 * 60 ? 0 : 1;
  return "p" + ((day - 1) * 2 + half);
}

function markDue() {
  var ids = ["buy"];
  for (var i = 0; i < 12; i++) ids.push("p" + i);
  var due = dueSlotId();
  ids.forEach(function (id) {
    var el = $(id);
    if (!el) return;
    if (id === due) el.classList.add("due");
    else el.classList.remove("due");
  });
}

function syncPrevLock() {
  var on = $("firstBuy").checked;
  $("prevPattern").disabled = on;
  $("prevHint").textContent = on
    ? "A first purchase is always a small spike. Last week does not change it."
    : "Last week's pattern changes the odds. Leave it on Unknown if you are not sure.";
}

function setEnabled(on) {
  ["sunday", "sundayBtn", "firstBuy", "prevPattern", "soldAt", "buy", "note", "importBtn", "deleteBtn", "importUrl", "saveBtn"].forEach(function (id) {
    $(id).disabled = !on;
  });
  for (var i = 0; i < 12; i++) $("p" + i).disabled = !on;
  if (on) syncPrevLock();
}

function renderWeekList() {
  var box = $("weekList");
  box.innerHTML = "";
  var weeks = state.weeks.slice().sort(byNewest);
  if (!weeks.length) {
    var empty = document.createElement("p");
    empty.className = "hint";
    empty.textContent = "No weeks yet.";
    box.appendChild(empty);
    return;
  }
  weeks.forEach(function (week) {
    var btn = document.createElement("button");
    btn.type = "button";
    btn.className = "week" + (week.id === state.activeId ? " on" : "");
    var date = document.createElement("span");
    date.className = "date";
    date.textContent = pretty(week.id);
    var marks = document.createElement("span");
    marks.className = "marks";
    var buy = week.prices[0];
    var best = null;
    week.prices.slice(1).forEach(function (v) {
      if (v != null && (best == null || v > best)) best = v;
    });
    function tag(text) {
      var el = document.createElement("span");
      el.className = "tag";
      el.textContent = text;
      return el;
    }
    if (buy != null) {
      marks.appendChild(tag("Buy"));
      var buyEl = document.createElement("b");
      buyEl.className = "buy";
      buyEl.textContent = "(" + buy + ")";
      marks.appendChild(buyEl);
    }
    if (buy != null && best != null) marks.appendChild(tag(":"));
    if (best != null) {
      marks.appendChild(tag("High"));
      var bestEl = document.createElement("b");
      bestEl.className = "best";
      bestEl.textContent = "(" + best + ")";
      marks.appendChild(bestEl);
    }
    btn.appendChild(date);
    btn.appendChild(marks);
    btn.addEventListener("click", function () { selectWeek(week.id); });
    box.appendChild(btn);
  });
}

function selectWeek(id) {
  if (id === state.activeId) return;
  readFormInto(active());
  state.activeId = id;
  fillForm(active());
  renderWeekList();
  saveSoon();
  schedulePredict();
}

function showNoWeek() {
  $("result").hidden = true;
  setEnabled(false);
  status("Add a week to start.");
}

function advice(rows, entered, firstBuy) {
  var all = rows[0];
  var odds = [];
  var seen = {};
  rows.forEach(function (row) {
    if (row.pattern_number === 4 || seen[row.pattern_number]) return;
    seen[row.pattern_number] = true;
    odds.push({
      n: row.pattern_number,
      name: PATTERN_NAME[row.pattern_number],
      p: row.category_total_probability
    });
  });
  odds.sort(function (a, b) { return b.p - a.p; });

  var last = -1;
  var best = null;
  for (var i = 2; i < 14; i++) {
    if (!Number.isFinite(entered[i])) continue;
    last = i;
    if (!best || entered[i] > best.price) best = { price: entered[i], idx: i };
  }

  if (last < 0) {
    return {
      tone: "wait",
      headline: "Need a sell price.",
      detail: "The buy price alone is a wide range. Check Nook's once on Monday morning.",
      odds: odds
    };
  }

  var future = [];
  for (var f = last + 1; f < 14; f++) {
    future.push({
      idx: f,
      name: slotName(f),
      min: all.prices[f].min,
      max: all.prices[f].max
    });
  }

  var current = entered[last];
  var currentName = slotName(last);
  var patternLine = odds.length
    ? odds[0].name + ", " + displayPercentage(odds[0].p) + "."
    : "";
  if (firstBuy) patternLine = "First purchase is always a small spike. " + patternLine;

  if (!future.length) {
    return {
      tone: "wait",
      headline: "Week is over. Best price was " + best.price + " on " + slotName(best.idx) + ".",
      detail: "Turnips spoil on Sunday morning. " + patternLine,
      odds: odds
    };
  }

  var next = future[0];
  var bestFuture = future[0];
  future.forEach(function (slot) {
    if (slot.max > bestFuture.max) bestFuture = slot;
  });

  if (bestFuture.max < best.price && next.max <= current) {
    var head;
    if (best.idx === last) {
      head = "Sell now. This is the high (" + current + " on " + currentName + ").";
    } else {
      head = "Sell now. The peak was " + best.price + " on " + slotName(best.idx) + ".";
    }
    return {
      tone: "sell",
      headline: head,
      detail: currentName.charAt(0).toUpperCase() + currentName.slice(1) + " is " + current +
        ". " + next.name.charAt(0).toUpperCase() + next.name.slice(1) +
        " only reaches " + next.min + "\u2013" + next.max + ". Waiting loses bells. " + patternLine,
      odds: odds
    };
  }

  if (next.max > current) {
    var later = bestFuture.idx !== next.idx
      ? " Later, " + bestFuture.name + " can still reach " + bestFuture.max + "."
      : "";
    return {
      tone: "wait",
      headline: "Wait. " + next.name.charAt(0).toUpperCase() + next.name.slice(1) + " can still reach " + next.max + ".",
      detail: "You are at " + current + " on " + currentName + "." + later + " " + patternLine,
      odds: odds
    };
  }

  return {
    tone: "wait",
    headline: "Hold for " + bestFuture.name + ". It can still reach " + bestFuture.max + ".",
    detail: "The next price (" + next.name + ") is only " + next.min + "\u2013" + next.max +
      ", so it may dip before that. You are at " + current + ". " + patternLine,
    odds: odds
  };
}

function renderBanner(model) {
  var box = $("banner");
  box.className = "banner " + (model.tone || "");
  box.innerHTML = "";
  var h = document.createElement("h2");
  h.textContent = model.headline;
  var p = document.createElement("p");
  p.textContent = model.detail;
  box.appendChild(h);
  box.appendChild(p);
}

function renderOdds(odds) {
  var box = $("odds");
  box.innerHTML = "";
  odds.forEach(function (o) {
    var row = document.createElement("div");
    row.className = "odd";
    var name = document.createElement("span");
    name.textContent = o.name;
    var track = document.createElement("div");
    track.className = "track";
    var fill = document.createElement("div");
    fill.className = "fill";
    var pct = Number.isFinite(o.p) ? Math.max(0, Math.min(1, o.p)) * 100 : 0;
    fill.style.width = pct + "%";
    track.appendChild(fill);
    var label = document.createElement("b");
    label.textContent = displayPercentage(o.p);
    row.appendChild(name);
    row.appendChild(track);
    row.appendChild(label);
    box.appendChild(row);
  });
}

function renderSlots(week, all) {
  var box = $("slots");
  box.innerHTML = "";
  var last = -1;
  for (var i = 0; i < 12; i++) if (week.prices[i + 1] != null) last = i;
  for (var s = 0; s < 12; s++) {
    var entered = week.prices[s + 1];
    var range = all.prices[s + 2];
    var card = document.createElement("div");
    var cls = "slot";
    if (entered != null) cls += " known";
    else if (s > last) cls += " future";
    else cls += " pastblank";
    if (Number(week.soldAt) === s + 1) cls += " sold";
    card.className = cls;
    var when = document.createElement("div");
    when.className = "when";
    when.textContent = SHORT[Math.floor(s / 2)] + " " + (s % 2 === 0 ? "8:00 AM" : "12:00 PM");
    var val = document.createElement("div");
    val.className = "val";
    val.textContent = entered != null ? String(entered) : (range.min + "\u2013" + range.max);
    card.appendChild(when);
    card.appendChild(val);
    box.appendChild(card);
  }
}

function td(text, cls) {
  var el = document.createElement("td");
  el.textContent = text;
  if (cls) el.className = cls;
  return el;
}

function renderTable(rows, buy) {
  var head = $("thead");
  var body = $("tbody");
  head.innerHTML = "";
  body.innerHTML = "";
  var hr = document.createElement("tr");
  var labels = ["Pattern", "Odds", "Path"];
  for (var i = 0; i < 6; i++) {
    labels.push(SHORT[i] + " AM");
    labels.push(SHORT[i] + " PM");
  }
  labels.push("Guaranteed min");
  labels.push("Potential max");
  labels.forEach(function (label, idx) {
    var th = document.createElement("th");
    th.textContent = label;
    if (idx === 0) th.className = "pat";
    hr.appendChild(th);
  });
  head.appendChild(hr);

  var counts = {};
  rows.forEach(function (row) {
    counts[row.pattern_number] = (counts[row.pattern_number] || 0) + 1;
  });
  var seen = {};
  rows.forEach(function (row) {
    var tr = document.createElement("tr");
    var pat = td(PATTERN_NAME[row.pattern_number] || "?", "pat");
    tr.appendChild(pat);
    if (!seen[row.pattern_number]) {
      seen[row.pattern_number] = true;
      var oddsCell = td(displayPercentage(row.category_total_probability));
      oddsCell.rowSpan = counts[row.pattern_number];
      tr.appendChild(oddsCell);
    }
    tr.appendChild(td(displayPercentage(row.probability)));
    row.prices.slice(2).forEach(function (day) {
      var text = day.min === day.max ? String(day.min) : (day.min + "\u2013" + day.max);
      tr.appendChild(td(text, priceClass(buy, day.max)));
    });
    tr.appendChild(td(Number.isFinite(row.weekGuaranteedMinimum) ? String(row.weekGuaranteedMinimum) : "\u2014", priceClass(buy, row.weekGuaranteedMinimum)));
    tr.appendChild(td(Number.isFinite(row.weekMax) ? String(row.weekMax) : "\u2014", priceClass(buy, row.weekMax)));
    body.appendChild(tr);
  });
}

function drawChart() {
  if (!lastChart) return;
  var canvas = $("chart");
  var dpr = window.devicePixelRatio || 1;
  var w = canvas.clientWidth || 800;
  var h = canvas.clientHeight || 240;
  canvas.width = Math.floor(w * dpr);
  canvas.height = Math.floor(h * dpr);
  var g = canvas.getContext("2d");
  g.setTransform(dpr, 0, 0, dpr, 0, 0);
  g.clearRect(0, 0, w, h);

  var entered = lastChart.entered;
  var all = lastChart.all;
  var buy = lastChart.buy;
  var n = 12;
  var disp = [];
  var last = -1;
  for (var i = 0; i < n; i++) {
    var known = Number.isFinite(entered[i + 2]) ? entered[i + 2] : null;
    if (known != null) last = i;
    var cell = all[i + 2];
    disp.push(known != null ? { min: known, max: known, known: known } : { min: cell.min, max: cell.max, known: null });
  }

  var lo = Infinity;
  var hi = 0;
  disp.forEach(function (v) {
    lo = Math.min(lo, v.min);
    hi = Math.max(hi, v.max);
  });
  if (Number.isFinite(buy)) {
    lo = Math.min(lo, buy);
    hi = Math.max(hi, buy);
  }
  if (!Number.isFinite(lo)) { lo = 0; hi = 100; }
  lo = Math.max(0, Math.floor((lo - 15) / 10) * 10);
  hi = Math.ceil((hi + 15) / 10) * 10;
  if (hi <= lo) hi = lo + 40;

  var padL = 48, padR = 12, padT = 16, padB = 44;
  var iw = w - padL - padR;
  var ih = h - padT - padB;
  function x(i) { return padL + (i + 0.5) * (iw / n); }
  function y(v) { return padT + (1 - (v - lo) / (hi - lo)) * ih; }

  g.font = "14px Segoe UI";
  g.fillStyle = "#4d6844";
  g.strokeStyle = "#d5e8bc";
  g.lineWidth = 1;
  var span = hi - lo;
  var step = span > 220 ? 50 : span > 90 ? 20 : 10;
  for (var bells = Math.ceil(lo / step) * step; bells <= hi; bells += step) {
    var yy = y(bells);
    g.beginPath();
    g.moveTo(padL, yy);
    g.lineTo(w - padR, yy);
    g.stroke();
    g.fillText(String(bells), 4, yy + 4);
  }

  if (Number.isFinite(buy) && buy >= lo && buy <= hi) {
    g.save();
    g.strokeStyle = "#d0892a";
    g.setLineDash([2, 4]);
    g.beginPath();
    g.moveTo(padL, y(buy));
    g.lineTo(w - padR, y(buy));
    g.stroke();
    g.restore();
  }

  g.beginPath();
  disp.forEach(function (v, i) {
    var px = x(i);
    var py = y(v.max);
    if (i === 0) g.moveTo(px, py);
    else g.lineTo(px, py);
  });
  for (var k = n - 1; k >= 0; k--) g.lineTo(x(k), y(disp[k].min));
  g.closePath();
  g.fillStyle = "rgba(90, 190, 70, 0.35)";
  g.fill();
  g.strokeStyle = "#3ea84a";
  g.stroke();

  if (last >= 0 && last < n - 1) {
    var xLine = (x(last) + x(last + 1)) / 2;
    g.save();
    g.strokeStyle = "#e0a106";
    g.setLineDash([3, 3]);
    g.beginPath();
    g.moveTo(xLine, padT);
    g.lineTo(xLine, padT + ih);
    g.stroke();
    g.restore();
  }

  g.strokeStyle = "#2f6b32";
  g.lineWidth = 2;
  g.beginPath();
  var pen = false;
  disp.forEach(function (v, i) {
    if (v.known == null) { pen = false; return; }
    if (!pen) g.moveTo(x(i), y(v.known));
    else g.lineTo(x(i), y(v.known));
    pen = true;
  });
  g.stroke();

  g.fillStyle = "#e09412";
  disp.forEach(function (v, i) {
    if (v.known == null) return;
    g.beginPath();
    g.arc(x(i), y(v.known), 4, 0, Math.PI * 2);
    g.fill();
  });

  g.fillStyle = "#4d6844";
  g.font = "13px Segoe UI";
  g.textAlign = "center";
  disp.forEach(function (_v, i) {
    g.fillText(i % 2 === 0 ? "AM" : "PM", x(i), h - 22);
    if (i % 2 === 0) g.fillText(SHORT[i / 2], x(i) + (iw / n) / 2, h - 6);
  });
  g.textAlign = "left";
}


function paidLine(sell, buy) {
  if (!Number.isFinite(buy) || !buy) return "";
  var times = Math.round((sell / buy) * 100) / 100;
  return " That is " + times.toFixed(2) + "× the " + buy + " you paid Daisy Mae.";
}

function bestBeforeIndex(entered, idx14) {
  var best = null;
  for (var i = 2; i < idx14; i++) {
    if (Number.isFinite(entered[i]) && (!best || entered[i] > best.price)) best = { price: entered[i], idx: i };
  }
  return best;
}

function bestAfterIndex(entered, idx14) {
  var best = null;
  for (var i = idx14 + 1; i < 14; i++) {
    if (Number.isFinite(entered[i]) && (!best || entered[i] > best.price)) best = { price: entered[i], idx: i };
  }
  return best;
}

function gradeSale(week, rows, entered) {
  var soldAt = week.soldAt;
  if (soldAt == null || soldAt === "" || soldAt === "hold") return null;
  var buy = Number.isFinite(entered[0]) ? entered[0] : NaN;

  if (soldAt === "spoil") {
    var seen = null;
    for (var s = 2; s < 14; s++) {
      if (Number.isFinite(entered[s]) && (!seen || entered[s] > seen.price)) seen = { price: entered[s], idx: s };
    }
    if (!seen) {
      return { tone: "late", headline: "Too late. They spoiled.", detail: "No sell prices were typed, so there is nothing to compare." };
    }
    return {
      tone: "late",
      headline: "Too late. They spoiled.",
      detail: "The best price you saw was " + seen.price + " on " + slotName(seen.idx) + "."
    };
  }

  var idx = Number(soldAt);
  var idx14 = idx + 1;
  if (!(idx >= 1 && idx <= 12) || !Number.isFinite(entered[idx14])) {
    var when = (idx >= 1 && idx <= 12) ? slotName(idx14) : "that time";
    return { tone: "wait", headline: "Add the price for " + when + ".", detail: "Then this can say if that sale was early, late, or right on time." };
  }

  var sell = entered[idx14];
  var name = slotName(idx14);
  var before = bestBeforeIndex(entered, idx14);
  var afterKnown = bestAfterIndex(entered, idx14);
  var futureMax = sell;
  var futureName = null;
  var all = rows && rows[0];
  if (all) {
    for (var f = idx14 + 1; f < 14; f++) {
      if (Number.isFinite(entered[f])) continue;
      if (all.prices[f].max > futureMax) {
        futureMax = all.prices[f].max;
        futureName = slotName(f);
      }
    }
  }

  var cap = function (n) { return n.charAt(0).toUpperCase() + n.slice(1); };

  if (afterKnown && afterKnown.price > sell + 10) {
    return {
      tone: "soon",
      headline: "Too soon.",
      detail: "You sold at " + sell + " on " + name + ". " + cap(slotName(afterKnown.idx)) + " got up to " + afterKnown.price + "." + paidLine(sell, buy)
    };
  }

  if (before && before.price > sell + 10 && (!afterKnown || afterKnown.price <= before.price)) {
    return {
      tone: "late",
      headline: "Too late.",
      detail: "The high was " + before.price + " on " + slotName(before.idx) + ". You sold at " + sell + " on " + name + "." + paidLine(sell, buy)
    };
  }

  if (futureName && futureMax > sell + 15) {
    return {
      tone: "soon",
      headline: "Maybe too soon.",
      detail: "You locked in " + sell + " on " + name + ". " + cap(futureName) + " can still reach " + futureMax + "." + paidLine(sell, buy)
    };
  }

  var extra = "";
  if (afterKnown && afterKnown.price < sell) {
    extra = " " + cap(slotName(afterKnown.idx)) + " was already down to " + afterKnown.price + ".";
  }
  var head = (!before || sell >= before.price) && (!afterKnown || sell >= afterKnown.price)
    ? "Good job! That was the best price!"
    : "Good job! Close to the best!";
  return {
    tone: "good",
    headline: head,
    detail: cap(name) + " at " + sell + " was the high this week." + extra + paidLine(sell, buy)
  };
}

function showFitError() {
  $("result").hidden = false;
  renderBanner({
    tone: "sell",
    headline: "These prices do not fit a pattern.",
    detail: "Check for a typo. If you time traveled backwards, do not mix prices from before the jump with prices from after it."
  });
  $("odds").innerHTML = "";
  $("slots").innerHTML = "";
  $("tbody").innerHTML = "";
  lastChart = null;
  status("");
}

function predict() {
  var token = ++gen;
  var week = active();
  if (!week) { showNoWeek(); return; }
  readFormInto(week);
  $("siteLink").href = permalink(week);
  if (!week.prices.some(function (v) { return v != null; })) {
    $("result").hidden = true;
    status("Type a price to see the week.");
    renderWeekList();
    return;
  }
  status("Working out the patterns\u2026");
  var weekId = week.id;
  window.setTimeout(function () {
    if (token !== gen) return;
    if (!active() || active().id !== weekId) return;
    var entered = prices14(week);
    var prev = week.firstBuy || !(week.previousPattern >= 0 && week.previousPattern <= 3)
      ? NaN
      : week.previousPattern;
    var realLog = console.log;
    console.log = function () {};
    var rows;
    try {
      var predictor = new Predictor(entered, !!week.firstBuy, prev);
      rows = predictor.analyze_possibilities();
    } catch (err) {
      console.log = realLog;
      showFitError();
      return;
    }
    console.log = realLog;
    if (!rows || !rows.length || !Number.isFinite(rows[0].weekGuaranteedMinimum) || rows[0].weekGuaranteedMinimum === Infinity) {
      showFitError();
      week.likelyPattern = -1;
      return;
    }
    var forecast = advice(rows, entered, week.firstBuy);
    var grade = gradeSale(week, rows, entered);
    var model = grade || forecast;
    week.likelyPattern = forecast.odds.length && forecast.odds[0].p >= 0.8 ? forecast.odds[0].n : -1;
    $("result").hidden = false;
    renderBanner(model);
    renderOdds(forecast.odds);
    renderSlots(week, rows[0]);
    renderTable(rows, Number.isFinite(entered[0]) ? entered[0] : rows[0].prices[0].min);
    lastChart = { all: rows[0].prices, entered: entered, buy: Number.isFinite(entered[0]) ? entered[0] : NaN };
    drawChart();
    var min = rows[0].weekGuaranteedMinimum;
    var max = rows[0].weekMax;
    $("legend").textContent = "Gold dots are prices you typed. Green is the range still possible. Faint dashes are your buy price. The gold dashed line is now. Guaranteed minimum " + min + ", potential maximum " + max + ".";
    status("");
    renderWeekList();
    saveSoon();
  }, 30);
}

function schedulePredict() {
  window.clearTimeout(predictTimer);
  predictTimer = window.setTimeout(predict, 400);
}

function saveSoon() {
  $("saveState").textContent = "Saving\u2026";
  window.clearTimeout(saveTimer);
  saveTimer = window.setTimeout(saveNow, 200);
}

function saveNow() {
  var week = active();
  if (week) readFormInto(week);
  return fetch("/api/state", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(state)
  }).then(function (res) {
    $("saveState").textContent = res.ok ? "Saved." : "Save failed";
    renderWeekList();
  }).catch(function () {
    $("saveState").textContent = "Save failed";
  });
}

var calCursor = null;

function syncSundayFace() {
  var btn = $("sundayBtn");
  if (!btn) return;
  var id = $("sunday").value;
  btn.textContent = id ? pretty(id) : "Pick a Sunday";
}

function renderCal() {
  var cal = $("cal");
  cal.innerHTML = "";
  var y = calCursor.getFullYear();
  var m = calCursor.getMonth();
  var head = document.createElement("div");
  head.className = "cal-head";
  var prev = document.createElement("button");
  prev.type = "button";
  prev.textContent = "\u2039";
  prev.addEventListener("click", function () {
    calCursor.setMonth(calCursor.getMonth() - 1);
    renderCal();
  });
  var title = document.createElement("span");
  title.textContent = calCursor.toLocaleDateString("en-US", { month: "long", year: "numeric" });
  var next = document.createElement("button");
  next.type = "button";
  next.textContent = "\u203a";
  next.addEventListener("click", function () {
    calCursor.setMonth(calCursor.getMonth() + 1);
    renderCal();
  });
  head.appendChild(prev);
  head.appendChild(title);
  head.appendChild(next);
  cal.appendChild(head);

  var names = document.createElement("div");
  names.className = "cal-names";
  ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"].forEach(function (n, i) {
    var s = document.createElement("span");
    s.textContent = n;
    if (i !== 0) s.className = "off";
    names.appendChild(s);
  });
  cal.appendChild(names);

  var grid = document.createElement("div");
  grid.className = "cal-grid";
  var start = new Date(y, m, 1).getDay();
  var days = new Date(y, m + 1, 0).getDate();
  var selected = $("sunday").value;
  for (var i = 0; i < start; i++) grid.appendChild(document.createElement("span"));
  for (var d = 1; d <= days; d++) {
    var dt = new Date(y, m, d);
    var id = fmt(dt);
    var sunday = dt.getDay() === 0;
    var cell = document.createElement("button");
    cell.type = "button";
    cell.textContent = String(d);
    cell.className = "day" + (sunday ? " sun" : " off") + (id === selected ? " pick" : "");
    if (!sunday) cell.disabled = true;
    else cell.addEventListener("click", (function (picked) {
      return function () {
        $("sunday").value = picked;
        $("cal").hidden = true;
        onSunday();
      };
    })(id));
    grid.appendChild(cell);
  }
  cal.appendChild(grid);
}

function toggleCal(e) {
  e.stopPropagation();
  if ($("sundayBtn").disabled) return;
  var cal = $("cal");
  if (!cal.hidden) { cal.hidden = true; return; }
  var id = $("sunday").value || thisSunday();
  var p = id.split("-").map(Number);
  calCursor = new Date(p[0], p[1] - 1, 1);
  renderCal();
  cal.hidden = false;
}

function onSunday() {
  var week = active();
  if (!week) return;
  var raw = $("sunday").value;
  var snapped = snapSunday(raw);
  if (!snapped) return;
  if (raw !== snapped) {
    $("sunday").value = snapped;
    status("Weeks start on Sunday.");
  }
  if (snapped !== week.id && state.weeks.some(function (w) { return w.id === snapped; })) {
    window.alert("You already have a week starting " + pretty(snapped) + ".");
    $("sunday").value = week.id;
    syncSundayFace();
    return;
  }
  week.id = snapped;
  state.activeId = snapped;
  document.title = "Turnip Prophet \u00b7 " + pretty(snapped);
  syncSundayFace();
  renderWeekList();
  saveSoon();
}

function addWeek() {
  if (active()) readFormInto(active());
  var latest = state.weeks.slice().sort(byNewest)[0];
  var sunday;
  if (!latest) sunday = thisSunday();
  else {
    var p = latest.id.split("-").map(Number);
    var dt = new Date(p[0], p[1] - 1, p[2]);
    dt.setDate(dt.getDate() + 7);
    sunday = fmt(dt);
  }
  if (state.weeks.some(function (w) { return w.id === sunday; })) {
    window.alert("You already have a week starting " + pretty(sunday) + ".");
    return;
  }
  var prev = latest && latest.likelyPattern >= 0 ? latest.likelyPattern : -1;
  var week = {
    id: sunday,
    firstBuy: false,
    previousPattern: prev,
    likelyPattern: -1,
    prices: [null, null, null, null, null, null, null, null, null, null, null, null, null],
    soldAt: null,
    note: "",
    updatedAt: new Date().toISOString()
  };
  state.weeks.push(week);
  state.activeId = week.id;
  setEnabled(true);
  fillForm(week);
  renderWeekList();
  saveNow();
  predict();
}

function deleteWeek() {
  var week = active();
  if (!week) return;
  if (!window.confirm("Delete the week of " + pretty(week.id) + "?")) return;
  state.weeks = state.weeks.filter(function (w) { return w.id !== week.id; });
  var left = state.weeks.slice().sort(byNewest);
  state.activeId = left.length ? left[0].id : null;
  if (!state.activeId) {
    renderWeekList();
    showNoWeek();
    saveNow();
    return;
  }
  fillForm(active());
  renderWeekList();
  saveNow();
  predict();
}

function importUrl() {
  var week = active();
  if (!week) return;
  var raw = $("importUrl").value.trim().replace(/\|+$/g, "");
  if (raw.indexOf("turnipprophet.io") === 0) raw = "https://" + raw;
  var url;
  try { url = new URL(raw); }
  catch (e) { status("That is not a link."); return; }
  var pricesParam = url.searchParams.get("prices");
  if (pricesParam == null) { status("No prices in that link."); return; }
  var parts = pricesParam.split(".").map(function (x) { return x === "" ? null : Number(x); });
  if (parts.length === 14 && parts[13] === null) parts.pop();
  while (parts.length < 13) parts.push(null);
  parts = parts.slice(0, 13);
  if (parts.some(function (v) { return v !== null && (!Number.isInteger(v) || v < 1 || v > 999); })) {
    status("A price in that link cannot be used.");
    return;
  }
  readFormInto(week);
  if (week.prices.some(function (v) { return v != null; }) && !window.confirm("Replace this week's prices with the link?")) return;
  week.prices = parts;
  week.firstBuy = url.searchParams.get("first") === "true";
  var pat = url.searchParams.get("pattern");
  week.previousPattern = pat == null || pat === "" || Number.isNaN(Number(pat)) ? -1 : Number(pat);
  fillForm(week);
  saveNow();
  predict();
  status("Imported.");
}

function buildPriceGrid() {
  var root = $("priceGrid");
  root.innerHTML = "";
  var buyRow = document.createElement("div");
  buyRow.className = "buy-row";
  var buyLabel = document.createElement("label");
  buyLabel.appendChild(document.createTextNode("Sunday Buy Price"));
  var line = document.createElement("div");
  line.className = "buy-line";
  var buy = document.createElement("input");
  buy.id = "buy";
  buy.type = "number";
  buy.min = "1";
  buy.max = "999";
  buy.step = "1";
  buy.inputMode = "numeric";
  var buyHint = document.createElement("p");
  buyHint.className = "hint";
  buyHint.innerHTML = "<span class=\"hours\">Daisy Mae</span> only sells on Sunday mornings. Her price is between <span class=\"hours\">90</span> and <span class=\"hours\">110</span>.";
  line.appendChild(buy);
  line.appendChild(buyHint);
  buyLabel.appendChild(line);
  buyRow.appendChild(buyLabel);
  root.appendChild(buyRow);

  var head = document.createElement("div");
  head.className = "price-head";
  ["", "8:00 AM\n– 11:59 AM", "12:00 PM\n– 10:00 PM"].forEach(function (text) {
    var span = document.createElement("span");
    text.split("\n").forEach(function (line, idx) {
      if (idx) span.appendChild(document.createElement("br"));
      span.appendChild(document.createTextNode(line));
    });
    head.appendChild(span);
  });
  root.appendChild(head);

  ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"].forEach(function (day, d) {
    var row = document.createElement("div");
    row.className = "price-row";
    var name = document.createElement("span");
    name.textContent = day;
    row.appendChild(name);
    for (var h = 0; h < 2; h++) {
      var input = document.createElement("input");
      input.type = "number";
      input.min = "1";
      input.max = "999";
      input.step = "1";
      input.id = "p" + (d * 2 + h);
      input.inputMode = "numeric";
      row.appendChild(input);
    }
    root.appendChild(row);
  });
  var hint = document.createElement("p");
  hint.className = "hint";
  hint.textContent = "The morning price is from 8:00 AM to 11:59 AM. A new price starts at noon and lasts until 10:00 PM. Do not erase the morning number.";
  root.appendChild(hint);
}

function bind() {
  buildPriceGrid();
  ["buy"].concat(Array.from({ length: 12 }, function (_, i) { return "p" + i; })).forEach(function (id) {
    $(id).addEventListener("input", function () { saveSoon(); schedulePredict(); });
  });
  $("firstBuy").addEventListener("change", function () { syncPrevLock(); saveSoon(); schedulePredict(); });
  $("prevPattern").addEventListener("change", function () { saveSoon(); schedulePredict(); });
  $("sundayBtn").addEventListener("click", toggleCal);
  $("cal").addEventListener("click", function (e) { e.stopPropagation(); });
  document.addEventListener("click", function () { $("cal").hidden = true; });
  $("note").addEventListener("input", saveSoon);
  $("soldAt").addEventListener("change", function () { saveSoon(); schedulePredict(); });
  $("saveBtn").addEventListener("click", function () { saveNow(); });
  $("newWeek").addEventListener("click", addWeek);
  $("deleteBtn").addEventListener("click", deleteWeek);
  $("importBtn").addEventListener("click", importUrl);
  $("importUrl").addEventListener("keydown", function (e) { if (e.key === "Enter") importUrl(); });
  $("quitBtn").addEventListener("click", function () {
    fetch("/api/shutdown", { method: "POST" }).finally(function () {
      document.body.innerHTML = '<p class="bye">Turnip Prophet has quit. You can close this window.</p>';
    });
  });
  window.addEventListener("resize", drawChart);
  markDue();
  window.setInterval(markDue, 30000);
}

function boot() {
  if (typeof Predictor !== "function") {
    status("Prediction engine failed to load.");
    return;
  }
  bind();
  fetch("/api/info").then(function (r) { return r.json(); }).then(function (info) {
    $("dataPath").textContent = info.root + "\\turnips.json";
  }).catch(function () {});
  fetch("/api/state").then(function (r) { return r.json(); }).then(function (data) {
    state = data && data.weeks ? data : { weeks: [], activeId: null };
    var seeded = false;
    if (!state.weeks.length) {
      state = JSON.parse(JSON.stringify(SEED));
      seeded = true;
    }
    state.weeks.forEach(function (week) {
      if (!week.prices || week.prices.length !== 13) {
        var next = [];
        for (var i = 0; i < 13; i++) next.push(week.prices && week.prices[i] != null ? week.prices[i] : null);
        week.prices = next;
      }
      if (week.previousPattern == null) week.previousPattern = -1;
      if (week.likelyPattern == null) week.likelyPattern = -1;
      if (week.note == null) week.note = "";
      if (week.soldAt === undefined) {
        if (week.id === "2026-10-18" && week.prices[10] === 191) week.soldAt = 10;
        else week.soldAt = null;
        seeded = true;
      }
    });
    if (!state.activeId || !state.weeks.some(function (w) { return w.id === state.activeId; })) {
      state.activeId = state.weeks.slice().sort(byNewest)[0].id;
    }
    setEnabled(true);
    fillForm(active());
    renderWeekList();
    if (seeded) saveNow();
    predict();
  }).catch(function () {
    status("Could not load saved weeks.");
  });
}

boot();
