package diffr

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/imrajdas/diffr/static"
	"github.com/spf13/cobra"
)

func openBrowser(rawURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "linux":
		cmd = exec.Command("xdg-open", rawURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", rawURL)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}

func Run(cmd *cobra.Command, args []string) {
	opts := Options{
		Exclude:          Exclude,
		IgnoreFile:       IgnoreFile,
		NoDefaultExclude: NoDefaultExclude,
		NoGitignore:      NoGitignore,
		IgnoreWhitespace: IgnoreWhitespace,
		IgnoreCase:       IgnoreCase,
		Context:          Context,
	}

	res, err := Compare(args[0], args[1], opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if Stdout || JSON || PatchFile != "" {
		if err := WriteCLI(res, Stdout, JSON, PatchFile); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
		if res.HasChanges() {
			os.Exit(1)
		}
		return
	}

	runWeb(res)
}

func runWeb(res *Result) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		renderPage(w, res)
	})
	mux.HandleFunc("/media/", func(w http.ResponseWriter, r *http.Request) {
		serveMedia(w, r, res)
	})

	listenAddr, serverURL, err := listenAndDisplay(Address, Port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	server := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server started at %s (listening on %s)\n", serverURL, listenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
			os.Exit(2)
		}
	}()

	if !NoOpen {
		fmt.Println("Opening browser...")
		if err := openBrowser(serverURL); err != nil {
			fmt.Fprintf(os.Stderr, "Error opening browser: %v\n", err)
		}
	} else {
		fmt.Println("Browser launch skipped (--no-open)")
	}

	<-signalCh
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error shutting down server: %v\n", err)
	}
}

type ImageView struct {
	RelPath  string
	Status   string
	Summary  string
	LeftURL  string
	RightURL string
	DiffURL  string
	HasLeft  bool
	HasRight bool
	HasDiff  bool
}

type ExtraView struct {
	RelPath string
	Status  string
	Summary string
	Kind    string
	Unified string
}

type PageData struct {
	Title    string
	Diff     string
	HasText  bool
	HasAny   bool
	Stats    Stats
	Images   []ImageView
	PDFs     []ExtraView
	Binaries []ExtraView
}

func renderPage(w http.ResponseWriter, res *Result) {
	tmpl, err := template.New("html").Parse(static.HTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := pageData(res)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func pageData(res *Result) PageData {
	data := PageData{
		Title:  "Diffr - A web-based content difference analyzer",
		Stats:  res.Stats,
		HasAny: res.HasChanges(),
	}

	var textParts []string
	for _, f := range res.Files {
		switch f.Kind {
		case KindText:
			if f.Unified != "" {
				textParts = append(textParts, f.Unified)
			}
		case KindImage:
			data.Images = append(data.Images, ImageView{
				RelPath:  f.RelPath,
				Status:   f.Status.String(),
				Summary:  f.Summary,
				LeftURL:  mediaURL("left", f.RelPath),
				RightURL: mediaURL("right", f.RelPath),
				DiffURL:  mediaURL("diff", f.RelPath),
				HasLeft:  f.LeftPath != "",
				HasRight: f.RightPath != "",
				HasDiff:  len(f.DiffPNG) > 0,
			})
		case KindPDF:
			data.PDFs = append(data.PDFs, ExtraView{
				RelPath: f.RelPath,
				Status:  f.Status.String(),
				Summary: f.Summary,
				Kind:    f.Kind.String(),
				Unified: f.Unified,
			})
		case KindBinary:
			data.Binaries = append(data.Binaries, ExtraView{
				RelPath: f.RelPath,
				Status:  f.Status.String(),
				Summary: f.Summary,
				Kind:    f.Kind.String(),
			})
		}
	}
	data.Diff = strings.Join(textParts, "")
	data.HasText = strings.TrimSpace(data.Diff) != ""
	return data
}

func mediaURL(side, rel string) string {
	parts := strings.Split(rel, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return "/media/" + side + "/" + strings.Join(parts, "/")
}

func serveMedia(w http.ResponseWriter, r *http.Request, res *Result) {
	rest := strings.TrimPrefix(r.URL.Path, "/media/")
	side, rel, ok := strings.Cut(rest, "/")
	if !ok || rel == "" {
		http.NotFound(w, r)
		return
	}
	rel = path.Clean(rel)
	if rel == ".." || strings.HasPrefix(rel, "../") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	switch side {
	case "diff":
		for _, f := range res.Files {
			if f.RelPath == rel && len(f.DiffPNG) > 0 {
				w.Header().Set("Content-Type", "image/png")
				_, _ = w.Write(f.DiffPNG)
				return
			}
		}
		http.NotFound(w, r)
	case "left":
		serveSideFile(w, r, res.LeftAbs, res.LeftDir, rel)
	case "right":
		serveSideFile(w, r, res.RightAbs, res.RightDir, rel)
	default:
		http.NotFound(w, r)
	}
}

func serveSideFile(w http.ResponseWriter, r *http.Request, root string, isDir bool, rel string) {
	target, err := safeJoin(root, rel, isDir)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, target)
}

func safeJoin(root, rel string, isDir bool) (string, error) {
	if !isDir {
		return root, nil
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	joined := filepath.Join(rootAbs, filepath.FromSlash(rel))
	abs, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	sep := string(os.PathSeparator)
	if abs != rootAbs && !strings.HasPrefix(abs, rootAbs+sep) {
		return "", fmt.Errorf("invalid path")
	}
	return abs, nil
}
