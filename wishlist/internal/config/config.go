package config

import (
    "log"
    "os"
)

var (
    Port               string
    DataSourceName     string
    SessionSecretKey   string
    IndexPath          string
    IndexTemplatePath  string
    LoginPath          string
    LoginTemplatePath  string
    LogoutPath         string
)

func Init() {
    Port = os.Getenv("PORT")
    if Port == "" {
        log.Fatal("Error: Environment variable PORT is not set.")
    }

    DataSourceName = os.Getenv("DATA_SOURCE_NAME")
    if DataSourceName == "" {
        log.Fatal("Error: Environment variable DATA_SOURCE_NAME is not set.")
    }

    SessionSecretKey = os.Getenv("SESSION_SECRET_KEY")
    if SessionSecretKey == "" {
        log.Fatal("Error: Environment variable SESSION_SECRET_KEY is not set.")
    }

    IndexPath = os.Getenv("INDEX_PATH")
    if IndexPath == "" {
        IndexPath = "/"
    }

    IndexTemplatePath = os.Getenv("INDEX_TEMPLATE_PATH")
    if IndexTemplatePath == "" {
        log.Fatal("Error: Environment variable INDEX_TEMPLATE_PATH is not set.")
    }

    LoginPath = os.Getenv("LOGIN_PATH")
    if LoginPath == "" {
        log.Fatal("Error: Environment variable LOGIN_PATH is not set.")
    }

    LoginTemplatePath = os.Getenv("LOGIN_TEMPLATE_PATH")
    if LoginTemplatePath == "" {
        log.Fatal("Error: Environment variable LOGIN_TEMPLATE_PATH is not set.")
    }

    LogoutPath = os.Getenv("LOGOUT_PATH")
    if LogoutPath == "" {
        log.Fatal("Error: Environment variable LOGOUT_PATH is not set.")
    }
}

