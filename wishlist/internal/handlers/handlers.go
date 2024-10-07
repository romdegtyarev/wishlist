package handlers

import (
    "log"
    "net/http"
    "time"
    "html/template"
    "github.com/gorilla/sessions"
    "golang.org/x/crypto/bcrypt"

    "wishlist/internal/config"
    "wishlist/internal/db"
)

type Session struct {
    UserID    int
    Username  string
    CreatedAt time.Time
    ExpiresAt time.Time
}

var store *sessions.CookieStore

func newSession(userID int, username string) *Session {
    return &Session{
        UserID:    userID,
        Username:  username,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
}




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
    err = tmpl.Execute(w, user)
    if err != nil {
        log.Println("Error executing template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")

        user, err := db.GetUser(username)
        if err != nil || user.PasswordHash == "" {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            return
        }
        err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
        if err != nil {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            return
        }

        session, err := store.Get(r, "usersession")
        if err != nil {
            log.Println("Error getting session:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        userSession := newSession(user.ID, user.Username)
        session.Values["user"] = userSession
        err = session.Save(r, w)
        if err != nil {
            log.Println("Error saving session:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, config.IndexPath, http.StatusFound)
        return
    }

    tmpl := template.Must(template.ParseFiles(config.LoginTemplatePath))
    err := tmpl.Execute(w, nil)
    if err != nil {
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
    err = session.Save(r, w);
    if err != nil {
        log.Println("Error saving session:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, config.LoginPath, http.StatusFound)
}

