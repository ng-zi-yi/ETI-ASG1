package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type User struct {
	UserID          int    `json:"UserID"`
	FirstName       string `json:"FirstName"`
	LastName        string `json:"LastName"`
	MobileNo        string `json:"MobileNo"`
	Email           string `json:"Email"`
	DriverLicenseNo string `json:"DriverLicenseNo"`
	CarPlateNo      string `json:"CarPlateNo"`
	UserType        string `json:"UserType"`
}

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

type enrolmentData struct {
	TripID int `json:"tripId"`
	UserID int `json:"userId"`
}

func main() {
outer:
	for {
		fmt.Println("===============================================")
		fmt.Println("Welcome to the Car-Pooling Console Application!")
		fmt.Println("1. Create User Account")
		fmt.Println("2. Login")
		fmt.Println("0. Exit")
		fmt.Print("Enter an option: ")

		var choice int
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			// user account creation
			fmt.Println("----Create User Account----")
			createAcc()
		case 2:
			// user login
			fmt.Println("----Login----")
			user, err := login()
			if err != nil {
				fmt.Println("Login failed:", err)
				return
			}
			//after login display passenger main menu
			if user.UserType == "Passenger" {
				passengerMainMenu(user)
			} else {
				carOwnerMainMenu(user)
			}
		case 0:
			break outer
		default:
			fmt.Println("Invalid option")
		}
		fmt.Scanln()
	}
}

func createAcc() {
	var user User
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter First Name: ")
	fmt.Scanf("%v", &user.FirstName)
	reader.ReadString('\n')
	fmt.Print("Enter Last Name: ")
	fmt.Scanf("%v", &user.LastName)
	reader.ReadString('\n')
	fmt.Print("Enter Mobile Number: ")
	fmt.Scanf("%v", &user.MobileNo)
	reader.ReadString('\n')
	fmt.Print("Enter Email: ")
	fmt.Scanf("%v", &user.Email)

	user.DriverLicenseNo = ""
	user.CarPlateNo = ""
	user.UserType = "Passenger"

	postBody, _ := json.Marshal(user)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5001/api/v1/users", bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 201 {
				fmt.Println("Account created successfully!")
			} else {
				fmt.Println("Error creating user account")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

}

func login() (*User, error) {
	var email string
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter Email: ")
	fmt.Scanf("%v", &email)

	// Perform login check
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5001/api/v1/users?email="+email, nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusOK {
				var user User
				err := json.NewDecoder(res.Body).Decode(&user)
				if err == nil {
					fmt.Printf("Welcome back, %s!\n", user.FirstName)
					return &user, nil
				} else {
					return nil, fmt.Errorf("Error decoding response: %v", err)
				}
			} else {
				return nil, fmt.Errorf("Wrong Email")
			}
		} else {
			return nil, fmt.Errorf("Error making request: %v", err)
		}
	} else {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
}

func passengerMainMenu(user *User) {
	for {
		var choice int
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("===============================================")
		fmt.Println("--------------Passenger Main Menu-------------")
		fmt.Println("1. Upgrade to Car Owner")
		fmt.Println("2. Update Profile Information")
		fmt.Println("3. Delete Account")
		fmt.Println("4. List all published trips")
		fmt.Println("5. View past trips")
		fmt.Println("0. Exit")
		reader.ReadString('\n')
		fmt.Print("Enter an option: ")

		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			// Upgrade to Car Owner
			fmt.Println("--------------Upgrade to Car Owner--------------")
			upgradeToCarOwner(user)
			carOwnerMainMenu(user)
		case 2:
			// Update Profile Information
			fmt.Println("----Update Profile Information----")
			updatePassengerProfile(user)
		case 3:
			// Delete Account
			fmt.Println("----Delete Account----")
			//deleteAccount(user)
			deleteUserAcc(user)
		case 4:
			// List all published trips
			fmt.Println("----List all Published Trips----")
			listPublishedTrips(user)
			// menu for passengers to search for trip or enrol trip
			managePassengerTripMenu(user)
		case 5:
			// View past trips
			fmt.Println("----View Past trips----")
			viewPastTripIDs(user)
		case 0:
			os.Exit(0)
		default:
			fmt.Println("Invalid option")
		}
	}
}

func upgradeToCarOwner(user *User) {
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter your Driver's License Number: ")
	fmt.Scanf("%v", &user.DriverLicenseNo)
	reader.ReadString('\n')
	fmt.Print("Enter your Car Plate Number: ")
	fmt.Scanf("%v", &user.CarPlateNo)

	user.UserType = "CarOwner"

	postBody, _ := json.Marshal(user)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:5001/api/v1/users", bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Upgrade successful, you are now a Car Owner!")
			} else {
				fmt.Println("Error upgrading to Car Owner")
			}
		} else {
			fmt.Println("Error making request", err)
		}
	} else {
		fmt.Println("Error creating request", err)
	}
}

func carOwnerMainMenu(user *User) {
	for {
		var choice int
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("===============================================")
		fmt.Println("--------------Car Owner Main Menu--------------")
		fmt.Println("1. Update Profile Information")
		fmt.Println("2. Delete Account")
		fmt.Println("3. Publish new trip")
		fmt.Println("4. List all my created trips")
		fmt.Println("0. Exit")
		reader.ReadString('\n')
		fmt.Print("Enter an option: ")

		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			// Update Profile Information
			fmt.Println("----Update Profile Information----")
			updateCarOwnerProfile(user)
		case 2:
			// Delete Account
			fmt.Println("----Delete Account----")
			deleteUserAcc(user)
		case 3:
			// Publish new trip
			fmt.Println("----Publish new trip----")
			publishNewTrip(user)
		case 4:
			// List all my created trips
			fmt.Println("----List all my created trips----")
			listCarOwnerTrips(user)
			//manage trip menu - start/cancel trip
			manageTripMenu()
		case 0:
			os.Exit(0)
		default:
			fmt.Println("Invalid option")
		}
	}
}

func updatePassengerProfile(user *User) {
	// Request updated information from the user
	var updatedUser User
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter updated First Name: ")
	fmt.Scanf("%v", &updatedUser.FirstName)
	reader.ReadString('\n')
	fmt.Print("Enter updated Last Name: ")
	fmt.Scanf("%v", &updatedUser.LastName)
	reader.ReadString('\n')
	fmt.Print("Enter updated Mobile Number: ")
	fmt.Scanf("%v", &updatedUser.MobileNo)
	reader.ReadString('\n')
	fmt.Print("Enter updated Email: ")
	fmt.Scanf("%v", &updatedUser.Email)

	// Perform the update by making a PUT request to the API
	postBody, _ := json.Marshal(updatedUser)
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:5001/api/v1/users/%d", user.UserID), bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				fmt.Println("Profile updated successfully!")
			} else {
				fmt.Println("Error updating profile")
			}
		} else {
			fmt.Println("Error making request", err)
		}
	} else {
		fmt.Println("Error creating request", err)
	}
}

func updateCarOwnerProfile(user *User) {
	// Request updated information from the user
	var updatedUser User
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter updated First Name: ")
	fmt.Scanf("%v", &updatedUser.FirstName)
	reader.ReadString('\n')
	fmt.Print("Enter updated Last Name: ")
	fmt.Scanf("%v", &updatedUser.LastName)
	reader.ReadString('\n')
	fmt.Print("Enter updated Mobile Number: ")
	fmt.Scanf("%v", &updatedUser.MobileNo)
	reader.ReadString('\n')
	fmt.Print("Enter updated Email: ")
	fmt.Scanf("%v", &updatedUser.Email)
	reader.ReadString('\n')
	fmt.Print("Enter updated Driver's License Number: ")
	fmt.Scanf("%v", &updatedUser.DriverLicenseNo)
	reader.ReadString('\n')
	fmt.Print("Enter updated Car Plate Number: ")
	fmt.Scanf("%v", &updatedUser.CarPlateNo)

	// Perform the update by making a PUT request to the API
	postBody, _ := json.Marshal(updatedUser)
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:5001/api/v1/carowners/%d", user.UserID), bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				fmt.Println("Profile updated successfully!")
			} else {
				fmt.Println("Error updating profile")
			}
		} else {
			fmt.Println("Error making request", err)
		}
	} else {
		fmt.Println("Error creating request", err)
	}
}

// / was displayAccAge function, modified to include deleteAccount
func deleteUserAcc(user *User) {
	// get account creation date
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:5001/api/v1/users/%d/accountCreationDate", user.UserID), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusOK {
				var accCreationDate string
				err := json.NewDecoder(res.Body).Decode(&accCreationDate)
				if err == nil {
					// Calculate the account age, parse accCreationDate
					layout := "2006-01-02 15:04:05"
					parsedTime, err := time.Parse(layout, accCreationDate)
					if err != nil {
						fmt.Println("Error parsing time:", err)
						return
					}

					// Get the current time
					currentTime := time.Now()

					// Calculate the age by subtracting parsedTime from currentTime
					accAge := currentTime.Sub(parsedTime).Hours() / 24 / 365

					// Display the account age for debugging
					fmt.Printf("Account creation datetime: %v\n", parsedTime)
					fmt.Printf("current datetime: %v\n", currentTime)
					fmt.Printf("Account age: %v\n", accAge)

					// Check if account age is greater than or equal to 1 year
					if accAge >= 1 {
						// Perform account deletion
						deleteAccount(user)
					} else {
						// Display message for unsuccessful deletion
						fmt.Println("Unsuccessful deletion, account can only be deleted after 1 year.")
					}
				} else {
					fmt.Println("Error decoding response:", err)
				}
			} else {
				fmt.Println("Error fetching account creation date")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}

// Add the deleteAccount function, called within deleteUserAcc to check for accAge
func deleteAccount(user *User) {
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:5001/api/v1/users/%d", user.UserID), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				fmt.Println("Account deleted successfully!")
				os.Exit(3)
			} else {
				fmt.Println("Error deleting account")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}

func publishNewTrip(user *User) {
	var trip Trip
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter a Pick-up Address: ")
	fmt.Scanf("%v", &trip.PickupAddr)
	reader.ReadString('\n')
	fmt.Print("Enter an Alternative Pick-up Address: ")
	fmt.Scanf("%v", &trip.AltPickupAddr)
	reader.ReadString('\n')
	fmt.Print("Enter a Start Travel Time (YYYY-MM-DD hh:mm:ss): ")
	trip.StartTravelTime, _ = reader.ReadString('\n')
	trip.StartTravelTime = strings.TrimSpace(trip.StartTravelTime)
	fmt.Print("Enter a Destination Address: ")
	fmt.Scanf("%v", &trip.DestAddr)
	reader.ReadString('\n')
	fmt.Print("Enter the number of passengers your car can accommodate: ")
	fmt.Scanf("%v", &trip.MaxPassengers)

	trip.UserID = user.UserID
	trip.FirstName = user.FirstName
	trip.VacanciesLeft = trip.MaxPassengers
	trip.TripStatus = "Waiting" //waiting refers to trip that has been published but not started

	postBody, _ := json.Marshal(trip)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5002/api/v1/trips", bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 201 {
				fmt.Println("New trip created successfully!")
			} else {
				fmt.Println("Error creating new trip")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

}

// Add the listPublishedTrips function, lists all trips with the watiting status in the passenger menu
func listPublishedTrips(user *User) {
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5002/api/v1/trips", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusOK {
				var trips []Trip
				err := json.NewDecoder(res.Body).Decode(&trips)
				if err == nil {
					fmt.Println("Published Trips:")
					for _, trip := range trips {
						fmt.Printf("\nTrip ID: %d\nDriver Name: %s\nPick-up Address: %s\nAlternative Pick-up Address: %s\nStart Travel Time: %s\nDestination Address: %s\nNumber of passengers car can accomodate: %d\nNumber of Vacancies Left: %d\nStatus: %s\n\n",
							trip.TripID, trip.FirstName, trip.PickupAddr, trip.AltPickupAddr, trip.StartTravelTime, trip.DestAddr, trip.MaxPassengers, trip.VacanciesLeft, trip.TripStatus)
					}
				} else {
					fmt.Println("Error decoding response:", err)
				}
			} else {
				fmt.Println("Error fetching published trips")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}

// manage passenger trip menu -> search for trip | enrol trip
func managePassengerTripMenu(user *User) {
	var choice int
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nManage your trips:")
	fmt.Println("1. Search for Trip")
	fmt.Println("2. Enrol Trip")
	fmt.Println("0. Go Back")
	reader.ReadString('\n')
	fmt.Print("Enter an option: ")
	fmt.Scanf("%d", &choice)

	switch choice {
	case 1:
		//search for an available trip -> all tripStatus waiting
		fmt.Println("----Search for Trip----")
		searchForTrip()
	case 2:
		//enrol trip
		fmt.Println("----Enrol Trip----")
		enrolInTrip(user)
	case 0:
		// Go back
	default:
		fmt.Println("Invalid option")
	}
}

// add the searchForTrip function, search for trips using destination address based on input
func searchForTrip() {
	var destination string
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter Destination Address: ")
	fmt.Scanf("%s", &destination)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:5002/api/v1/trips/search?destination=%s", destination), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusOK {
				var searchResults []Trip
				err := json.NewDecoder(res.Body).Decode(&searchResults)
				if err == nil {
					fmt.Println("Search Results:")
					for _, trip := range searchResults {
						fmt.Printf("\nTrip ID: %d\nDriver Name: %s\nPick-up Address: %s\nAlternative Pick-up Address: %s\nStart Travel Time: %s\nDestination Address: %s\nNumber of passengers car can accomodate: %d\nNumber of Vacancies Left: %d\nStatus: %s\n\n",
							trip.TripID, trip.FirstName, trip.PickupAddr, trip.AltPickupAddr, trip.StartTravelTime, trip.DestAddr, trip.MaxPassengers, trip.VacanciesLeft, trip.TripStatus)
					}
				} else {
					fmt.Println("Error decoding response:", err)
				}
			} else {
				fmt.Println("No trips found")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}

// Add the enrolInTrip function
func enrolInTrip(user *User) {
	var enrolmentData struct {
		TripID int `json:"tripId"`
		UserID int `json:"userId"`
	}
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter Trip ID: ")
	fmt.Scanf("%d", &enrolmentData.TripID)

	// Get the UserID from the currently logged-in user
	enrolmentData.UserID = user.UserID

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5002/api/v1/trips/enrol", bytes.NewBuffer([]byte(fmt.Sprintf(`{"tripId":%d,"userId":%d}`, enrolmentData.TripID, enrolmentData.UserID)))); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusCreated {
				fmt.Println("Enroled successfully!")
			} else {
				fmt.Println("You can't enrol for this trip as there are either no vacancies or date and time conflicts with your enroled trips")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}

// view past trips
func viewPastTripIDs(user *User) {
	// Get all enroled trip IDs
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:5002/api/v1/trips/enroltrips/%d", user.UserID), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusOK {
				var enrolDatas []enrolmentData
				err := json.NewDecoder(res.Body).Decode(&enrolDatas)
				if err == nil {
					for _, enrolData := range enrolDatas {
						//trip trip id = enroldata trip id then use trip id to get from trip table
						fetchAndPrintTripDetails(enrolData.TripID)
					}
				} else {
					fmt.Println("Error decoding response:", err)
				}
			} else {
				fmt.Println("Error fetching created trips")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}

func fetchAndPrintTripDetails(tripID int) {
	// Fetch detailed information about the trip from the trip table
	//debugging fmt.Printf("within fetch print Trip ID: %d\n", tripID)
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:5002/api/v1/trips/%d", tripID), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusOK {
				var trip Trip
				err := json.NewDecoder(res.Body).Decode(&trip)
				if err == nil {
					// Print detailed information about the trip
					fmt.Printf("\nTrip ID: %d\nDriver Name: %s\nPick-up Address: %s\nAlternative Pick-up Address: %s\nStart Travel Time: %s\nDestination Address: %s\nNumber of passengers car can accomodate: %d\nNumber of Vacancies Left: %d\nStatus: %s\n\n",
						trip.TripID, trip.FirstName, trip.PickupAddr, trip.AltPickupAddr, trip.StartTravelTime, trip.DestAddr, trip.MaxPassengers, trip.VacanciesLeft, trip.TripStatus)
				} else {
					fmt.Println("Error decoding response:", err)
				}
			} else {
				fmt.Println("Error fetching trip details")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}

// list all trips created by specific car owner
// Add the listCarOwnerTrips function
func listCarOwnerTrips(user *User) {
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:5002/api/v1/trips/carowners/%d", user.UserID), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusOK {
				var trips []Trip
				err := json.NewDecoder(res.Body).Decode(&trips)
				if err == nil {
					fmt.Println("Your Created Trips:")
					for _, trip := range trips {
						fmt.Printf("\nTrip ID: %d\nDriver Name: %s\nPick-up Address: %s\nAlternative Pick-up Address: %s\nStart Travel Time: %s\nDestination Address: %s\nNumber of passengers car can accomodate: %d\nNumber of Vacancies Left: %d\nStatus: %s\n\n",
							trip.TripID, trip.FirstName, trip.PickupAddr, trip.AltPickupAddr, trip.StartTravelTime, trip.DestAddr, trip.MaxPassengers, trip.VacanciesLeft, trip.TripStatus)
					}
				} else {
					fmt.Println("Error decoding response:", err)
				}
			} else {
				fmt.Println("Error fetching created trips")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}

func manageTripMenu() {
	var choice int
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nManage your trips:")
	fmt.Println("1. Start Trip")
	fmt.Println("2. Cancel Trip")
	fmt.Println("0. Go Back")
	reader.ReadString('\n')
	fmt.Print("Enter an option: ")
	fmt.Scanf("%d", &choice)

	switch choice {
	case 1:
		//start trip
		fmt.Println("----Start Trip----")
		startTrip()
	case 2:
		//cancel trip if 30mins before start time
		fmt.Println("----Cancel Trip----")
		//cancelTrip()
		selectCancelTrip()
	case 0:
		// Go back
	default:
		fmt.Println("Invalid option")
	}
}

// Add the startTrip function
func startTrip() {
	var tripID int
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter the Trip ID of the trip you want to start: ")
	fmt.Scanf("%d", &tripID)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:5002/api/v1/trips/start/%d", tripID), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusAccepted {
				fmt.Println("Trip started successfully!")
			} else {
				fmt.Println("Error starting trip")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}

// Add the cancelTrip function, within selectCancelTrip
func cancelTrip() {
	var tripID int
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter the same Trip ID again to confirm: ") //lol
	fmt.Scanf("%d", &tripID)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:5002/api/v1/trips/cancel/%d", tripID), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusAccepted {
				fmt.Println("Your trip has been cancelled successfully!")
			} else {
				fmt.Println("Error cancelling trip")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}

func selectCancelTrip() {
	var tripID int
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter the Trip ID of the trip you want to cancel: ")
	fmt.Scanf("%d", &tripID)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:5002/api/v1/trips/start-time/%d", tripID), nil); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()

			if res.StatusCode == http.StatusOK {
				var response map[string]string
				err := json.NewDecoder(res.Body).Decode(&response)
				if err == nil {
					startTimeStr := response["start_time"]
					//fmt.Printf("Start Time: %s\n", startTimeStr)

					// Parse the start time string
					startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
					if err != nil {
						fmt.Println("Error parsing start time:", err)
						return
					}

					// Calculate cancelBeforeTime (30 minutes before the start time)
					cancelBeforeTime := startTime.Add(-30 * time.Minute)
					currentTime := time.Now()
					//fmt.Printf("Cancel Before Time: %s\n", cancelBeforeTime.Format("2006-01-02 15:04:05"))
					//fmt.Printf("Current Time: %s\n", currentTime.Format("2006-01-02 15:04:05"))

					//if current time is 30 mins before scheduled time then cancel trip
					if currentTime.Before(cancelBeforeTime) {
						cancelTrip()
					} else {
						fmt.Println("Cancel Unsuccessful, trip can only be cancelled 30 mins before the scheduled time")
					}

				} else {
					fmt.Println("Error decoding response:", err)
				}
			} else {
				fmt.Println("Error fetching start time")
			}
		} else {
			fmt.Println("Error making request:", err)
		}
	} else {
		fmt.Println("Error creating request:", err)
	}
}
