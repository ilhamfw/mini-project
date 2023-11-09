package handler

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"rental-games/config"
	"rental-games/entity"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

const MaxDepositAmount = 100000.0

// @Summary Register a new user
// @Description Register a new user with the provided details.
// @ID register-user
// @Accept json
// @Produce json
// @Param user body entity.User true "User object to be registered"
// @Success 201 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users/register [post]
// RegisterUser handler untuk endpoint /users/register
func RegisterUser(c echo.Context) error {
	log := logrus.New()

	// Dapatkan data dari request
	user := new(entity.User)
	if err := c.Bind(user); err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("Invalid request data")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid request data"})
	}

	// Validasi data
	validator := validator.New() // Inisialisasi validator
	if err := validator.Struct(user); err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("Validation error")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Validation error"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("Error hashing password")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Error hashing password"})
	}

	if err := sendRegistrationEmail(user.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to send registration email"})
	}
	user.Password = string(hashedPassword)

	// Simpan data user ke database
	db := config.DB
	if err := db.Create(&user).Error; err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("Failed to create user")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to create user"})
	}

	// Kirim response
	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "User created successfully", "user": user})
}

func sendRegistrationEmail(recipientEmail string) error {
	// Konfigurasi email
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Membuat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", "ilham.fw18@gmail.com") // Ganti dengan alamat email pengirim
	m.SetHeader("To", recipientEmail)           // Menggunakan alamat email penerima dari parameter fungsi
	m.SetHeader("Subject", "Selamat bergabung!")
	m.SetBody("text/plain", "Selamat, pendaftaran Anda berhasil!")

	// Mengirim email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	// Email terkirim dengan sukses
	return nil
}

// @Summary Login user
// @Description Login a user with the provided email and password.
// @ID login-user
// @Accept json
// @Produce json
// @Param user body entity.User true "User object with email and password"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users/login [post]

// LoginUser handler untuk endpoint /users/login
func LoginUser(c echo.Context) error {
	log := logrus.New()

	// Dapatkan data dari request
	requestUser := new(entity.User)
	if err := c.Bind(requestUser); err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("Invalid request data")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid request data"})
	}

	// Dapatkan data user dari database berdasarkan email
	db := config.DB
	user := new(entity.User)
	if err := db.Where("email = ?", requestUser.Email).First(user).Error; err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("User not found")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "User not found"})
	}

	// Bandingkan password dengan hash yang ada di database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestUser.Password))
	if err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("Invalid password")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Invalid password"})
	}

	// Generate token JWT
	token, err := generateToken(user)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Error("Error generating token")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Error generating token"})
	}

	// Kirim response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
}

// @Summary Generate JWT Token
// @Description Generate a JWT token for the authenticated user.
// @ID generate-token
// @Param user body entity.User true "User object with valid credentials"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /generate-token [post]

func generateToken(user *entity.User) (string, error) {
	// Buat token JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		// Jika token valid, Anda dapat melanjutkan ke handler berikutnya
		c.Set("user", claims)
		return next(c)
	}
}
