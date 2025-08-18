package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	Nc *nats.Conn
}

func NewNats() Nats {
	nc, _ := nats.Connect(nats.DefaultURL)
	return Nats{
		Nc: nc,
	}
}

type natsResponse struct {
	Status     string          `json:"status"`
	Data       json.RawMessage `json:"data,omitempty"`
	Code       string          `json:"code,omitempty"`
	Message    string          `json:"message,omitempty"`
	HttpStatus int             `json:"httpStatus,omitempty"`
}

type natsRequest struct {
	Token string          `json:"token,omitempty"`
	Query string          `json:"query,omitempty"`
	Body  json.RawMessage `json:"body,omitempty"`
}

func (n Nats) do(subject string, req natsRequest) (ServiceResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return ServiceResponse{StatusCode: 500, Body: []byte("marshal error")}, err
	}
	msg, err := n.Nc.Request(subject, data, 20*time.Second)
	if err != nil {
		log.Printf("NATS request error: %v", err)
		return ServiceResponse{StatusCode: 503, Body: []byte(`{"status":"error","message":"service unavailable"}`)}, err
	}
	var resp natsResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return ServiceResponse{StatusCode: 500, Body: []byte(`{"status":"error","message":"invalid response"}`)}, err
	}
	status := resp.HttpStatus
	if status == 0 {
		if resp.Status == "ok" {
			status = 200
		} else {
			status = 400
		}
	}
	if resp.Status == "ok" {
		if len(resp.Data) == 0 {
			return ServiceResponse{StatusCode: status, Body: []byte(`{"status":"ok"}`)}, nil
		}
		return ServiceResponse{StatusCode: status, Body: resp.Data}, nil
	}
	// error case
	body, _ := json.Marshal(map[string]string{
		"status":  "error",
		"code":    resp.Code,
		"message": resp.Message,
	})
	return ServiceResponse{StatusCode: status, Body: body}, fmt.Errorf("%s", resp.Message)
}

// Auth
func (n Nats) SignUp(body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("signup.request", natsRequest{Body: json.RawMessage(b)})
}
func (n Nats) Login(body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("login.request", natsRequest{Body: json.RawMessage(b)})
}
func (n Nats) Verify(body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("verify.request", natsRequest{Body: json.RawMessage(b)})
}
func (n Nats) Logout(token string) (ServiceResponse, error) {
	return n.do("logout.request", natsRequest{Token: token})
}

// Admin
func (n Nats) ShowProfessors(token string) (ServiceResponse, error) {
	return n.do("professors.list", natsRequest{Token: token})
}
func (n Nats) GetTerms(token string) (ServiceResponse, error) {
	return n.do("terms.get", natsRequest{Token: token})
}
func (n Nats) ShowAllUsers(token string, query string) (ServiceResponse, error) {
	return n.do("users.list", natsRequest{Token: token, Query: query})
}
func (n Nats) InsertLesson(token string, body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("lessons.insert", natsRequest{Token: token, Body: json.RawMessage(b)})
}
func (n Nats) ShowAllLessons(token string) (ServiceResponse, error) {
	return n.do("lessons.list", natsRequest{Token: token})
}
func (n Nats) DeleteLesson(token string, body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("lesson.delete", natsRequest{Token: token, Body: json.RawMessage(b)})
}
func (n Nats) InsertClass(token string, body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("class.insert", natsRequest{Token: token, Body: json.RawMessage(b)})
}
func (n Nats) ShowClasses(token string, query string) (ServiceResponse, error) {
	return n.do("classes.show", natsRequest{Token: token, Query: query})
}
func (n Nats) DeleteClass(token string, body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("class.delete", natsRequest{Token: token, Body: json.RawMessage(b)})
}
func (n Nats) ShowUsersByRole(token string, body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("users.byRole", natsRequest{Token: token, Body: json.RawMessage(b)})
}
func (n Nats) AddStudent(token string, body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("student.add", natsRequest{Token: token, Body: json.RawMessage(b)})
}
func (n Nats) AddProfessor(token string, body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("professor.add", natsRequest{Token: token, Body: json.RawMessage(b)})
}

// Professor
func (n Nats) AddMark(token string, body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("mark.add", natsRequest{Token: token, Body: json.RawMessage(b)})
}
func (n Nats) ShowStudentsForProfessor(token string) (ServiceResponse, error) {
	return n.do("professor.students", natsRequest{Token: token})
}

// Student
func (n Nats) AddStudentUnit(token string, body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("student.unit.add", natsRequest{Token: token, Body: json.RawMessage(b)})
}
func (n Nats) PickedUnits(token string) (ServiceResponse, error) {
	return n.do("student.units", natsRequest{Token: token})
}
func (n Nats) DelStudentUnit(token string, body io.Reader) (ServiceResponse, error) {
	b, _ := io.ReadAll(body)
	return n.do("student.unit.delete", natsRequest{Token: token, Body: json.RawMessage(b)})
}
