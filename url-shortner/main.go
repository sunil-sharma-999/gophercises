package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type MapUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func main() {

	// initialize
	var urlMaps []MapUrl

	// parse flags
	filePath := flag.String("input", "", "YAML file")
	flag.Parse()

	// check flag values
	if *filePath == "" || !strings.HasSuffix(*filePath, ".yaml") {
		log.Fatal("Provide YAML file")
	}

	// open YAML file
	file, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	// unmarshal
	err = yaml.Unmarshal([]byte(file), &urlMaps)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	for _, urlMap := range urlMaps {
		mux.HandleFunc(urlMap.Path, func(w http.ResponseWriter, r *http.Request) {
			tmpl, err := template.ParseFiles("./assets/page.html")
			if err != nil {
				log.Fatal(err)
			}
			countCookie, err := r.Cookie("count")
			if errors.Is(err, http.ErrNoCookie) {
				http.SetCookie(w, &http.Cookie{
					Name:  "count",
					Value: "5",
				})
				tmpl.Execute(w, nil)
			} else {
				count, err := strconv.Atoi(countCookie.Value)
				if err != nil {
					http.Error(w, "cookie not found", http.StatusBadRequest)
				} else {
					if count <= 0 {
						http.SetCookie(w, &http.Cookie{
							Name:    "count",
							Value:   "",
							Expires: time.Unix(0, 0).Add(-1 * time.Hour),
						})
						w.Header().Add("HX-Redirect", urlMap.Url)
						w.WriteHeader(http.StatusOK)
						return
					} else {
						http.SetCookie(w, &http.Cookie{
							Name:  "count",
							Value: fmt.Sprint(count - 1),
						})
						tmpl.Execute(w, nil)
					}
				}
			}

		})
	}
	fmt.Println("Server started at: \nhttp://localhost:9000\nhttp://192.168.1.91:9000")
	err = http.ListenAndServe(":9000", mux)
	if err != nil {
		log.Fatal(err)
	}

}
