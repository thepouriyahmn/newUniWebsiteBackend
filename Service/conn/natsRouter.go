package restful

import (
	"UniWebsite/auth"
	"UniWebsite/bussinessLogic"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/nats-io/nats.go"
)

type HandlerFunc func(m *nats.Msg)

type Router struct {
	Nc    *nats.Conn
	Logic bussinessLogic.Bussinesslogic
}

type natsRequest struct {
	Token string          `json:"token,omitempty"`
	Query string          `json:"query,omitempty"`
	Body  json.RawMessage `json:"body,omitempty"`
}

type natsResponse struct {
	Status     string      `json:"status"`
	Data       interface{} `json:"data,omitempty"`
	Code       string      `json:"code,omitempty"`
	Message    string      `json:"message,omitempty"`
	HttpStatus int         `json:"httpStatus,omitempty"`
}

func respondOK(msg *nats.Msg, data interface{}, status int) {
	if status == 0 {
		status = 200
	}
	res := natsResponse{Status: "ok", Data: data, HttpStatus: status}
	b, _ := json.Marshal(res)
	_ = msg.Respond(b)
}

func respondErr(msg *nats.Msg, httpStatus int, code, message string) {
	res := natsResponse{Status: "error", Code: code, Message: message, HttpStatus: httpStatus}
	b, _ := json.Marshal(res)
	_ = msg.Respond(b)
}

func NewNats(logic bussinessLogic.Bussinesslogic) *Router {
	nc, _ := nats.Connect(nats.DefaultURL)
	return &Router{Nc: nc, Logic: logic}
}

func (r *Router) Run() {
	fmt.Println("run in nats")
	// Auth
	r.Nc.Subscribe("signup.request", r.handleSignup)
	r.Nc.Subscribe("login.request", r.handleLogin)
	r.Nc.Subscribe("verify.request", r.handleVerify)
	r.Nc.Subscribe("logout.request", r.handleLogout)

	// Admin
	r.Nc.Subscribe("professors.list", r.handleProfessorsList)
	r.Nc.Subscribe("terms.get", r.handleTermsGet)
	r.Nc.Subscribe("users.list", r.handleUsersList)
	r.Nc.Subscribe("lessons.insert", r.handleLessonInsert)
	r.Nc.Subscribe("lessons.list", r.handleLessonsList)
	r.Nc.Subscribe("lesson.delete", r.handleLessonDelete)
	r.Nc.Subscribe("class.insert", r.handleClassInsert)
	r.Nc.Subscribe("classes.show", r.handleClassesShow)
	r.Nc.Subscribe("class.delete", r.handleClassDelete)
	r.Nc.Subscribe("users.byRole", r.handleUsersByRole)
	r.Nc.Subscribe("student.add", r.handleStudentAdd)
	r.Nc.Subscribe("professor.add", r.handleProfessorAdd)

	// Professor
	r.Nc.Subscribe("mark.add", r.handleMarkAdd)
	r.Nc.Subscribe("professor.students", r.handleProfessorStudents)

	// Student
	r.Nc.Subscribe("student.unit.add", r.handleStudentUnitAdd)
	r.Nc.Subscribe("student.units", r.handleStudentUnits)
	r.Nc.Subscribe("student.unit.delete", r.handleStudentUnitDelete)

	select {}
}

func (r *Router) handleSignup(msg *nats.Msg) {
	type User struct {
		Username      string `json:"username"`
		Password      string `json:"password"`
		StudentRole   bool   `json:"studentRole"`
		ProfessorRole bool   `json:"professorRole"`
		Email         string `json:"email"`
	}
	var req natsRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		log.Println("json error:", err)
		respondErr(msg, 400, "bad_request", "invalid json")
		return
	}
	var user User
	if err := json.Unmarshal(req.Body, &user); err != nil {
		log.Println("error:", err)
		respondErr(msg, 400, "bad_request", "invalid json body")
		return
	}
	if err := r.Logic.SignUp(user.Username, user.Password, user.Email, user.StudentRole, user.ProfessorRole); err != nil {
		log.Println(" error:", err)
		respondErr(msg, 400, "signup_failed", "signup failed")
		return
	}
	respondOK(msg, map[string]string{"message": "signup successful"}, 200)
}

func (r *Router) handleLogin(msg *nats.Msg) {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req natsRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		respondErr(msg, 400, "bad_request", "invalid json")
		return
	}
	var body Req
	if err := json.Unmarshal(req.Body, &body); err != nil {
		respondErr(msg, 400, "bad_request", "invalid json body")
		return
	}
	id, err := r.Logic.Login(body.Username, body.Password)
	if err != nil {
		respondErr(msg, 401, "unauthorized", "invalid credentials")
		return
	}
	respondOK(msg, map[string]int{"id": id}, 200)
}

func (r *Router) handleVerify(msg *nats.Msg) {
	type Req struct {
		Id   int    `json:"id"`
		Code string `json:"code"`
	}
	var req natsRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		respondErr(msg, 400, "bad_request", "invalid json")
		return
	}
	var body Req
	if err := json.Unmarshal(req.Body, &body); err != nil {
		respondErr(msg, 400, "bad_request", "invalid json body")
		return
	}
	token, err := r.Logic.Verify(body.Id, body.Code)
	if err != nil {
		respondErr(msg, 401, "unauthorized", "code not valid")
		return
	}
	respondOK(msg, map[string]string{"token": token}, 200)
}

func (r *Router) handleLogout(msg *nats.Msg) {
	var req natsRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		respondErr(msg, 400, "bad_request", "invalid json")
		return
	}
	if req.Token == "" {
		respondErr(msg, 401, "unauthorized", "missing token")
		return
	}
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	rdb.Set(req.Token, "blocked", 5*time.Minute)
	respondOK(msg, nil, 200)
}

// Admin handlers (role 1)
func (r *Router) requireAdmin(token string) (*auth.Claims, error) {
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	if !auth.ClaimsHasRole(claims, "1") {
		return nil, fmt.Errorf("forbidden")
	}
	return claims, nil
}

func (r *Router) requireProfessor(token string) (*auth.Claims, error) {
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	if !auth.ClaimsHasRole(claims, "3") {
		return nil, fmt.Errorf("forbidden")
	}
	return claims, nil
}

func (r *Router) requireStudent(token string) (*auth.Claims, error) {
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	if !auth.ClaimsHasRole(claims, "2") {
		return nil, fmt.Errorf("forbidden")
	}
	return claims, nil
}

func (r *Router) handleProfessorsList(msg *nats.Msg) {
	var req natsRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		respondErr(msg, 400, "bad_request", "invalid json")
		return
	}
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	res, err := r.Logic.GetAllProfessors()
	if err != nil {
		respondErr(msg, 500, "server_error", "failed to fetch")
		return
	}
	respondOK(msg, res, 200)
}

func (r *Router) handleTermsGet(msg *nats.Msg) {
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	res, err := r.Logic.GetAllTerms()
	if err != nil {
		respondErr(msg, 500, "server_error", "failed to fetch terms")
		return
	}
	respondOK(msg, res, 200)
}

func (r *Router) handleUsersList(msg *nats.Msg) {
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	res, err := r.Logic.GetAllUsers(req.Query)
	if err != nil {
		respondErr(msg, 500, "server_error", "failed to get users")
		return
	}
	respondOK(msg, res, 200)
}

func (r *Router) handleLessonInsert(msg *nats.Msg) {
	type Body struct {
		LessonName string `json:"lessonName"`
		LessonUnit int    `json:"lessonUnit"`
	}
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	var b Body
	if err := json.Unmarshal(req.Body, &b); err != nil {
		respondErr(msg, 400, "bad_request", "invalid body")
		return
	}
	if err := r.Logic.InsertLesson(b.LessonName, b.LessonUnit); err != nil {
		respondErr(msg, 500, "server_error", "failed to insert lesson")
		return
	}
	respondOK(msg, nil, 201)
}

func (r *Router) handleLessonsList(msg *nats.Msg) {
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	res, err := r.Logic.GetAllLessons()
	if err != nil {
		respondErr(msg, 500, "server_error", "failed to get lessons")
		return
	}
	respondOK(msg, res, 200)
}

func (r *Router) handleLessonDelete(msg *nats.Msg) {
	type Body struct {
		LessonName string `json:"lessonName"`
	}
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	var b Body
	if err := json.Unmarshal(req.Body, &b); err != nil {
		respondErr(msg, 400, "bad_request", "invalid body")
		return
	}
	if err := r.Logic.DeleteLesson(b.LessonName); err != nil {
		respondErr(msg, 500, "server_error", "failed to delete lesson")
		return
	}
	respondOK(msg, nil, 200)
}

func (r *Router) handleClassInsert(msg *nats.Msg) {
	type Body struct {
		LessonName    string `json:"lessonName"`
		ProfessorName string `json:"professorName"`
		Capacity      int    `json:"capacity"`
		ClassNum      int    `json:"classNumber"`
		Date          string `json:"date"`
		Term          string `json:"term"`
	}
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	var b Body
	if err := json.Unmarshal(req.Body, &b); err != nil {
		respondErr(msg, 400, "bad_request", "invalid body")
		return
	}
	if err := r.Logic.InsertClass(b.LessonName, b.ProfessorName, b.Date, b.Term, b.Capacity, b.ClassNum); err != nil {
		respondErr(msg, 500, "server_error", "failed to insert class")
		return
	}
	respondOK(msg, nil, 201)
}

func (r *Router) handleClassesShow(msg *nats.Msg) {
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	res, err := r.Logic.GetAllClassesByTerm(req.Query)
	if err != nil {
		respondErr(msg, 500, "server_error", "failed to get classes")
		return
	}
	respondOK(msg, res, 200)
}

func (r *Router) handleClassDelete(msg *nats.Msg) {
	type Body struct {
		Id int `json:"id"`
	}
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	var b Body
	if err := json.Unmarshal(req.Body, &b); err != nil {
		respondErr(msg, 400, "bad_request", "invalid body")
		return
	}
	if err := r.Logic.DeleteClass(b.Id); err != nil {
		respondErr(msg, 500, "server_error", "failed to delete class")
		return
	}
	respondOK(msg, nil, 200)
}

func (r *Router) handleUsersByRole(msg *nats.Msg) {
	type Body struct {
		RoleId int `json:"roleId"`
	}
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	var b Body
	if err := json.Unmarshal(req.Body, &b); err != nil {
		respondErr(msg, 400, "bad_request", "invalid body")
		return
	}
	res, err := r.Logic.GetUsersByRole(b.RoleId)
	if err != nil {
		respondErr(msg, 500, "server_error", "failed to get users by role")
		return
	}
	respondOK(msg, res, 200)
}

func (r *Router) handleStudentAdd(msg *nats.Msg) {
	type Body struct {
		Id int `json:"id"`
	}
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	var b Body
	if err := json.Unmarshal(req.Body, &b); err != nil {
		respondErr(msg, 400, "bad_request", "invalid body")
		return
	}
	if err := r.Logic.AddStudent(b.Id); err != nil {
		respondErr(msg, 500, "server_error", "failed to add student")
		return
	}
	respondOK(msg, nil, 201)
}

func (r *Router) handleProfessorAdd(msg *nats.Msg) {
	type Body struct {
		Id int `json:"id"`
	}
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireAdmin(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	var b Body
	if err := json.Unmarshal(req.Body, &b); err != nil {
		respondErr(msg, 400, "bad_request", "invalid body")
		return
	}
	if err := r.Logic.AddProfessor(b.Id); err != nil {
		respondErr(msg, 500, "server_error", "failed to add professor")
		return
	}
	respondOK(msg, nil, 201)
}

// Professor handlers
func (r *Router) handleMarkAdd(msg *nats.Msg) {
	type Req struct {
		Mark    int `json:"mark"`
		UserId  int `json:"userId"`
		ClassId int `json:"classId"`
	}
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	if _, err := r.requireProfessor(req.Token); err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	var b Req
	if err := json.Unmarshal(req.Body, &b); err != nil {
		respondErr(msg, 400, "bad_request", "invalid body")
		return
	}
	if err := r.Logic.AddMark(b.UserId, b.ClassId, b.Mark); err != nil {
		respondErr(msg, 500, "server_error", "failed to add mark")
		return
	}
	respondOK(msg, nil, 200)
}

func (r *Router) handleProfessorStudents(msg *nats.Msg) {
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	claims, err := r.requireProfessor(req.Token)
	if err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	res, err := r.Logic.GetStudentsForProfessor(claims.Id)
	if err != nil {
		respondErr(msg, 500, "server_error", "failed to get students")
		return
	}
	respondOK(msg, res, 200)
}

// Student handlers
func (r *Router) handleStudentUnitAdd(msg *nats.Msg) {
	type Body struct {
		Id int `json:"id"`
	}
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	claims, err := r.requireStudent(req.Token)
	if err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	var b Body
	if err := json.Unmarshal(req.Body, &b); err != nil {
		respondErr(msg, 400, "bad_request", "invalid body")
		return
	}
	if err := r.Logic.AddStudentUnit(claims.Id, b.Id); err != nil {
		respondErr(msg, 500, "server_error", "failed to add unit")
		return
	}
	respondOK(msg, nil, 201)
}

func (r *Router) handleStudentUnits(msg *nats.Msg) {
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	claims, err := r.requireStudent(req.Token)
	if err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	res, err := r.Logic.GetClassesByUserId(claims.Id)
	if err != nil {
		respondErr(msg, 500, "server_error", "failed to get units")
		return
	}
	respondOK(msg, res, 200)
}

func (r *Router) handleStudentUnitDelete(msg *nats.Msg) {
	type Body struct {
		Id int `json:"id"`
	}
	var req natsRequest
	_ = json.Unmarshal(msg.Data, &req)
	claims, err := r.requireStudent(req.Token)
	if err != nil {
		respondErr(msg, 403, "forbidden", "forbidden")
		return
	}
	var b Body
	if err := json.Unmarshal(req.Body, &b); err != nil {
		respondErr(msg, 400, "bad_request", "invalid body")
		return
	}
	if err := r.Logic.RemoveStudentUnit(b.Id, claims.Id); err != nil {
		respondErr(msg, 500, "server_error", "failed to delete unit")
		return
	}
	respondOK(msg, nil, 200)
}
