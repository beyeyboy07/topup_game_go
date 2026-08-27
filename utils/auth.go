package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-this-secret-in-production"
	}
	return []byte(secret)
}

func CreateToken(userID uint, role string) (string, error) {
	header := encodeTokenPart(map[string]string{"alg": "HS256", "typ": "JWT"})
	payload := encodeTokenPart(map[string]interface{}{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})
	unsigned := header + "." + payload
	return unsigned + "." + sign(unsigned), nil
}

func ParseToken(tokenString string) (uint, string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 || !hmac.Equal([]byte(parts[2]), []byte(sign(parts[0]+"."+parts[1]))) {
		return 0, "", errors.New("invalid token")
	}
	var claims struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role"`
		Exp    int64  `json:"exp"`
	}
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || json.Unmarshal(decoded, &claims) != nil || claims.UserID == 0 || claims.Role == "" || time.Now().Unix() >= claims.Exp {
		return 0, "", errors.New("invalid token claims")
	}
	return claims.UserID, claims.Role, nil
}

func encodeTokenPart(value interface{}) string {
	data, _ := json.Marshal(value)
	return base64.RawURLEncoding.EncodeToString(data)
}

func sign(value string) string {
	hash := hmac.New(sha256.New, jwtSecret())
	hash.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}
