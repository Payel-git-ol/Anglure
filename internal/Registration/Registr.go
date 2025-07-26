package Registration

import (
	"html/template"
	"log"
	"net/http"
)

func HandleRegistr(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("web/templates/RegistrationTemplates/Registr.html")
		if err != nil {
			log.Fatal("Ошибка прогрузки шаблона Registr")
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		email := r.PostFormValue("useremail")
		password := r.PostFormValue("password")

		// Сохраняем email и password во временные куки (на 10 минут, например)
		http.SetCookie(w, &http.Cookie{
			Name:   "reg_email",
			Value:  email,
			Path:   "/",
			MaxAge: 600,
		})

		http.SetCookie(w, &http.Cookie{
			Name:   "reg_password",
			Value:  password,
			Path:   "/",
			MaxAge: 600,
		})

		http.Redirect(w, r, "/usname", http.StatusSeeOther)
	}
}
