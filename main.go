package main

import (
	"UniWebsite/bussinessLogic"
	"UniWebsite/cache"
	"UniWebsite/databases"
	"UniWebsite/protocols"
	"UniWebsite/restful"
	"UniWebsite/verification"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var Icache cache.ICache
	var Idatabase databases.IDatabase
	var IProtocol protocols.Protocols
	var Iverify verification.ISendVerificationCode
	useCache := "redis"
	useDatabase := "Mysql"
	useProtocol := "http"
	VerifyType := "email"
	redis := cache.NewRedis("localhost:6379")
	if useCache == "redis" {
		Icache = redis
	}
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
	http := protocols.NewHttp(Idatabase, Iverify, Icache)
	if useProtocol == "http" {
		IProtocol = http
	}

	logic := bussinessLogic.NewBussinessLogic(IProtocol)
	r := restful.NewRestFul(logic)
	r.Run()

}
