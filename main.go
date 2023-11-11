// Listing 3.1
package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func getData(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		log.Printf("Unable to retrieve data from %v: %v", url, err)
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	welcomePath := filepath.Join("templates", "welcome.gohtml")
	welcomeTemplate, err := template.ParseFiles(welcomePath)
	if err != nil {
		log.Printf("Unable to parse welcome template file: %v", err)
		http.Error(w, "Unable to parse welcome template file", http.StatusInternalServerError)
		return
	}
	err = welcomeTemplate.Execute(w, nil)
	if err != nil {
		log.Printf("Unable to execute welcome template: %v", err)
		http.Error(w, "Unable to execute welcome template", http.StatusInternalServerError)
		return
	}

	url := r.FormValue("url")
	fmt.Printf("\nValue of url entered by user: %v", url)

	body, err := getData(url)
	if err != nil {
		log.Printf("Unable to get body of requested web page: %v", err)
		http.Error(w, "Unable get body of requested web page", http.StatusInternalServerError)
		return
	}

	outputPath := filepath.Join("webpages", "web-page-body.html")
	file, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Unable to create file for web page data: %v", err)
		http.Error(w, "Unable to create file for web page data", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	bytes, err := file.Write(body)
	if err != nil {
		log.Printf("Unable to write to html file: %v", err)
		http.Error(w, "Unable to write to html file", http.StatusInternalServerError)
		return
	}
	fmt.Printf("\nWrote %d bytes to file\n", bytes)
	return
}

func main() {
	http.HandleFunc("/", pageHandler)
	fmt.Println("Starting the web server on localhost port 3000, localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Printf("Unable to start web server or server crash: %v", err)
		return
	}
}
