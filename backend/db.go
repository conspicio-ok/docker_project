package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func readSecret(name string) string {
	data, err := os.ReadFile("/run/secrets/" + name)
	if err != nil {
		log.Fatal("Impossible de lire le secret :", err)
	}
	return strings.TrimSpace(string(data))
}

func InitDB() {
	password := readSecret("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	dbname := os.Getenv("MYSQL_DATABASE")
	user := os.Getenv("MYSQL_USER")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, password, host, dbname)
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur connexion DB :", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("DB injoignable :", err)
	}
	log.Println("Connexion DB OK")
}