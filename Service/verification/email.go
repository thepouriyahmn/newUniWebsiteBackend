package verification

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type Email struct{}

func NewEmail() Email {
	return Email{}
}
func (e Email) SendCode(reciever string) (string, error) {
	fmt.Println("using email")
	code, err := SendEmailVerificationCode(reciever)
	if err != nil {
		return "", err
	}
	return code, nil
}
func SendEmailVerificationCode(reciever string) (string, error) {
	code, err := GenerateSecureCode()
	if err != nil {
		fmt.Printf("reading error: %v", err)
		//panic(err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "pouriyahmn@gmail.com")
	m.SetHeader("To", reciever)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", "Your verification code is: "+code)

	d := gomail.NewDialer("smtp.gmail.com", 587, "pouriyahmn@gmail.com", "yezs zujy czwx xiew")

	err = d.DialAndSend(m)
	if err != nil {
		return "", err
	}

	return code, nil
}
