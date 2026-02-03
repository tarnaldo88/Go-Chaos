package server

import (
	"net/http"
)

const adminPage = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Go Chaos Admin</title>
  <style>
    body { font-family: sans-serif; margin: 24px; }
    h1 { margin-bottom: 8px; }
    textarea { width: 100%; height: 320px; font-family: monospace; }
    button { margin-right: 8px; padding: 8px 12px; }
    .row { margin: 12px 0; }
    .status { margin-top: 8px; font-size: 0.9rem; }
  </style>
</head>
<body>
  <h1>Go Chaos Admin</h1>
  <div class="row">
    <button onclick="loadConfig()">Load</button>
    <button onclick="saveConfig()">Save</button>
  </div>
  <textarea id="cfg"></textarea>
  <div class="status" id="status"></div>

<script>
async function loadConfig() {
  const res = await fetch('/admin/config');
  if (!res.ok) {
    setStatus('Load failed: ' + res.status);
    return;
  }
  const text = await res.text();
  document.getElementById('cfg').value = text;
  setStatus('Loaded');
}

async function saveConfig() {
  const body = document.getElementById('cfg').value;
  const res = await fetch('/admin/config', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-yaml' },
    body
  });
  if (!res.ok) {
    const err = await res.text();
    setStatus('Save failed: ' + err);
    return;
  }
  setStatus('Saved');
}

function setStatus(msg) {
  document.getElementById('status').innerText = msg;
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
