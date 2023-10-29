package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"urlshort"
)

func main() {
	filename := flag.String("file", "../paths.yml", "a yaml file that serves as a URL store")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	content, ext := openFile(*filename)
	var handler http.Handler
	var err error

	switch ext {
	case ".yml":
		handler, err = urlshort.YAMLHandler(content, mapHandler)
		if err != nil {
			panic(err)
		}
	case ".json":
		handler, err = urlshort.JSONHandler(content, mapHandler)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func openFile(name string) ([]byte, string) {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return content, filepath.Ext(name)
}
