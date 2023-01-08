package main

import (
	"io"
	"net/http"
)

type handlerA struct{}

func (hand *handlerA) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Welcome handleA")
}

func handler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Welcome")
}

func main() {

	//using handle with handlerFunc
	http.Handle("/", http.HandlerFunc(handler))
	//using handle with struct
	http.Handle("/handleA", &handlerA{})
	//using handlefunc
	http.HandleFunc("/handleB", func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, "Welcome handleB")
	})

	http.ListenAndServe(":8080", nil)
}
