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

type Trip struct {
	TripID          int    `json:"tripId"`
	UserID          int    `json:"userId"`
	FirstName       string `json:"firstName"`
	PickupAddr      string `json:"pickupAddr"`
	AltPickupAddr   string `json:"altPickupAddr"`
	StartTravelTime string `json:"startTravelTime"`
	DestAddr        string `json:"destAddr"`
	MaxPassengers   int    `json:"maxPassengers"`
	VacanciesLeft   int    `json:"vacanciesLeft"`
	TripStatus      string `json:"tripStatus"`
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

	router.HandleFunc("/api/v1/trips", publishNewTripHandler).Methods("POST")
	router.HandleFunc("/api/v1/trips", listPublishedTripsHandler).Methods("GET")
	router.HandleFunc("/api/v1/trips/carowners/{userID}", listCarOwnerTripsHandler).Methods("GET")
	router.HandleFunc("/api/v1/trips/start/{tripID}", startTripHandler).Methods("PUT")
	router.HandleFunc("/api/v1/trips/start-time/{tripID}", getStartTimeHandler).Methods("GET")
	router.HandleFunc("/api/v1/trips/cancel/{tripID}", cancelTripHandler).Methods("PUT")
	router.HandleFunc("/api/v1/trips/search", searchTripsHandler).Methods("GET")
	router.HandleFunc("/api/v1/trips/enrol", enrolInTripHandler).Methods("POST")
	router.HandleFunc("/api/v1/trips/enroltrips/{userID}", listPassengerEnroledTripIDsHandler).Methods("GET")
	router.HandleFunc("/api/v1/trips/{tripID}", getTripByIDHandler).Methods("GET")

	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))
}

// Add the publishNewTripHandler function
func publishNewTripHandler(w http.ResponseWriter, r *http.Request) {
	var newTrip Trip
	err := json.NewDecoder(r.Body).Decode(&newTrip)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the new trip into the database
	stmt, err := db.Prepare("INSERT INTO Trip (UserID, FirstName, PickupAddr, AltPickupAddr, StartTravelTime, DestAddr, MaxPassengers, VacanciesLeft, TripStatus) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(newTrip.UserID, newTrip.FirstName, newTrip.PickupAddr, newTrip.AltPickupAddr, newTrip.StartTravelTime, newTrip.DestAddr, newTrip.MaxPassengers, newTrip.MaxPassengers, newTrip.TripStatus)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Trip published successfully!")
}

// Add the listPublishedTripsHandler function
func listPublishedTripsHandler(w http.ResponseWriter, r *http.Request) {
	trips, err := listPublishedTripsFromDB()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert trips to JSON
	response, err := json.Marshal(trips)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Add the listPublishedTripsFromDB function to retrieve trips from the database
func listPublishedTripsFromDB() ([]Trip, error) {
	rows, err := db.Query("SELECT TripID, UserID, FirstName, PickupAddr, AltPickupAddr, StartTravelTime, DestAddr, MaxPassengers,VacanciesLeft, TripStatus FROM Trip WHERE TripStatus = 'Waiting'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []Trip
	for rows.Next() {
		var trip Trip
		if err := rows.Scan(&trip.TripID, &trip.UserID, &trip.FirstName, &trip.PickupAddr, &trip.AltPickupAddr, &trip.StartTravelTime, &trip.DestAddr, &trip.MaxPassengers, &trip.VacanciesLeft, &trip.TripStatus); err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}
	return trips, nil
}

// Add the listCarOwnerTripsHandler function
func listCarOwnerTripsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	trips, err := listCarOwnerTripsFromDB(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert trips to JSON
	response, err := json.Marshal(trips)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Add the listCarOwnerTripsFromDB function to retrieve car owner's trips from the database
func listCarOwnerTripsFromDB(userID string) ([]Trip, error) {
	rows, err := db.Query("SELECT TripID, UserID, FirstName, PickupAddr, AltPickupAddr, StartTravelTime, DestAddr, MaxPassengers, VacanciesLeft, TripStatus FROM Trip WHERE UserID = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []Trip
	for rows.Next() {
		var trip Trip
		if err := rows.Scan(&trip.TripID, &trip.UserID, &trip.FirstName, &trip.PickupAddr, &trip.AltPickupAddr, &trip.StartTravelTime, &trip.DestAddr, &trip.MaxPassengers, &trip.VacanciesLeft, &trip.TripStatus); err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}
	return trips, nil
}

// Add the startTripHandler function
func startTripHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tripID := params["tripID"]

	// Update the trip status to "Started" in the database
	stmt, err := db.Prepare("UPDATE Trip SET TripStatus = 'Started' WHERE TripID = ?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(tripID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Trip started successfully!")
}

// Add the getStartTimeHandler function, to calculate beforeCancelTime for checking (part of cancel trip feature)
func getStartTimeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tripID := params["tripID"]

	// Retrieve the start time of the trip from the database
	var startTime string
	err := db.QueryRow("SELECT StartTravelTime FROM Trip WHERE TripID = ?", tripID).Scan(&startTime)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return the start time as a JSON response
	response := map[string]string{"start_time": startTime}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Add the cancelTripHandler function, updates trip status to cancelled if cancel 30mins b4 scheduled time(check is in admin.go)
func cancelTripHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tripID := params["tripID"]

	// Update the trip status to "Started" in the database
	stmt, err := db.Prepare("UPDATE Trip SET TripStatus = 'Cancelled' WHERE TripID = ?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(tripID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Trip cancelled successfully!")
}

// Add the searchTripsHandler function
func searchTripsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the destination address from the query parameters
	destination := r.URL.Query().Get("destination")

	// Query the database to search for trips based on the destination address
	rows, err := db.Query("SELECT TripID, UserID, FirstName, PickupAddr, AltPickupAddr, StartTravelTime, DestAddr, MaxPassengers, VacanciesLeft, TripStatus FROM Trip WHERE DestAddr LIKE ?", "%"+destination+"%")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the search results
	var searchResults []Trip

	// Iterate through the rows and populate the search results slice
	for rows.Next() {
		var trip Trip
		if err := rows.Scan(&trip.TripID, &trip.UserID, &trip.FirstName, &trip.PickupAddr, &trip.AltPickupAddr, &trip.StartTravelTime, &trip.DestAddr, &trip.MaxPassengers, &trip.VacanciesLeft, &trip.TripStatus); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		searchResults = append(searchResults, trip)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Encode the search results as JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResults)
}

// Add the enrolInTripHandler function
func enrolInTripHandler(w http.ResponseWriter, r *http.Request) {
	var enrolmentData struct {
		TripID int `json:"tripId"`
		UserID int `json:"userId"`
	}
	// Decode enrolment data from the request body
	err := json.NewDecoder(r.Body).Decode(&enrolmentData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the user is already enroled in another trip at the same time
	conflictingTripID, err := checkDateTimeConflict(enrolmentData.UserID, enrolmentData.TripID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if conflictingTripID != 0 {
		http.Error(w, fmt.Sprintf("Unsuccessful enrolment as there is a conflict with tripID %d at the same time.", conflictingTripID), http.StatusBadRequest)
		return
	}

	// Check if there are vacancies in the selected trip
	vacancies, err := checkVacancies(enrolmentData.TripID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if vacancies <= 0 {
		http.Error(w, "No vacancies available for this trip", http.StatusBadRequest)
		return
	}

	// Enrol the user in the trip
	stmt, err := db.Prepare("INSERT INTO TripEnrolment (TripID, UserID) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(enrolmentData.TripID, enrolmentData.UserID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Decrement the number of vacancies in the selected trip
	_, err = db.Exec("UPDATE Trip SET VacanciesLeft = VacanciesLeft - 1 WHERE TripID = ?", enrolmentData.TripID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Enroled successfully!")
}

// Add the checkDateTimeConflict function
func checkDateTimeConflict(userID, tripID int) (int, error) {
	var conflictingTripID int

	// Query the database to check for conflicting trips
	err := db.QueryRow("SELECT t.TripID FROM TripEnrolment te JOIN Trip t ON te.TripID = t.TripID WHERE te.UserID = ? AND t.StartTravelTime = (SELECT StartTravelTime FROM Trip WHERE TripID = ?)", userID, tripID).Scan(&conflictingTripID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No conflicting trips found
			return 0, nil
		}
		return 0, err
	}

	return conflictingTripID, nil
}

// Add the checkVacancies function
func checkVacancies(tripID int) (int, error) {
	var vacancies int

	// Query the database to check the number of vacancies in the selected trip
	err := db.QueryRow("SELECT VacanciesLeft FROM Trip WHERE TripID = ?", tripID).Scan(&vacancies)
	if err != nil {
		return 0, err
	}

	return vacancies, nil
}

type enrolmentData struct {
	TripID int `json:"tripId"`
	UserID int `json:"userId"`
}

func listPassengerEnroledTripIDsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	enrolDatas, err := listPassengerEnroledTripIDsFromDB(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert trips to JSON
	response, err := json.Marshal(enrolDatas)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Add the listPassengerEnroledTripIDsFromDB function to retrieve passenger's tripids from the database
func listPassengerEnroledTripIDsFromDB(userID string) ([]enrolmentData, error) {
	rows, err := db.Query("SELECT TripID, UserID FROM TripEnrolment WHERE UserID = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrolDatas []enrolmentData
	for rows.Next() {
		var enrolData enrolmentData
		if err := rows.Scan(&enrolData.TripID, &enrolData.UserID); err != nil {
			return nil, err
		}
		enrolDatas = append(enrolDatas, enrolData)
	}
	return enrolDatas, nil
}

// Add the getTripByIDHandler function
func getTripByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripID, err := strconv.Atoi(vars["tripID"])
	if err != nil {
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	// Fetch the trip from the database based on tripID
	trip, err := getTripByID(tripID)
	if err != nil {
		http.Error(w, "Error fetching trip", http.StatusInternalServerError)
		return
	}

	// Return the trip details as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trip)
}

func getTripByID(tripID int) (*Trip, error) {
	// Perform a SELECT query to fetch the trip details based on tripID
	row := db.QueryRow("SELECT TripID, UserID, FirstName, PickupAddr, AltPickupAddr, StartTravelTime, DestAddr, MaxPassengers, VacanciesLeft, TripStatus FROM Trip WHERE TripID = ?", tripID)

	var trip Trip
	err := row.Scan(&trip.TripID, &trip.UserID, &trip.FirstName, &trip.PickupAddr, &trip.AltPickupAddr, &trip.StartTravelTime, &trip.DestAddr, &trip.MaxPassengers, &trip.VacanciesLeft, &trip.TripStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Trip not found")
		}
		return nil, err
	}

	return &trip, nil
}
