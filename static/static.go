package static

// HTML is the template for the HTML page
const HTML = `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta charset="utf-8" />
    <!-- Make sure to load the highlight.js CSS file before the Diff2Html CSS file -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.1/styles/github.min.css" />
    <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/diff2html/bundles/css/diff2html.min.css" />
    <!-- Include Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <!-- Include Bootstrap Icons (for GitHub icon) -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons/font/bootstrap-icons.css" rel="stylesheet">
    <!-- Include custom CSS -->
	<style>
		.navbar {
			background-color: #010b18; /* Dark blue color */
		}
		.navbar-brand {
			color: #fff;
			font-weight: bold;
		}
		.github-star {
			font-size: 24px;
			color: #fff;
			margin-right: 20px;
		}
		
		#myDiffElement {
			margin-left: 10%;
			margin-right: 10%;
		}
	</style>
    <script type="text/javascript" src="https://cdn.jsdelivr.net/npm/diff2html/bundles/js/diff2html-ui.min.js"></script>
    <script async defer src="https://buttons.github.io/buttons.js"></script>
</head>
<body>
    <!-- Navigation Bar -->
    <nav class="navbar">
        <div class="container">
            <div class="navbar-brand">{{.Title}}</div>
            <div class="justify-content-end">
                <a class="github-button"
                   href="https://github.com/imrajdas/diffr"
                   data-icon="octicon-star"
                   data-size="large"
                   data-show-count="true"
                   aria-label="Star ntkme/github-buttons on GitHub">Star</a>
            </div>
        </div>
    </nav>
    </br>
    <div id="myDiffElement" class="container-left"></div>
    <div id="diff-data" hidden>{{.Diff}}</div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <script>
		var diffString = document.getElementById("diff-data").textContent;
		document.addEventListener('DOMContentLoaded', function () {
			var targetElement = document.getElementById('myDiffElement');
			var configuration = {
				drawFileList: true,
				fileListToggle: true,
				fileListStartVisible: true,
				fileContentToggle: true,
				matching: 'lines',
				outputFormat: 'side-by-side',
				synchronisedScroll: true,
				highlight: true,
				highlightLanguages: true,
				renderNothingWhenEmpty: false,
			};
			var diff2htmlUi = new Diff2HtmlUI(targetElement, diffString, configuration);
			diff2htmlUi.draw();
			diff2htmlUi.highlightCode();
		
			// Dark mode toggle functionality
			const darkModeToggle = document.getElementById('darkModeToggle');
			const body = document.body;
		
			darkModeToggle.addEventListener('click', () => {
				body.classList.toggle('dark-mode');
			});
		});
	</script>
</body>
</html>
`
