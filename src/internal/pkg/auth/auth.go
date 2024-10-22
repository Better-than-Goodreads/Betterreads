package auth


import (
    "os"
    "fmt"
    "time"
    "strconv"
    "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret          string
	jwtExpirationHours int
)

func init() {
	jwtSecret = os.Getenv("JWT_SECRET")
	jwtExpirationHoursStr := os.Getenv("JWT_DURATION_HOURS")

	if jwtExpirationHoursStr != "" {
		jwtExpirationHours, _ = strconv.Atoi(jwtExpirationHoursStr)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

type Claims struct {
    IsAuthor bool `json:"is_author"`
    UserId    string `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateToken(userId string, isAuthor bool) (string, error) {
    fmt.Printf("is author: %v", isAuthor)
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is not set")
	}
	if jwtExpirationHours == 0 {
		return "", fmt.Errorf("JWT_DURATION_HOURS environment variable is not set or invalid")
	}

	claims := Claims{
        IsAuthor: isAuthor,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Beterreads-monke-crack",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sign, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error signing JWT token: %w", err)
	}

	return sign, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing JWT token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}
