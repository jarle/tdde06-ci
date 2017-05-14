package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	databaseUser = "postgres"
	databaseHost = "localhost"
	databaseName = "tdde06aaaa"
)

type Task struct {
	Id     int
	Name   string
	Done   bool
	ListId int
}

type CreateTaskRequest struct {
	Name   string `json:"name"`
	ListId int    `json:"list_id"`
}

type List struct {
	Id   int
	Name string
}

type ListCreateRequest struct {
	Name string
}

type ListCreateResponse struct {
	Id int
}

func CheckFatal(err error, w http.ResponseWriter) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type Database struct {
	Db *sql.DB
}

// getURLParameter takes out the first URL parameter from the path
// path should be formated as /type/param.
// It returns a parameter representing a string.
func getURLParameter(path string) *string {
	// TODO: Handle multiple parameters
	param := strings.Split(path, "/")
	if len(param) == 3 {
		return &param[2]
	} else {
		return nil
	}
}

// getLists retrieves all the lists from the database.
// It returns a slice of List structs.
func getLists(db *sql.DB, w http.ResponseWriter) []List {
	rows, err := db.Query("select * from list")
	CheckFatal(err, w)
	// Retrieve all lists from the query
	res := make([]List, 0)
	for rows.Next() {
		list := List{}
		err := rows.Scan(&list.Id, &list.Name)
		CheckFatal(err, w)
		res = append(res, list)
	}

	return res
}

// getTass retrieves all the tasks from the database.
// It returns a slice of Task structs.
func getTasks(db *sql.DB, listId int, w http.ResponseWriter) []Task {
	// Query the database for all tasks that references the specified list
	rows, err := db.Query("select * from task where list=$1", listId)
	CheckFatal(err, w)

	// Retrieve all tasks from the query
	res := make([]Task, 0)
	for rows.Next() {
		var name string
		var id, list int
		var done bool
		err := rows.Scan(&id, &name, &done, &list)
		CheckFatal(err, w)
		res = append(res, Task{Id: id, Name: name, Done: done, ListId: list})
	}

	return res
}

// insertList adds a list to the database with listName as its name.
// It returns the Id of the list.
func insertList(db *sql.DB, listName string, w http.ResponseWriter) int {
	var listId int
	err := db.QueryRow("insert into list (name) values ($1) returning id", listName).Scan(&listId)
	CheckFatal(err, w)

	return listId
}

// insertTask adds a task to the database.
// taskName specifies the name of the task, and listId the list that it belongs to.
func insertTask(db *sql.DB, taskName string, listId int, w http.ResponseWriter) {
	_, err := db.Exec("insert into task (name, list) values ($1, $2)", taskName, listId)
	// Handle non-existing list id
	CheckFatal(err, w)
}

// listHandler manages requests with regards to the lists.
// A GET request to /list will retrieve all the lists.
// A GET request to /list/<id> will retrieve all the tasks of the list with id <id>.
func (db *Database) listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Handle GET Request
		param := getURLParameter(r.URL.Path)

		// If no parameter exists, retrieve all lists
		if param == nil || *param == "" {
			// Retrieve lists
			list := getLists(db.Db, w)
			json.NewEncoder(w).Encode(&list)
		} else {
			// Get the list id from the parameter
			listId, err := strconv.Atoi(*param)
			CheckFatal(err, w)

			// Retrieve tasks and send them back
			tasks := getTasks(db.Db, listId, w)
			json.NewEncoder(w).Encode(&tasks)
		}
	} else if r.Method == "POST" {
		// Parse the request and create a new list
		body, err := ioutil.ReadAll(r.Body)
		CheckFatal(err, w)
		listRequest := ListCreateRequest{}
		err = json.Unmarshal(body, &listRequest)
		CheckFatal(err, w)

		listResponse := ListCreateResponse{}
		listResponse.Id = insertList(db.Db, listRequest.Name, w)

		json.NewEncoder(w).Encode(&listResponse)
	}
}

// taskHandler manages requests with regards to the tasks.
// A POST request to /task will create a new task with the name and list provided
// in the Post Body. The Body should be in the format
// {"name": "taskName", "list_id": 123}
func (db *Database) taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		CheckFatal(err, w)
		taskRequest := CreateTaskRequest{}
		err = json.Unmarshal(body, &taskRequest)
		CheckFatal(err, w)

		insertTask(db.Db, taskRequest.Name, taskRequest.ListId, w)

		fmt.Fprintf(w, "OK")
	}
}

// ConnextDB connects to a postgres database.
// it returns a database handle.
func ConnectDb() *sql.DB {
	// TODO: Refactor the database config
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable", databaseUser, databaseHost, databaseName))
	if err != nil {
		log.Fatal(err)
	}
	if db == nil {
		log.Fatal("Woops, this db is nil, not good for you!")
	}
	return db
}

// Handlers retrieves all handlers for the server.
func Handlers() *http.ServeMux {
	db := Database{Db: ConnectDb()}
	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(db.listHandler))
	mux.Handle("/list/", http.HandlerFunc(db.listHandler))
	mux.Handle("/task", http.HandlerFunc(db.taskHandler))
	return mux
}

func main() {
	// Listen on port 5050
	err := http.ListenAndServe(":8080", Handlers())
	if err != nil {
		log.Fatal(err)
	}
}
