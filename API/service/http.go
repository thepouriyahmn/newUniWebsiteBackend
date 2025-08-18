package service

import (
	"io"
	"net/http"
	"net/url"
)

type Http struct {
	serviceURL string
}

func NewHttp(serviceURL string) Http {
	return Http{serviceURL: serviceURL}
}

func (h Http) do(method, path, token, rawQuery string, body io.Reader) (ServiceResponse, error) {
	u := url.URL{Scheme: "http", Host: h.serviceURL, Path: path, RawQuery: rawQuery}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return ServiceResponse{StatusCode: http.StatusInternalServerError, Body: []byte("request build error")}, err
	}
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ServiceResponse{StatusCode: http.StatusServiceUnavailable, Body: []byte("service unavailable")}, err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	return ServiceResponse{StatusCode: resp.StatusCode, Body: bodyBytes}, nil
}

// Auth
func (h Http) SignUp(body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/signUp", "", "", body)
}
func (h Http) Login(body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/login", "", "", body)
}
func (h Http) Verify(body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/verify", "", "", body)
}
func (h Http) Logout(token string) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/logout", token, "", nil)
}

// Admin
func (h Http) ShowProfessors(token string) (ServiceResponse, error) {
	return h.do(http.MethodGet, "/showProfessors", token, "", nil)
}
func (h Http) GetTerms(token string) (ServiceResponse, error) {
	return h.do(http.MethodGet, "/getTerms", token, "", nil)
}
func (h Http) ShowAllUsers(token string, query string) (ServiceResponse, error) {
	return h.do(http.MethodGet, "/showAllUsers", token, query, nil)
}
func (h Http) InsertLesson(token string, body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/insertLesson", token, "", body)
}
func (h Http) ShowAllLessons(token string) (ServiceResponse, error) {
	return h.do(http.MethodGet, "/showAllLessons", token, "", nil)
}
func (h Http) DeleteLesson(token string, body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/deleteLesson", token, "", body)
}
func (h Http) InsertClass(token string, body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/insertClass", token, "", body)
}
func (h Http) ShowClasses(token string, query string) (ServiceResponse, error) {
	return h.do(http.MethodGet, "/showClasses", token, query, nil)
}
func (h Http) DeleteClass(token string, body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/deleteClass", token, "", body)
}
func (h Http) ShowUsersByRole(token string, body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/showUsersByRole", token, "", body)
}
func (h Http) AddStudent(token string, body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/addStudent", token, "", body)
}
func (h Http) AddProfessor(token string, body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/addProfessor", token, "", body)
}

// Professor
func (h Http) AddMark(token string, body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/addMark", token, "", body)
}
func (h Http) ShowStudentsForProfessor(token string) (ServiceResponse, error) {
	return h.do(http.MethodGet, "/showStudentsForProfessor", token, "", nil)
}

// Student
func (h Http) AddStudentUnit(token string, body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/add", token, "", body)
}
func (h Http) PickedUnits(token string) (ServiceResponse, error) {
	return h.do(http.MethodGet, "/pickedUnits", token, "", nil)
}
func (h Http) DelStudentUnit(token string, body io.Reader) (ServiceResponse, error) {
	return h.do(http.MethodPost, "/delStudentUnit", token, "", body)
}
