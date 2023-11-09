package handler

import (
	"fmt"
	"net/http"
	"rental-games/config"
	"rental-games/entity"
	"os"
	"context"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	xendit "github.com/xendit/xendit-go/v3"
    invoice "github.com/xendit/xendit-go/v3/invoice"
)

func DepositAmount(c echo.Context) error {
	log := logrus.New()

	// Dapatkan token dari header Authorization
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		log.Error("Authorization token is required")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Authorization token is required"})
	}

	// Validasi token JWT
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	userId := claims["id"].(float64)
	fmt.Println(userId)

	if err != nil || !jwtToken.Valid {
		log.Error("Invalid or expired token")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Invalid or expired token"})
	}
	fmt.Println(userId)
	// Dapatkan data dari request
	depositData := struct {
		Deposit float32 `json:"deposit"`
	}{}
	if err := c.Bind(&depositData); err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("Invalid request data")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid request data"})
	}
	// Dapatkan data user dari database berdasarkan ID
	db := config.DB
	user := new(entity.User)
	if err := db.Where("id = ?", userId).First(user).Error; err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("User not found")
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "User not found"})
	}

	// Lakukan validasi depositAmount
	if depositData.Deposit <= 0 {
		log.Error("Invalid deposit amount")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid deposit amount"})
	}

	// Pastikan deposit tidak melebihi batasan tertentu
	if depositData.Deposit > MaxDepositAmount {
		log.Error("Deposit amount exceeds maximum limit")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Deposit amount exceeds maximum limit"})
	}

	// Kirim permintaan ke API Xendit
	if err := sendDepositToXendit(c, userId, depositData.Deposit); err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("Gagal mengirim deposit ke Xendit")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Gagal mengirim deposit ke Xendit"})
	}

	

	// Lakukan update deposit amount pada data user
	user.Deposit += float64(depositData.Deposit)

	if err := db.Save(user).Error; err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("Failed to update deposit amount")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to update deposit amount"})
	}

	// Kirim response dengan data user yang telah diupdate
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Deposit successful",
		"user":    user,
	})
}

func sendDepositToXendit(c echo.Context, userId float64, depositAmount float32) error {
	judulInvoice := fmt.Sprintf("Invoice order user id = %v", userId)
	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(judulInvoice, depositAmount)

	xenditClient := xendit.NewClient(os.Getenv("XENDIT_APIKEY"))

	resp, r, err := xenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", err.Error())

		b, _ := json.Marshal(err.FullError())
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from CreateInvoice: Invoice
	fmt.Fprintf(os.Stdout, "Response from InvoiceApi.CreateInvoice: %v\n", resp)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "sukses xendit",
		"respon":  resp,
	})
}
