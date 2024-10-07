package db

import (
    "database/sql"
    "log"
    "golang.org/x/crypto/bcrypt"
    "github.com/lib/pq"
    "encoding/gob"

    "wishlist/internal/config"
)

var database *sql.DB

func createTables() error {
    userTableQuery := `
    CREATE TABLE IF NOT EXISTS userstable (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        passwordhash VARCHAR(255) NOT NULL,
        creationdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        lastupdateddate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        nick VARCHAR(50),
        photoid INT, --TODO
        link VARCHAR(255),
        following INT[], --TODO
        followers INT[] --TODO
    );`

    wishListTableQuery := `
    CREATE TABLE IF NOT EXISTS wishlisttable (
        id SERIAL PRIMARY KEY,
        userid INT REFERENCES userstable(id),
        name VARCHAR(100) NOT NULL,
        text TEXT,
        creationdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        lastupdateddate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        isprivate BOOLEAN DEFAULT TRUE
    );`

    itemsListTableQuery := `
    CREATE TABLE IF NOT EXISTS itemslisttable (
        id SERIAL PRIMARY KEY,
        wishlistid INT REFERENCES wishlisttable(id),
        userid INT REFERENCES userstable(id),
        status BOOLEAN DEFAULT FALSE,
        bookeduserid INT --TODO
    );`

    _, err := database.Exec(userTableQuery)
    if err != nil {
        return err
    }

    _, err = database.Exec(wishListTableQuery)
    if err != nil {
        return err
    }

    _, err = database.Exec(itemsListTableQuery)
    return err
}

// Database
func Init(dataSourceName string) {
    var err error

    database, err = sql.Open("postgres", config.DataSourceName)
    if err != nil {
        log.Fatal(err)
    }

    err = database.Ping()
    if err != nil {
        log.Fatal(err)
    }

    err = createTables()
    if err != nil {
        log.Fatalf("Error creating tables: %v", err)
    }

    gob.Register(UsersTable{})
}

func Close() {
    err := database.Close()
    if err != nil {
        log.Fatal(err)
    }
}

// User
func AddUser(username string, password string) (int, error) {
    var id int

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return 0, err
    }

    err = database.QueryRow("INSERT INTO userstable (username, passwordhash) VALUES ($1, $2) RETURNING id", username, passwordHash).Scan(&id)
    return id, err
}

func GetUser(username string) (UsersTable, error) {
    var user UsersTable

    err := database.QueryRow("SELECT id, username, passwordhash, creationdate, lastupdateddate, nick, photoid, link, following, followers FROM userstable WHERE username = $1",
        username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreationDate, &user.LastUpdatedDate, &user.Nick, &user.PhotoID, &user.Link, pq.Array(&user.Following), pq.Array(&user.Followers))
    if err != nil {
        return user, err
    }
    return user, nil
}

func DeleteUser(userID int) error {
    _, err := database.Exec("DELETE FROM userstable WHERE id = $1", userID)
    return err
}

func UpdateUserPhoto(userID int, photoID int) error {
    _, err := database.Exec("UPDATE userstable SET photoid = $1, lastupdateddate = CURRENT_TIMESTAMP WHERE id = $2", photoID, userID)
    return err
}

func UpdateUserNick(userID int, nick string) error {
    _, err := database.Exec("UPDATE userstable SET nick = $1, lastupdateddate = CURRENT_TIMESTAMP WHERE id = $2", nick, userID)
    return err
}

func UpdateUserLink(userID int, link string) error {
    _, err := database.Exec("UPDATE userstable SET link = $1, lastupdateddate = CURRENT_TIMESTAMP WHERE id = $2", link, userID)
    return err
}

func AddUserFollower(userID int, followerID int) error {
    _, err := database.Exec("UPDATE userstable SET followers = array_append(followers, $1), lastupdateddate = CURRENT_TIMESTAMP WHERE id = $2", followerID, userID)
    return err
}

func AddUserFollowing(userID int, followingID int) error {
    _, err := database.Exec("UPDATE userstable SET following = array_append(following, $1), lastupdateddate = CURRENT_TIMESTAMP WHERE id = $2", followingID, userID)
    return err
}

// Wishlist
func AddWishList(userID int, name string, text string, isPrivate bool) (int, error) {
    var id int

    err := database.QueryRow("INSERT INTO wishlisttable (userid, name, text, isprivate) VALUES ($1, $2, $3, $4) RETURNING id", userID, name, text, isPrivate).Scan(&id)
    return id, err
}

func GetWishList(wishListID int) (WishListTable, error) {
    var wishList WishListTable

    err := database.QueryRow("SELECT id, userid, name, text, creationdate, lastupdateddate, isprivate FROM wishlisttable WHERE id = $1",
        wishListID).Scan(&wishList.ID, &wishList.UserID, &wishList.Name, &wishList.Text, &wishList.CreationDate, &wishList.LastUpdatedDate, &wishList.IsPrivate)
    return wishList, err
}

func DeleteWishList(wishListID int) error {
    _, err := database.Exec("DELETE FROM wishlisttable WHERE id = $1", wishListID)
    return err
}

// Wishlist item
func AddItemToWishList(wishListID int, userID int) (int, error) {
    var id int

    err := database.QueryRow("INSERT INTO itemslisttable (wishlistid, userid) VALUES ($1, $2) RETURNING id", wishListID, userID).Scan(&id)
    return id, err
}

func GetItem(itemID int) (ItemsListTable, error) {
    var item ItemsListTable

    err := database.QueryRow("SELECT id, wishlistid, userid, status, bookeduserid FROM itemslisttable WHERE id = $1",
        itemID).Scan(&item.ID, &item.WishListID, &item.UserID, &item.Status, &item.BookedUserID)
    return item, err
}

func DeleteItemFromWishList(itemID int) error {
    var wishlistID int

    err := database.QueryRow("SELECT wishlistid FROM itemslisttable WHERE id = $1", itemID).Scan(&wishlistID)
    if err != nil {
        return err
    }

    _, err = database.Exec("DELETE FROM itemslisttable WHERE id = $1", itemID)
    if err != nil {
        return err
    }

    _, err = database.Exec("UPDATE wishlisttable SET lastupdateddate = CURRENT_TIMESTAMP WHERE id = $1", wishlistID)
    return err
}

func BookItem(itemID int, userID int) error {
    var wishlistID int

    err := database.QueryRow("SELECT wishlistid FROM itemslisttable WHERE id = $1", itemID).Scan(&wishlistID)
    if err != nil {
        return err
    }

    _, err = database.Exec("UPDATE itemslisttable SET status = TRUE, bookeduserid = $1 WHERE id = $2", userID, itemID)
    if err != nil {
        return err
    }

    _, err = database.Exec("UPDATE wishlisttable SET lastupdateddate = CURRENT_TIMESTAMP WHERE id = $1", wishlistID)
    return err
}

func UnbookItem(itemID int) error {
    var wishlistID int

    err := database.QueryRow("SELECT wishlistid FROM itemslisttable WHERE id = $1", itemID).Scan(&wishlistID)
    if err != nil {
        return err
    }

    _, err = database.Exec("UPDATE itemslisttable SET status = FALSE, bookeduserid = NULL WHERE id = $1", itemID)
    if err != nil {
        return err
    }

    _, err = database.Exec("UPDATE wishlisttable SET lastupdateddate = CURRENT_TIMESTAMP WHERE id = $1", wishlistID)
    return err
}

