package bussinessLogic

import (
	"errors"
	"fmt"
)

type Bussinesslogic struct {
	ICache      ICache
	IDatabase   IDatabase
	IVerify     ISendVerificationCode
	IValidation IPassValidation
}
type IbussinessLogic interface {
	SignUp(username, email, password string, professorRole, studentRole bool) error
}

func NewBussinessLogic(database IDatabase, cache ICache, verify ISendVerificationCode, passValidation IPassValidation) Bussinesslogic {
	return Bussinesslogic{
		//IProtocol: protocol,
		IDatabase:   database,
		ICache:      cache,
		IVerify:     verify,
		IValidation: passValidation,
	}

}

type ISendVerificationCode interface {
	SendCode(reciever string) (string, error)
}
type ICache interface {
	CacheTerms(terms []string)
	GetCacheValue(key string) (string, error)
}
type IPassValidation interface {
	IsValidPassword(password string) bool
}

type IDatabase interface {
	// متدهای احراز هویت و ثبت‌نام
	CheackUserByUsernameAndEmail(ClientUsername, ClientEmail string) error
	InsertUser(username, pass, email string, studentRole, professorRole bool) error
	CheackUserByUserNameAndPassword(username, pass string) (int, string, error)
	GetRole(id int) ([]string, string, error)

	// متدهای مدیریتی و آموزشی
	GetAllProfessors() ([]Professor, error)
	AddProfessorById(userId int) error
	AddStudent(userId int) error
	GetAllUsers(input string) ([]User, error)
	InsertLesson(lessonName string, lessonUnit int) error
	InsertClass(lessonName, professorName, date, term string, capacity, classNumber int) error
	GetAllClassesByTerm(term string) ([]Classes, error)
	DeleteClass(classId int) error
	GetAllTerms() ([]string, error)
	DeleteLesson(lessonName string) error
	GetAllLessons() ([]Lesson, error)
	GetUsersByRole(roleId int) ([]User, error)
	AddStudentById(userId int) error
	AddMark(userId, classId int, mark int) error
	GetStudentsForProfessor(professorId int) ([]Student2, error)
	RemoveStudentUnit(classid int, userid int) error
	GetClassesByUserId(userID int) ([]StudentClasses, error)
	InsertUnitForStudent(userid int, classid int) error
}

func (b Bussinesslogic) SignUp(username, password, email string, professorRole, studentRole bool) error {
	fmt.Println(password)
	valid := b.IValidation.IsValidPassword(password)
	if !valid {
		fmt.Println("pass validation wrong")
		return errors.New("")
	}
	err := b.IDatabase.CheackUserByUsernameAndEmail(username, email)
	if err != nil {
		fmt.Println("username,email validation wrong")
		return err
	}

	err = b.IDatabase.InsertUser(username, password, email, studentRole, professorRole)
	if err != nil {
		fmt.Println("insert wrong")
		return err
	}
	return nil
}
