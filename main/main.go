package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"urlshort"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "urlstore"
)

func main() {
	useDb := flag.Bool("database", false, "a database that serves as a URL store")
	filename := flag.String("file", "../paths.yml", "a yaml file that serves as a URL store")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var handler http.Handler
	var err error

	if *useDb {
		content := getURLsFromDB()
		handler, err = urlshort.DBHandler(content, mapHandler)
		if err != nil {
			panic(err)
		}
	} else {
		content, ext := openFile(*filename)

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

func getURLsFromDB() *sql.Rows {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	rows, err := db.Query("SELECT * FROM urlstore")
	if err != nil {
		panic(err)
	}
	return rows
}
