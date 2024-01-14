package main

import (
	"net/http"

	"nyiyui.ca/halation/web"
)

func main() {
	s := web.NewServer()
	http.ListenAndServe(":8080", s)
}
