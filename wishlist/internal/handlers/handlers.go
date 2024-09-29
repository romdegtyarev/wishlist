package handlers

import (
    "log"
    "net/http"
    "html/template"
    "github.com/gorilla/sessions"
    "golang.org/x/crypto/bcrypt"

    "wishlist/internal/config"
    "wishlist/internal/db"
)

var store *sessions.CookieStore

func InitHandlers() {
    store = sessions.NewCookieStore([]byte(config.SessionSecretKey))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "usersession")
    if err != nil {
        log.Println("Error getting session:", err)
        http.Redirect(w, r, config.LoginPath, http.StatusFound)
        return
    }

    user, ok := session.Values["user"].(db.UsersTable)
    if !ok {
        http.Redirect(w, r, config.LoginPath, http.StatusFound)
        return
    }

    tmpl := template.Must(template.ParseFiles(config.IndexTemplatePath))
    if err := tmpl.Execute(w, user); err != nil {
        log.Println("Error executing template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")

        storedPassword, err := db.GetUserFromDB(username)
        if err != nil || storedPassword == "" {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            return
        }

        session, err := store.Get(r, "usersession")
        if err != nil {
            log.Println("Error getting session:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        user := db.UsersTable{Username: username}
        session.Values["user"] = user
        if err := session.Save(r, w); err != nil {
            log.Println("Error saving session:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, config.IndexPath, http.StatusFound)
        return
    }

    tmpl := template.Must(template.ParseFiles(config.LoginTemplatePath))
    if err := tmpl.Execute(w, nil); err != nil {
        log.Println("Error executing template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "usersession")
    if err != nil {
        log.Println("Error getting session:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    delete(session.Values, "user")
    if err := session.Save(r, w); err != nil {
        log.Println("Error saving session:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, config.LoginPath, http.StatusFound)
}

