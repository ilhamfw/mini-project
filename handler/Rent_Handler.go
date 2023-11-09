package handler

import (
	"fmt"
	"net/http"
	"rental-games/entity"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"time"
)

// @Summary Rent a PlayStation
// @Description Rent a PlayStation based on the provided parameters.
// @ID rent-console
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization Token"
// @Param console_id formData integer true "ID of the PlayStation to rent"
// @Param rental_date formData string true "Rental date (YYYY-MM-DD)"
// @Param return_date formData string true "Return date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Success response"
// @Router /rent [post]
func GetAvailableConsoles(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Authorization token
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Authorization token is required"})
		}

		claims := jwt.MapClaims{}
		jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})

		if err != nil || !jwtToken.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Invalid or expired token"})
		}

		// Menggunakan GORM untuk mengambil daftar PS yang tersedia
		var consoles []entity.RentalPlaystation
		if err := db.Where("Availability = ?", "Tersedia").Find(&consoles).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Gagal mengambil data PS"})
		}

		return c.JSON(http.StatusOK, consoles)
	}
}

// @Summary Get Available Consoles
// @Description Get a list of available PlayStation consoles.
// @ID get-available-consoles
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization Token"
func RentConsole(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mengambil data pengguna dari token (pastikan Anda telah mengimplementasikan fitur otentikasi dan token)
		user := c.Get("user")
		userId := user.(jwt.MapClaims)["id"].(float64)
		/*
			{
				"console_id": 1
				"rental_date": 2023-10-10
				"return_date": 2023-10-10
			}
		*/

		var userRequestOrder entity.UserRequest
		if err := c.Bind(&userRequestOrder); err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Konsol tidak ditemukan atau tidak tersedia"})
		}
		
		// Menggunakan GORM untuk mengubah status ketersediaan konsol dalam database
		var console entity.RentalPlaystation
		if err := db.Where("ID = ? AND Availability = ?", userRequestOrder.ConsoleID, "Tersedia").First(&console).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Konsol tidak ditemukan atau tidak tersedia"})
		}

		// Menghitung biaya sewa berdasarkan tanggal
		// rentalDateStr := c.FormValue("rental_date")
		// returnDateStr := c.FormValue("return_date")

		rentalDate, err := time.Parse("2006-01-02", userRequestOrder.RentalDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Format tanggal sewa tidak valid"})
		}

		returnDate, err := time.Parse("2006-01-02", userRequestOrder.ReturnDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Format tanggal pengembalian tidak valid"})
		}

		// Menghitung selisih hari antara tanggal sewa dan tanggal pengembalian
		days := returnDate.Sub(rentalDate).Hours() / 24

		fmt.Printf("harga rental per hari: %v\n", console.RentalCosts)
		fmt.Printf("total hari: %v\n",days)
		// Menghitung biaya berdasarkan harga sewa per hari
		rentalCosts := console.RentalCosts * float64(days)
		fmt.Printf("Subtotal: %v\n", rentalCosts)

		// Ubah status ketersediaan menjadi "Tidak Tersedia"
		if err := db.Model(&console).Update("Availability", "Tidak Tersedia").Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Gagal melakukan penyewaan"})
		}

		// Membuat entri HistoryRental
		historyRental := entity.HistoryRental{
			UserID:     int(userId),    // Menggunakan ID pengguna yang telah diambil dari token
			RentalID:   console.ID, // Menggunakan ID konsol yang telah ditemukan
			RentalDate: userRequestOrder.RentalDate,
			ReturnDate: userRequestOrder.ReturnDate,
			RentCost:   rentalCosts, // Menggunakan biaya sewa yang telah dihitung
			Status:     "Tersedia",  // Ubah status sesuai kebutuhan
		}

		// Simpan entri HistoryRental ke dalam basis data
		if err := db.Create(&historyRental).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Gagal menyimpan riwayat penyewaan"})
		}

		// Return response dengan biaya sewa
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "PS berhasil disewa", "rental_cost": rentalCosts})
		// return c.JSON(http.StatusOK, map[string]interface{}{"message": "PS berhasil disewa"})
	}
}
