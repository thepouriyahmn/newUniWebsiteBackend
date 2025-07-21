package bussinessLogic

import "net/http"

type IDatabase interface {
	CheackUserByUsernameAndEmail(ClientUsername, ClientEmail string) error
	InsertUser(username, pass, email string, studentRole, professorRole bool) error
	CheackUserByUserNameAndPassword(username, pass string) (int, string, error)
	GetRole(id int) ([]string, string, error)
}
type IProtocol interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Verify(w http.ResponseWriter, r *http.Request)
}
type AuthBussinessLogic struct {
	IProtocol IProtocol
	IDatabase IDatabase
}

func NewBussinessLogic(protocol IProtocol, database IDatabase) AuthBussinessLogic {
	return AuthBussinessLogic{
		IProtocol: protocol,
		IDatabase: database,
	}

}
