package Complite

import (
	"html/template"
	"log"
	"net/http"
)

type UserRegister struct {
	Email    string
	Password string
	Name     string
	Region   string
	ID       uint
}

var email string
var password string

func HandleRegistr(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("template/Registr.html")
		if err != nil {
			log.Fatal("Ошибка прогрузки шаблона ")
		}
		tmpl.Execute(w, nil)
	}

	if r.Method == "POST" {
		email = r.PostFormValue("useremail")
		password = r.PostFormValue("password")

		http.Redirect(w, r, "/usname", http.StatusSeeOther)
		return
	}
}
