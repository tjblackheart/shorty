package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tjblackheart/shorty/pkg/models"
	"github.com/tjblackheart/shorty/pkg/models/sqlite"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

type application struct {
	config struct {
		dsn string
	}
	users *sqlite.UserModel
}

var app application

func init() {
	flag.StringVar(&app.config.dsn, "dsn", "./var/shorty.db", "Database DSN")
	flag.Parse()
}

func main() {
	db, err := openDB(app.config.dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	app.users = &sqlite.UserModel{DB: db}
	if err = app.users.InitTable(); err != nil {
		log.Fatalln(err)
	}

	var email string
	fmt.Print("Email: ")
	fmt.Scanf("%s", &email)

	if err = validateEmail(email); err != nil {
		log.Fatalln("Invalid email.")
	}

	fmt.Print("Password: ")
	pwd, _ := terminal.ReadPassword(0) // silence input
	fmt.Println()                      // newline

	if err = validatePassword(string(pwd)); err != nil {
		log.Fatalln(err)
	}

	fmt.Print("Confirm password: ")
	confirm, _ := terminal.ReadPassword(0)
	fmt.Println()

	if string(pwd) != string(confirm) {
		log.Fatalln("Passwords do not match.")
	}

	hash, err := bcrypt.GenerateFromPassword(pwd, 12)
	if err != nil {
		log.Fatalln(err)
	}

	u := &models.User{
		Email: email,
		Pass:  string(hash),
	}

	err = app.users.Create(u)
	if err != nil {
		if err == models.ErrUnique {
			log.Fatalln("The user already exists.")
		}

		log.Fatalln(err)
	}

	fmt.Printf("User %s successfully created.\n", email)
}

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func validateEmail(email string) error {
	rxEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(email) > 254 || rxEmail.MatchString(email) == false {
		return errors.New("create_user: not a valid email address")
	}

	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return errors.New("create_user: password cannot be empty")
	}

	if len(password) < 8 {
		return errors.New("create_user: password length too short")
	}

	return nil
}
