package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// Will be implicitly converted into an http.HandlerFunc because of the return type
	// declared above and since it matches the signature, https://golang.org/pkg/net/http/#HandlerFunc
	return func(w http.ResponseWriter, r *http.Request) {
		// If a path is matched, redirect to it
		path := r.URL.Path

		// If the path is mapped to a URL, redirect to it
		// otherwise use the fallback
		if url, keyFound := pathsToUrls[path]; keyFound {
			http.Redirect(w, r, url, http.StatusFound)
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// Parse the YAML into this slice
	var pairs []pathURLPair
	err := yaml.Unmarshal(yamlBytes, &pairs)
	if err != nil {
		return nil, err
	}

	// Slice of structs to a map so we can reuse the handler
	// defined above
	pathsToUrls := make(map[string]string)
	for _, pair := range pairs {
		pathsToUrls[pair.Path] = pair.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}

// A struct for parsing YAML to path/url pairs
type pathURLPair struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
