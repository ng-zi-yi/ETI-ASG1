package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	UserID          int    `json:"userId"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	MobileNo        string `json:"mobileNo"`
	Email           string `json:"email"`
	DriverLicenseNo string `json:"driverLicenseNo"`
	CarPlateNo      string `json:"carPlateNo"`
	UserType        string `json:"userType"`
}

var (
	db  *sql.DB
	err error
)

func dB() {
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpooling_db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}

func main() {
	dB()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/users", createUserHandler).Methods("POST")
	router.HandleFunc("/api/v1/users", getUserByEmailHandler).Methods("GET")
	router.HandleFunc("/api/v1/users", upgradeToCarOwnerHandler).Methods("PUT")
	router.HandleFunc("/api/v1/users/{userID}", updatePassengerProfileHandler).Methods("PUT")
	router.HandleFunc("/api/v1/carowners/{userID}", updateCarOwnerProfileHandler).Methods("PUT")
	router.HandleFunc("/api/v1/users/{userID}/accountCreationDate", getAccCreationDateHandler).Methods("GET")
	router.HandleFunc("/api/v1/users/{userID}", deleteUserAccHandler).Methods("DELETE")

	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the new user into the database
	stmt, err := db.Prepare("INSERT INTO User (FirstName, LastName, MobileNo, Email, DriverLicenseNo, CarPlateNo, UserType) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUser.FirstName, newUser.LastName, newUser.MobileNo, newUser.Email, newUser.DriverLicenseNo, newUser.CarPlateNo, newUser.UserType)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Account created successfully")
}

func getUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	var user User
	err := db.QueryRow("SELECT UserID, FirstName, LastName, MobileNo, Email, DriverLicenseNo, CarPlateNo, UserType FROM User WHERE Email = ?", email).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.MobileNo, &user.Email, &user.DriverLicenseNo, &user.CarPlateNo, &user.UserType)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond with user information
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func upgradeToCarOwnerHandler(w http.ResponseWriter, r *http.Request) {
	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the user's information in the database
	stmt, err := db.Prepare("UPDATE User SET DriverLicenseNo=?, CarPlateNo=?, UserType=? WHERE UserID=?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedUser.DriverLicenseNo, updatedUser.CarPlateNo, updatedUser.UserType, updatedUser.UserID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Upgrade successful, user is now a Car Owner!")
}

func updatePassengerProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from the request URL
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the user's information in the database
	stmt, err := db.Prepare("UPDATE User SET FirstName=?, LastName=?, MobileNo=?, Email=? WHERE UserID=?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedUser.FirstName, updatedUser.LastName, updatedUser.MobileNo, updatedUser.Email, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Passenger profile updated successfully!")
}

func updateCarOwnerProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from the request URL
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the user's information in the database
	stmt, err := db.Prepare("UPDATE User SET FirstName=?, LastName=?, MobileNo=?, Email=?, DriverLicenseNo=?, CarPlateNo=? WHERE UserID=?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedUser.FirstName, updatedUser.LastName, updatedUser.MobileNo, updatedUser.Email, updatedUser.DriverLicenseNo, updatedUser.CarPlateNo, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Car Owner profile updated successfully!")
}

// Add the getAccCreationDateHandler function
func getAccCreationDateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from the request URL
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch the account creation date from the database
	var accCreationDate string
	err = db.QueryRow("SELECT AccCreationDate FROM User WHERE UserID = ?", userID).Scan(&accCreationDate)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond with the account creation date
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accCreationDate)
}

// Add the deleteUserHandler function
func deleteUserAccHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from the request URL
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Delete the user from the database
	stmt, err := db.Prepare("DELETE FROM User WHERE UserID=?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User Account deleted successfully!")
}
