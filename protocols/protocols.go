package protocols

import "net/http"

type Protocols interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Verify(w http.ResponseWriter, r *http.Request)
	GetAllProfessors(w http.ResponseWriter, r *http.Request)
	AddProfessor(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	InsertLesson(w http.ResponseWriter, r *http.Request)
	DeleteLesson(w http.ResponseWriter, r *http.Request)
	GetAllLessons(w http.ResponseWriter, r *http.Request)
	GetUsersByRole(w http.ResponseWriter, r *http.Request)
	AddMark(w http.ResponseWriter, r *http.Request)
	GetStudentsForProfessor(w http.ResponseWriter, r *http.Request)
	AddStudentUnit(w http.ResponseWriter, r *http.Request)
	DelStudentUnit(w http.ResponseWriter, r *http.Request)
}
