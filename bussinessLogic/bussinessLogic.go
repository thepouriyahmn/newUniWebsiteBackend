package bussinessLogic

import (
	"UniWebsite/auth"
	"errors"
	"fmt"
	"time"
)

type Bussinesslogic struct {
	ICache      ICache
	IDatabase   IDatabase
	IVerify     ISendVerificationCode
	IValidation IPassValidation
}

type IbussinessLogic interface {
	SignUp(username, email, password string, professorRole, studentRole bool) error
	Verify(id int, code string) (string, error)
	Login(username, pass string) (int, error)
	GetAllProfessors() ([]Professor, error)
	AddProfessor(userID int) error
	GetAllUsers(input string) ([]User, error)
	InsertLesson(lessonName string, lessonUnit int) error
	DeleteLesson(lessonName string) error
	GetAllLessons() ([]Lesson, error)
	GetUsersByRole(roleID int) ([]User, error)
	AddMark(userID, classID int, mark int) error
	GetStudentsForProfessor(professorID int) ([]Student2, error)
	AddStudentUnit(userID, classID int) error
	RemoveStudentUnit(classID, userID int) error
	InsertClass(lessonName, professorName, date, term string, capacity, classNumber int) error
	GetAllClassesByTerm(term string) ([]Classes, error)
	DeleteClass(classID int) error
	AddStudent(userID int) error
	GetClassesByUserId(userID int) ([]StudentClasses, error)
	GetAllTerms() ([]string, error)
}

func NewBussinessLogic(database IDatabase, cache ICache, verify ISendVerificationCode, passValidation IPassValidation) Bussinesslogic {
	return Bussinesslogic{
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

func (b Bussinesslogic) Login(username, pass string) (int, error) {

	id, email, err := b.IDatabase.CheackUserByUserNameAndPassword(username, pass)
	if err != nil {

		return 0, err
	}
	code, err := b.IVerify.SendCode(email)
	if err != nil {
		fmt.Printf("reading err: %v")
		return 0, err
	}
	mu.Lock()
	verificationCodes[id] = CodeInfo{
		Code:      code,
		CreatedAt: time.Now(),
	}
	mu.Unlock()
	return id, nil

}
func (b Bussinesslogic) Verify(id int, code string) (string, error) {
	roleSlice, username, err := b.IDatabase.GetRole(id)
	fmt.Println("roleslice: ", roleSlice)
	if err != nil {

		return "", err
	}
	fmt.Println(verificationCodes)
	mu.Lock()
	userInfo, ok := verificationCodes[id]
	mu.Unlock()

	fmt.Println("map : ", verificationCodes[id])
	if !ok {

		return "", err
	}

	if time.Since(userInfo.CreatedAt) > 2*time.Minute {
		mu.Lock()
		delete(verificationCodes, id)
		mu.Unlock()

		return "", err
	}
	fmt.Println("clientCode: ", code, "userInfo: ", userInfo.Code)
	if code != userInfo.Code {

		return "", err
	}

	tokenStr := auth.GenerateJWT(id, username, roleSlice)
	fmt.Println("token created: ", tokenStr)

	mu.Lock()
	delete(verificationCodes, id)
	mu.Unlock()
	return tokenStr, nil
}

// SignUp business logic
func (b Bussinesslogic) SignUp(username, password, email string, professorRole, studentRole bool) error {
	fmt.Println(password)
	valid := b.IValidation.IsValidPassword(password)
	if !valid {

		return errors.New("invalid password")
	}
	err := b.IDatabase.CheackUserByUsernameAndEmail(username, email)
	if err != nil {

		return err
	}

	err = b.IDatabase.InsertUser(username, password, email, studentRole, professorRole)
	if err != nil {

		return err
	}
	return nil
}

// GetAllProfessors business logic
func (b Bussinesslogic) GetAllProfessors() ([]Professor, error) {
	professors, err := b.IDatabase.GetAllProfessors()
	if err != nil {
		return nil, err
	}
	return professors, nil
}

// AddProfessor business logic
func (b Bussinesslogic) AddProfessor(userID int) error {
	err := b.IDatabase.AddProfessorById(userID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllUsers business logic
func (b Bussinesslogic) GetAllUsers(input string) ([]User, error) {
	users, err := b.IDatabase.GetAllUsers(input)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// InsertLesson business logic
func (b Bussinesslogic) InsertLesson(lessonName string, lessonUnit int) error {
	err := b.IDatabase.InsertLesson(lessonName, lessonUnit)
	if err != nil {
		return err
	}
	return nil
}

// DeleteLesson business logic
func (b Bussinesslogic) DeleteLesson(lessonName string) error {
	err := b.IDatabase.DeleteLesson(lessonName)
	if err != nil {
		return err
	}
	return nil
}

// GetAllLessons business logic
func (b Bussinesslogic) GetAllLessons() ([]Lesson, error) {
	lessons, err := b.IDatabase.GetAllLessons()
	if err != nil {
		return nil, err
	}
	return lessons, nil
}

// GetUsersByRole business logic
func (b Bussinesslogic) GetUsersByRole(roleID int) ([]User, error) {
	users, err := b.IDatabase.GetUsersByRole(roleID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// AddMark business logic
func (b Bussinesslogic) AddMark(userID, classID int, mark int) error {
	err := b.IDatabase.AddMark(userID, classID, mark)
	if err != nil {
		return err
	}
	return nil
}

// GetStudentsForProfessor business logic
func (b Bussinesslogic) GetStudentsForProfessor(professorID int) ([]Student2, error) {
	students, err := b.IDatabase.GetStudentsForProfessor(professorID)
	if err != nil {
		return nil, err
	}
	return students, nil
}

// AddStudentUnit business logic
func (b Bussinesslogic) AddStudentUnit(userID, classID int) error {
	err := b.IDatabase.InsertUnitForStudent(userID, classID)
	if err != nil {
		return err
	}
	return nil
}

// RemoveStudentUnit business logic
func (b Bussinesslogic) RemoveStudentUnit(classID, userID int) error {
	err := b.IDatabase.RemoveStudentUnit(classID, userID)
	if err != nil {
		return err
	}
	return nil
}

// InsertClass business logic
func (b Bussinesslogic) InsertClass(lessonName, professorName, date, term string, capacity, classNumber int) error {
	err := b.IDatabase.InsertClass(lessonName, professorName, date, term, capacity, classNumber)
	if err != nil {
		return err
	}
	return nil
}

// GetAllClassesByTerm business logic
func (b Bussinesslogic) GetAllClassesByTerm(term string) ([]Classes, error) {
	classes, err := b.IDatabase.GetAllClassesByTerm(term)
	if err != nil {
		return nil, err
	}
	return classes, nil
}

// DeleteClass business logic
func (b Bussinesslogic) DeleteClass(classID int) error {
	err := b.IDatabase.DeleteClass(classID)
	if err != nil {
		return err
	}
	return nil
}

// AddStudent business logic
func (b Bussinesslogic) AddStudent(userID int) error {
	err := b.IDatabase.AddStudentById(userID)
	if err != nil {
		return err
	}
	return nil
}

// GetClassesByUserId business logic
func (b Bussinesslogic) GetClassesByUserId(userID int) ([]StudentClasses, error) {
	classes, err := b.IDatabase.GetClassesByUserId(userID)
	if err != nil {
		return nil, err
	}
	return classes, nil
}

// GetAllTerms business logic
func (b Bussinesslogic) GetAllTerms() ([]string, error) {
	terms, err := b.IDatabase.GetAllTerms()
	if err != nil {
		return nil, err
	}
	return terms, nil
}
