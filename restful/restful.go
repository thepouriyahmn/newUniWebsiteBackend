package restful

import (
	"UniWebsite/auth"
	"UniWebsite/bussinessLogic"

	"fmt"
	"net/http"
)

type Restful struct {
	AuthBussinessLogic bussinessLogic.AuthBussinessLogic
}

func NewRestFul(authLogic bussinessLogic.AuthBussinessLogic) Restful {
	return Restful{
		AuthBussinessLogic: authLogic,
	}
}

func (rest Restful) Run() {
	http.HandleFunc("/signUp", rest.SignUp)

	http.HandleFunc("/login", rest.Login1)
	http.HandleFunc("/verify", rest.Verify)
	http.HandleFunc("/showProfessors", auth.AdminJwtMiddleware(rest.showProfessors))
	http.HandleFunc("/addProfessor", auth.AdminJwtMiddleware(rest.addProfessor))
	http.HandleFunc("/showAllUsers", auth.AdminJwtMiddleware(rest.showAllUsers))
	http.HandleFunc("/insertLesson", auth.AdminJwtMiddleware(rest.insertLesson))
	http.HandleFunc("/showClasses", auth.AdminJwtMiddleware(rest.showClasses))
	http.HandleFunc("/insertClass", auth.AdminJwtMiddleware(rest.insertClass))
	http.HandleFunc("/deleteClass", auth.AdminJwtMiddleware(rest.deleteClass))
	http.HandleFunc("/deleteLesson", auth.AdminJwtMiddleware(rest.deleteLesson))
	http.HandleFunc("/showAllLessons", auth.AdminJwtMiddleware(rest.showAllLessons))
	http.HandleFunc("/showUsersByRole", auth.AdminJwtMiddleware(rest.showUsersByRole))
	http.HandleFunc("/addMark", auth.ProfessorjwtMiddleware3(rest.addMark))
	http.HandleFunc("/showStudentsForProfessor", auth.ProfessorjwtMiddleware3(rest.showStudentsForProfessor))
	http.HandleFunc("/add", auth.StudentJwtMiddleware(rest.addStudentUnit))
	http.HandleFunc("/delStudentUnit", auth.StudentJwtMiddleware(rest.delStudentUnit))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("reding error: %v", err)
		panic(err)
	}
}
func (rest Restful) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// پاسخ به preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return

	}

	rest.AuthBussinessLogic.IProtocol.SignUp(w, r)
}
func (rest Restful) Login1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// پاسخ به preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return

	}
	fmt.Println("r.body is: ", r.Body, r.Method, r.Host)

	rest.AuthBussinessLogic.IProtocol.Login(w, r)
}
func (rest Restful) Verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// پاسخ به preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return

	}

	rest.AuthBussinessLogic.IProtocol.Verify(w, r)
}
func (rest Restful) showProfessors(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.GetAllProfessors(w, r)
}

func (rest Restful) addProfessor(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.AddProfessor(w, r)
}

func (rest Restful) showAllUsers(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.GetAllUsers(w, r)
}

func (rest Restful) insertLesson(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.InsertLesson(w, r)
}
func (rest Restful) insertClass(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.InsertClass(w, r)
}
func (rest Restful) showClasses(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.ShowClasses(w, r)
}

func (rest Restful) deleteLesson(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.DeleteLesson(w, r)
}

func (rest Restful) showAllLessons(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.GetAllLessons(w, r)
}

func (rest Restful) showUsersByRole(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.GetUsersByRole(w, r)
}

func (rest Restful) addMark(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.AddMark(w, r)
}

func (rest Restful) showStudentsForProfessor(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.GetStudentsForProfessor(w, r)
}

func (rest Restful) addStudentUnit(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.AddStudentUnit(w, r)
}

func (rest Restful) delStudentUnit(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.DelStudentUnit(w, r)
}
func (rest Restful) deleteClass(w http.ResponseWriter, r *http.Request) {
	rest.AuthBussinessLogic.IProtocol.DeleteClass(w, r)
}
