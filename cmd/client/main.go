package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	res, err := http.Get("http://localhost:8080/file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	file, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(file))

	res, err = http.Get("http://localhost:8080/download")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	file, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("file_copy.txt", file, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
