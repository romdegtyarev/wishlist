package db

import (
    "time"
)

type UsersTable struct {
    ID              int
    Username        string
    PasswordHash    string
    CreationDate    time.Time
    LastUpdatedDate time.Time
    Nick            string
    PhotoID         int
    Link            string
    Following       []int
    Followers       []int
}

type WishListTable struct {
    ID              int
    UserID          int
    Name            string
    Text            string
    CreationDate    time.Time
    LastUpdatedDate time.Time
    IsPrivate       bool
}

type ItemsListTable struct {
    ID           int
    WishListID   int
    UserID       int
    Status       bool
    BookedUserID int
}

