package myapp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	Id        int
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email`
	CreatedAt time.Time `json:"created_at"`
}
type handlerA struct{}

var userMap map[int]*User
var lastID int

func handlerB(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, err)
		return
	}
	user, ok := userMap[id]

	if !ok {
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, "No user ID: ", id)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(res, string(data))

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

func createHandlerB(res http.ResponseWriter, req *http.Request) {
	user := new(User)
	err := json.NewDecoder(req.Body).Decode(user)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, err)
		return
	}

	//Created User
	lastID++
	user.Id = lastID
	user.CreatedAt = time.Now()
	userMap[user.Id] = user

	res.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(res, string(data))

}

// Make a new handler
func NewHttpHander() http.Handler {

	userMap = make(map[int]*User)
	lastID = 0

	mux := mux.NewRouter()
	//using handle with handlerFunc
	mux.HandleFunc("/", handler)
	//using handle with struct
	mux.Handle("/handleA", &handlerA{})
	//using handlefunc
	mux.HandleFunc("/handleB/{id:[0-9]+}", handlerB)
	mux.HandleFunc("/handleB", createHandlerB).Methods("POST")

	mux.Handle("/public/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/upload", uploadHandler)

	return mux
}
