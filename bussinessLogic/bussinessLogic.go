package bussinessLogic

type Bussinesslogic struct {
	ICache    ICache
	IDatabase IDatabase
	IVerify   ISendVerificationCode
}

func NewBussinessLogic(database IDatabase, cache ICache, verify ISendVerificationCode) Bussinesslogic {
	return Bussinesslogic{
		//IProtocol: protocol,
		IDatabase: database,
		ICache:    cache,
		IVerify:   verify,
	}

}

type ISendVerificationCode interface {
	SendCode(reciever string) (string, error)
}
type ICache interface {
	CacheTerms(terms []string)
	GetCacheValue(key string) (string, error)
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

// type IProtocol interface {
// 	ShowPickedUnitsForStudent(w http.ResponseWriter, r *http.Request)
// 	SignUp(w http.ResponseWriter, r *http.Request)
// 	Login(w http.ResponseWriter, r *http.Request)
// 	Verify(w http.ResponseWriter, r *http.Request)
// 	GetAllProfessors(w http.ResponseWriter, r *http.Request)
// 	GetTerms(w http.ResponseWriter, r *http.Request)
// 	AddProfessor(w http.ResponseWriter, r *http.Request)
// 	GetAllUsers(w http.ResponseWriter, r *http.Request)
// 	InsertLesson(w http.ResponseWriter, r *http.Request)
// 	InsertClass(w http.ResponseWriter, r *http.Request)
// 	DeleteLesson(w http.ResponseWriter, r *http.Request)
// 	GetAllLessons(w http.ResponseWriter, r *http.Request)
// 	GetUsersByRole(w http.ResponseWriter, r *http.Request)
// 	AddMark(w http.ResponseWriter, r *http.Request)
// 	GetStudentsForProfessor(w http.ResponseWriter, r *http.Request)
// 	AddStudentUnit(w http.ResponseWriter, r *http.Request)
// 	DelStudentUnit(w http.ResponseWriter, r *http.Request)
// 	ShowClasses(w http.ResponseWriter, r *http.Request)
// 	DeleteClass(w http.ResponseWriter, r *http.Request)
// 	AddStudent(w http.ResponseWriter, r *http.Request)
// }
