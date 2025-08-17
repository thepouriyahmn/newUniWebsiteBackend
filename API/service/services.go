package service

import "io"

type Services interface {
	SignUp(body io.Reader) error
}
