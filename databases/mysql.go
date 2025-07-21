package databases

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // این خط را حتما اضافه کن
	"golang.org/x/crypto/bcrypt"
)

type Mysql struct {
	db *sql.DB
}

func NewMysql(dsn string) (Mysql, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return Mysql{}, err
	}
	return Mysql{db: db}, nil
}
func (m Mysql) CheackUserByUsernameAndEmail(ClientUsername, ClientEmail string) error {
	var userslice []string
	var emailSlice []string
	var user, email string
	rows, err := m.db.Query("SELECT username,email FROM users")
	if err != nil {
		fmt.Printf("reding error: %v", err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(&user, &email)
		if err != nil {
			fmt.Println(err)
		}
		userslice = append(userslice, user)
		emailSlice = append(emailSlice, email)
	}
	for _, v := range userslice {
		if v == ClientUsername {
			fmt.Printf("username already exist")
			return errors.New("")
		}
	}
	for _, v := range emailSlice {
		if v == ClientEmail {
			return errors.New("")
		}
	}
	return nil
}
func (m Mysql) InsertUser(username, pass, email string, studentRole, professorRole bool) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("reding error: %v", err)
		return err
		//panic(err)
	}
	fmt.Println("email: ", email)
	stmt, err := m.db.Prepare("INSERT INTO users(`username`,`password`,`claim_student`,`claim_professor`, `email`) VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Printf("reding error: %v", err)
		return err
		//panic(err)
	}
	_, err = stmt.Exec(username, hashedPassword, studentRole, professorRole, email)
	if err != nil {
		fmt.Printf("reding error: %v", err)
		return err
		//panic(err)
	}
	return nil
}
func (m Mysql) CheackUserByUserNameAndPassword(username, pass string) (int, string, error) {
	fmt.Println("enter")
	var usernameDB, passwordDB, email string
	var id int
	row := m.db.QueryRow("SELECT username,password,ID,email FROM users WHERE username = ?", username)
	err := row.Scan(&usernameDB, &passwordDB, &id, &email)
	fmt.Println("dbUsername: ", usernameDB, "dbpass: ", passwordDB)
	if err != nil {
		fmt.Printf("scan err: %v", err)
		return 0, "", errors.New("username or password is incorrect")
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(pass))
	if err != nil {
		fmt.Printf("reading error: %v", err)
		fmt.Printf("scan err: %v", err)
		return 0, "", errors.New("username or password is incorrect")
	}
	fmt.Println("not error")
	return id, email, nil
}
func (m Mysql) GetRole(id int) ([]string, string, error) {
	var roleslice []string
	var role string
	row, err := m.db.Query("SELECT role_id from user_roles where user_id = ?", id)
	if err != nil {
		fmt.Printf("readingg error: %v", err)
		return []string{}, "", errors.New("not found")
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&role)
		if err != nil {
			fmt.Printf("reading errorr: %v", err)
			//panic(err)
		}
		roleslice = append(roleslice, role)
	}
	var username string
	row2 := m.db.QueryRow("SELECT username from users where ID = ?", id)

	err = row2.Scan(&username)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return []string{}, "", errors.New("")
	}
	return roleslice, username, nil
}
