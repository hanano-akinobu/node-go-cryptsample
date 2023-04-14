package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Message struct {
	Message string `json:"message"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/", func(c echo.Context) error {
		var m Message
		if err := c.Bind(&m); err != nil {
			println(err.Error())
		}
		println(m.Message)

		privateKeyFile := "../private-key.pem"

		// Read the private key file
		privateKeyBytes, err := ioutil.ReadFile(privateKeyFile)
		if err != nil {
			log.Fatalf("Error reading private key file: %v", err)
		}

		// Parse the private key from the PEM block
		privateKeyPem, _ := pem.Decode(privateKeyBytes)
		if privateKeyPem == nil {
			log.Fatalf("Error decoding private key PEM block")
		}

		privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyPem.Bytes)
		if err != nil {
			log.Fatalf("Error parsing private key: %v", err)
		}

		// Decode the base64-encoded encrypted message
		encryptedMessage, err := base64.StdEncoding.DecodeString(m.Message)
		if err != nil {
			log.Fatalf("Error decoding base64 message: %v", err)
		}

		// Decrypt the message using RSA-OAEP with SHA-256
		label := []byte("") // Optional label, set to an empty byte slice if not used
		decryptedMessage, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedMessage, label)
		if err != nil {
			log.Fatalf("Error decrypting message: %v", err)
		}

		// Now you have the decrypted message
		fmt.Printf("Decrypted message: %s\n", decryptedMessage)

		return c.JSON(http.StatusOK, Message{Message: "done"})
	})
	e.Logger.Fatal(e.Start(":1323"))
}
