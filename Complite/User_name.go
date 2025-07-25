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

func HandleUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("template/Username.html")
		if err != nil {
			log.Fatal("Ошибка прогрузки шаблона")
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		emailCookie, err := r.Cookie("reg_email")
		if err != nil {
			http.Error(w, "Email cookie not found", http.StatusBadRequest)
			return
		}
		passwordCookie, err := r.Cookie("reg_password")
		if err != nil {
			http.Error(w, "Password cookie not found", http.StatusBadRequest)
			return
		}

		name := r.PostFormValue("username")
		region := r.PostFormValue("region")

		newUser := UserRegister{
			Email:    emailCookie.Value,
			Password: passwordCookie.Value,
			Name:     name,
			Region:   region,
		}

		result := db.Create(&newUser)
		if result.Error != nil {
			log.Printf("Ошибка при сохранении пользователя: %v", result.Error)
			http.Error(w, "Ошибка регистрации: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Основные куки
		http.SetCookie(w, &http.Cookie{
			Name:   "user_email",
			Value:  newUser.Email,
			Path:   "/",
			MaxAge: 86400,
		})
		http.SetCookie(w, &http.Cookie{
			Name:   "user_name",
			Value:  newUser.Name,
			Path:   "/",
			MaxAge: 86400,
		})

		// Удаляем временные куки
		http.SetCookie(w, &http.Cookie{Name: "reg_email", Value: "", Path: "/", MaxAge: -1})
		http.SetCookie(w, &http.Cookie{Name: "reg_password", Value: "", Path: "/", MaxAge: -1})

		http.Redirect(w, r, "/chat", http.StatusSeeOther)
	}
}
