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
   - Passengers can search for trips based on destination and enrol in available trips with vacancies and no datetime clashes with their other enroled trips
     
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
   - Step 1: Create an Account and input First Name, Last Name, Mobile Number, Email Address<br>
             (Select option 1 from the start menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/9558b9d6-a0a9-4c2b-8e9a-cf2a7ca4b3b3)
     
   - Step 2: Login using newly created account's email<br>
             (Select option 2 from the start menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/115624e2-9ad3-426a-863e-b70b42391b60)
     
   - Step 3: Delete Account (Unsuccessful Deletion Message as Account isn't 1 year old)<br>
             (Select option 3 from the Passenger Main Menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/ebc92025-7544-406c-96ae-1baa6020270f)
     
   - Step 4: Exit and Log back in with 'del@gmail.com'<br>
             (Select option 0 'Exit', run 'go run admin.go', select option 2 'Login')<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/f1e3ce9e-66eb-492f-bf26-b998ec915d44)
      
   - Step 5: Delete Del's Account (Successful Deletion as Account is older than 1 year)<br>
             (Select option 3 from the Passenger Main Menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/3790b377-7456-4cf0-b734-88f2d93caedc)
  
   - Step 6: Log back in with created account's email<br>
             (Run 'go run admin.go', select option 2 'Login')<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/1b392340-4046-4d39-90c8-ddda9d144bea)
     
   - Step 7: Update Profile Information and input new profile information<br>
             (Select option 2 from the Passenger Main Menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/6e6e3214-7565-40c9-9dae-ae3262ec37c7)
     
   - Step 8: Upgrade to Car Owner and input Driver's License Number and Car Plate Number. Main menu changes to Car Owner Main Menu<br>
             (Select option 1 from the Passenger Main Menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/78a3d154-e359-46eb-b183-bd93f5cc6e7f)
     
   - Step 9: Update Profile Information and input new profile information<br>
             (Select option 1 from the Car Owner Main Menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/06a08a53-62ff-4bb0-b0ef-3149ca0016d4)
     
   - Step 9: Publish new trip and input trip details<br>
             (Select option 3 from the Car Owner Main Menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/3410fab7-d0c6-4420-be9a-53f856f35431)
     
   - Step 10: Publish another new trip and input trip details<br>
             (Select option 3 from the Car Owner Main Menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/63176064-e580-4b94-82ee-e19fc46025b7)

   - Step 11: List all created trips<br>
             (Select option 4 from the Car Owner Main Menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/07e740c8-75ab-4542-aa25-dff4bd263dde)
     
   - Step 12: Start Trip and input TripID of trip that you want to start<br>
             (Select option 1 from the manage trip menu -> from list all created trips)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/a01afb4c-6dd9-4b61-b69d-bd909acac409)
     
   - Step 13: Cancel Trip and input TripID of trip that you want to cancel<br>
             (Select option 4 from Car Owner Main Menu, then Select option 2 from the manage trip menu -> from list all created trips)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/e24da407-3cc2-43bb-bc07-891d44e52136)
     
   - Step 14: Exit/Logout and Log back in using 'rob@gmail.com'<br>
             (select option 0 to exit, run 'go run admin.go', select option 2 'Login')<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/9ce67981-a0d3-4790-a1a7-69972079576f)
     
   - Step 15: List all available/published trips (only trips with 'Waiting' status will be shown)<br>
             (Select option 4 from the Passenger main menu)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/9357cdfa-ecb0-4d49-9dea-d9bbad0177b7) <br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/845f6c7f-d384-498c-b8ca-f780850666e8)
    
   - Step 16: Search for Trip by Destination Address<br>
             (Select option 1 and search 'plaza')<br>
 ![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/079b6fef-4a89-4a90-8c3b-345d5bc36995)
    
   - Step 17: List all published trips and Select Enrol Trip, input TripID 9 as the trip you want to enrol<br>
             (Select option 4 from Passenger main menu then option 2 from List all Published Trips)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/2bf09dce-1df9-447b-99c6-606723809231)

   - Step 18: Enrol Trip again using TripID 10 to show the unsuccessful enrolment message due to date and time conflicts with user's enroled trips<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/5fde0144-a803-4915-ab65-c7198ac99b9d)
  
   - Step 19: Enrol Trip again using TripID 8 to show the unsuccessful enrolment message due to no vacancies available<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/4ce3c4b4-03cc-44d3-b078-50af75fc80b8)
  
   - Step 20: Exit/Logout and Log back in using 'johnlim@gmail.com' (As John is a Car Owner after Log in it shows Car Owner Main Menu)<br>
             (select option 0 to exit, run 'go run admin.go', select option 2 'Login')<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/ad4d71b7-0a51-482f-b20d-b4110df5153e)
     
   - Step 21: List all created trips and Cancel Trip with TripID 6 to show the unsuccessful cancelation message as trip can only be cancelled 30 minutes before the scheduled time<br>
             (Select option 4 from Car Owner Main Menu, Select option 2 from the manage trip menu -> from list all created trips)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/c2c81e79-a8f5-4100-8a8f-33a41399c6a6)

   - Step 22: Exit/Logout and Log back in using 'bach@gmail.com' to demonstrate delete using Car Owner Account<br>
             (select option 1 to exit, run 'go run admin.go', select option 2 'Login')<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/9c122aa0-f1c1-40bf-8d06-676b3d78a67d)

   - Step 23: Select Delete Account, after deletion the program will exit<br>
             (select option 2 from Car Owner Main menu to Delete)<br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/051c933b-c33b-4309-a147-6b4deded5376)

- From the User Table in MySQL you can see that the deleted users Del and Bach are gone <br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/c58dcb85-b2df-4f96-8243-82b953026f23)
- The updated Trip table <br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/d26478c3-d959-4846-85f7-883b1a77e720)
- The updated TripEnrolment table <br>
![image](https://github.com/ng-zi-yi/ETI-ASG1/assets/93900155/f867f641-a0a7-4899-be44-6f16fc5afe2b)




