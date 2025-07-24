package databases

import (
	"UniWebsite/bussinessLogic"
	"database/sql"
	"errors"
	"fmt"
	"log"

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
	fmt.Println("roleslice in db: ", roleslice)
	return roleslice, username, nil
}
func (m Mysql) GetAllProfessors() ([]bussinessLogic.Professor, error) {
	var professor bussinessLogic.Professor
	var professorSlice []bussinessLogic.Professor

	rows, err := m.db.Query("SELECT users.username,user_roles.user_id FROM user_roles INNER JOIN users ON user_roles.user_id=users.ID WHERE user_roles.role_id = ?", 3)
	if err != nil {
		fmt.Printf("reading query error: %v", err)
		return []bussinessLogic.Professor{}, errors.New("")
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&professor.Name, &professor.Id)
		if err != nil {
			fmt.Printf("reading scan error: %v", err)
			return []bussinessLogic.Professor{}, errors.New("")
		}
		professorSlice = append(professorSlice, professor)

	}
	return professorSlice, nil
}

func (m Mysql) AddProfessorById(userId int) error {
	var roleId int
	rows, err := m.db.Query("SELECT role_id FROM user_roles WHERE user_id = ?", userId)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}

	for rows.Next() {
		err = rows.Scan(&roleId)
		if err != nil {
			fmt.Printf("reading error: %v", err)
			return err
		}
		if roleId == 3 {
			return errors.New("")
		}
	}

	_, err = m.db.Exec("INSERT INTO user_roles(user_id,role_id) VALUES (?, ?)", userId, 3)
	if err != nil {

		fmt.Printf("reading error: %v", err)
		return err
	}
	return nil
}

func (m Mysql) AddStudent(userId int) error {
	_, err := m.db.Exec("INSERT INTO user_roles(user_id,role_id) VALUES (?, ?)", userId, 2)
	return err
}

func (m Mysql) GetAllUsers(input string) ([]bussinessLogic.User, error) {
	var users []bussinessLogic.User
	rows, err := m.db.Query("SELECT username,claim_student,claim_professor,ID FROM users WHERE username LIKE ?", "%"+input+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u bussinessLogic.User
		err = rows.Scan(&u.Username, &u.StudentRole, &u.StudentRole, &u.Id)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (m Mysql) InsertLesson(lessonName string, lessonUnit int) error {
	_, err := m.db.Exec("INSERT INTO lessons(lesson_name,lesson_unit) VALUES (?,?)", lessonName, lessonUnit)
	return err
}

func (m Mysql) DeleteLesson(lessonName string) error {
	_, err := m.db.Exec("DELETE FROM lessons WHERE lesson_name = ?", lessonName)
	return err
}

func (m Mysql) GetAllLessons() ([]bussinessLogic.Lesson, error) {
	var lessons []bussinessLogic.Lesson
	rows, err := m.db.Query("SELECT lesson_id,lesson_name,lesson_unit FROM lessons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var l bussinessLogic.Lesson
		err = rows.Scan(&l.Id, &l.LessonName, &l.LessonUnit)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (m Mysql) GetUsersByRole(roleId int) ([]bussinessLogic.User, error) {
	var users []bussinessLogic.User
	rows, err := m.db.Query("SELECT users.username,user_roles.user_id FROM user_roles INNER JOIN users ON user_roles.user_id=users.ID WHERE user_roles.role_id = ?", roleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u bussinessLogic.User
		err = rows.Scan(&u.Username, &u.Id)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (m Mysql) AddMark(userId, classId int, mark int) error {
	fmt.Println("userid: ", userId)
	fmt.Println("classid: ", classId)
	fmt.Println("mark: ", mark)
	_, err := m.db.Exec("UPDATE users_classes SET mark = ? WHERE user_class_id = ? AND class_id = ?", mark, userId, classId)
	return err
}

func (m Mysql) GetStudentsForProfessor(professorId int) ([]bussinessLogic.Student2, error) {
	var students []bussinessLogic.Student2
	fmt.Println("professor id: ", professorId)
	rows, err := m.db.Query(`SELECT u.username, v.lesson_name, v.class_number, v.class_time, v.user_class_id, v.class_id, v.mark FROM users_classes_view v JOIN users u ON v.user_class_id = u.id WHERE v.username = (SELECT username FROM users WHERE ID = ?)`, professorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s bussinessLogic.Student2
		err = rows.Scan(&s.Name, &s.Lesson, &s.Class, &s.Date, &s.UserId, &s.ClassId, &s.Mark)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	fmt.Println("student for pro: ", students)
	return students, nil
}

func (m Mysql) RemoveStudentUnit(classid int, userid int) error {
	stmt, err := m.db.Prepare("DELETE FROM users_classes WHERE class_id = ? AND user_class_id = ?")
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(classid, userid)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	return nil
}
func (m Mysql) InsertClass(lessonName, professorName, date string, capacity, classNumber int) error {
	var professorId int
	professorRow := m.db.QueryRow("SELECT user_roles.user_id FROM user_roles INNER JOIN users ON user_roles.user_id=users.ID where users.username = ?", professorName)
	err := professorRow.Scan(&professorId)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	var lessonId int
	lessonRow := m.db.QueryRow("SELECT lesson_id FROM lessons where lesson_name = ?", lessonName)
	err = lessonRow.Scan(&lessonId)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}

	_, err = m.db.Exec("INSERT INTO classes(`lesson_id`,`professor_id`,`class_number`,`capacity`,`class_time`)VALUES(?,?,?,?,?)", lessonId, professorId, classNumber, capacity, date)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	return nil
}
func (m Mysql) GetAllClasses() ([]bussinessLogic.Classes, error) {
	var registered int
	var classesSlice []bussinessLogic.Classes

	rows, err := m.db.Query("SELECT `lesson_unit`,`lesson_name`,`username`,`class_id`,`class_number`,`capacity`,`class_time` FROM classes_view LIMIT 100")
	if err != nil {
		panic(err)

	}
	defer rows.Close()
	var classes bussinessLogic.Classes
	for rows.Next() {

		err = rows.Scan(&classes.LessonUnit, &classes.LessonName, &classes.ProfessorName, &classes.Id, &classes.ClassNumber, &classes.Capacity, &classes.Date)
		if err != nil {
			fmt.Printf("reading error: %v", err)
			return []bussinessLogic.Classes{}, err
		}
		err = m.db.QueryRow("SELECT COUNT(*) FROM users_classes WHERE class_id = ?", classes.Id).Scan(&registered)
		if err != nil {
			fmt.Printf("reading error: %v", err)
			return []bussinessLogic.Classes{}, err
		}
		fmt.Println("class id: ", classes.Id)

		LeftCapacity := classes.Capacity - registered
		fmt.Println(registered, LeftCapacity, classes.Capacity)
		if LeftCapacity == 0 {
			LeftCapacity = -1
		}

		if LeftCapacity == 0 {
			fmt.Println("full")

		}
		classes.Capacity = LeftCapacity

		classesSlice = append(classesSlice, classes)
		fmt.Println(classesSlice)

	}
	return classesSlice, nil

}
func (m Mysql) DeleteClass(classId int) error {
	stmt, err := m.db.Prepare("DELETE FROM classes WHERE class_id = ?")
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(classId)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	return nil
}
func (m Mysql) AddStudentById(userId int) error {
	var roleId int
	rows, err := m.db.Query("SELECT role_id FROM user_roles WHERE user_id = ?", userId)
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&roleId)
		if err != nil {
			panic(err)
		}
		if roleId == 2 {
			return errors.New("")
		}
	}

	_, err = m.db.Exec("INSERT INTO user_roles(user_id,role_id) VALUES (?, ?)", userId, 2)
	if err != nil {

		return err
	}
	return nil
}
func (m Mysql) GetClassesByUserId(userID int) ([]bussinessLogic.StudentClasses, error) {
	var lesson []bussinessLogic.StudentClasses
	var nullableMark sql.NullInt64
	fmt.Println("userid: ", userID)
	var registered int
	rows, err := m.db.Query("SELECT `lesson_unit`,`lesson_name`,`username`,`class_id`,`class_number`,`capacity`,`class_time`,`mark`FROM users_classes_view WHERE user_class_id = ?", userID)
	if err != nil {
		log.Println("Prepare statement error:", err)

		return []bussinessLogic.StudentClasses{}, err
	}
	defer rows.Close()
	var classes2 bussinessLogic.StudentClasses
	for rows.Next() {

		err = rows.Scan(&classes2.LessonUnit, &classes2.LessonName, &classes2.ProfessorName, &classes2.Id, &classes2.ClassNumber, &classes2.Capacity, &classes2.Date, &nullableMark)

		if err != nil {
			log.Println("scan statement error: ", err)

			return []bussinessLogic.StudentClasses{}, err
		}
		if nullableMark.Valid {
			classes2.Mark = int(nullableMark.Int64)
		} else {
			classes2.Mark = -1 // یا مقدار پیش‌فرضی که می‌خوای بذاری
		}
		err = m.db.QueryRow("SELECT COUNT(*) FROM users_classes WHERE class_id = ?", classes2.Id).Scan(&registered)
		if err != nil {
			panic(err)
		}
		LeftCapacity := classes2.Capacity - registered
		if LeftCapacity == 0 {
			classes2.Capacity = -1
		}
		lesson = append(lesson, classes2)

	}
	fmt.Println("classes: ", lesson)
	return lesson, nil
}
func (m Mysql) InsertUnitForStudent(userid, classid int) error {

	var registered, capacity int
	err := m.db.QueryRow("SELECT COUNT(*) FROM users_classes WHERE class_id = ?", classid).Scan(&registered)
	if err != nil {
		return err
	}

	err = m.db.QueryRow("SELECT capacity FROM classes WHERE class_id = ?", classid).Scan(&capacity)
	if err != nil {
		fmt.Printf("reading errorL %v:", err)

	}

	if capacity-registered == 0 {

		return err
	}
	row, err := m.db.Query("SELECT `class_id`FROM users_classes WHERE `user_class_id` = ?", userid)
	if err != nil {
		fmt.Printf("reading errorL %v:", err)

	}
	defer row.Close()
	var id int

	var lessonId []int
	for row.Next() {
		err = row.Scan(&id)
		if err != nil {
			break
		}
		lessonId = append(lessonId, id)

	}
	fmt.Println(lessonId)
	for _, v := range lessonId {
		if v == classid {
			return err
		}
	}

	_, err = m.db.Exec("INSERT INTO users_classes(`user_class_id`,`class_id`) VALUE (?,?)", userid, classid)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	return nil
}
