package databases

import (
	"testing"
)

func SetupDatabase(t *testing.T) Mysql {
	db, err := NewMysql("root:newpassword@tcp(localhost:3306)/hellodb")
	if err != nil {
		t.Errorf("cannot connect to DB: %v", err)
	}

	return db

}

func TestInsertUser(t *testing.T) {
	db := SetupDatabase(t)
	username := "paradise"
	password := "123qweQWE"
	email := "test@example.com"

	// تست InsertUser
	err := db.InsertUser(username, password, email, true, false)
	if err != nil {
		t.Errorf("InsertUser failed: %v", err)
	}

	err = db.CheackUserByUsernameAndEmail(username, "newemail@example.com")
	if err == nil {
		t.Errorf("CheackUserByUsernameAndEmail should return error for duplicate username")
	}

	err = db.CheackUserByUsernameAndEmail("newuser", email)
	if err == nil {
		t.Errorf("CheackUserByUsernameAndEmail should return error for duplicate email")
	}

	err = db.CheackUserByUsernameAndEmail("anotheruser", "another@example.com")
	if err != nil {
		t.Errorf("CheackUserByUsernameAndEmail failed unexpectedly: %v", err)
	}
	_, err = db.db.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		t.Errorf("delete failed unexpectedly: %v", err)
	}

}
