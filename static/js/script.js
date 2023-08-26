var diffString = document.getElementById("diff-data").getAttribute("data-diff");
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
