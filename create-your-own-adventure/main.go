package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	err_handle "github.com/sunil-sharma-999/gophercises/err"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

type AdventureOption struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type AdventureValueJSON struct {
	Title   string            `json:"title"`
	Story   []string          `json:"story"`
	Options []AdventureOption `json:"options"`
}

type AdventureJSON map[string]AdventureValueJSON

func main() {
	storyPath := flag.String("f", "", "Provide story in JSON format")
	server := flag.Bool("s", false, "Play in Web")
	flag.Parse()

	if (*storyPath == "" || !strings.HasSuffix(*storyPath, ".json")) && !*server {
		log.Fatal("Provide story in JSON format")
	}

	file, err := os.ReadFile(*storyPath)
	err_handle.HandleError(err)

	dataJSON := AdventureJSON{}

	err = json.Unmarshal(file, &dataJSON)
	err_handle.HandleError(err)
	if !*server {
		cliGame(dataJSON)
	} else {
		webGame(dataJSON)
	}

}

func cliGame(dataJSON AdventureJSON) {
	key := "intro"
	for key != "" {
	scan:
		CallClear()
		options := dataJSON[key].Options

		title := dataJSON[key].Title
		story := dataJSON[key].Story

		fmt.Printf("\n%v\n", title)
		fmt.Println()
		for _, line := range story {
			fmt.Printf("\t%v\n\n", line)
		}
		fmt.Println()

		if len(options) == 0 {
			break
		}

		fmt.Println("Select:")
		for i, option := range options {
			fmt.Printf("%d: %v\n", i+1, option.Text)
		}

		scannedValue := ""
		_, err := fmt.Scan(&scannedValue)
		if err != nil {
			goto scan
		}
		selectedOption, err := strconv.Atoi(strings.TrimSpace(scannedValue))
		if err != nil {
			goto scan
		}

		if selectedOption-1 >= 0 && selectedOption-1 < len(options) {
			key = strings.TrimSpace(options[selectedOption-1].Arc)
		} else {
			goto scan
		}

	}
}

func webGame(dataJSON AdventureJSON) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/arc/intro", http.StatusPermanentRedirect)
	})

	mux.HandleFunc("/arc/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(""))
		}

		arcName := strings.TrimPrefix(r.URL.Path, "/arc/")
		value, ok := dataJSON[arcName]

		if !ok {
			http.NotFound(w, r)
		} else {
			_, err := json.Marshal(value)
			err_handle.HandleResponseErr(err, w, http.StatusInternalServerError)

			tpl, err := template.ParseFiles("./static/index.html")
			err_handle.HandleResponseErr(err, w, http.StatusInternalServerError)

			err = tpl.Execute(w, value)
			err_handle.HandleResponseErr(err, w, http.StatusInternalServerError)

		}
	})

	fmt.Println("Server started on: http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	err_handle.HandleError(err)

}
