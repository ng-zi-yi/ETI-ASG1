# ETI-ASG1

## Design Consideration of my microservices
1. User Microservice
   - Manages user accounts, user information, and account-related operations
   - Utilizes the User Table in the 'carpooling_db' MySQL database to store user information
   - Features include account creation, login, profile update, allowing passengers to upgrade to car owner account, and account deletion

2. Trip Microservice
   - Handles the creation, management, and status of carpooling trips
   - Utilizes the Trip Table in the 'carpooling_db' MySQL database to store trip-related data
   - Features include trip creation, listing all available trips (Passenger side), listing all created trips (Car Owner side), starting a trip, cancelling a trip, searching for a trip, enroling in a trip, and viewing past trips
   - Car Owners can publish new trips, specifying details such as pickup locations, destination, timings, maximum passengers
   - Car Owners can start or cancel trips based on specific conditions from viewing all of their created trips
   - Passengers can search for all available trips based on destination and enrol in available trips with vacancies and no datetime clashes with their other enroled trips
     
3.  Overall
    - Communication between microservices is implemented via HTTP RESTful APIs
    - Data Consistency is acheived through the use of the relational MySQL database 'carpooling_db'
    - Console Applcation developed in Go supports both passenger and car owner functionalities

## Architecture Diagram
![eti asg 1 architecture diagram](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/b7fa8f73-00fc-4ab8-86fd-89ddc9d2e435)
   - Both User and Trip microservices share a common MySQL database 'carpooling_db'
   - The User Microservice uses the User table in the database
   - The Trip Microservice uses the Trip table and TripEnrolment table in the database
   - Each microservice communicates with its respective table
   - The Console Application 'admin.go' interacts with both microservices through their APIs

## Instructions for setting up and running my microservices
1. Database Setup
   - Execute the SQL script 'carpooling_db.sql' to create the required database and tables (User, Trip, TripEnrolment)


2. Microservices Setup
   - After cloning the repository, navigate to the root directory each microservice and run 'go run user.go' and 'go run trip.go'.
   - User microservice runs on port 5001, and Trip microservice runs on port 5002.
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/7b598199-9572-49a8-90a1-379ec2ef2f45)
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/58d43d8d-b165-4bf2-b1d7-6a478c7c0f83)

3. Console Application (Demo)
   - At the console application 'admin.go' run 'go run admin.go' to interact with the microservices
   - Throughout the demo, can check if there have been updates or changes to the data in the MySQL database. While the console application provides real-time interactions, can manually inspect the database for changes.
   - Step 1: Create an Account and input First Name, Last Name, Mobile Number, Email Address
             (Select option 1 from the start menu)
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/9558b9d6-a0a9-4c2b-8e9a-cf2a7ca4b3b3)
     
   - Step 2: Login using newly created account's email<br>
             (Select option 2 from the start menu)
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/115624e2-9ad3-426a-863e-b70b42391b60)
     
   - Step 3: Delete Account (Unsuccessful Deletion Message as Account isn't 1 year old)<br>
             (Select option 3 from the Passenger Main Menu)
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/ebc92025-7544-406c-96ae-1baa6020270f)
     
   - Step 4: Exit and Log back in with 'del@gmail.com'
             (Select option 0 'Exit', run 'go run admin.go', select option 2 'Login')<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/f1e3ce9e-66eb-492f-bf26-b998ec915d44)
      
   - Step 5: Delete Del's Account (Successful Deletion as Account is older than 1 year)<br>
             (Select option 3 from the Passenger Main Menu)
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/3790b377-7456-4cf0-b734-88f2d93caedc)
  
   - Step 6: Log back in with created account's email<br>
             (Run 'go run admin.go', select option 2 'Login')
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/1b392340-4046-4d39-90c8-ddda9d144bea)
     
   - Step 7: Update Profile Information and input new profile information
             (Select option 2 from the Passenger Main Menu)
     
   - Step 8: Upgrade to Car Owner and input Driver's License Number and Car Plate Number
             (Select option 1 from the Passenger Main Menu)
     
   - Step 9: Update Profile Information and input new profile information
             (Select option 1 from the Car Owner Main Menu)
     
   - Step 9: Publish new trip and input trip details
             (Select option 3 from the Car Owner Main Menu)
     
   - Step 10: Publish another new trip and input trip details
             (Select option 3 from the Car Owner Main Menu)

   - Step 11: List all created trips
             (Select option 4 from the Car Owner Main Menu)
     
   - Step 12: Start Trip and input TripID of trip that you want to start
             (Select option 1 from the manage trip menu -> from list all created trips)
     
   - Step 13: Cancel Trip and input TripID of trip that you want to cancel
             (Select option 2 from the manage trip menu -> from list all created trips)
     
   - Step 14: Exit/Logout and Log back in using 'rob@gmail.com'
             (select option 1 to exit, run 'go run admin.go', select option 2 'Login')
     
   - Step 15: List all available/published trips (only trips with 'Waiting' status will be shown)
             (Select option 4 from the Passenger main menu)
     
   - Step 16: Search for Trip by Destination Address
             (Select option 1 and search 'plaza' or 'ee')
     
   - Step 17: List all published trips and Select Enrol Trip, input TripID 9 as the trip you want to enrol (
             (Select option 4 from Passenger main menu then option 2 from List all Published Trips)

   - Step 18: Enrol Trip again using TripID 10 to show the unsuccessful enrolment message due to date and time conflicts with user's enroled trips
  
   - Step 19: Enrol Trip again using TripID 8 to show the unsuccessful enrolment message due to no vacancies available
  
   - Step 20: Exit/Logout and Log back in using 'johnlim@gmail.com' (As John is a Car Owner after Log in it shows Car Owner Main Menu)
             (select option 1 to exit, run 'go run admin.go', select option 2 'Login')
     
   - Step 21: List all created trips and Cancel Trip with TripID 6 to show the unsuccessful cancelation message as trip can only be cancelled 30 minutes before the scheduled time
             (Select option 4 from Car Owner Main Menu, Select option 2 from the manage trip menu -> from list all created trips)

   - Step 22: Exit/Logout and Log back in using 'bach@gmail.com' to demonstrate delete using Car Owner Account
             (select option 1 to exit, run 'go run admin.go', select option 2 'Login')

   - Step 23: Select Delete Account, after deletion the program will exit
             (select option 2 from Car Owner Main menu to Delete)





