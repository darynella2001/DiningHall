package main

import (
	"net/http"
	"io"
)
type Food struct {
	Id               int              `json:"id"`
	Name             string           `json:"name"`
	PreparationTime  int              `json:"preparation-time"`
	Complexity       int              `json:"complexity"`
	CookingApparatus string           `json:"cooking-apparatus"`
}
type Order struct {
	Id       int    `json:"id"`
	Items    []int  `json:"items"`
	Priority int    `json:"priority"`
	MaxWait  int    `json:"maxWait"`
}
func main() {
	http.HandleFunc("/", servePage)
	http.ListenAndServe(":8080", nil)
}

func servePage(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "Hello world!")
}