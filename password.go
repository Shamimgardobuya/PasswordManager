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
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"io"
	"bytes"
	"strconv"


)

type User struct {
	Email string
	Password string
}


type UserDB struct {
	id int
	email string
	password string
}

type PasswordKeeper struct {
	Item string
	Password string

}

type DataDb struct {
	UserId string
	Item string
	Password string


}
type ValidatePasswordError struct {
	password string
	message string
}
type Users []User

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
    pwd, err := os.Getwd()
	if err != nil {
				panic(err)
	}
	err = godotenv.Load(filepath.Join(pwd, ".env"))
	if err != nil {
				log.Fatal("Error loading .env file")
	}
	
	
	route := mux.NewRouter()

	route.HandleFunc("/users/create", createUserHandler).Methods("POST")
    route.HandleFunc("/password/{item}/{password}", getPasswordHandler).Methods("GET")
    route.HandleFunc("/password", loginMiddleware(createPasswordHandler)).Methods("POST")
	route.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
	
		var dbDriver  = "mysql" 
		var dbUser  = os.Getenv("DB_USER")
		var dbPass   = os.Getenv("DB_PASSWORD")
		var dbName   = os.Getenv("DB_NAME")
	
		fmt.Println(dbUser)
		db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
		if err != nil {
		panic(err.Error())
		}
		defer db.Close()

		data, err := getAllUsers(db)
		if err != nil {
            
		  fmt.Println(err)
		}
		fmt.Fprintln(w, data)
		
	} )
	
    
	
    fmt.Println("Server listening on port http://localhost:8070")
    if err := http.ListenAndServe(":8070", route); err != nil {
        fmt.Println("Error starting server:", err)
    }


	
}

func getPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var dbDriver  = "mysql" 
	var dbUser  = os.Getenv("DB_USER")
	var dbPass   = os.Getenv("DB_PASSWORD")
	var dbName   = os.Getenv("DB_NAME")

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
      panic(err.Error())
    }
    defer db.Close()

	vars := mux.Vars(r)
    passwordStr := vars["password"]
	itemStr := vars["item"]

    fmt.Println("is it present ",passwordStr)
	data, err1 := getPassword(db, passwordStr, itemStr)
	if err1 != nil {
		return 
	}
	
    fmt.Println("Password of item" , itemStr, "has been retrieved successfully", data.Password)
	


}

func createUserHandler( w http.ResponseWriter, r *http.Request) {
	var dbDriver  = "mysql" 
	var dbUser  = os.Getenv("DB_USER")
	var dbPass   = os.Getenv("DB_PASSWORD")
	var dbName   = os.Getenv("DB_NAME")

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
      panic(err.Error())
    }
    defer db.Close()

	var user User

	json.NewDecoder(r.Body).Decode(&user)

	CreateUser(db, user.Email, user.Password)
	
	w.WriteHeader(http.StatusCreated)
	fmt.Println(w, "User created successfully")
	
}

func addBodyToRequest(r *http.Request, data []byte) {
	r.Body = io.NopCloser(bytes.NewReader(data))
	r.ContentLength = int64(len(data))
	r.Body.Close()
}

func loginMiddleware(f func(w http.ResponseWriter, r *http.Request)) (func(http.ResponseWriter, *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dbDriver  = "mysql" 
		var dbUser  = os.Getenv("DB_USER")
		var dbPass   = os.Getenv("DB_PASSWORD")
		var dbName   = os.Getenv("DB_NAME")
		
		type InputData struct {
			Email string 
			Item string
			Password string
		}


		var input InputData

	
		json.NewDecoder(r.Body).Decode(&input)

		db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
		if err != nil {
		panic(err.Error())
		}
    	defer db.Close()
		createdUser, err1 :=  GetUser(db, input.Email)
		if err1 != nil {
			fmt.Println(err1)
		}
		s2 := strconv.Itoa(createdUser.id)
		
		data := map[string]string{
			"item" : input.Item,
			"password" : input.Password,
			"UserId" : s2,
		}
	
		newBo, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		
		fmt.Println("whole data", string(newBo))
		newBody := []byte(string(newBo))
	
		addBodyToRequest(r,newBody)
		
		f(w, r)
	}

}

func createPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var password DataDb
	var dbDriver  = "mysql" 
	var dbUser  = os.Getenv("DB_USER")
	var dbPass   = os.Getenv("DB_PASSWORD")
	var dbName   = os.Getenv("DB_NAME")

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
      panic(err.Error())
    }
    defer db.Close()

	json.NewDecoder(r.Body).Decode(&password)
	fmt.Println("password", password)

	CreatePassword(db, password.UserId, password.Item , password.Password)
	w.WriteHeader(http.StatusCreated)
	fmt.Println("Password created successfully")
}




func getAllUsers(db *sql.DB) (Users, error) {
	var userss Users
	rows, err := db.Query("SELECT Users.email , Users.password FROM Users")
	
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Email, &user.Password); err != nil {
            return nil, err
        }
		userss = append(userss, user)
	}
	return userss, nil
}

func CreateUser(db *sql.DB, email, password string) error{
	query := "INSERT INTO Users (email, password ) VALUES(?, ?)"
	_, err := db.Exec(query, email, password)
	if err != nil {
		return err
	}
	return nil

}

func getPassword(db *sql.DB, alpha_password, item string) (*DataDb, error){
	query := "Select id from Users where password = ?"
	row := db.QueryRow(query, alpha_password)

	var user string
	err := row.Scan(&user)
	
	if err != nil {
        return nil, err
    }
	query2 := "Select userId, item, password from Password where item = ?"
	
    password := &DataDb{}
	row1 := db.QueryRow(query2, item)


	err1 := row1.Scan(&password.UserId, &password.Item, &password.Password)
	
	if err1 != nil {
        return nil, err
    }
	return password,nil
}
func CreatePassword(db *sql.DB, userId, item, password string) error {
	
	passErr := validateInput(item, password)
	if passErr, ok := passErr.(ValidatePasswordError); ok {
		fmt.Printf("Password storage failed: %s\n", passErr.message)
	}else {
		result := storesecret(item, password)
		fmt.Println("Your password for site ",result.Item, "is saved with this password \n", result.Password )

	}
	query := "INSERT INTO Password (userId, item, password ) VALUES(?, ?, ?)"
	_, err := db.Exec(query, userId , item, password)
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

