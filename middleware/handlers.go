package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"go-postgres/models" // models package where User schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to access the environment variables
	"strconv"  // used to convert string to int64

	"github.com/gorilla/mux" // used to access the params in the url
	// used to access the environment variables
	_ "github.com/lib/pq" // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	MESSAGE string `json:"message,omitempty"`
}

// cretae connection to the database
func createConnection() *sql.DB {
	// load the environment variables
	// err := godotenv.Load(".env")

	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	// open the connection to the database
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")
	return db
}

// CreateUser func is used to create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty user of models.User
	var user models.User
	fmt.Println("User: ", user)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("unable to decode the request body: %v", err)
	}

	// call insert user func and pass the user object
	insertID := insertUser(user)

	// format a response object
	res := response{
		ID:      insertID,
		MESSAGE: "User created successfully",
	}

	// encode the response object into json
	json.NewEncoder(w).Encode(res)
}

//  GetUser func is used to get a user by id
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get the user id from the params. key is "id"
	params := mux.Vars(r)

	// convert the string id to int64
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("unable to convert the id to int: %v", err)
	}

	// call the get user func and pass the id
	user, err := getUser(int64(id))

	if err != nil {
		log.Fatalf("unable to get the user: %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(user)
}

// GetAllUsers func is used to get all users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get all the users in the db
	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("unable to get the users: %v", err)
	}

	// send all the users as response
	json.NewEncoder(w).Encode(users)

}

// UpdateUser update user's detail in the postgres db
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Get the user id from the params. key is "id"
	params := mux.Vars(r)

	// convert the string id to int64
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("unable to convert the id to int: %v", err)
	}

	// create an empty user of type models.User
	var user models.User

	// decode the json request to user
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("unable to decode the request body: %v", err)
	}

	// call update user to update the user
	updatedRows := updateUser(int64(id), user)

	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)

	// format the response object
	res := response{
		ID:      int64(id),
		MESSAGE: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// DeleteUser delete user's detail in the postgres db
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deleteUser, convert the int to int64
	deletedRows := deleteUser(int64(id))

	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		MESSAGE: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
// insert one user in the DB
func insertUser(user models.User) int64 {
	// open the db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := "INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid"

	// insered id will be stored in the id variable
	var id int64

	// excecute the query
	// the Scan function will fill the id variable with the id of the inserted user
	err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)

	if err != nil {
		log.Fatalf("unable to execute the query: %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the id of the inserted user
	return id

}

// get one user from the DB by its userid
func getUser(id int64) (models.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a user of models.User type
	var user models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users WHERE userid=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to user
	err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return user, err

}

// get all users from the DB
func getAllUsers() ([]models.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create an empty array of type models.User
	var users []models.User

	// create the select sql query
	sqlStatement := "SELECT * FROM users"

	// execute the query
	rows, err := db.Query(sqlStatement)

	// unmarshal the rows into the users
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location)
		if err != nil {
			log.Fatalf("unable to scan the row: %v", err)
		}
		users = append(users, user)
	}

	// return the users array
	return users, err
}

// update user in the DB
func updateUser(id int64, user models.User) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := "UPDATE users SET name = $2, location = $3, age = $4 WHERE userid = $1"

	// execute the query
	res, err := db.Exec(sqlStatement, id, user.Name, user.Location, user.Age)

	if err != nil {
		log.Fatalf("unable to execute the query: %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete user from the DB
func deleteUser(id int64) int64 {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM users WHERE userid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("unable to execute the query: %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
