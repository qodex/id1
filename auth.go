package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

func auth(id string, cmd Command) bool {
	isOwner := cmd.Key.Id == id
	isPublicGet := cmd.Key.Pub && (cmd.Op == Get || cmd.Op == List)
	authorized := isOwner || isPublicGet
	isNewIdClaim := !authorized && (cmd.Op == Set && cmd.Key.Pub && cmd.Key.Last == "key")
	exists := idExists(cmd.Key.Id)
	authorized = authorized || (isNewIdClaim && !exists)
	return authorized
}

func idExists(id string) bool {
	pubKey := KK(id, "pub", "key")
	data, err := CmdGet(pubKey).Exec()
	return err == nil && len(data) > 0
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func validateToken(token, secret string) (Claims, error) {
	claims := Claims{}
	if jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}); err != nil {
		return claims, err
	} else if !jwtToken.Valid {
		return claims, fmt.Errorf("invalid token")
	} else {
		return claims, nil
	}
}
