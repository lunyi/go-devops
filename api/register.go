package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"tgs-devops/utils"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	uName := strings.Trim(r.FormValue("username"), " ")
	email := strings.Trim(r.FormValue("emailsignup"), " ")
	pwd := strings.Trim(r.FormValue("passwordsignup"), " ")
	confirmPwd := strings.Trim(r.FormValue("passwordconfirm"), " ")

	_uName, _email, _pwd, _confirmPwd := false, false, false, false
	_uName = !utils.IsEmpty(uName)
	_email = !utils.IsEmpty(email)
	_pwd = !utils.IsEmpty(pwd)
	_confirmPwd = !utils.IsEmpty(confirmPwd)

	if _uName && _email && _pwd && _confirmPwd {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 8)
		password := string(hashedPassword)

		stmt, err := utils.DB.Prepare("insert into dbo.users values ($1, $2, $3)")
		if err != nil {
			log.Println("db error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := stmt.Exec(uName, email, password)

		log.Println("res:", res)

		if err != nil {
			fmt.Println("db error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		SetCookie(uName, w)
		http.Redirect(w, r, "/index", 302)
	}
}
