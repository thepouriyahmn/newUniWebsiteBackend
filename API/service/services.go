package service

import "io"

// ServiceResponse represents a normalized transport response
type ServiceResponse struct {
	StatusCode int
	Body       []byte
}

type Services interface {
	// Auth
	SignUp(body io.Reader) (ServiceResponse, error)
	Login(body io.Reader) (ServiceResponse, error)
	Verify(body io.Reader) (ServiceResponse, error)
	Logout(token string) (ServiceResponse, error)

	// Admin
	ShowProfessors(token string) (ServiceResponse, error)
	GetTerms(token string) (ServiceResponse, error)
	ShowAllUsers(token string, query string) (ServiceResponse, error)
	InsertLesson(token string, body io.Reader) (ServiceResponse, error)
	ShowAllLessons(token string) (ServiceResponse, error)
	DeleteLesson(token string, body io.Reader) (ServiceResponse, error)
	InsertClass(token string, body io.Reader) (ServiceResponse, error)
	ShowClasses(token string, query string) (ServiceResponse, error)
	DeleteClass(token string, body io.Reader) (ServiceResponse, error)
	ShowUsersByRole(token string, body io.Reader) (ServiceResponse, error)
	AddStudent(token string, body io.Reader) (ServiceResponse, error)
	AddProfessor(token string, body io.Reader) (ServiceResponse, error)

	// Professor
	AddMark(token string, body io.Reader) (ServiceResponse, error)
	ShowStudentsForProfessor(token string) (ServiceResponse, error)

	// Student
	AddStudentUnit(token string, body io.Reader) (ServiceResponse, error)
	PickedUnits(token string) (ServiceResponse, error)
	DelStudentUnit(token string, body io.Reader) (ServiceResponse, error)
}
