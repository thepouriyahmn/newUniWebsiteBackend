package restful

import (
	"UniWebsite/auth"
	"UniWebsite/bussinessLogic"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

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
	http.HandleFunc("/login", auth.CheackOriginMiddleWare(rest.Login))
	http.HandleFunc("/logout", auth.CheackOriginMiddleWare(rest.Logout))
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
		fmt.Printf("reading error: %v", err)
		panic(err)
	}
}

func (rest Restful) showProfessors(w http.ResponseWriter, r *http.Request) {
	rest.GetAllProfessors(w, r)
}

func (rest Restful) getTerms(w http.ResponseWriter, r *http.Request) {
	rest.GetTerms(w, r)
}

func (rest Restful) showAllUsers(w http.ResponseWriter, r *http.Request) {
	rest.GetAllUsers(w, r)
}

func (rest Restful) insertLesson(w http.ResponseWriter, r *http.Request) {
	rest.InsertLesson(w, r)
}

func (rest Restful) showClasses(w http.ResponseWriter, r *http.Request) {
	rest.ShowClasses(w, r)
}

func (rest Restful) insertClass(w http.ResponseWriter, r *http.Request) {
	rest.InsertClass(w, r)
}

func (rest Restful) deleteClass(w http.ResponseWriter, r *http.Request) {
	rest.DeleteClass(w, r)
}

func (rest Restful) deleteLesson(w http.ResponseWriter, r *http.Request) {
	rest.DeleteLesson(w, r)
}

func (rest Restful) showAllLessons(w http.ResponseWriter, r *http.Request) {
	rest.GetAllLessons(w, r)
}

func (rest Restful) showUsersByRole(w http.ResponseWriter, r *http.Request) {
	rest.GetUsersByRole(w, r)
}

func (rest Restful) addStudent(w http.ResponseWriter, r *http.Request) {
	rest.AddStudent(w, r)
}

func (rest Restful) addProfessor(w http.ResponseWriter, r *http.Request) {
	var req struct{ Id int }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := rest.Bussinesslogic.IDatabase.AddProfessorById(req.Id)
	if err != nil {
		http.Error(w, "Failed to add professor", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Professor Handlers
func (rest Restful) addMark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	type Req struct {
		Mark    int `json:"mark"`
		UserId  int `json:"userId"`
		ClassId int `json:"classId"`
	}
	var req Req

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println("received req: ", req)
	err = rest.Bussinesslogic.IDatabase.AddMark(req.UserId, req.ClassId, req.Mark)
	if err != nil {
		http.Error(w, "Failed to add mark", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rest Restful) showStudentsForProfessor(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)

	students, err := rest.Bussinesslogic.IDatabase.GetStudentsForProfessor(userID)
	if err != nil {
		http.Error(w, "Failed to get students for professor", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(students)
}

// Student Handlers
func (rest Restful) addStudentUnit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	var req struct{ Id int }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("reading error: %v:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println("req id: ", req.Id)
	err := rest.Bussinesslogic.IDatabase.InsertUnitForStudent(userID, req.Id)
	if err != nil {
		http.Error(w, "Failed to add student unit", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (rest Restful) pickedUnits(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	classesSlice, err := rest.Bussinesslogic.IDatabase.GetClassesByUserId(userID)
	if err != nil {
		http.Error(w, "Failed to get student units", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(classesSlice)
	if err != nil {
		http.Error(w, "Failed to show class", http.StatusInternalServerError)
		return
	}
}

func (rest Restful) delStudentUnit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	var req struct{ Id int }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := rest.Bussinesslogic.IDatabase.RemoveStudentUnit(req.Id, userID)
	if err != nil {
		http.Error(w, "Failed to delete student unit", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Business Logic Methods Implementation
func (rest Restful) GetAllProfessors(w http.ResponseWriter, r *http.Request) {
	professorSlice, err := rest.Bussinesslogic.IDatabase.GetAllProfessors()
	if err != nil {
		fmt.Println("http err")
		http.Error(w, "professors not found in db", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(professorSlice)
	if err != nil {
		http.Error(w, "professors not found in db", http.StatusInternalServerError)
	}
}

func (rest Restful) AddProfessor(w http.ResponseWriter, r *http.Request) {
	var req struct{ Id int }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := rest.Bussinesslogic.IDatabase.AddProfessorById(req.Id)
	if err != nil {
		http.Error(w, "Failed to add professor", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (rest Restful) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	fmt.Println("input:", input)
	users, err := rest.Bussinesslogic.IDatabase.GetAllUsers(input)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
}

func (rest Restful) GetUsersByRole(w http.ResponseWriter, r *http.Request) {
	var req struct{ RoleId int }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	users, err := rest.Bussinesslogic.IDatabase.GetUsersByRole(req.RoleId)
	if err != nil {
		http.Error(w, "Failed to get users by role", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (rest Restful) InsertLesson(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LessonName string `json:"lessonName"`
		LessonUnit int    `json:"lessonUnit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := rest.Bussinesslogic.IDatabase.InsertLesson(req.LessonName, req.LessonUnit)
	if err != nil {
		http.Error(w, "Failed to insert lesson", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (rest Restful) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	var req struct{ LessonName string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := rest.Bussinesslogic.IDatabase.DeleteLesson(req.LessonName)
	if err != nil {
		http.Error(w, "Failed to delete lesson", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rest Restful) GetAllLessons(w http.ResponseWriter, r *http.Request) {
	lessons, err := rest.Bussinesslogic.IDatabase.GetAllLessons()
	if err != nil {
		http.Error(w, "Failed to get lessons", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(lessons)
}

func (rest Restful) AddMark(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Mark    int `json:"mark"`
		UserId  int `json:"userId"`
		ClassId int `json:"classId"`
	}
	var req Req

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println("received req: ", req)
	err = rest.Bussinesslogic.IDatabase.AddMark(req.UserId, req.ClassId, req.Mark)
	if err != nil {
		http.Error(w, "Failed to add mark", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rest Restful) GetStudentsForProfessor(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)

	students, err := rest.Bussinesslogic.IDatabase.GetStudentsForProfessor(userID)
	if err != nil {
		http.Error(w, "Failed to get students for professor", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(students)
}

func (rest Restful) AddStudentUnit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	var req struct{ Id int }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("reading error: %v:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println("req id: ", req.Id)
	err := rest.Bussinesslogic.IDatabase.InsertUnitForStudent(userID, req.Id)
	if err != nil {
		http.Error(w, "Failed to add student unit", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (rest Restful) DelStudentUnit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	var req struct{ Id int }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := rest.Bussinesslogic.IDatabase.RemoveStudentUnit(req.Id, userID)
	if err != nil {
		http.Error(w, "Failed to delete student unit", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rest Restful) InsertClass(w http.ResponseWriter, r *http.Request) {
	type Lesson struct {
		LessonName    string `json:"lessonName"`
		ProfessorName string `json:"professorName"`
		Capacity      int    `json:"capacity"`
		ClassNum      int    `json:"classNumber"`
		Date          string `json:"date"`
		Term          string `json:"term"`
	}
	var lesson Lesson

	err := json.NewDecoder(r.Body).Decode(&lesson)
	if err != nil {
		fmt.Printf("reading error: %v:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	fmt.Println("Received Lesson:", lesson)
	err = rest.Bussinesslogic.IDatabase.InsertClass(lesson.LessonName, lesson.ProfessorName, lesson.Date, lesson.Term, lesson.Capacity, lesson.ClassNum)
	if err != nil {
		http.Error(w, "Failed to add class", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (rest Restful) ShowClasses(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	fmt.Println("inputtttttt: ", input)
	classesSlice, err := rest.Bussinesslogic.IDatabase.GetAllClassesByTerm(input)
	if err != nil {
		http.Error(w, "Failed to get classes", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(classesSlice)
	if err != nil {
		http.Error(w, "Failed to show class", http.StatusInternalServerError)
		return
	}
}

func (rest Restful) DeleteClass(w http.ResponseWriter, r *http.Request) {
	type Class struct {
		Id int `json:"id"`
	}
	var class Class
	err := json.NewDecoder(r.Body).Decode(&class)
	if err != nil {
		fmt.Printf("reading error: %v:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	err = rest.Bussinesslogic.IDatabase.DeleteClass(class.Id)
	if err != nil {
		http.Error(w, "Failed to delete class", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rest Restful) AddStudent(w http.ResponseWriter, r *http.Request) {
	type studentRequest struct {
		Id int `json:"id"`
	}
	var student studentRequest

	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		fmt.Printf("reading error: %v:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	err = rest.Bussinesslogic.IDatabase.AddStudentById(student.Id)
	if err != nil {
		http.Error(w, "Failed to add student", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (rest Restful) ShowPickedUnitsForStudent(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	classesSlice, err := rest.Bussinesslogic.IDatabase.GetClassesByUserId(userID)
	if err != nil {
		http.Error(w, "Failed to get student units", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(classesSlice)
	if err != nil {
		http.Error(w, "Failed to show class", http.StatusInternalServerError)
		return
	}
}

func (rest Restful) GetTerms(w http.ResponseWriter, r *http.Request) {
	terms, err := rest.Bussinesslogic.ICache.GetCacheValue("terms")
	if err != nil {
		fmt.Println("err in cache")
		terms, err := rest.Bussinesslogic.IDatabase.GetAllTerms()
		if err != nil {
			http.Error(w, "Failed to get terms", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(terms)

		rest.Bussinesslogic.ICache.CacheTerms(terms)

		if err != nil {
			http.Error(w, "Failed to show terms", http.StatusInternalServerError)
		}
		return
	}
	var termSlice []string
	err = json.Unmarshal([]byte(terms), &termSlice)
	if err != nil {
		fmt.Println("Error unmarshaling:", err)
		return
	}
	err = json.NewEncoder(w).Encode(termSlice)
	if err != nil {
		http.Error(w, "Failed", http.StatusInternalServerError)
		return
	}
}

// Authentication Methods
func (rest Restful) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type User struct {
		Username      string `json:"username"`
		Password      string `json:"password"`
		StudentRole   bool   `json:"studentRole"`
		ProfessorRole bool   `json:"professorRole"`
		Email         string `json:"email"`
	}
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if !auth.IsValidPassword(user.Password) {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	err = rest.Bussinesslogic.IDatabase.CheackUserByUsernameAndEmail(user.Username, user.Email)
	if err != nil {
		http.Error(w, "Username or email already exists", http.StatusConflict)
		return
	}

	err = rest.Bussinesslogic.IDatabase.InsertUser(user.Username, user.Password, user.Email, user.StudentRole, user.ProfessorRole)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type CodeInfo struct {
	Code      string
	CreatedAt time.Time
}

var mu sync.Mutex
var verificationCodes = make(map[int]CodeInfo)

func (rest Restful) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	fmt.Println("RAW BODY:", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	type ClaimedUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var claimedUser ClaimedUser

	err = json.NewDecoder(r.Body).Decode(&claimedUser)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	id, email, err := rest.Bussinesslogic.IDatabase.CheackUserByUserNameAndPassword(claimedUser.Username, claimedUser.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Println("ok")
	fmt.Println("enter send code")
	code, err := rest.Bussinesslogic.IVerify.SendCode(email)
	fmt.Println("still fine")

	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	mu.Lock()
	verificationCodes[id] = CodeInfo{
		Code:      code,
		CreatedAt: time.Now(),
	}
	mu.Unlock()
	fmt.Println(verificationCodes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("id: ", id)
	err = json.NewEncoder(w).Encode(map[string]int{"id": id})
	if err != nil {
		fmt.Printf("reading error: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (rest Restful) Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type ClientInfo struct {
		Id   int    `json:"id"`
		Code string `json:"code"`
	}
	var clientinfo ClientInfo
	err := json.NewDecoder(r.Body).Decode(&clientinfo)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		fmt.Println("reading error:", err)
		return
	}

	roleSlice, username, err := rest.Bussinesslogic.IDatabase.GetRole(clientinfo.Id)
	fmt.Println("roleslice: ", roleSlice)
	if err != nil {
		http.Error(w, "Failed to get user role", http.StatusInternalServerError)
		return
	}
	fmt.Println(verificationCodes)
	mu.Lock()
	userInfo, ok := verificationCodes[clientinfo.Id]
	mu.Unlock()

	fmt.Println("map : ", verificationCodes[clientinfo.Id])
	if !ok {
		http.Error(w, "Unauthorized: code not found", http.StatusUnauthorized)
		return
	}

	if time.Since(userInfo.CreatedAt) > 2*time.Minute {
		mu.Lock()
		delete(verificationCodes, clientinfo.Id)
		mu.Unlock()
		http.Error(w, "Code expired", http.StatusGatewayTimeout)
		return
	}
	fmt.Println("clientCode: ", clientinfo.Code, "userInfo: ", userInfo.Code)
	if clientinfo.Code != userInfo.Code {
		http.Error(w, "Invalid code", http.StatusUnauthorized)
		return
	}

	tokenStr := auth.GenerateJWT(clientinfo.Id, username, roleSlice)
	fmt.Println("token created: ", tokenStr)

	mu.Lock()
	delete(verificationCodes, clientinfo.Id)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
	if err != nil {
		http.Error(w, "Failed to respond with token", http.StatusInternalServerError)
		return
	}
}

func (rest Restful) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	rdb.Set(tokenStr, "blocked", 5*time.Minute)
	w.WriteHeader(http.StatusOK)
}
