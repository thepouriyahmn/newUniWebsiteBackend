package main

import (
	"UniWebsite/auth"
	"UniWebsite/bussinessLogic"
	"UniWebsite/cache"
	"UniWebsite/databases"
	"UniWebsite/restful"
	"UniWebsite/verification"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var Icache cache.ICache
	var Idatabase databases.IDatabase
	var IPassValidation auth.IPassValidation

	var Iverify verification.ISendVerificationCode
	useCache := "redis"
	useDatabase := "Mysql"
	usePassValidation := "regex"

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
	regex := auth.NewRegex()
	if usePassValidation == "regex" {
		IPassValidation = regex
	}

	logic := bussinessLogic.NewBussinessLogic(Idatabase, Icache, Iverify, IPassValidation)
	r := restful.NewRestFul(logic)
	r.Run()

}
