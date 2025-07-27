package main

import (
	"Angular/internal/Registration"
	"Angular/internal/chat"
	"fmt"
	"net/http"
)

func main() {
	st := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", st))

	http.HandleFunc("/", Registration.HandleRegistr)

	http.HandleFunc("/usname", Registration.HandleUsername)

	http.HandleFunc("/chat", chat.HandleChat)

	fmt.Println("       /\\\n      /  \\\n     / /\\ \\\n    / /  \\ \\\n   / /    \\ \\\n  / /------\\ \\\n /_/        \\_\\   ")
	fmt.Println("    ANGLURE")
	fmt.Println("http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

//git add .
//git commit -m "      "
//git push --force
//go run cmd/main.go --files-dir=./Complite тут пока не точно
