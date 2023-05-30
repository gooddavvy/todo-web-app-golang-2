package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	// interfaceMod "github.com/gooddavvy/todo-web-app-golang-2/interface"
)

var (
	port        = "1055"
	list        []TodoItem
	currentRole Role
)

type Role string
type TodoListCtrlType = func([]TodoItem, TodoItem) (string, string)

const (
	Add     Role = "add"
	Replace Role = "replace"
	Delete  Role = "delete"
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

func addFn(todoItems []TodoItem, todoItem TodoItem) (string, string) {
	currentRole = "add"
	todoItems = append(todoItems, todoItem)

	updatedData, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the updated JSON data back to the file
	if err := ioutil.WriteFile("todos.json", updatedData, 0644); err != nil {
		log.Fatal(err)
	}

	return "Successfully added!", ""
}

func deleteFn(todoItems []TodoItem, todoItem TodoItem) (string, string) {
	currentRole = "delete"
	i := 0
	for i <= len(todoItems) {
		if todoItems[i] == todoItem {
			todoItems = append(todoItems[:i], todoItems[(i+1):]...)
			break
		} else {
			i++
		}
	}

	updatedData, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the updated JSON data back to the file
	if err := ioutil.WriteFile("todos.json", updatedData, 0644); err != nil {
		log.Fatal(err)
	}

	return "Successfully deleted.", ""
}

func replaceFn(todoItems []TodoItem, todoItem TodoItem, newTodoItem TodoItem) (string, string) {
	currentRole = "replace"
	i := 0
	for i <= len(todoItems) {
		if todoItems[i] == todoItem {
			todoItems[i] = newTodoItem
			break
		} else {
			i++
		}
	}

	updatedData, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the updated JSON data back to the file
	if err := ioutil.WriteFile("todos.json", updatedData, 0644); err != nil {
		log.Fatal(err)
	}

	return "Successfully replaced.", ""
}

func TodoListCtrls() (
	TodoListCtrlType,
	TodoListCtrlType,
	func([]TodoItem, TodoItem, TodoItem) (string, string),
) {
	return addFn, deleteFn, replaceFn
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
	queries := r.URL.Query()

	title := queries.Get("title")
	desc := queries.Get("desc")
	dueDate := queries.Get("due-date")
	completed := queries.Get("completed")
	role := queries.Get("role")

	add, delete, _ := TodoListCtrls()
	newTodoItem := TodoItem{
		Title:     title,
		Desc:      desc,
		DueDate:   dueDate,
		Completed: completed,
	}
	success := true
	Error := ""

	if role == "add" {
		_, err := add(list, newTodoItem)
		if err != "" {
			success = false
			Error = err
		}
	} else if role == "remove" {
		_, err := delete(list, newTodoItem)
		if err != "" {
			success = false
			Error = err
		}
	}

	message := ""
	typeOfMessage := ""

	if role == "add" && success {
		message = "Successfully added todo item to the list!"
		typeOfMessage = "success"
	} else if role == "remove" && success {
		message = "Successfully removed todo item from the list!"
		typeOfMessage = "success"
	}

	if role == "add" && !success {
		message = "Error: " + Error
		typeOfMessage = "unsuccess"
	} else if role == "remove" && !success {
		message = "Error: " + Error
		typeOfMessage = "unsuccess"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CtrlMessage{Message: message, Type: typeOfMessage})
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
