package handlers

const IndexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Azzurro Technology Platform</title>
    <style>
        :root { --primary: #0056b3; --bg: #f4f6f9; --card: #ffffff; --text: #333; --border: #ddd; }
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; background: var(--bg); color: var(--text); margin: 0; padding: 20px; line-height: 1.6; }
        .container { max-width: 900px; margin: 0 auto; }
        header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; border-bottom: 2px solid var(--border); padding-bottom: 10px; }
        h1 { margin: 0; font-size: 1.5rem; color: var(--primary); }
        .controls { display: flex; gap: 10px; }
        button { background: var(--primary); color: white; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer; font-size: 0.9rem; }
        button:hover { opacity: 0.9; }
        button.secondary { background: #6c757d; }
        .input-group { display: flex; gap: 10px; margin-bottom: 30px; background: var(--card); padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); }
        input[type="text"], textarea { width: 100%; padding: 10px; border: 1px solid var(--border); border-radius: 4px; font-family: inherit; }
        textarea { height: 80px; resize: vertical; }
        .source-list { display: grid; gap: 20px; }
        .source-card { background: var(--card); border: 1px solid var(--border); border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.05); }
        .card-header { padding: 15px 20px; background: #fafafa; border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; }
        .card-header a { text-decoration: none; color: var(--primary); font-weight: bold; font-size: 1.1rem; }
        .card-body { padding: 20px; }
        details { margin-top: 10px; }
        summary { cursor: pointer; color: #555; font-weight: 500; outline: none; }
        summary:hover { color: var(--primary); }
        .note-box { background: #fffbe6; border-left: 4px solid #ffc107; padding: 10px; margin-top: 10px; font-size: 0.9rem; }
        .meta { font-size: 0.8rem; color: #888; margin-top: 5px; }
        .empty-state { text-align: center; color: #888; margin-top: 50px; }
        .modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: none; align-items: center; justify-content: center; z-index: 1000; }
        .modal { background: white; padding: 25px; border-radius: 8px; width: 90%; max-width: 500px; }
        .modal h2 { margin-top: 0; }
        .modal textarea { height: 150px; margin: 10px 0; }
        .modal-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 15px; }
    </style>
</head>
<body>
<div class="container">
    <header>
        <h1>Azzurro Platform</h1>
        <div class="controls">
            <button onclick="exportConfig()" title="Copy configuration to clipboard">Export Config</button>
            <button class="secondary" onclick="openImportModal()">Import Config</button>
        </div>
    </header>
    <div class="input-group">
        <form id="addForm" method="POST" action="/add" style="display:flex; flex-direction:column; width:100%;">
            <label style="font-weight:bold; margin-bottom:5px;">Add New Source URL</label>
            <input type="text" name="url" placeholder="https://example.com/article" required autocomplete="off">
            <textarea name="notes" placeholder="Initial notes (optional)" style="margin-top:10px;"></textarea>
            <button type="submit" style="margin-top:10px; align-self:flex-start;">Generate Article Tag</button>
        </form>
    </div>
    <div id="sourcesList" class="source-list">
        {{if .Sources}}
            {{range .Sources}}
            <article class="source-card">
                <div class="card-header">
                    <a href="{{.URL}}" target="_blank">{{if .Title}}{{.Title}}{{else}}{{.URL}}{{end}}</a>
                    <span class="meta">{{.CreatedAt.Format "Jan 02, 2006"}}</span>
                </div>
                <div class="card-body">
                    {{if .Summary}}
                    <details>
                        <summary>View Summary</summary>
                        <p>{{.Summary}}</p>
                    </details>
                    {{end}}
                    {{if .Notes}}
                    <div class="note-box">
                        <strong>Notes:</strong><br>
                        {{.Notes | nl2br}}
                    </div>
                    {{end}}
                    <details>
                        <summary>Details / Metadata</summary>
                        <p>ID: {{.ID}}</p>
                        <p>URL: {{.URL}}</p>
                    </details>
                </div>
            </article>
            {{end}}
        {{else}}
            <div class="empty-state">
                <p>No sources added yet. Enter a URL above to begin.</p>
            </div>
        {{end}}
    </div>
</div>
<div id="importModal" class="modal-overlay">
    <div class="modal">
        <h2>Import Configuration</h2>
        <p>Paste the Base64 configuration string here:</p>
        <textarea id="importData"></textarea>
        <div class="modal-actions">
            <button class="secondary" onclick="closeImportModal()">Cancel</button>
            <button onclick="performImport()">Import</button>
        </div>
    </div>
</div>
<script>
    function nl2br(str) { return str.replace(/\n/g, '<br>'); }
    function exportConfig() {
        fetch('/api/export').then(res => res.text()).then(data => {
            navigator.clipboard.writeText(data).then(() => alert('Configuration copied to clipboard!'));
        }).catch(err => alert('Error exporting config: ' + err));
    }
    function openImportModal() { document.getElementById('importModal').style.display = 'flex'; }
    function closeImportModal() { document.getElementById('importModal').style.display = 'none'; }
    function performImport() {
        const data = document.getElementById('importData').value.trim();
        if (!data) return alert('Please paste a configuration string.');
        window.location.href = '/?config=' + encodeURIComponent(data);
    }
    const urlParams = new URLSearchParams(window.location.search);
    const configParam = urlParams.get('config');
    if (configParam) {
        fetch('/api/import', { method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({config: configParam}) })
        .then(res => { if(res.ok) { const newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname; window.history.pushState({}, document.title, newUrl); window.location.reload(); } else { alert('Failed to import configuration.'); } })
        .catch(err => alert('Import error: ' + err));
    }
</script>
</body>
</html>`
