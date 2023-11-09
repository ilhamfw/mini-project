package entity

// User represents a user of the system
type User struct {
	ID       int     `json:"id"`
	Email    string  `json:"email" validate:"required"`
	Password string  `json:"password" validate:"required"`
	Deposit  float64 `json:"deposit"`
}

type UserRequest struct {
	ConsoleID  int    `json:"console_id"`
	RentalDate string `json:"rental_date"`
	ReturnDate string `json:"return_date"`
}

// RentalPlaystation represents a rental PlayStation
type RentalPlaystation struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Availability string  `json:"availability"` // Nilai yang mungkin: "Tersedia" atau "Tidak Tersedia"
	RentalCosts  float64 `json:"rentalcosts"`
	Category     string  `json:"category"`
}

// HistoryRental represents a rental history
type HistoryRental struct {
	ID         int     `json:"id" gorm:"primaryKey"`
	UserID     int     `json:"user_id"`
	RentalID   int     `json:"rental_id"`
	RentalDate string  `json:"rental_date"`
	ReturnDate string  `json:"return_date"`
	RentCost   float64 `json:"rent_cost" gorm:"type:numeric(10,2)"`
	Status     string  `json:"status"`
}
