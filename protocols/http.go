package protocols

import (
	"UniWebsite/auth"
	"UniWebsite/databases"
	"UniWebsite/verification"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Http struct {
	Database         databases.IDatabase
	VerificationCode verification.ISendVerificationCode
}

func NewHttp(database databases.IDatabase, verifyType verification.ISendVerificationCode) Http {
	return Http{Database: database,
		VerificationCode: verifyType,
	}
}
func (h Http) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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
			fmt.Printf("reding error: %v", err)
			//panic(err)
		}
		if !auth.IsValidPassword(user.Password) {
			http.Error(w, "invalid password", http.StatusBadRequest)
			return
		}
		err = h.Database.CheackUserByUsernameAndEmail(user.Username, user.Email)
		if err != nil {
			http.Error(w, "Username or email already exists", http.StatusConflict)
		}

		err = h.Database.InsertUser(user.Username, user.Password, user.Email, user.StudentRole, user.ProfessorRole)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type CodeInfo struct {
	Code      string
	CreatedAt time.Time
}

var mu sync.Mutex // prevent Race Condition
var verificationCodes = make(map[int]CodeInfo)

func (h Http) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
	}
	fmt.Println("RAW BODY:", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if r.Method == "POST" {
		fmt.Println("r.body isس: ", r.Body)
		type ClaimedUser struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var claimedUser ClaimedUser

		err := json.NewDecoder(r.Body).Decode(&claimedUser)

		if err != nil {
			fmt.Printf("reading error: %v", err)
			//panic(err)
		}

		id, email, err := h.Database.CheackUserByUserNameAndPassword(claimedUser.Username, claimedUser.Password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		//err = json.NewEncoder(w).Encode(map[string]string{"token": token})
		//if err != nil {
		//	panic(err)
		//}

		fmt.Println("ok")

		fmt.Println("enter send code")
		code, err := h.VerificationCode.SendCode(email)
		fmt.Println("still fine")

		if err != nil {
			http.Error(w, "service unavaible", http.StatusServiceUnavailable)
		}
		mu.Lock()
		//convert to int

		verificationCodes[id] = CodeInfo{
			Code:      code,
			CreatedAt: time.Now(),
		}

		mu.Unlock()
		fmt.Println(verificationCodes)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // اول این
		fmt.Println("id: ", id)
		err = json.NewEncoder(w).Encode(map[string]int{"id": id})
		if err != nil {
			fmt.Printf("reading error: %v", err)
		}
		fmt.Println("ok2")

	}
}
func (h Http) Verify(w http.ResponseWriter, r *http.Request) {
	type ClientInfo struct {
		Id   int
		Code string
	}
	var clientinfo ClientInfo
	err := json.NewDecoder(r.Body).Decode(&clientinfo)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		fmt.Println("reading error:", err)
		return
	}

	roleSlice, username, err := h.Database.GetRole(clientinfo.Id)
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
