package main

import (
	"UniWebsite/bussinessLogic"
	"UniWebsite/databases"
	"UniWebsite/protocols"
	"UniWebsite/restful"
	"UniWebsite/verification"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var Idatabase databases.IDatabase
	var IProtocol protocols.Protocols
	var Iverify verification.ISendVerificationCode
	useDatabase := "Mysql"
	useProtocol := "http"
	VerifyType := "email"
	email := verification.NewEmail()
	if VerifyType == "email" {
		Iverify = email
	}

	Mysql, err := databases.NewMysql("root:newpassword@tcp(localhost:3306)/hellodb")
	if err != nil {
		fmt.Printf("reding error: %v", err)
		os.Exit(1)
	}
	if useDatabase == "Mysql" {
		Idatabase = Mysql
	}
	if Idatabase == nil {
		fmt.Println("Database is nil. Exiting.")
		os.Exit(1)
	}
	http := protocols.NewHttp(Idatabase, Iverify)
	if useProtocol == "http" {
		IProtocol = http
	}
	if IProtocol == nil {
		fmt.Println("Protocol is nil. Exiting.")
		os.Exit(1)
	}

	logic := bussinessLogic.NewBussinessLogic(IProtocol, Idatabase)
	r := restful.NewRestFul(logic)
	r.Run()

}
