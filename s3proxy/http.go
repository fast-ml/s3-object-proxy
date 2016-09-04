package s3proxy

import ( 
	"log"
	"net/http"
	"strings"
)

// ImageProxy HTTP handling
const dir = "/proxy/"

func (i *ImageProxy) handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	if !strings.HasPrefix(r.URL.Path, dir) {
		http.NotFound(w, r)
		return
	}
	path := strings.TrimLeft(r.URL.Path, dir)
	resp, err := i.GetObject(path)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(resp)
}

// ListenAndServe wraps http.ListenAndServe so you can test the connection and logic of your permission lookup easily.
func (i *ImageProxy) ListenAndServe(port string) error {
	http.HandleFunc("/", i.handleRequest)
	log.Printf("Listening at %s%s\nBucket: %s\nRegion: %s", port, dir, i.bucket, *i.s3.Config.Region)
	return http.ListenAndServe(port, nil)
}
