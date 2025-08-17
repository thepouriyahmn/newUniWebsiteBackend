package verification

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type ISendVerificationCode interface {
	SendCode(reciever string) (string, error)
}

func GenerateSecureCode() (string, error) {
	max := big.NewInt(1000000) // ۶ رقم
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
