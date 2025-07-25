package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	// Gera 32 bytes aleatórios
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		panic(err)
	}

	// Codifica em base64 para ficar legível e seguro
	encoded := base64.StdEncoding.EncodeToString(secret)
	fmt.Printf("JWT_SECRET=%s\n", encoded)
}
