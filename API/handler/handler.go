package handler

import (
	"api/middlewear"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	serviceURL string
}

func NewHandler(url string) Handler {
	return Handler{
		serviceURL: url,
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

func (h Handler) ProxySignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	resp, err := http.Post("http://"+h.serviceURL+"/signUp", "application/json", r.Body)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	fmt.Println("done")
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyLogin called")

	// Read the request body for logging
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	fmt.Println("Request body:", string(reqBodyBytes))

	// Create a new reader for the request body since we consumed it
	reqBody := bytes.NewReader(reqBodyBytes)

	resp, err := http.Post("http://"+h.serviceURL+"/login", "application/json", reqBody)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Read the response body to get the token for logging
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	fmt.Println("Login response:", string(respBodyBytes))

	w.WriteHeader(resp.StatusCode)
	w.Write(respBodyBytes)
}

func (h Handler) ProxyLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyLogout called")
	resp, err := http.Post("http://"+h.serviceURL+"/logout", "application/json", r.Body)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyVerify(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyVerify called")
	resp, err := http.Post("http://"+h.serviceURL+"/verify", "application/json", r.Body)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Read the response body to get the token for logging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	fmt.Println("token response:", string(bodyBytes))

	w.WriteHeader(resp.StatusCode)
	w.Write(bodyBytes)
}

func (h Handler) ProxyShowProfessors(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyShowProfessors called")

	// Create a new request with headers
	req, err := http.NewRequest("GET", "http://"+h.serviceURL+"/showProfessors", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyGetTerms(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyGetTerms called")

	// Create a new request with headers
	req, err := http.NewRequest("GET", "http://"+h.serviceURL+"/getTerms", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyShowAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyShowAllUsers called")
	url := "http://" + h.serviceURL + "/showAllUsers"
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	// Create a new request with headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyInsertLesson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyInsertLesson called")

	// Read the request body
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Create a new request with headers
	req, err := http.NewRequest("POST", "http://"+h.serviceURL+"/insertLesson", bytes.NewReader(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyShowClasses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyShowClasses called")
	url := "http://" + h.serviceURL + "/showClasses"
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	// Create a new request with headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyInsertClass(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyInsertClass called")

	// Read the request body
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Create a new request with headers
	req, err := http.NewRequest("POST", "http://"+h.serviceURL+"/insertClass", bytes.NewReader(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyDeleteClass(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyDeleteClass called")

	// Read the request body
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Create a new request with headers
	req, err := http.NewRequest("POST", "http://"+h.serviceURL+"/deleteClass", bytes.NewReader(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyDeleteLesson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyDeleteLesson called")

	// Read the request body
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Create a new request with headers
	req, err := http.NewRequest("POST", "http://"+h.serviceURL+"/deleteLesson", bytes.NewReader(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyShowAllLessons(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyShowAllLessons called")

	// Create a new request with headers
	req, err := http.NewRequest("GET", "http://"+h.serviceURL+"/showAllLessons", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyShowUsersByRole(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyShowUsersByRole called")

	// Read the request body
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Create a new request with headers
	req, err := http.NewRequest("POST", "http://"+h.serviceURL+"/showUsersByRole", bytes.NewReader(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyAddStudent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyAddStudent called")

	// Read the request body
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Create a new request with headers
	req, err := http.NewRequest("POST", "http://"+h.serviceURL+"/addStudent", bytes.NewReader(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyAddProfessor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyAddProfessor called")

	// Read the request body
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Create a new request with headers
	req, err := http.NewRequest("POST", "http://"+h.serviceURL+"/addProfessor", bytes.NewReader(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyAddMark(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyAddMark called")

	// Read the request body
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Create a new request with headers
	req, err := http.NewRequest("POST", "http://"+h.serviceURL+"/addMark", bytes.NewReader(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyShowStudentsForProfessor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyShowStudentsForProfessor called")

	// Create a new request with headers
	req, err := http.NewRequest("GET", "http://"+h.serviceURL+"/showStudentsForProfessor", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyAddStudentUnit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyAddStudentUnit called")

	// Read the request body
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Create a new request with headers
	req, err := http.NewRequest("POST", "http://"+h.serviceURL+"/add", bytes.NewReader(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyPickedUnits(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyPickedUnits called")

	// Create a new request with headers
	req, err := http.NewRequest("GET", "http://"+h.serviceURL+"/pickedUnits", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h Handler) ProxyDelStudentUnit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyDelStudentUnit called")

	// Read the request body
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Create a new request with headers
	req, err := http.NewRequest("POST", "http://"+h.serviceURL+"/delStudentUnit", bytes.NewReader(reqBodyBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Copy important headers from the original request
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
		fmt.Println("Authorization header set:", authHeader)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err is: ", err)
		http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
