package utils

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func DashboardInformation(tokenString string) (string, error) {
	var dashboardData string
	fmt.Println("Dashboard Information:", tokenString)
	dashboardData = "Dashboard Information"
	return dashboardData, nil
}

func ExtractUserIDFromToken(tokenString string) (string, error) {
	claims := &jwt.MapClaims{}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid JWT token")
	}

	// Extract userID from claims
	userID, ok := (*claims)["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return userID, nil
}
