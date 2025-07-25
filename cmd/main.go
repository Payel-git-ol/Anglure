package main

import (
	"Angular/Complite"
	"fmt"
	"net/http"
)

func main() {
	st := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", st))

	http.HandleFunc("/", Complite.HandleRegistr)

	http.HandleFunc("/usname", Complite.HandleUsername)

	http.HandleFunc("/chat", Complite.HandleChat)

	fmt.Println("       /\\\n      /  \\\n     / /\\ \\\n    / /  \\ \\\n   / /    \\ \\\n  / /------\\ \\\n /_/        \\_\\   ")
	fmt.Println("    ANGLURE")
	fmt.Println("http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

//git commit -m "      "
//git push --force
