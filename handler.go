package main

import (
	"goWebServerTutorial/myapp"
	"net/http"
)

func main() {

	http.ListenAndServe(":8080", myapp.NewHttpHander())
}
