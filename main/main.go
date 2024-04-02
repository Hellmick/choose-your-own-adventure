package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

type Option struct {
	Text string `json:text`
	Arc  string `json:arc`
}

type StoryPage struct {
	Title     string   `json:"title"`
	StoryText []string `json:"story"`
	Options   []Option `json:"options"`
}

func readJSON(filename string) ([]byte, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	doc, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func parseJSON(jsonBytes []byte) (map[string]StoryPage, error) {

	story := &map[string]StoryPage{}

	if err := json.Unmarshal(jsonBytes, &story); err != nil {
		return nil, err
	}

	return *story, nil
}

func createMux(story map[string]StoryPage) (http.Handler, error) {

	mux := http.NewServeMux()

	template, err := template.ParseFiles("template.html")
	if err != nil {
		return nil, err
	}

	for page, storyArc := range story {
		page := page
		storyArc := storyArc
		mux.HandleFunc("/"+page, func(w http.ResponseWriter, r *http.Request) {
			err := template.Execute(w, storyArc)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // Set the status code
				http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
				return // Return to prevent further writes to the response writer
			}
		})
	}

	return mux, nil
}

func main() {

	filename := flag.String("f", "story.json", "specify location of the story file")
	port := flag.String("p", "8080", "specify port to serve")
	flag.Parse()

	doc, err := readJSON(*filename)
	if err != nil {
		panic(err)
	}

	story, err := parseJSON(doc)
	if err != nil {
		panic(err)
	}

	handler, err := createMux(story)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Serving on port %s", *port)
	http.ListenAndServe(":"+*port, handler)

}
