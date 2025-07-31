package restful

import (
	"UniWebsite/auth"
	"UniWebsite/bussinessLogic"
	"time"

	"fmt"
	"net/http"

	"github.com/go-redis/redis"
)

type Restful struct {
	Bussinesslogic bussinessLogic.Bussinesslogic
}

func NewRestFul(authLogic bussinessLogic.Bussinesslogic) Restful {
	return Restful{
		Bussinesslogic: authLogic,
	}
}

func (rest Restful) Run() {
	http.HandleFunc("/signUp", auth.CheackOriginMiddleWare(rest.SignUp))

	http.HandleFunc("/login", auth.CheackOriginMiddleWare(rest.Login1))
	http.HandleFunc("/logout", auth.CheackOriginMiddleWare(rest.logout))
	http.HandleFunc("/verify", auth.CheackOriginMiddleWare(rest.Verify))
	http.HandleFunc("/showProfessors", auth.AdminJwtMiddleware(rest.showProfessors))
	http.HandleFunc("/getTerms", auth.AdminJwtMiddleware(rest.getTerms))

	http.HandleFunc("/showAllUsers", auth.AdminJwtMiddleware(rest.showAllUsers))
	http.HandleFunc("/insertLesson", auth.AdminJwtMiddleware(rest.insertLesson))
	http.HandleFunc("/showClasses", auth.NormalJwtmiddleWare(rest.showClasses))
	http.HandleFunc("/insertClass", auth.AdminJwtMiddleware(rest.insertClass))
	http.HandleFunc("/deleteClass", auth.AdminJwtMiddleware(rest.deleteClass))
	http.HandleFunc("/deleteLesson", auth.AdminJwtMiddleware(rest.deleteLesson))
	http.HandleFunc("/showAllLessons", auth.AdminJwtMiddleware(rest.showAllLessons))
	http.HandleFunc("/showUsersByRole", auth.AdminJwtMiddleware(rest.showUsersByRole))
	http.HandleFunc("/addStudent", auth.AdminJwtMiddleware(rest.addStudent))
	http.HandleFunc("/addProfessor", auth.AdminJwtMiddleware(rest.addProfessor))
	http.HandleFunc("/addMark", auth.ProfessorjwtMiddleware3(rest.addMark))
	http.HandleFunc("/showStudentsForProfessor", auth.ProfessorjwtMiddleware3(rest.showStudentsForProfessor))
	http.HandleFunc("/add", auth.StudentJwtMiddleware(rest.addStudentUnit))
	http.HandleFunc("/pickedUnits", auth.StudentJwtMiddleware(rest.pickedUnits))
	http.HandleFunc("/delStudentUnit", auth.StudentJwtMiddleware(rest.delStudentUnit))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("reding error: %v", err)
		panic(err)
	}
}
func (rest Restful) SignUp(w http.ResponseWriter, r *http.Request) {

	rest.Bussinesslogic.IProtocol.SignUp(w, r)
}
func (rest Restful) Login1(w http.ResponseWriter, r *http.Request) {

	rest.Bussinesslogic.IProtocol.Login(w, r)
}
func (rest Restful) Verify(w http.ResponseWriter, r *http.Request) {

	rest.Bussinesslogic.IProtocol.Verify(w, r)
}
func (rest Restful) showProfessors(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.GetAllProfessors(w, r)
}

func (rest Restful) addProfessor(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.AddProfessor(w, r)
}

func (rest Restful) showAllUsers(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.GetAllUsers(w, r)
}

func (rest Restful) insertLesson(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.InsertLesson(w, r)
}
func (rest Restful) insertClass(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.InsertClass(w, r)
}
func (rest Restful) showClasses(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.ShowClasses(w, r)
}

func (rest Restful) deleteLesson(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.DeleteLesson(w, r)
}

func (rest Restful) showAllLessons(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.GetAllLessons(w, r)
}

func (rest Restful) showUsersByRole(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.GetUsersByRole(w, r)
}

func (rest Restful) addMark(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.AddMark(w, r)
}

func (rest Restful) showStudentsForProfessor(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.GetStudentsForProfessor(w, r)
}

func (rest Restful) addStudentUnit(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.AddStudentUnit(w, r)
}

func (rest Restful) delStudentUnit(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.DelStudentUnit(w, r)
}
func (rest Restful) deleteClass(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.DeleteClass(w, r)
}
func (rest Restful) addStudent(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.AddStudent(w, r)
}
func (rest Restful) pickedUnits(w http.ResponseWriter, r *http.Request) {
	rest.Bussinesslogic.IProtocol.ShowPickedUnitsForStudent(w, r)
}
func (rest Restful) logout(w http.ResponseWriter, r *http.Request) {

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	rdb.Set(tokenStr, "blocked", 5*time.Minute)
}
func (rest Restful) getTerms(w http.ResponseWriter, r *http.Request) {

	rest.Bussinesslogic.IProtocol.GetTerms(w, r)
}
