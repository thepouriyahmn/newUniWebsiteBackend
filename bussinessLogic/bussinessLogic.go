package bussinessLogic

import "net/http"

type IDatabase interface {
	//auth methodes
	CheackUserByUsernameAndEmail(ClientUsername, ClientEmail string) error
	InsertUser(username, pass, email string, studentRole, professorRole bool) error
	CheackUserByUserNameAndPassword(username, pass string) (int, string, error)
	GetRole(id int) ([]string, string, error)
	//management methodes
	GetAllProfessors() ([]Professor, error)
	AddProfessorById(userId int) error
	AddStudent(userId int) error
	GetAllUsers() ([]User, error)
	InsertLesson(lessonName string, lessonUnit int) error
	InsertClass(lessonName, professorName, date string, capacity, classNumber int) error
	DeleteLesson(lessonName string) error
	GetAllLessons() ([]Lesson, error)
	GetUsersByRole(roleId int) ([]User, error)
	AddMark(userId, classId int, mark *int) error
	GetStudentsForProfessor(professorId int) ([]Student, error)
	RemoveStudentUnit(id int) error
	GetAllClasses() ([]Classes, error)
	DeleteClass(classId int) error
	AddStudentById(userId int) error
}
type IProtocol interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Verify(w http.ResponseWriter, r *http.Request)
	GetAllProfessors(w http.ResponseWriter, r *http.Request)
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
type AuthBussinessLogic struct {
	IProtocol IProtocol
	IDatabase IDatabase
}

func NewBussinessLogic(protocol IProtocol, database IDatabase) AuthBussinessLogic {
	return AuthBussinessLogic{
		IProtocol: protocol,
		IDatabase: database,
	}

}

func (b AuthBussinessLogic) ShowProfessors() ([]Professor, error) {
	return b.IDatabase.GetAllProfessors()
}

func (b AuthBussinessLogic) AddProfessor(userId int) error {
	return b.IDatabase.AddProfessorById(userId)
}

func (b AuthBussinessLogic) ShowAllUsers() ([]User, error) {
	return b.IDatabase.GetAllUsers()
}

func (b AuthBussinessLogic) InsertLesson(lessonName string, lessonUnit int) error {
	return b.IDatabase.InsertLesson(lessonName, lessonUnit)
}

func (b AuthBussinessLogic) DeleteLesson(lessonName string) error {
	return b.IDatabase.DeleteLesson(lessonName)
}

func (b AuthBussinessLogic) ShowAllLessons() ([]Lesson, error) {
	return b.IDatabase.GetAllLessons()
}

func (b AuthBussinessLogic) ShowUsersByRole(roleId int) ([]User, error) {
	return b.IDatabase.GetUsersByRole(roleId)
}

func (b AuthBussinessLogic) AddMark(userId, classId int, mark *int) error {
	return b.IDatabase.AddMark(userId, classId, mark)
}

func (b AuthBussinessLogic) ShowStudentsForProfessor(professorId int) ([]Student, error) {
	return b.IDatabase.GetStudentsForProfessor(professorId)
}

func (b AuthBussinessLogic) RemoveStudentUnit(id int) error {
	return b.IDatabase.RemoveStudentUnit(id)
}
