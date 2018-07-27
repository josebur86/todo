package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var GlobalTasks = []Task{
	Task{0, "Test Task 1", time.Time{}, false},
	Task{1, "Test Task 2", time.Time{}, false},
	Task{2, "Test Task 3", time.Time{}, false},
	Task{3, "Test Task 4", time.Time{}, false},
	Task{4, "Test Task 5", time.Time{}, false},
	Task{5, "Test Task 6", time.Time{}, false},
	Task{6, "Test Task 7", time.Time{}, false},
}

func GetAllTodosHandler(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(GlobalTasks)
}

func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId, err := strconv.Atoi(vars["todoId"])
	if err != nil {
		log.Fatal("Unable to fetch todo")
	}
	task := GlobalTasks[taskId]
	json.NewEncoder(w).Encode(task)
}

func PutTodoHandler(w http.ResponseWriter, r *http.Request) {
	var updatedTask = Task{}
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		log.Fatal("Unable to parse request body: ", err)
	}

	vars := mux.Vars(r)

	taskId, err := strconv.Atoi(vars["todoId"])
	if err != nil {
		log.Fatal("Unable to fetch todo")
	}

	if taskId != updatedTask.Id {
		log.Fatal("Invalid task id does not match")
	}

	GlobalTasks[taskId] = updatedTask

	json.NewEncoder(w).Encode(updatedTask)
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTask = Task{}
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Fatal("Unable to parse request body: ", err)
	}

	id := len(GlobalTasks)
	GlobalTasks = append(GlobalTasks, newTask)

	GlobalTasks[id].Id = id

	json.NewEncoder(w).Encode(GlobalTasks[id])
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/api/todos/").Handler(http.HandlerFunc(GetAllTodosHandler))
	router.Methods("GET").Path("/api/todos/{todoId}/").Handler(http.HandlerFunc(GetTodoHandler))
	router.Methods("PUT").Path("/api/todos/{todoId}/").Handler(http.HandlerFunc(PutTodoHandler))
	router.Methods("POST").Path("/api/todos/").Handler(http.HandlerFunc(CreateTodoHandler))
	log.Fatal(http.ListenAndServe(":5000", router))
}
