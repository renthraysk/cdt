package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/renthraysk/cdt"
)

func main() {

	u := url.URL{
		Scheme: "https",
		Host:   "localhost:8080",
		Path:   "/demo.txt",
	}

	compendium, err := cdt.NewCompendium(3, u.Path, "")
	if err != nil {
		log.Fatalf("NewCompendium failed: %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle(u.Path, cdt.NewSelfPack(compendium))

	fmt.Printf("Use: curl -X PUT --data-binary @<file> %s\n", u.String())
	s := http.Server{
		Addr:    u.Host,
		Handler: mux,
	}

	if err := s.ListenAndServeTLS(appRel("/keys/localhost.pem"), appRel("/keys/localhost-key.pem")); err != nil {
		fmt.Println("Server error:", err)
	}
	s.Shutdown(context.Background())
}

func appRel(name string) string {
	path, _ := os.Executable()
	return filepath.Join(filepath.Dir(path), name)
}
