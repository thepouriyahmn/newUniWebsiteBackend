package main

import (
	"UniWebsite/auth"
	"UniWebsite/bussinessLogic"
	"UniWebsite/cache"
	restful "UniWebsite/conn"
	"UniWebsite/databases"
	messagebrokers "UniWebsite/messageBrokers"

	"UniWebsite/verification"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var IconnService restful.Service
	var Icache cache.ICache
	//var Idatabase databases.IDatabase
	var SubDatabase databases.SubDatabase
	var mainDatabase databases.MainDatabase
	var IPassValidation auth.IPassValidation
	var IToken auth.IToken
	var Iverify verification.ISendVerificationCode
	var IDatabaseMessageBroker messagebrokers.IMessageBroker
	useCache := "redis"
	useSubDatabase := "mysql"
	useMainDatabase := "mongo"
	usePassValidation := "regex"
	useToken := "jwt"
	VerifyType := "email"
	serviceConn := "nats"
	databaseBroker := "nats"

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

	mongo := databases.NewMongo("mongodb://localhost:27017")
	switch useSubDatabase {
	case "mysql":
		SubDatabase = Mysql
	case "mongo":
		SubDatabase = mongo
	default:
		fmt.Println("no ReadDatabase has been chosen")
		panic("no ReadDatabase selected")
	}
	switch useMainDatabase {
	case "mysql":
		mainDatabase = Mysql

	case "mongo":
		mainDatabase = mongo

	default:
		fmt.Println("no WriteDatabase has been chosen")
		panic("")
	}

	regex := auth.NewRegex()
	if usePassValidation == "regex" {
		IPassValidation = regex
	}
	jwt := auth.NewJwt()
	if useToken == "jwt" {
		IToken = jwt
	}

	// Initialize message broker with SubDatabase
	natsDatabase := messagebrokers.NewNats()
	natsDatabase.ISubDatabase = SubDatabase

	switch databaseBroker {
	case "nats":
		IDatabaseMessageBroker = natsDatabase

	default:
		fmt.Println("you have not choose valid databaseBroker")
		panic("")
	}

	logic := bussinessLogic.NewBussinessLogic(mainDatabase, SubDatabase, IDatabaseMessageBroker, Icache, Iverify, IPassValidation, IToken)
	r := restful.NewRestFul(logic)
	nats := restful.NewNats(logic)

	go IDatabaseMessageBroker.Run()

	// if serviceConn == "nats" {
	// //	nats.Run()
	// 	IconnService = nats
	// } else {
	// 	r.Run()
	// }
	switch serviceConn {
	case "nats":
		IconnService = nats
	case "http":
		IconnService = r
	default:
		fmt.Println("you have not choose valid connServicer")
		panic("")
	}
	// run
	IconnService.Run()

}
