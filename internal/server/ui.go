package server

import "net/http"

const adminPage = `<!doctype html>
<html>
<head>
  <link rel="icon" href="/images/GoChaosLogo.ico" type="image/x-icon">

  <meta charset="utf-8">
  <title>Go Chaos Admin</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <style>
    :root {
      --bg: #0f1710;
      --card: #161b22;
      --muted: #9aa4b2;
      --text: #e6edf3;
      --accent: #2ea043;
      --accent2: #58a6ff;
      --danger: #f85149;
      --border: #2d333b;
      --input: #0b0f14;
	  .header {
		display: flex;
		align-items: center;
		justify-content: left;
		gap: 12px;
		margin-bottom: 6px;
	  }
	  .logo {
		width: 216px;
		height: 216px;
		object-fit: contain;
	   }

    }
    body {
      margin: 0; padding: 24px;
      font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, sans-serif;
      background: radial-gradient(1200px 800px at 10% -10%, #1b2230, transparent),
                  radial-gradient(1200px 800px at 110% 10%, #1a2b1f, transparent),
                  var(--bg);
      color: var(--text);
    }
    h1 { margin: 0 0 6px; font-size: 24px; }
    p { margin: 0 0 20px; color: var(--muted); }
    .grid { display: grid; gap: 16px; grid-template-columns: 1fr; }
    @media (min-width: 900px) { .grid { grid-template-columns: 1fr 1fr; } }
    .card {
      background: var(--card);
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 16px;
      box-shadow: 0 6px 24px rgba(0,0,0,0.25);
    }
    label { display: block; font-size: 12px; color: var(--muted); margin-bottom: 6px; }
    input, textarea, select {
      width: 100%;
      background: var(--input);
      border: 1px solid var(--border);
      color: var(--text);
      padding: 8px 10px;
      border-radius: 8px;
      outline: none;
    }
    textarea { min-height: 140px; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
    .row { display: grid; gap: 12px; grid-template-columns: 1fr 1fr; }
    .row3 { display: grid; gap: 12px; grid-template-columns: 1fr 1fr 1fr; }
    .switch { display: flex; align-items: center; gap: 10px; }
    .switch input { width: 18px; height: 18px; }
    .range { display: flex; align-items: center; gap: 8px; }
    .range input[type="range"] { flex: 1; }
    .btns { display: flex; gap: 8px; flex-wrap: wrap; }
    button {
      background: var(--accent);
      border: none;
      color: #08110a;
      padding: 8px 12px;
      border-radius: 8px;
      cursor: pointer;
      font-weight: 600;
    }
    button.secondary { background: var(--accent2); color: #08111a; }
    button.ghost { background: transparent; color: var(--text); border: 1px solid var(--border); }
    .status { margin-top: 8px; font-size: 12px; color: var(--muted); }
    .error { color: var(--danger); }
    .pill { display:inline-block; padding:2px 8px; border:1px solid var(--border); border-radius:999px; font-size:11px; color:var(--muted); }
  </style>
</head>
<body>
	<div class="header">
	<img src="/images/GoChaosLogo.png" alt="Go Chaos" class="logo">
	<h1>Go Chaos Admin</h1>
	</div>
	<p>Live update chaos settings without restarting the proxy.</p>

  <div class="grid">
    <div class="card">
      <div class="row">
        <div>
          <label>Listen Address</label>
          <input id="listen_addr" placeholder=":8080">
        </div>
        <div>
          <label>Target URL</label>
          <input id="target_url" placeholder="http://localhost:9000">
        </div>
      </div>

      <div style="margin-top:12px" class="row3">
        <div>
          <label>Error Rate (0-1)</label>
          <input id="error_rate" type="number" min="0" max="1" step="0.01" value="0">
        </div>
        <div>
          <label>Disconnect Rate (0-1)</label>
          <input id="disconnect_rate" type="number" min="0" max="1" step="0.01" value="0">
        </div>
        <div>
          <label>Upstream Timeout Rate (0-1)</label>
          <input id="upstream_timeout_rate" type="number" min="0" max="1" step="0.01" value="0">
        </div>
      </div>

      <div style="margin-top:12px" class="row3">
        <div>
          <label>DNS Failure Rate (0-1)</label>
          <input id="dns_failure_rate" type="number" min="0" max="1" step="0.01" value="0">
        </div>
        <div>
          <label>Latency Min (ms)</label>
          <input id="latency_min_ms" type="number" min="0" step="10" value="0">
        </div>
        <div>
          <label>Latency Max (ms)</label>
          <input id="latency_max_ms" type="number" min="0" step="10" value="0">
        </div>
      </div>

      <div style="margin-top:12px" class="row">
        <div>
          <label>Include Paths (one per line)</label>
          <textarea id="include_paths" placeholder="/api/"></textarea>
        </div>
        <div>
          <label>Exclude Paths (one per line)</label>
          <textarea id="exclude_paths" placeholder="/healthz&#10;/admin/"></textarea>
        </div>
      </div>

      <div class="btns" style="margin-top:12px">
        <button onclick="loadConfig()">Load</button>
        <button class="secondary" onclick="saveConfig()">Save</button>
        <button class="ghost" onclick="resetToDefaults()">Defaults</button>
      </div>
      <div class="status" id="status"></div>
    </div>

    <div class="card">
      <div class="pill">Raw YAML</div>
      <textarea id="raw"></textarea>
      <div class="btns" style="margin-top:12px">
        <button class="ghost" onclick="syncFromForm()">Sync From Form</button>
        <button class="ghost" onclick="syncToForm()">Sync To Form</button>
      </div>
      <div class="status" id="rawStatus"></div>
    </div>
  </div>

<script>
const $ = (id) => document.getElementById(id);

function setStatus(msg, isError=false) {
  const el = $("status");
  el.textContent = msg;
  el.className = "status" + (isError ? " error" : "");
}

function setRawStatus(msg, isError=false) {
  const el = $("rawStatus");
  el.textContent = msg;
  el.className = "status" + (isError ? " error" : "");
}

function readLines(id) {
  return $(id).value.split("\n").map(s => s.trim()).filter(Boolean);
}

function toConfigObj() {
  return {
    listen_addr: $("listen_addr").value.trim(),
    target_url: $("target_url").value.trim(),
    chaos: {
      error_rate: Number($("error_rate").value),
      disconnect_rate: Number($("disconnect_rate").value),
      upstream_timeout_rate: Number($("upstream_timeout_rate").value),
      dns_failure_rate: Number($("dns_failure_rate").value),
      latency_min_ms: Number($("latency_min_ms").value),
      latency_max_ms: Number($("latency_max_ms").value),
      include_paths: readLines("include_paths"),
      exclude_paths: readLines("exclude_paths")
    }
  };
}

function validateConfig(cfg) {
  const errs = [];
  if (!cfg.listen_addr) errs.push("listen_addr required");
  if (!cfg.target_url) errs.push("target_url required");
  const rates = ["error_rate","disconnect_rate","upstream_timeout_rate","dns_failure_rate"];
  for (const r of rates) {
    const v = cfg.chaos[r];
    if (v < 0 || v > 1 || Number.isNaN(v)) errs.push(r + " must be 0-1");
  }
  if (cfg.chaos.latency_min_ms < 0 || cfg.chaos.latency_max_ms < 0) {
    errs.push("latency min/max must be >= 0");
  }
  if (cfg.chaos.latency_max_ms < cfg.chaos.latency_min_ms) {
    errs.push("latency max must be >= latency min");
  }
  const allPaths = [...cfg.chaos.include_paths, ...cfg.chaos.exclude_paths];
  for (const p of allPaths) {
    if (!p.startsWith("/")) errs.push("paths must start with '/'");
  }
  return errs;
}

async function loadConfig() {
  const res = await fetch("/admin/config");
  if (!res.ok) { setStatus("Load failed: " + res.status, true); return; }
  const text = await res.text();
  $("raw").value = text;
  syncToForm();
  setStatus("Loaded");
}

async function saveConfig() {
  syncFromForm();
  const cfg = toConfigObj();
  const errs = validateConfig(cfg);
  if (errs.length) { setStatus(errs.join("; "), true); return; }
  const res = await fetch("/admin/config", {
    method: "POST",
    headers: { "Content-Type": "application/x-yaml" },
    body: $("raw").value
  });
  if (!res.ok) {
    const err = await res.text();
    setStatus("Save failed: " + err, true);
    return;
  }
  setStatus("Saved");
}

function resetToDefaults() {
  $("listen_addr").value = ":8080";
  $("target_url").value = "http://localhost:9000";
  $("error_rate").value = 0;
  $("disconnect_rate").value = 0;
  $("upstream_timeout_rate").value = 0;
  $("dns_failure_rate").value = 0;
  $("latency_min_ms").value = 0;
  $("latency_max_ms").value = 0;
  $("include_paths").value = "/api/";
  $("exclude_paths").value = "/healthz\n/admin/";
  syncFromForm();
  setStatus("Defaults loaded");
}

function syncFromForm() {
  const cfg = toConfigObj();
  // simple YAML output
  const yaml =
    'listen_addr: "' + cfg.listen_addr + '"\n' +
    'target_url: "' + cfg.target_url + '"\n' +
    'chaos:\n' +
    '  error_rate: ' + cfg.chaos.error_rate + '\n' +
    '  disconnect_rate: ' + cfg.chaos.disconnect_rate + '\n' +
    '  upstream_timeout_rate: ' + cfg.chaos.upstream_timeout_rate + '\n' +
    '  dns_failure_rate: ' + cfg.chaos.dns_failure_rate + '\n' +
    '  latency_min_ms: ' + cfg.chaos.latency_min_ms + '\n' +
    '  latency_max_ms: ' + cfg.chaos.latency_max_ms + '\n' +
    '  include_paths:\n' +
    cfg.chaos.include_paths.map(p => '    - "' + p + '"').join('\n') + '\n' +
    '  exclude_paths:\n' +
    cfg.chaos.exclude_paths.map(p => '    - "' + p + '"').join('\n');

  $("raw").value = yaml;
  setRawStatus("Synced from form");
}

function syncToForm() {
  // basic YAML parse: use the server GET output as source of truth
  // We'll just load it in a naive way by calling /admin/config and then manual edits.
  // For full YAML parsing in the browser, add a JS YAML library.
  setRawStatus("Edit raw YAML directly or use form + Sync");
}

loadConfig();
</script>
</body>
</html>`

func (s *Server) handleAdminUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(adminPage))
}
