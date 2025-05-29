package main

import (
	"fmt"
	"github.com/go-passwd/validator"
	"errors"
	"database/sql"
	"os"
	_ "github.com/go-sql-driver/mysql"
    "path/filepath"
	"github.com/joho/godotenv"  // Load environment variables from .env file
	"log"

)


type User struct {
	email string
	password string
}


type UserDB struct {
	id int
	email string
	password string
}
//Passwords table , item , password and userId foreign key
//DB CONNECTION

type PasswordKeeper struct {
	Item string
	Password string


}
type ValidatePasswordError struct {
	password string
	message string
}

func validatePassword(password string) error {
	passwordValidator := validator.New(
		validator.MinLength(8, errors.New("password must be at least 8 characters long")),
		validator.MaxLength(16, errors.New("password cannot exceed 16 characters")),
		validator.ContainsAtLeast("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 1, errors.New("password must contain at least one letter")),
		validator.ContainsAtLeast("0123456789", 1, errors.New("password must contain at least one digit")),
		validator.ContainsOnly("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!(@}~", errors.New("password can only contain letters and digits")),
	)
	err := passwordValidator.Validate(password)

	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}

func (e  ValidatePasswordError) Error() string {
	return fmt.Sprintf("%s Error: %s", e.password, e.message)

}
func main() {
	var item, password string
	var user User

	fmt.Println("This is a password manager cli tool, use it to securely store your passwords")
    fmt.Println("Enter your credentials")

	fmt.Println("Enter email")
	fmt.Scan(&user.email)
	fmt.Println("Enter your alpha password that will be used to access your other passwords")
	fmt.Scan(&user.password)

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

	fmt.Println("You have entered email", dbUser, dbPass, dbName)
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
      panic(err.Error())
    }
    defer db.Close()
    CreateUser(db, user.email, user.password)
	fmt.Println("Excellent! you can begin using the app")
	
	


	fmt.Println("Enter an item")
	fmt.Scan(&item)
	fmt.Println("You have entered item", item)
	fmt.Println("Enter its password")
	fmt.Scan(&password)
	passErr := validateInput(item, password)
	if passErr, ok := passErr.(ValidatePasswordError); ok {
		fmt.Printf("Password storage failed: %s\n", passErr.message)
	}else {
		result := storesecret(item, password)
		fmt.Println("Your password for site ",result.Item, "is saved with this password \n", result.Password )


	}
	
	err2 := CreatePassword(db,  user.email, item , user.password)
	if err2 != nil {
	    fmt.Println(err2)
	}
	// fmt.Println("Success for storing your first password", created_password)
	
}


func CreateUser(db *sql.DB, email, password string) error{
	query := "INSERT INTO Users (email, password ) VALUES(?, ?)"
	_, err := db.Exec(query, email, password)
	if err != nil {
		return err
	}
	return nil

}
func CreatePassword(db *sql.DB, email, item, password string) error {
	createdUser, err1 :=  GetUser(db, email)
	if err1 != nil {
		fmt.Println(err1)
	}
    
	query := "INSERT INTO Password (userId, item, password ) VALUES(?, ?, ?)"
	_, err := db.Exec(query, createdUser.id , item, password)
	if err != nil {
		return err
	}
	return nil

}
func GetUser(db *sql.DB, email string) (*UserDB, error) {
    query := "SELECT * FROM Users WHERE email = ?"
    row := db.QueryRow(query, email)

    user := &UserDB{}
    err := row.Scan(&user.id, &user.email, &user.password)
    if err != nil {
        return nil, err
    }
    return user, nil
}





func storesecret(item string, new_password string) PasswordKeeper {
	value := PasswordKeeper{
		Item : item,
		Password: new_password,

	}
    return value
}
func validateInput(item,new_password string) error {
	error1 := validatePassword(new_password)
	fmt.Println(error1)
	if error1 != nil {
		ans := ValidatePasswordError{
			password: new_password,
			message: error1.Error(),
		}
		return ans

	}
	
	return nil



}

