package err_handle

import (
	"fmt"
	"log"
	"net/http"
)

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func HandleCustomErr(err error, message string) {
	if err != nil {
		if message != "" {
			fmt.Println(message)
		} else {
			fmt.Println(err)
		}
	}
}

func HandleResponseErr(err error, w http.ResponseWriter, status int) {
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
}
