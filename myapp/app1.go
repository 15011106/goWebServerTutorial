package myapp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	res.WriteHeader(http.StatusCreated)
	fmt.Fprintf(res, string(data))

}

func handler(res http.ResponseWriter, req *http.Request) {

	name := req.URL.Query().Get("name")

	if name == "" {
		name = "world"
	}

	fmt.Fprintf(res, "Hello %s!", name)
}

func uploadHandler(res http.ResponseWriter, req *http.Request) {
	uploadFile, header, err := req.FormFile("uploadFile")

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, err)
		return
	}

	defer uploadFile.Close()

	dirname := "./upload"
	os.Mkdir(dirname, 0777)

	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
	file, err := os.Create(filepath)

	defer file.Close()

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, err)
		return
	}

	io.Copy(file, uploadFile)
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, filepath)
}

func NewHttpHander() http.Handler {

	mux := http.NewServeMux()
	//using handle with handlerFunc
	mux.HandleFunc("/", handler)
	//using handle with struct
	mux.Handle("/handleA", &handlerA{})
	//using handlefunc
	mux.HandleFunc("/handleB", handlerB)

	mux.Handle("/public/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/upload", uploadHandler)

	return mux
}
