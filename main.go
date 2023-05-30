package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	// interfaceMod "github.com/gooddavvy/todo-web-app-golang/interface"
	"io/ioutil"
)

var (
	port = "1055"
	list []TodoItem
)

type TodoItem struct {
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	DueDate   string `json:"due_date"`
	Completed string `json:"completed"`
	ID        int    `json:"id"`
}

type CtrlMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func getJson() {
	// Read JSON file
	data, err := ioutil.ReadFile("todos.json")
	if err != nil {
		panic(err)
	}

	// Unmarshal JSON data into slice of struct
	err = json.Unmarshal(data, &list)
	if err != nil {
		panic(err)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Println(err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func TodoList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func TodoListCtrl(w http.ResponseWriter, r *http.Request) {
	/* queries := r.URL.Query()

	title := queries.Get("title")
	desc := queries.Get("desc")
	dueDate := queries.Get("due-date")
	completed := queries.Get("completed")
	role := queries.Get("role")

	if role == "add" {
		list = append(list, TodoItem{
			Title:     title,
			Desc:      desc,
			DueDate:   dueDate,
			Completed: completed,
			ID:        len(list) + 1,
		})

		updatedData, err := json.Marshal(list)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile("todos.json", updatedData, 0644)
		if err != nil {
			panic(err)
		}
	} else if role == "remove" {
		i := 0
		for i < len(list) {
			item := TodoItem{Title: title, Desc: desc, DueDate: dueDate, Completed: completed}
			if list[i] == item {
				list = append(list[:i], list[i+1:]...)
			}
			i++
		}
	}

	message := ""
	typeOfMessage := ""
	if role == "add" {
		message = "Successfully added todo item to the list!"
		typeOfMessage = "success"
	} else if role == "remove" {
		message = "Successfully removed todo item from the list!"
		typeOfMessage = "success"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CtrlMessage{Message: message, Type: typeOfMessage}) */
}

func main() {
	getJson()

	http.HandleFunc("/", Home)
	http.HandleFunc("/api/todoList", TodoList)
	http.HandleFunc("/api/todoListCtrl", TodoListCtrl)
	// http.HandleFunc("/api/todoListEdit", todoListEdit)

	fmt.Println("Server started on port " + port)
	http.ListenAndServe(":"+port, nil)
}
