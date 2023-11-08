package urlshort

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

type Entry struct {
	Path string
	URL  string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := pathsToUrls[r.URL.Path]
		if url != "" {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

func DBHandler(rows *sql.Rows, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseROWS(rows)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

func parseROWS(rows *sql.Rows) ([]Entry, error) {
	var entries []Entry
	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.Path, &entry.URL)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func parseYAML(yml []byte) ([]Entry, error) {
	var entries []Entry
	err := yaml.Unmarshal(yml, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func parseJSON(js []byte) ([]Entry, error) {
	var entries []Entry
	err := json.Unmarshal(js, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func buildMap(entries []Entry) map[string]string {
	pathMap := make(map[string]string)
	for _, entry := range entries {
		pathMap[entry.Path] = entry.URL
	}
	return pathMap
}
