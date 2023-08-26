package diffr

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	dir1 string
	dir2 string
)

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func RunWebServer(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Errorf("Error: Usage: \n diffr /path/to/dir1 /path/to/dir2")
		return
	}

	dir1 = args[0]
	dir2 = args[1]

	serverURL := fmt.Sprintf("%s:%d", Address, Port)

	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := &http.Server{Addr: fmt.Sprintf(":%d", Port)}

	// Channel to receive signals (e.g., interrupt or termination)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server started at %s\n", serverURL)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
			os.Exit(1)
		}
	}()

	fmt.Println("Opening browser...")
	err := openBrowser(serverURL)
	if err != nil {
		fmt.Println("Error opening browser:", err)
	}

	// Wait for a termination signal
	<-signalCh

	fmt.Println("Shutting down server...")
	err = server.Shutdown(nil)
	if err != nil {
		fmt.Println("Error shutting down server:", err)
	}
}

type PageData struct {
	Title string
	Diff  string
}

func handler(w http.ResponseWriter, r *http.Request) {
	var (
		wg        sync.WaitGroup
		finalStr  = ""
		diffChan  = make(chan string)
		errorChan = make(chan error)
	)

	go func() {
		for diff := range diffChan {
			finalStr += diff
		}
	}()

	go func() {
		for err := range errorChan {
			fmt.Printf("error: %v", err)
		}
	}()

	wg.Add(1)
	go CompareDirectories(dir1, dir2, diffChan, errorChan, &wg)
	wg.Wait()

	close(diffChan)
	close(errorChan)

	tmpl, err := template.ParseFiles("static/templates/template.html")
	if err != nil {
		fmt.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title: "Diffr - A web-based content difference analyzer",
		Diff:  finalStr,
	}

	// Execute the template with the data and write the output to the response writer
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
