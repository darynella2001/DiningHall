package main

import (
	"net/http"
	"io"
)

func main() {
	http.HandleFunc("/", servePage)
	http.ListenAndServe(":8080", nil)
}

func servePage(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "Hello world!")
}