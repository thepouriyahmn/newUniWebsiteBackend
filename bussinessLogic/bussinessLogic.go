package bussinessLogic

import "net/http"

type Bussinesslogic struct {
	IProtocol IProtocol
}

func NewBussinessLogic(protocol IProtocol) Bussinesslogic {
	return Bussinesslogic{
		IProtocol: protocol,
	}

}

type IProtocol interface {
	ShowPickedUnitsForStudent(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Verify(w http.ResponseWriter, r *http.Request)
	GetAllProfessors(w http.ResponseWriter, r *http.Request)
	GetTerms(w http.ResponseWriter, r *http.Request)
	AddProfessor(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	InsertLesson(w http.ResponseWriter, r *http.Request)
	InsertClass(w http.ResponseWriter, r *http.Request)
	DeleteLesson(w http.ResponseWriter, r *http.Request)
	GetAllLessons(w http.ResponseWriter, r *http.Request)
	GetUsersByRole(w http.ResponseWriter, r *http.Request)
	AddMark(w http.ResponseWriter, r *http.Request)
	GetStudentsForProfessor(w http.ResponseWriter, r *http.Request)
	AddStudentUnit(w http.ResponseWriter, r *http.Request)
	DelStudentUnit(w http.ResponseWriter, r *http.Request)
	ShowClasses(w http.ResponseWriter, r *http.Request)
	DeleteClass(w http.ResponseWriter, r *http.Request)
	AddStudent(w http.ResponseWriter, r *http.Request)
}
