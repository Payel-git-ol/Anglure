package Complite

import (
	"html/template"
	"log"
	"net/http"
)

var nameus string
var region string

func HandleUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("template/Username.html")
		if err != nil {
			log.Fatal("Ошибка прогрузки шаблона ")
		}
		tmpl.Execute(w, nil)
	}

	if r.Method == "POST" {

		nameus = r.PostFormValue("username")
		region = r.PostFormValue("region")

		newUser := UserRegister{
			Email:    email,
			Password: password,
			Name:     nameus,
			Region:   region,
		}

		result := db.Create(&newUser)
		if result.Error != nil {
			log.Printf("Ошибка при сохранении пользователя: %v", result.Error)
			http.Error(w, "Ошибка регистрации: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}
}
