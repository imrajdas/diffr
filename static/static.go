package static

// HTML is the template for the HTML page
const HTML = `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta charset="utf-8" />
    <link id="hljs-light" rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.1/styles/github.min.css" />
    <link id="hljs-dark" rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.1/styles/github-dark.min.css" disabled />
    <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/diff2html/bundles/css/diff2html.min.css" />
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons/font/bootstrap-icons.css" rel="stylesheet">
	<style>
		:root {
			--bg: #ffffff;
			--fg: #212529;
			--muted: #6c757d;
			--panel: #f4f6f8;
			--border: #e2e8f0;
			--nav: #010b18;
		}
		body {
			background: var(--bg);
			color: var(--fg);
		}
		body.dark-mode {
			--bg: #0f1419;
			--fg: #e6edf3;
			--muted: #9aa4b2;
			--panel: #1c2333;
			--border: #2d3a4f;
			--nav: #010b18;
		}
		.navbar { background-color: var(--nav); }
		.navbar-brand { color: #fff; font-weight: bold; }
		.content { margin-left: 10%; margin-right: 10%; margin-top: 1.5rem; padding-bottom: 3rem; }
		#myDiffElement { margin-bottom: 2rem; }
		.stats-bar, .toolbar, .media-card {
			background: var(--panel);
			border: 1px solid var(--border);
			border-radius: 6px;
		}
		.stats-bar { padding: 0.75rem 1rem; margin-bottom: 1rem; font-size: 0.95rem; }
		.toolbar {
			display: flex;
			flex-wrap: wrap;
			gap: 0.6rem;
			align-items: center;
			padding: 0.75rem 1rem;
			margin-bottom: 1.5rem;
		}
		.toolbar input[type="search"] { min-width: 220px; flex: 1; }
		.media-card { padding: 1rem; margin-bottom: 1rem; }
		.img-compare-frame {
			position: relative;
			overflow: hidden;
			border: 1px solid var(--border);
			background: #111;
			cursor: ew-resize;
			user-select: none;
		}
		.img-compare-frame img { display: block; width: 100%; height: auto; pointer-events: none; }
		.img-compare-left-wrap {
			position: absolute;
			top: 0; left: 0; bottom: 0;
			width: 50%;
			overflow: hidden;
		}
		.img-compare-left-wrap img { max-width: none; }
		.img-compare-handle {
			position: absolute;
			top: 0; bottom: 0;
			left: 50%;
			width: 2px;
			background: #fff;
			transform: translateX(-50%);
			pointer-events: none;
		}
		.img-compare-label {
			position: absolute;
			top: 8px;
			padding: 2px 8px;
			font-size: 0.75rem;
			background: rgba(0,0,0,0.55);
			color: #fff;
			pointer-events: none;
		}
		.img-compare-label.left { left: 8px; }
		.img-compare-label.right { right: 8px; }
		.img-pixel-diff { max-width: 100%; border: 1px solid var(--border); margin-top: 0.75rem; }
		.pdf-unified {
			max-height: 240px;
			overflow: auto;
			background: var(--bg);
			border: 1px solid var(--border);
			padding: 0.75rem;
			font-size: 0.8rem;
			white-space: pre-wrap;
		}
		.text-muted, .muted { color: var(--muted) !important; }
		body.dark-mode .table { color: var(--fg); --bs-table-bg: transparent; }
		body.dark-mode .form-control, body.dark-mode .form-select {
			background: #0f1419;
			color: var(--fg);
			border-color: var(--border);
		}
		body.dark-mode .btn-outline-secondary { color: #d0d7de; border-color: var(--border); }
		body.dark-mode .btn-outline-secondary.active, body.dark-mode .btn-outline-secondary:hover {
			background: #2d3a4f;
			color: #fff;
		}
		body.dark-mode .d2h-wrapper { color: var(--fg); }
		body.dark-mode .d2h-file-header { background-color: #1c2333; color: var(--fg); border-color: var(--border); }
		body.dark-mode .d2h-file-wrapper { border-color: var(--border); }
		body.dark-mode .d2h-file-name { color: var(--fg); }
		body.dark-mode .d2h-code-side-linenumber,
		body.dark-mode .d2h-code-linenumber { background-color: #161b22; color: #8b949e; border-color: var(--border); }
		body.dark-mode .d2h-ins { background-color: #1f3d2a; }
		body.dark-mode .d2h-del { background-color: #3d1f1f; }
		body.dark-mode .d2h-code-line del, body.dark-mode .d2h-code-side-line del { background-color: #5c2020; }
		body.dark-mode .d2h-code-line ins, body.dark-mode .d2h-code-side-line ins { background-color: #1b5c32; }
		body.dark-mode .d2h-emptyplaceholder, body.dark-mode .d2h-code-side-emptyplaceholder { background: #161b22; }
		body.dark-mode .d2h-file-list-wrapper { background: var(--panel); color: var(--fg); }
		#filterEmpty { display: none; }
		.nav-actions { display: flex; align-items: center; gap: 0.75rem; }
	</style>
    <script type="text/javascript" src="https://cdn.jsdelivr.net/npm/diff2html/bundles/js/diff2html-ui.min.js"></script>
    <script async defer src="https://buttons.github.io/buttons.js"></script>
</head>
<body>
    <nav class="navbar">
        <div class="container">
            <div class="navbar-brand">{{.Title}}</div>
            <div class="nav-actions">
                <button type="button" id="darkModeToggle" class="btn btn-sm btn-outline-light" title="Toggle dark mode">
                    <i class="bi bi-moon"></i> Dark
                </button>
                <a class="github-button"
                   href="https://github.com/imrajdas/diffr"
                   data-icon="octicon-star"
                   data-size="large"
                   data-show-count="true"
                   aria-label="Star imrajdas/diffr on GitHub">Star</a>
            </div>
        </div>
    </nav>
    <div class="content">
        <div class="stats-bar">
            <strong>{{.Stats.Changed}}</strong> changed ·
            <strong>{{.Stats.Added}}</strong> added ·
            <strong>{{.Stats.Removed}}</strong> removed ·
            <strong>{{.Stats.Identical}}</strong> identical
            <span class="muted">· {{.Stats.Elapsed}}</span>
        </div>
        <div class="toolbar">
            <input id="filterQuery" class="form-control form-control-sm" type="search" placeholder="Filter by filename or extension" />
            <select id="filterStatus" class="form-select form-select-sm" style="max-width: 160px;">
                <option value="">All statuses</option>
                <option value="changed">Changed</option>
                <option value="added">Added</option>
                <option value="removed">Removed</option>
            </select>
            <div class="btn-group btn-group-sm" role="group" aria-label="Diff format">
                <button type="button" id="formatSide" class="btn btn-outline-secondary">Side-by-side</button>
                <button type="button" id="formatUnified" class="btn btn-outline-secondary">Unified</button>
            </div>
        </div>
        {{if not .HasAny}}
        <div class="alert alert-success">No differences found.</div>
        {{end}}
        <div id="filterEmpty" class="alert alert-secondary">No files match this filter.</div>

        {{if .Images}}
        <div data-section="images">
        <h5 class="mt-4">Images</h5>
        {{range .Images}}
        <section class="media-card filterable" data-path="{{.RelPath}}" data-status="{{.Status}}">
            <div class="d-flex justify-content-between align-items-center mb-2">
                <strong>{{.RelPath}}</strong>
                <span class="badge text-bg-secondary">{{.Status}}</span>
            </div>
            <p class="mb-3 muted">{{.Summary}}</p>
            {{if and .HasLeft .HasRight}}
            <div class="img-compare">
                <div class="img-compare-frame">
                    <img class="img-compare-right" src="{{.RightURL}}" alt="right {{.RelPath}}">
                    <div class="img-compare-left-wrap">
                        <img class="img-compare-left" src="{{.LeftURL}}" alt="left {{.RelPath}}">
                    </div>
                    <div class="img-compare-handle"></div>
                    <span class="img-compare-label left">Left</span>
                    <span class="img-compare-label right">Right</span>
                </div>
                <input class="form-range mt-2" type="range" min="0" max="100" value="50" aria-label="Image compare slider">
            </div>
            {{else if .HasLeft}}
            <img src="{{.LeftURL}}" alt="left {{.RelPath}}" class="img-pixel-diff">
            {{else if .HasRight}}
            <img src="{{.RightURL}}" alt="right {{.RelPath}}" class="img-pixel-diff">
            {{end}}
            {{if .HasDiff}}
            <details class="mt-2">
                <summary>Pixel diff</summary>
                <img src="{{.DiffURL}}" alt="diff {{.RelPath}}" class="img-pixel-diff">
            </details>
            {{end}}
        </section>
        {{end}}
        </div>
        {{end}}

        {{if .PDFs}}
        <div data-section="pdfs">
        <h5 class="mt-4">PDFs</h5>
        <table class="table table-sm">
            <thead><tr><th>File</th><th>Status</th><th>Summary</th></tr></thead>
            <tbody>
            {{range .PDFs}}
            <tr class="filterable" data-path="{{.RelPath}}" data-status="{{.Status}}">
                <td>{{.RelPath}}</td>
                <td>{{.Status}}</td>
                <td>
                    {{.Summary}}
                    {{if .Unified}}
                    <details class="mt-1"><summary>Extracted text diff</summary><pre class="pdf-unified">{{.Unified}}</pre></details>
                    {{end}}
                </td>
            </tr>
            {{end}}
            </tbody>
        </table>
        </div>
        {{end}}

        {{if .Binaries}}
        <div data-section="binaries">
        <h5 class="mt-4">Binaries</h5>
        <table class="table table-sm">
            <thead><tr><th>File</th><th>Status</th><th>Summary</th></tr></thead>
            <tbody>
            {{range .Binaries}}
            <tr class="filterable" data-path="{{.RelPath}}" data-status="{{.Status}}">
                <td>{{.RelPath}}</td>
                <td>{{.Status}}</td>
                <td>{{.Summary}}</td>
            </tr>
            {{end}}
            </tbody>
        </table>
        </div>
        {{end}}

        <div id="textDiffSection" data-section="text">
            <div id="myDiffElement"></div>
        </div>
        <div id="diff-data" hidden>{{.Diff}}</div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script>
		var diffString = "";
		var diff2htmlUi = null;
		var outputFormat = "side-by-side";

		function extOf(path) {
			var base = (path || "").split("/").pop();
			var i = base.lastIndexOf(".");
			return i >= 0 ? base.slice(i + 1).toLowerCase() : "";
		}

		function applyFilter() {
			var q = (document.getElementById("filterQuery").value || "").trim().toLowerCase();
			var status = document.getElementById("filterStatus").value;
			var qExt = q.replace(/^\./, "");
			var anyVisible = false;
			var total = 0;
			document.querySelectorAll(".filterable").forEach(function (el) {
				total += 1;
				var path = (el.getAttribute("data-path") || "").toLowerCase();
				var st = el.getAttribute("data-status") || "";
				var ext = extOf(path);
				var matchQ = !q || path.indexOf(q) !== -1 || ext === qExt;
				var matchS = !status || st === status;
				var show = matchQ && matchS;
				el.hidden = !show;
				if (show) anyVisible = true;
			});

			document.querySelectorAll("[data-section]").forEach(function (sec) {
				var vis = sec.querySelectorAll(".filterable:not([hidden])").length;
				sec.hidden = vis === 0;
				if (vis > 0) anyVisible = true;
			});

			var empty = document.getElementById("filterEmpty");
			empty.style.display = (!anyVisible && total > 0) ? "block" : "none";
		}

		function tagDiffWrappers() {
			document.querySelectorAll(".d2h-file-wrapper").forEach(function (el) {
				el.classList.add("filterable");
				var nameEl = el.querySelector(".d2h-file-name");
				var path = nameEl ? nameEl.textContent.trim() : "";
				el.setAttribute("data-path", path);
				var tag = (el.querySelector(".d2h-tag") || {}).textContent || "";
				tag = tag.toLowerCase();
				var st = "changed";
				if (tag.indexOf("added") !== -1) st = "added";
				else if (tag.indexOf("deleted") !== -1) st = "removed";
				el.setAttribute("data-status", st);
			});
			document.querySelectorAll(".d2h-file-list-line").forEach(function (el) {
				el.classList.add("filterable");
				var nameEl = el.querySelector(".d2h-file-name");
				el.setAttribute("data-path", nameEl ? nameEl.textContent.trim() : "");
				var tag = (el.querySelector(".d2h-tag") || {}).textContent || "";
				tag = tag.toLowerCase();
				var st = "changed";
				if (tag.indexOf("added") !== -1) st = "added";
				else if (tag.indexOf("deleted") !== -1) st = "removed";
				el.setAttribute("data-status", st);
			});
		}

		function drawDiff() {
			var target = document.getElementById("myDiffElement");
			if (!diffString || !diffString.trim()) {
				target.innerHTML = "";
				return;
			}
			target.innerHTML = "";
			diff2htmlUi = new Diff2HtmlUI(target, diffString, {
				drawFileList: true,
				fileListToggle: true,
				fileListStartVisible: true,
				fileContentToggle: true,
				matching: "lines",
				outputFormat: outputFormat,
				synchronisedScroll: true,
				highlight: true,
				highlightLanguages: true,
				renderNothingWhenEmpty: true,
			});
			diff2htmlUi.draw();
			diff2htmlUi.highlightCode();
			tagDiffWrappers();
			applyFilter();
		}

		function setFormat(fmt) {
			outputFormat = fmt;
			try { localStorage.setItem("diffr-format", fmt); } catch (e) {}
			document.getElementById("formatSide").classList.toggle("active", fmt === "side-by-side");
			document.getElementById("formatUnified").classList.toggle("active", fmt === "line-by-line");
			drawDiff();
		}

		function setDark(on) {
			document.body.classList.toggle("dark-mode", on);
			document.getElementById("hljs-light").disabled = on;
			document.getElementById("hljs-dark").disabled = !on;
			var btn = document.getElementById("darkModeToggle");
			btn.innerHTML = on
				? '<i class="bi bi-sun"></i> Light'
				: '<i class="bi bi-moon"></i> Dark';
			try { localStorage.setItem("diffr-dark", on ? "1" : "0"); } catch (e) {}
			if (diff2htmlUi) {
				diff2htmlUi.highlightCode();
			}
		}

		function setComparePct(root, pct) {
			pct = Math.max(0, Math.min(100, pct));
			var wrap = root.querySelector(".img-compare-left-wrap");
			var handle = root.querySelector(".img-compare-handle");
			var range = root.querySelector("input[type=range]");
			wrap.style.width = pct + "%";
			handle.style.left = pct + "%";
			if (range) range.value = String(pct);
		}

		function syncCompareWidths(root) {
			var frame = root.querySelector(".img-compare-frame");
			var leftImg = root.querySelector(".img-compare-left");
			if (frame && leftImg) {
				leftImg.style.width = frame.clientWidth + "px";
			}
		}

		function initImageComparisons() {
			document.querySelectorAll(".img-compare").forEach(function (root) {
				var leftImg = root.querySelector(".img-compare-left");
				var rightImg = root.querySelector(".img-compare-right");
				function sync() { syncCompareWidths(root); }
				if (leftImg) leftImg.addEventListener("load", sync);
				if (rightImg) rightImg.addEventListener("load", sync);
				sync();
				setComparePct(root, 50);
				var frame = root.querySelector(".img-compare-frame");
				var range = root.querySelector("input[type=range]");
				var dragging = false;
				function fromEvent(e) {
					var rect = frame.getBoundingClientRect();
					var x = (e.touches ? e.touches[0].clientX : e.clientX) - rect.left;
					setComparePct(root, (x / rect.width) * 100);
				}
				frame.addEventListener("pointerdown", function (e) {
					dragging = true;
					frame.setPointerCapture(e.pointerId);
					fromEvent(e);
				});
				frame.addEventListener("pointermove", function (e) {
					if (dragging) fromEvent(e);
				});
				frame.addEventListener("pointerup", function () { dragging = false; });
				if (range) {
					range.addEventListener("input", function () { setComparePct(root, Number(range.value)); });
				}
			});
			window.addEventListener("resize", function () {
				document.querySelectorAll(".img-compare").forEach(syncCompareWidths);
			});
		}

		document.addEventListener("DOMContentLoaded", function () {
			diffString = document.getElementById("diff-data").textContent;
			var dark = false;
			var fmt = "side-by-side";
			try {
				dark = localStorage.getItem("diffr-dark") === "1";
				fmt = localStorage.getItem("diffr-format") || fmt;
			} catch (e) {}
			setDark(dark);
			outputFormat = fmt;
			document.getElementById("formatSide").classList.toggle("active", fmt === "side-by-side");
			document.getElementById("formatUnified").classList.toggle("active", fmt === "line-by-line");
			document.getElementById("formatSide").addEventListener("click", function () { setFormat("side-by-side"); });
			document.getElementById("formatUnified").addEventListener("click", function () { setFormat("line-by-line"); });
			document.getElementById("darkModeToggle").addEventListener("click", function () {
				setDark(!document.body.classList.contains("dark-mode"));
			});
			document.getElementById("filterQuery").addEventListener("input", applyFilter);
			document.getElementById("filterStatus").addEventListener("change", applyFilter);
			drawDiff();
			initImageComparisons();
			applyFilter();
		});
	</script>
</body>
</html>
`
