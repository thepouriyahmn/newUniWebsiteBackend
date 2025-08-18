package handler

import (
	"api/middlewear"
	"api/service"
	"fmt"
	"net/http"
)

type Handler struct {
	serviceURL string
	service    service.Services
}

func NewHandler(url string, service service.Services) Handler {
	return Handler{
		serviceURL: url,
		service:    service,
	}
}
func (h Handler) RunApi() {
	http.HandleFunc("/signUp", middlewear.CheackOriginMiddleWare(h.ProxySignUp))
	http.HandleFunc("/login", middlewear.CheackOriginMiddleWare(h.ProxyLogin))
	http.HandleFunc("/logout", middlewear.CheackOriginMiddleWare(h.ProxyLogout))
	http.HandleFunc("/verify", middlewear.CheackOriginMiddleWare(h.ProxyVerify))
	http.HandleFunc("/showProfessors", middlewear.CheackOriginMiddleWare(h.ProxyShowProfessors))
	http.HandleFunc("/getTerms", middlewear.CheackOriginMiddleWare(h.ProxyGetTerms))

	http.HandleFunc("/showAllUsers", middlewear.CheackOriginMiddleWare(h.ProxyShowAllUsers))
	http.HandleFunc("/insertLesson", middlewear.CheackOriginMiddleWare(h.ProxyInsertLesson))
	http.HandleFunc("/showClasses", middlewear.CheackOriginMiddleWare(h.ProxyShowClasses))
	http.HandleFunc("/insertClass", middlewear.CheackOriginMiddleWare(h.ProxyInsertClass))
	http.HandleFunc("/deleteClass", middlewear.CheackOriginMiddleWare(h.ProxyDeleteClass))
	http.HandleFunc("/deleteLesson", middlewear.CheackOriginMiddleWare(h.ProxyDeleteLesson))
	http.HandleFunc("/showAllLessons", middlewear.CheackOriginMiddleWare(h.ProxyShowAllLessons))
	http.HandleFunc("/showUsersByRole", middlewear.CheackOriginMiddleWare(h.ProxyShowUsersByRole))
	http.HandleFunc("/addStudent", middlewear.CheackOriginMiddleWare(h.ProxyAddStudent))
	http.HandleFunc("/addProfessor", middlewear.CheackOriginMiddleWare(h.ProxyAddProfessor))
	http.HandleFunc("/addMark", middlewear.CheackOriginMiddleWare(h.ProxyAddMark))
	http.HandleFunc("/showStudentsForProfessor", middlewear.CheackOriginMiddleWare(h.ProxyShowStudentsForProfessor))
	http.HandleFunc("/add", middlewear.CheackOriginMiddleWare(h.ProxyAddStudentUnit))
	http.HandleFunc("/pickedUnits", middlewear.CheackOriginMiddleWare(h.ProxyPickedUnits))
	http.HandleFunc("/delStudentUnit", middlewear.CheackOriginMiddleWare(h.ProxyDelStudentUnit))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("reading error: %v ", err)
		panic(err)
	}

}

func writeServiceResponse(w http.ResponseWriter, sr service.ServiceResponse) {
	w.Header().Set("Content-Type", "application/json")
	status := sr.StatusCode
	if status == 0 {
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
	if len(sr.Body) > 0 {
		_, _ = w.Write(sr.Body)
	}
}

func (h Handler) ProxySignUp(w http.ResponseWriter, r *http.Request) {
	sr, _ := h.service.SignUp(r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyLogin(w http.ResponseWriter, r *http.Request) {
	sr, _ := h.service.Login(r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyLogout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.Logout(token)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyVerify(w http.ResponseWriter, r *http.Request) {
	sr, _ := h.service.Verify(r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyShowProfessors(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.ShowProfessors(token)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyGetTerms(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.GetTerms(token)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyShowAllUsers(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.ShowAllUsers(token, r.URL.RawQuery)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyInsertLesson(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.InsertLesson(token, r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyShowClasses(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.ShowClasses(token, r.URL.RawQuery)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyInsertClass(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.InsertClass(token, r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyDeleteClass(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.DeleteClass(token, r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyDeleteLesson(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.DeleteLesson(token, r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyShowAllLessons(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.ShowAllLessons(token)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyShowUsersByRole(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.ShowUsersByRole(token, r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyAddStudent(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.AddStudent(token, r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyAddProfessor(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.AddProfessor(token, r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyAddMark(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.AddMark(token, r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyShowStudentsForProfessor(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.ShowStudentsForProfessor(token)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyAddStudentUnit(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.AddStudentUnit(token, r.Body)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyPickedUnits(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.PickedUnits(token)
	writeServiceResponse(w, sr)
}

func (h Handler) ProxyDelStudentUnit(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sr, _ := h.service.DelStudentUnit(token, r.Body)
	writeServiceResponse(w, sr)
}
