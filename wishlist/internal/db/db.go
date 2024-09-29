package db

import (
    "database/sql"
    "log"
    "golang.org/x/crypto/bcrypt"
    _ "github.com/lib/pq"
    "encoding/gob"

    "wishlist/internal/config"
)

var database *sql.DB

func InitDB(dataSourceName string) {
    var err error
    database, err = sql.Open("postgres", config.DataSourceName)
    if err != nil {
        log.Fatal(err)
    }

    if err = database.Ping(); err != nil {
        log.Fatal(err)
    }

    err = createUsersTable()
    if err != nil {
        log.Fatalf("Error creating users table: %v", err)
    }

    gob.Register(UsersTable{})
}

func CloseDB() {
    if err := database.Close(); err != nil {
        log.Fatal(err)
    }
}

func GetUserFromDB(username string) (string, error) {
    var passwordHash string
    err := database.QueryRow("SELECT passwordhash FROM userstable WHERE username = $1", username).Scan(&passwordHash)
    if err != nil {
        return "", err
    }
    return passwordHash, nil
}

func AddUserToDB(username string, password string) error {
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    _, err = database.Exec("INSERT INTO userstable (username, passwordhash) VALUES ($1, $2)", username, passwordHash)
    if err != nil {
        return err
    }

    return nil
}

func createUsersTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS userstable (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        passwordhash VARCHAR(255) NOT NULL
    );`
    _, err := database.Exec(query)
    return err
}

