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
	var IToken auth.IToken
	var Iverify verification.ISendVerificationCode
	useCache := "redis"
	useDatabase := "Mysql"
	usePassValidation := "regex"
	useToken := "jwt"
	VerifyType := "email"
	serviceConn := "nats"

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
	jwt := auth.NewJwt()
	if useToken == "jwt" {
		IToken = jwt
	}

	logic := bussinessLogic.NewBussinessLogic(Idatabase, Icache, Iverify, IPassValidation, IToken)
	r := restful.NewRestFul(logic)
	nats := restful.NewNats(logic)
	if serviceConn == "nats" {
		nats.Run()
	} else {
		r.Run()
	}

}
