package databases

import "UniWebsite/bussinessLogic"

type IDatabase interface {
	// متدهای احراز هویت و ثبت‌نام
	CheackUserByUsernameAndEmail(ClientUsername, ClientEmail string) error
	InsertUser(username, pass, email string, studentRole, professorRole bool) error
	CheackUserByUserNameAndPassword(username, pass string) (int, string, error)
	GetRole(id int) ([]string, string, error)

	// متدهای مدیریتی و آموزشی
	GetAllProfessors() ([]bussinessLogic.Professor, error)
	AddProfessor(userId int) error
	AddStudent(userId int) error
	GetAllUsers() ([]bussinessLogic.User, error)
	InsertLesson(lessonName string, lessonUnit int) error
	DeleteLesson(lessonName string) error
	GetAllLessons() ([]bussinessLogic.Lesson, error)
	GetUsersByRole(roleId int) ([]bussinessLogic.User, error)
	AddMark(userId, classId int, mark *int) error
	GetStudentsForProfessor(professorId int) ([]bussinessLogic.Student, error)
	RemoveStudentUnit(id int) error
}
