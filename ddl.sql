CREATE TABLE users (
    ID SERIAL PRIMARY KEY,
    Email VARCHAR(255) NOT NULL UNIQUE,
    Password VARCHAR(255) NOT NULL,
    Deposit FLOAT DEFAULT 0
);

CREATE TABLE rental_playstations (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Availability VARCHAR(255) NOT NULL CHECK (Availability IN ('Tersedia', 'Tidak Tersedia')),
    rental_costs DECIMAL(10, 2) NOT NULL,
    Category VARCHAR(255) NOT NULL
);


CREATE TABLE history_rentals (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(ID),
  rental_id INT REFERENCES rental_playstations(ID),
  rental_date DATE,
  return_date DATE,
  rent_cost DECIMAL(10, 2), 
  status VARCHAR(255)
);


INSERT INTO rental_playstations (Name, Availability, rental_costs, Category)
VALUES ('PS3', 'Tersedia', 5.99, 'PlayStation 3'),
       ('PS4', 'Tersedia', 7.99, 'PlayStation 4'),
       ('PS5', 'Tersedia', 9.99, 'PlayStation 5');




