package main

import (
    "log"
    "net/http"

    "wishlist/internal/config"
    "wishlist/internal/db"
    "wishlist/internal/handlers"
)

func main() {
    config.Init()
    db.Init(config.DataSourceName)
    defer db.Close()
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

