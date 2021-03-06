package main

import "github.com/dgrijalva/jwt-go"
import "net/http"
import "encoding/json"

import "fmt"

func validateToken(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("x-Myheader", "token")
	res.Write([]byte("validate"))
}

func handleToken(res http.ResponseWriter, req *http.Request) {
	type MyCustomClaims struct {
		Admin bool   `json:"admin"`
		Name  string `json:"name"`
		jwt.StandardClaims
	}

	claims := MyCustomClaims{
		true,
		"Gokul",
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "Pulsecheck",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte("secret"))

	res.Header().Set("Authorization", fmt.Sprintf("Bearer %v", tokenString))

	tok, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tok.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if claims, ok := tok.Claims.(*MyCustomClaims); ok && tok.Valid {
		fmt.Println(claims.Admin, claims.Name)
	} else {
		fmt.Println(err)
	}

	res.Write([]byte(tokenString))
}

func handleHome(res http.ResponseWriter, req *http.Request) {
	type Product struct {
		Name string
	}
	prod := Product{
		Name: "Under Construction",
	}
	payload, _ := json.Marshal(prod)
	res.Write([]byte(payload))
}

func main() {
	http.HandleFunc("/validate", validateToken)
	http.HandleFunc("/token", handleToken)
	http.HandleFunc("/", handleHome)
	http.ListenAndServe(":8080", nil)
}
