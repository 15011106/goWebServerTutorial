package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email`
	CreatedAt time.Time `json:"created_at"`
}
type handlerA struct{}

func handlerB(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Welcome handleB")
}

func (hand *handlerA) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	user := new(User)
	err := json.NewDecoder(req.Body).Decode(user)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Bad request: ", err)
		return
	}
	user.CreatedAt = time.Now()

	data, _ := json.Marshal(user)
	res.Header().Add("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, string(data))

}

func handler(res http.ResponseWriter, req *http.Request) {

	name := req.URL.Query().Get("name")

	if name == "" {
		name = "world"
	}

	fmt.Fprintf(res, "Hello %s!", name)
}

func NewHttpHander() http.Handler {

	mux := http.NewServeMux()
	//using handle with handlerFunc
	mux.HandleFunc("/", handler)
	//using handle with struct
	mux.Handle("/handleA", &handlerA{})
	//using handlefunc
	mux.HandleFunc("/handleB", handlerB)
	return mux
}
