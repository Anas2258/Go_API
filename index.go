package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// const (
// 	DB_USER     = "postgres"
// 	DB_PASSWORD = "anas@2258"
// 	DB_PASSWORD     = "movies"
// )

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	checkErr(err)

	return os.Getenv(key)
}

// DB set up
func setupDB() *sql.DB {
	userName := goDotEnvVariable("DB_USER")
	password := goDotEnvVariable("DB_PASSWORD")
	dbName := goDotEnvVariable("DB_NAME")

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", userName, password, dbName)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

type Movie struct {
	MovieID   string `json:"movieid"`
	MovieName string `json:"moviename"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Movie `json:"data"`
	Message string  `json:"message"`
}

func main() {

	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	// Get all movies
	router.HandleFunc("/movies/", GetMovies).Methods("GET")

	// Create a movie
	router.HandleFunc("/movies/", CreateMovie).Methods("POST")

	// Update a movie
	router.HandleFunc("/movies/{movieid}", UpdateMovie).Methods("POST")

	// Delete a specific movie by the movieID
	router.HandleFunc("/movies/{movieid}", DeleteMovie).Methods("DELETE")

	// Delete all movies
	router.HandleFunc("/movies/", DeleteMovies).Methods("DELETE")

	// serve the app
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting movies...")

	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM movies")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var movies []Movie

	// Foreach movie
	for rows.Next() {
		var id int
		var movieID string
		var movieName string

		err = rows.Scan(&id, &movieID, &movieName)

		// check errors
		checkErr(err)

		movies = append(movies, Movie{MovieID: movieID, MovieName: movieName})
	}

	var response = JsonResponse{Type: "success", Data: movies}

	json.NewEncoder(w).Encode(response)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	movieID := r.FormValue("movieid")
	movieName := r.FormValue("moviename")

	var response = JsonResponse{}

	if movieID == "" || movieName == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID or movieName parameter."}
	} else {
		db := setupDB()

		printMessage("Inserting movie into DB")

		fmt.Println("Inserting new movie with ID: " + movieID + " and name: " + movieName)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO movies(movieID, movieName) VALUES($1, $2) returning id;", movieID, movieName).Scan(&lastInsertID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieName := r.FormValue("moviename")

	movieID := params["movieid"]

	var response = JsonResponse{}

	if movieID == "" || movieName == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	} else {
		db := setupDB()

		printMessage("Updating movie name from DB")

		_, err := db.Exec("UPDATE movies SET MovieName = $2 WHERE movieID = $1", movieID, movieName)
		// _, err := db.Exec("DELETE FROM movies where movieID = $1", movieID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been updated successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	movieID := params["movieid"]

	var response = JsonResponse{}

	if movieID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting movie from DB")

		_, err := db.Exec("DELETE FROM movies where movieID = $1", movieID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all movies...")

	_, err := db.Exec("DELETE FROM movies")

	// check errors
	checkErr(err)

	printMessage("All movies have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All movies have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
