package restful

import (
	"UniWebsite/bussinessLogic"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

// نوع هندلر
type HandlerFunc func(m *nats.Msg)

// Router
type Router struct {
	Nc *nats.Conn
	// handlers map[string]nats.MsgHandler
	Logic bussinessLogic.Bussinesslogic
}

func NewNats(logic bussinessLogic.Bussinesslogic) *Router {
	nc, _ := nats.Connect(nats.DefaultURL)
	return &Router{
		Nc: nc,
		// handlers: make(map[string]nats.MsgHandler),
		Logic: logic,
	}
}

// ثبت هندلر
// func (r *Router) Handle(subject string, handler nats.MsgHandler) {
// 	r.handlers[subject] = handler
// }

func (r *Router) Run() {
	// فقط یک بار subscribe
	_, err := r.Nc.Subscribe("signup.request", r.handleSignup)
	if err != nil {
		log.Println("subscribe error:", err)
	}
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

	var user User
	if err := json.Unmarshal(msg.Data, &user); err != nil {
		log.Println("json error:", err)
		_ = msg.Respond([]byte(`{"status":"error","msg":"invalid json"}`))
		return
	}

	// فراخوانی لاجیک
	if err := r.Logic.SignUp(user.Username, user.Password, user.Email, user.StudentRole, user.ProfessorRole); err != nil {
		log.Println("signup error:", err)
		_ = msg.Respond([]byte(`{"status":"error","msg":"signup failed"}`))
		return
	}

	// جواب موفقیت
	_ = msg.Respond([]byte(`{"status":"ok"}`))

}
