package main

import (
    "log"
    "net/http"

    "wishlist/internal/config"
    "wishlist/internal/db"
    "wishlist/internal/handlers"
)

func main() {
    config.InitConfig()
    db.InitDB(config.DataSourceName)
    defer db.CloseDB()
    handlers.InitHandlers()

    http.HandleFunc(config.IndexPath, handlers.IndexHandler)
    http.HandleFunc(config.LoginPath, handlers.LoginHandler)

    username := "test"
    password := "test"
    db.AddUserToDB(username, password)
    log.Println("User successfully added!")

    log.Printf("Server started on :%s\n", config.Port)
    if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
        log.Fatal(err)
    }
}

