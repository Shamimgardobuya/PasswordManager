package database

import (
	"database/sql"
	"os"
	"log"
	_ "github.com/go-sql-driver/mysql"
    "path/filepath"
	"github.com/joho/godotenv"  


)

func InitDb() {
	pwd, err := os.Getwd()
	if err != nil {
				panic(err)
	}
	err = godotenv.Load(filepath.Join(pwd, ".env"))
	if err != nil {
				log.Fatal("Error loading .env file")
	}
	
	var dbDriver  = "mysql" 
	var dbUser  = os.Getenv("DB_USER")
	var dbPass   = os.Getenv("DB_PASSWORD")
	var dbName   = os.Getenv("DB_NAME")

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
      panic(err.Error())
    }
    defer db.Close()

	createTableUsers := `CREATE TABLE IF NOT EXISTS Users (
		id INTEGER NOT NULL AUTO_INCREMENT,
		email VARCHAR(50) NOT NULL,
		UNIQUE (email),
		PRIMARY KEY (id)
	);`
	_, err = db.Exec(createTableUsers)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	createTablePassword := `CREATE TABLE IF NOT EXISTS Password (
		id INTEGER NOT NULL AUTO_INCREMENT,
		userId INTEGER,
		item VARCHAR(20) NOT NULL,
		password  VARCHAR(50)  NOT NULL,
		PRIMARY KEY (id), 
		FOREIGN KEY (userId) REFERENCES Users(id)
	);`
	_, err = db.Exec(createTablePassword)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

}