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

type StoryArc struct {
	Page struct {
		Title    string   `json:"title"`
		StoryArc []string `json:"story"`
		Options  []struct {
			Text string `json:text`
			Arc  string `json:arc`
		} `json:"options"`
	}
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

func parseJSON(jsonBytes []byte) (map[string]StoryArc, error) {

	story := &map[string]StoryArc{}

	if err := json.Unmarshal(jsonBytes, &story); err != nil {
		return nil, err
	}

	return *story, nil
}

func createMux(story map[string]StoryArc) (http.Handler, error) {

	mux := http.NewServeMux()

	template, err := template.ParseFiles("template.html")
	if err != nil {
		return nil, err
	}

	for page, storyArc := range story {
		http.HandleFunc("/"+page, func(w http.ResponseWriter, r *http.Request) {
			template.Execute(w, storyArc)
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
