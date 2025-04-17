package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

func auth(id string, cmd Command) bool {
	k := K(cmd.Key)
	isOwner := k.Id == id
	isPublicGet := k.Pub && (cmd.Op == CmdGet || cmd.Op == CmdList)
	authorized := isOwner || isPublicGet
	newIdClaim := !authorized && cmd.Op == CmdSet && k.Pub && k.Last == "key"
	authorized = authorized || (newIdClaim && !idExists(k.Id))
	return authorized
}

func idExists(id string) bool {
	getPubKeyCmd := NewCommand(CmdGet, fmt.Sprintf("%s/pub/key", id), map[string]string{}, []byte{})
	data, err := getPubKeyCmd.Exec()
	return err == nil && len(data) > 0
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func parseToken(token, secret string) (Claims, error) {
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
