package main

import (
	"fmt"
	"net/http"
)

// func CheckFile(folder string) []string {

// 	files, err := ioutil.ReadDir(folder)

// }

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/file":
			fmt.Fprint(w, "file.txt")
		case "/download":
			http.ServeFile(w, r, "file.txt")
		default:
			http.NotFound(w, r)
		}
	})

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
