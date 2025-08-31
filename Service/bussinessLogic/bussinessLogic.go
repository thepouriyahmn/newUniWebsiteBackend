package bussinessLogic

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Bussinesslogic struct {
	Ibroker       IMessageBroker
	ICache        ICache
	ISubDatabase  SubDatabase
	IMainDatabase MainDatabase
	IVerify       ISendVerificationCode
	IValidation   IPassValidation
	IToken        IToken
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

func NewBussinessLogic(writeDatabase MainDatabase, readDatabase SubDatabase, Ibroker IMessageBroker, cache ICache, verify ISendVerificationCode, passValidation IPassValidation, token IToken) Bussinesslogic {
	return Bussinesslogic{
		Ibroker:       Ibroker,
		ISubDatabase:  readDatabase,
		IMainDatabase: writeDatabase,
		ICache:        cache,
		IVerify:       verify,
		IValidation:   passValidation,
		IToken:        token,
	}
}

type IMessageBroker interface {
	// Subscribe(topic string, cb func()) error
	Publish(subject string, req json.RawMessage) error
}

type ISendVerificationCode interface {
	SendCode(reciever string) (string, error)
}

// type MainDatabase interface {
// 	InsertUser(username, pass, email string, studentRole, professorRole bool) error
// }
// type SubDatabase interface {
// 	GetAllUsers(input string) ([]User, error)
// }

type ICache interface {
	CacheTerms(terms []string)
	GetCacheValue(key string) (string, error)
}

type IPassValidation interface {
	IsValidPassword(password string) bool
}
type IToken interface {
	GenerateToken(id int, username string, roleSlice []string) string
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
	GetAllProfessors() ([]Professor, error)
	GetAllUsers(input string) ([]User, error)
	GetAllLessons() ([]Lesson, error)
	GetAllClassesByTerm(term string) ([]Classes, error)
	GetAllTerms() ([]string, error)
	GetUsersByRole(roleId int) ([]User, error)
	GetStudentsForProfessor(professorId int) ([]Student2, error)
	GetClassesByUserId(userID int) ([]StudentClasses, error)

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

func (b Bussinesslogic) Login(username, pass string) (int, error) {

	id, email, err := b.ISubDatabase.CheackUserByUserNameAndPassword(username, pass)
	if err != nil {

		return 0, err
	}
	code, err := b.IVerify.SendCode(email)
	if err != nil {
		fmt.Printf("reading err: %v", err)
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
	roleSlice, username, err := b.ISubDatabase.GetRole(id)
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

	tokenStr := b.IToken.GenerateToken(id, username, roleSlice)
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

	// Check for duplicates from MongoDB (write database)
	// err := b.IMainDatabase.CheackUserByUsernameAndEmail(username, email)
	// if err != nil {
	// 	fmt.Printf("reading error: %v", err)
	// 	return err
	// }

	// Insert user to MongoDB (write database)
	err := b.IMainDatabase.InsertUser(username, password, email, studentRole, professorRole)
	if err != nil {
		fmt.Printf("reading errorr: %v", err)
		return err
	}

	// Publish event to NATS for MySQL sync
	var users Users
	users.Email = email
	users.Password = password
	users.Username = username
	users.ProfessorRole = professorRole
	users.StudentRole = studentRole
	data, err := json.Marshal(users)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	err = b.Ibroker.Publish("user.created", data)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return nil
	}

	return nil
}

// GetAllProfessors business logic
func (b Bussinesslogic) GetAllProfessors() ([]Professor, error) {
	professors, err := b.ISubDatabase.GetAllProfessors()
	if err != nil {
		return nil, err
	}
	return professors, nil
}

// AddProfessor business logic
func (b Bussinesslogic) AddProfessor(userID int) error {
	// Add to MongoDB (write database) - CQRS pattern
	err := b.IMainDatabase.AddProfessorById(userID)
	if err != nil {
		return err
	}

	// Publish event to NATS for MySQL sync
	professor := Professor{
		Id: userID,
	}
	data, err := json.Marshal(professor)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	err = b.Ibroker.Publish("professor.added", data)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return nil
	}

	return nil
}

// GetAllUsers business logic
func (b Bussinesslogic) GetAllUsers(input string) ([]User, error) {
	users, err := b.ISubDatabase.GetAllUsers(input)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// InsertLesson business logic
func (b Bussinesslogic) InsertLesson(lessonName string, lessonUnit int) error {
	// Insert to MongoDB (write database) - CQRS pattern
	err := b.IMainDatabase.InsertLesson(lessonName, lessonUnit)
	if err != nil {
		return err
	}

	// Publish event to NATS for MySQL sync
	lesson := Lesson{
		LessonName: lessonName,
		LessonUnit: lessonUnit,
	}
	data, err := json.Marshal(lesson)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	err = b.Ibroker.Publish("lesson.created", data)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return nil
	}

	return nil
}

// DeleteLesson business logic
func (b Bussinesslogic) DeleteLesson(lessonName string) error {
	// Delete from MongoDB (write database) - CQRS pattern
	err := b.IMainDatabase.DeleteLesson(lessonName)
	if err != nil {
		return err
	}

	// Publish event to NATS for MySQL sync
	lesson := Lesson{
		LessonName: lessonName,
	}
	data, err := json.Marshal(lesson)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	err = b.Ibroker.Publish("lesson.deleted", data)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return nil
	}

	return nil
}

// GetAllLessons business logic
func (b Bussinesslogic) GetAllLessons() ([]Lesson, error) {
	lessons, err := b.ISubDatabase.GetAllLessons()
	if err != nil {
		return nil, err
	}
	return lessons, nil
}

// GetUsersByRole business logic
func (b Bussinesslogic) GetUsersByRole(roleID int) ([]User, error) {
	users, err := b.ISubDatabase.GetUsersByRole(roleID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// AddMark business logic
func (b Bussinesslogic) AddMark(userID, classID int, mark int) error {
	// Add to MongoDB (write database) - CQRS pattern
	err := b.IMainDatabase.AddMark(userID, classID, mark)
	if err != nil {
		return err
	}

	// Publish event to NATS for MySQL sync
	markData := struct {
		UserId  int `json:"userId"`
		ClassId int `json:"classId"`
		Mark    int `json:"mark"`
	}{
		UserId:  userID,
		ClassId: classID,
		Mark:    mark,
	}
	data, err := json.Marshal(markData)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	err = b.Ibroker.Publish("mark.added", data)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return nil
	}

	return nil
}

// GetStudentsForProfessor business logic
func (b Bussinesslogic) GetStudentsForProfessor(professorID int) ([]Student2, error) {
	students, err := b.ISubDatabase.GetStudentsForProfessor(professorID)
	if err != nil {
		return nil, err
	}
	return students, nil
}

// AddStudentUnit business logic
func (b Bussinesslogic) AddStudentUnit(userID, classID int) error {
	// Add to MongoDB (write database) - CQRS pattern
	err := b.IMainDatabase.InsertUnitForStudent(userID, classID)
	if err != nil {
		return err
	}

	// Publish event to NATS for MySQL sync
	unitData := struct {
		UserId  int `json:"userId"`
		ClassId int `json:"classId"`
	}{
		UserId:  userID,
		ClassId: classID,
	}
	data, err := json.Marshal(unitData)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	err = b.Ibroker.Publish("student.unit.added", data)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return nil
	}

	return nil
}

// RemoveStudentUnit business logic
func (b Bussinesslogic) RemoveStudentUnit(classID, userID int) error {
	// Remove from MongoDB (write database) - CQRS pattern
	err := b.IMainDatabase.RemoveStudentUnit(classID, userID)
	if err != nil {
		return err
	}

	// Publish event to NATS for MySQL sync
	unitData := struct {
		UserId  int `json:"userId"`
		ClassId int `json:"classId"`
	}{
		UserId:  userID,
		ClassId: classID,
	}
	data, err := json.Marshal(unitData)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	err = b.Ibroker.Publish("student.unit.removed", data)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return nil
	}

	return nil
}

// InsertClass business logic
func (b Bussinesslogic) InsertClass(lessonName, professorName, date, term string, capacity, classNumber int) error {
	// Insert to MongoDB (write database) - CQRS pattern
	err := b.IMainDatabase.InsertClass(lessonName, professorName, date, term, capacity, classNumber)
	if err != nil {
		return err
	}

	// Publish event to NATS for MySQL sync
	class := Classes{
		LessonName:    lessonName,
		ProfessorName: professorName,
		Date:          date,
		Term:          term,
		Capacity:      capacity,
		ClassNumber:   classNumber,
	}
	data, err := json.Marshal(class)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	err = b.Ibroker.Publish("class.created", data)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return nil
	}

	return nil
}

// GetAllClassesByTerm business logic
func (b Bussinesslogic) GetAllClassesByTerm(term string) ([]Classes, error) {
	fmt.Println("term: ", term)
	classes, err := b.ISubDatabase.GetAllClassesByTerm(term)
	if err != nil {
		return nil, err
	}
	return classes, nil
}

// DeleteClass business logic
func (b Bussinesslogic) DeleteClass(classID int) error {
	// Delete from MongoDB (write database) - CQRS pattern
	err := b.IMainDatabase.DeleteClass(classID)
	if err != nil {
		return err
	}

	// Publish event to NATS for MySQL sync
	class := Classes{
		Id: classID,
	}
	data, err := json.Marshal(class)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	err = b.Ibroker.Publish("class.deleted", data)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return nil
	}

	return nil
}

// AddStudent business logic
func (b Bussinesslogic) AddStudent(userID int) error {
	// Add to MongoDB (write database) - CQRS pattern
	err := b.IMainDatabase.AddStudentById(userID)
	if err != nil {
		return err
	}

	// Publish event to NATS for MySQL sync
	student := struct {
		UserId int `json:"userId"`
	}{
		UserId: userID,
	}
	data, err := json.Marshal(student)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	err = b.Ibroker.Publish("student.added", data)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return nil
	}

	return nil
}

// GetClassesByUserId business logic
func (b Bussinesslogic) GetClassesByUserId(userID int) ([]StudentClasses, error) {
	classes, err := b.ISubDatabase.GetClassesByUserId(userID)
	if err != nil {
		return nil, err
	}
	return classes, nil
}

// GetAllTerms business logic
func (b Bussinesslogic) GetAllTerms() ([]string, error) {
	terms, err := b.ISubDatabase.GetAllTerms()
	if err != nil {
		return nil, err
	}
	return terms, nil
}
