package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"os"
	"strconv"

	"net/http"
)

type Todos struct {
	toRequest  []TodoItem
	requested  []TodoItem
	onDelivery []TodoItem
	delivered  []TodoItem
}

type TodoItem struct {
	Id      int
	State   string
	Content string
}

var (
	dbUser     = os.Getenv("DB_USER_TODO_APP")     // todo
	dbPassword = os.Getenv("DB_PASSWORD_TODO_APP") // todoPass
	dbName     = os.Getenv("DB_NAME_TODO_APP")     // todoDb
	todoItems  []TodoItem
	todos      map[string][]TodoItem
	states     = [...]string{"toRequest", "requested", "onDelivery", "delivered"}
	id         int
)

func main() {
	serveFile := func(w http.ResponseWriter, r *http.Request) {
		path := "../todo-app" + r.URL.Path
		fmt.Println(path)
		http.ServeFile(w, r, path)
	}
	// util.InitDB(dbUser + ":" + dbPassword +
	// "@tcp(localhost:3306)/" + dbName)
	fillDummieData()
	organizeTodos()
	fmt.Println(todos)
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/scripts/", serveFile)
	http.HandleFunc("/styles/", serveFile)
	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		id++
		w.Write([]byte(strconv.Itoa(id)))
	})
	port := ":8080"
	fmt.Printf("Listen on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func fillDummieData() {
	for i := 0; i < 10; i++ {

		todoItems = append(todoItems, TodoItem{
			Id:      i,
			State:   states[rand.Intn(len(states))],
			Content: "blabla content" + strconv.Itoa(i),
		})
		id = i
	}
}

func organizeTodos() {
	if todos == nil {
		todos = make(map[string][]TodoItem)
	}
	for _, item := range todoItems {
		fmt.Println(item)
		todos[item.State] = append(todos[item.State], item)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	check(err)
	err = t.Execute(w, todos)
	check(err)
}
