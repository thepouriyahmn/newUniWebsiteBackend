package restful

import (
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
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("reding error: %v", err)
		//	panic(err)
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
