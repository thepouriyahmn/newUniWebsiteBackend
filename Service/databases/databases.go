package databases

import "UniWebsite/bussinessLogic"

type IDatabase interface {
	// متدهای احراز هویت و ثبت‌نام
	// CheackUserByUsernameAndEmail(ClientUsername, ClientEmail string) error
	// InsertUser(username, pass, email string, studentRole, professorRole bool) error
	// CheackUserByUserNameAndPassword(username, pass string) (int, string, error)
	// GetRole(id int) ([]string, string, error)

	// // متدهای مدیریتی و آموزشی
	// GetAllProfessors() ([]bussinessLogic.Professor, error)
	// AddProfessorById(userId int) error
	// AddStudent(userId int) error
	// GetAllUsers(input string) ([]bussinessLogic.User, error)
	// InsertLesson(lessonName string, lessonUnit int) error
	// InsertClass(lessonName, professorName, date, term string, capacity, classNumber int) error
	// GetAllClassesByTerm(term string) ([]bussinessLogic.Classes, error)
	// DeleteClass(classId int) error
	// GetAllTerms() ([]string, error)
	// DeleteLesson(lessonName string) error
	// GetAllLessons() ([]bussinessLogic.Lesson, error)
	// GetUsersByRole(roleId int) ([]bussinessLogic.User, error)
	// AddStudentById(userId int) error
	// AddMark(userId, classId int, mark int) error
	// GetStudentsForProfessor(professorId int) ([]bussinessLogic.Student2, error)
	// RemoveStudentUnit(classid int, userid int) error
	// GetClassesByUserId(userID int) ([]bussinessLogic.StudentClasses, error)
	// InsertUnitForStudent(userid int, classid int) error
	MainDatabase
	SubDatabase
}

type MainDatabase interface {
	// Authentication methods (available in MongoDB)
	CheackUserByUsernameAndEmail(ClientUsername, ClientEmail string) error
	InsertUser(username, pass, email string, studentRole, professorRole bool) error
	CheackUserByUserNameAndPassword(username, pass string) (int, string, error)
	GetRole(id int) ([]string, string, error)

	// Write operations for CQRS (available in MongoDB)
	InsertLesson(lessonName string, lessonUnit int) error
	InsertClass(lessonName, professorName, date, term string, capacity, classNumber int) error
	DeleteClass(classId int) error
	DeleteLesson(lessonName string) error
	AddProfessorById(userId int) error
	AddStudentById(userId int) error
	AddMark(userId, classId int, mark int) error
	RemoveStudentUnit(classid int, userid int) error
	InsertUnitForStudent(userid int, classid int) error
}

type SubDatabase interface {
	// Read operations for CQRS (available in MySQL)
	GetAllProfessors() ([]bussinessLogic.Professor, error)
	GetAllUsers(input string) ([]bussinessLogic.User, error)
	GetAllLessons() ([]bussinessLogic.Lesson, error)
	GetAllClassesByTerm(term string) ([]bussinessLogic.Classes, error)
	GetAllTerms() ([]string, error)
	GetUsersByRole(roleId int) ([]bussinessLogic.User, error)
	GetStudentsForProfessor(professorId int) ([]bussinessLogic.Student2, error)
	GetClassesByUserId(userID int) ([]bussinessLogic.StudentClasses, error)

	// Sync operations (called by NATS handlers)
	CheackUserByUsernameAndEmail(ClientUsername, ClientEmail string) error
	InsertUser(username, pass, email string, studentRole, professorRole bool) error
	CheackUserByUserNameAndPassword(username, pass string) (int, string, error)
	GetRole(id int) ([]string, string, error)
	InsertLesson(lessonName string, lessonUnit int) error
	InsertClass(lessonName, professorName, date, term string, capacity, classNumber int) error
	DeleteClass(classId int) error
	DeleteLesson(lessonName string) error
	AddProfessorById(userId int) error
	AddStudentById(userId int) error
	AddMark(userId, classId int, mark int) error
	RemoveStudentUnit(classid int, userid int) error
	InsertUnitForStudent(userid int, classid int) error
}
