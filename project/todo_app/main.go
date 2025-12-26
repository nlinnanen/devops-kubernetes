package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	templates *template.Template

	imageURL  = os.Getenv("IMAGE_URL")  // required
	imagePath = os.Getenv("IMAGE_PATH") // default: ./static/image.jpg

	mu           sync.Mutex
	lastDownload time.Time
)

func downloadImage(url, path string) error {
	tmp := path + ".tmp"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &httpError{status: resp.StatusCode}
	}

	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(f, resp.Body)
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.Remove(tmp)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return closeErr
	}

	return os.Rename(tmp, path)
}

type httpError struct{ status int }

func (e *httpError) Error() string { return http.StatusText(e.status) }

func ensureFreshImage() {
	if imageURL == "" {
		return
	}
	if imagePath == "" {
		imagePath = "./static/image.jpg"
	}

	now := time.Now()

	// Quick check under lock: do we need to download?
	mu.Lock()
	need := lastDownload.IsZero() || now.Sub(lastDownload) >= 10*time.Minute
	mu.Unlock()

	if !need {
		return
	}

	// Do the download outside the lock so requests aren't blocked.
	if err := downloadImage(imageURL, imagePath); err != nil {
		log.Printf("image download failed: %v", err)
		return
	}

	// Only update timestamp on success.
	mu.Lock()
	lastDownload = now
	mu.Unlock()

	log.Printf("image downloaded to %s at %s", imagePath, now.Format(time.RFC3339))
}

func handler(w http.ResponseWriter, r *http.Request) {
	go ensureFreshImage()

	err := templates.ExecuteTemplate(w, "index.html", map[string]any{
		"ImagePath": "/static/image.jpg",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", handler)

	log.Println("Server started on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
