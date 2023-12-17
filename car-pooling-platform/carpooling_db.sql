CREATE database IF NOT EXISTS carpooling_db;

USE carpooling_db;
CREATE TABLE User (
UserID int auto_increment PRIMARY KEY,
FirstName varchar (20) NOT NULL,
LastName varchar (20) NOT NULL,
MobileNo varchar (8) NOT NULL, 
Email varchar (50) NOT NULL, 
DriverLicenseNo varchar (10),
CarPlateNo varchar (8),
UserType ENUM('Passenger', 'CarOwner') NOT NULL,
AccCreationDate timestamp default current_timestamp
);

CREATE TABLE Trip (
TripID int auto_increment PRIMARY KEY,
UserID int,
FirstName varchar (20) NOT NULL,
PickupAddr varchar (100) NOT NULL,
AltPickupAddr varchar (100),
StartTravelTime varchar (100) NOT NULL,
DestAddr varchar (100) NOT NULL,
MaxPassengers int NOT NULL,
VacanciesLeft int NOT NULL,
TripStatus ENUM('Waiting', 'Started', 'Cancelled'), 
TripCreationDate timestamp default current_timestamp,
foreign key (UserID) references User(UserID)
);

CREATE TABLE TripEnrolment (
    EnrolmentID int AUTO_INCREMENT PRIMARY KEY,
    TripID int,
    UserID int,
    FOREIGN KEY (TripID) REFERENCES Trip(TripID),
    FOREIGN KEY (UserID) REFERENCES User(UserID)
);

INSERT INTO User (UserID, FirstName, LastName, MobileNo, Email, DriverLicenseNo, CarPlateNo, UserType, AccCreationDate)
VALUES(1,"Sally", "Tan", "98765432", "sallytan@gmail.com", "S9876543A", "SAL9876A", "CarOwner", "2019-01-12 12:00:00"),
(2,"Bill", "Wong", "87654321", "billwong@gmail.com", "", "", "Passenger", "2019-03-18 03:00:00"),
(3,"John", "Lim", "76543210", "johnlim@gmail.com", "S7654321B", "SJO7654H", "CarOwner", "2020-08-19 08:20:00"),
(4,"Del", "Lee", "11223344", "del@gmail.com", "", "", "Passenger", "2020-10-11 13:00:00"),
(5,"Bach", "Sim", "22334455", "bach@gmail.com", "S2233445B", "SBA2233C", "CarOwner", "2020-11-19 20:20:00"),
(6,"Robert", "Wee", "67291234", "rob@gmail.com", "", "", "Passenger", "2021-01-24 15:00:00");


INSERT INTO Trip (TripID, UserID, FirstName, PickupAddr, AltPickupAddr, StartTravelTime, DestAddr, MaxPassengers, VacanciesLeft, TripStatus, TripCreationDate)
VALUES(1 ,1, "Sally", "Hougang1", "Hougang2", "2019-03-11 10:30:00", "Nex", 4, 3, "Started", "2019-03-02 10:00:00"),
(2, 3, "John", "Kranji1", "Kranji2", "2019-05-23 11:30:00", "JurongLake", 3, 1, "Started", "2019-05-13 20:00:00"),
(3, 3, "John", "Changi1", "Changi2", "2020-01-25 12:30:00", "JCube", 4, 3, "Started", "2020-01-21 14:00:00"),
(4, 1, "Sally", "Clementi1", "Clementi2", "2020-02-28 13:30:00", "YewTee1", 4, 2, "Started", "2020-02-21 09:00:00"),
(5, 1, "Sally", "Geylang1", "Geylang2", "20121-11-25 14:30:00", "Pioneer1", 4, 3, "Started", "2021-11-18 08:00:00");

INSERT INTO TripEnrolment (EnrolmentID, TripID, UserID)
VALUES(1, 1, 2),
(2, 2, 2),
(3, 2, 6),
(4, 3, 2),
(5, 4, 2),
(6, 4, 6),
(7, 5, 2);