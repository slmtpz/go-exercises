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
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.String()
		if toUrl, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, toUrl, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
	return fn
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
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)

	return MapHandler(pathMap, fallback), nil
}

type YamlInput []struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYAML(yamlInput []byte) (YamlInput, error) {
	var yi YamlInput
	err := yaml.Unmarshal(yamlInput, &yi)
	return yi, err
}

func buildMap(yi YamlInput) map[string]string {
	pathsToUrls := make(map[string]string)

	for _, pathUrl := range yi {
		pathsToUrls[pathUrl.Path] = pathUrl.Url
	}

	return pathsToUrls
}

