package messagebrokers

import (
	"UniWebsite/bussinessLogic"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type SubDatabase interface {
	// Read operations for CQRS
	GetAllProfessors() ([]bussinessLogic.Professor, error)
	GetAllUsers(input string) ([]bussinessLogic.User, error)
	GetAllLessons() ([]bussinessLogic.Lesson, error)
	GetAllClassesByTerm(term string) ([]bussinessLogic.Classes, error)
	GetAllTerms() ([]string, error)
	GetUsersByRole(roleId int) ([]bussinessLogic.User, error)
	GetStudentsForProfessor(professorId int) ([]bussinessLogic.Student2, error)
	GetClassesByUserId(userID int) ([]bussinessLogic.StudentClasses, error)

	// Sync operations (called by NATS handlers)
	CheackUserByUsernameAndEmail(ClientUsername, ClientEmail string) error
	InsertUser(username, pass, email string, studentRole, professorRole bool) error
	CheackUserByUserNameAndPassword(username, pass string) (int, string, error)
	GetRole(id int) ([]string, string, error)
	InsertLesson(lessonName string, lessonUnit int) error
	InsertClass(lessonName, professorName, date, term string, capacity, classNumber int) error
	DeleteClass(classId int) error
	DeleteLesson(lessonName string) error
	AddProfessorById(userId int) error
	AddStudentById(userId int) error
	AddMark(userId, classId int, mark int) error
	RemoveStudentUnit(classid int, userid int) error
	InsertUnitForStudent(userid int, classid int) error
}

type Nats struct {
	Nc           *nats.Conn
	ISubDatabase SubDatabase
}

func NewNats() Nats {

	nc, _ := nats.Connect(nats.DefaultURL)
	return Nats{
		Nc: nc,
	}

}

func (n Nats) Publish(subject string, req json.RawMessage) error {
	// Check if NATS connection is valid
	if n.Nc == nil || !n.Nc.IsConnected() {
		return fmt.Errorf("NATS connection is not available or disconnected")
	}

	data, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return err
	}

	err = n.Nc.Publish(subject, data)
	if err != nil {
		fmt.Printf("Error publishing to NATS: %v\n", err)
		return err
	}

	fmt.Printf("✅ Event published to NATS: %s\n", subject)
	return nil
}

//	func (n Nats) Subscribe(topic string, cb func()) error {
//		_, err := n.Nc.Subscribe(topic, cb)
//		if err != nil {
//			fmt.Printf("reading error: %v", err)
//			return err
//		}
//		return nil
//	}
func (n Nats) Run() {
	// Check if NATS connection is valid
	if n.Nc == nil || !n.Nc.IsConnected() {
		fmt.Println("❌ Cannot start NATS service: connection is not available")
		return
	}

	fmt.Println("🚀 Starting NATS service...")

	// Subscribe to all CQRS events
	subs := []struct {
		subject string
		handler nats.MsgHandler
	}{
		{"user.created", n.handleUserCreated},
		{"lesson.created", n.handleLessonCreated},
		{"lesson.deleted", n.handleLessonDeleted},
		{"class.created", n.handleClassCreated},
		{"class.deleted", n.handleClassDeleted},
		{"professor.added", n.handleProfessorAdded},
		{"student.added", n.handleStudentAdded},
		{"mark.added", n.handleMarkAdded},
		{"student.unit.added", n.handleStudentUnitAdded},
		{"student.unit.removed", n.handleStudentUnitRemoved},
	}

	// Subscribe to all events
	for _, sub := range subs {
		_, err := n.Nc.Subscribe(sub.subject, sub.handler)
		if err != nil {
			fmt.Printf("❌ Error subscribing to %s: %v\n", sub.subject, err)
		} else {
			fmt.Printf("✅ Subscribed to %s\n", sub.subject)
		}
	}

	fmt.Println("✅ SubDatabase is running in NATS")
	fmt.Println("📡 Waiting for events...")

	// Keep the service running
	select {}
}

// Event handlers for CQRS pattern
func (n Nats) handleUserCreated(msg *nats.Msg) {
	var users bussinessLogic.Users
	if err := json.Unmarshal(msg.Data, &users); err != nil {
		fmt.Printf("Error unmarshaling user data: %v", err)
		return
	}

	// Insert to MySQL (read database)
	if err := n.ISubDatabase.InsertUser(users.Username, users.Password, users.Email, users.StudentRole, users.ProfessorRole); err != nil {
		fmt.Printf("Error inserting user to MySQL: %v", err)
	} else {
		fmt.Printf("User %s successfully synced to MySQL\n", users.Username)
	}
}

func (n Nats) handleLessonCreated(msg *nats.Msg) {
	var lesson bussinessLogic.Lesson
	if err := json.Unmarshal(msg.Data, &lesson); err != nil {
		fmt.Printf("Error unmarshaling lesson data: %v", err)
		return
	}

	// Insert to MySQL (read database)
	if err := n.ISubDatabase.InsertLesson(lesson.LessonName, lesson.LessonUnit); err != nil {
		fmt.Printf("Error inserting lesson to MySQL: %v", err)
	} else {
		fmt.Printf("Lesson %s successfully synced to MySQL\n", lesson.LessonName)
	}
}

func (n Nats) handleLessonDeleted(msg *nats.Msg) {
	var lesson bussinessLogic.Lesson
	if err := json.Unmarshal(msg.Data, &lesson); err != nil {
		fmt.Printf("Error unmarshaling lesson data: %v", err)
		return
	}

	// Delete from MySQL (read database)
	if err := n.ISubDatabase.DeleteLesson(lesson.LessonName); err != nil {
		fmt.Printf("Error deleting lesson from MySQL: %v", err)
	} else {
		fmt.Printf("Lesson %s successfully deleted from MySQL\n", lesson.LessonName)
	}
}

func (n Nats) handleClassCreated(msg *nats.Msg) {
	var class bussinessLogic.Classes
	if err := json.Unmarshal(msg.Data, &class); err != nil {
		fmt.Printf("Error unmarshaling class data: %v", err)
		return
	}

	// Insert to MySQL (read database)
	if err := n.ISubDatabase.InsertClass(class.LessonName, class.ProfessorName, class.Date, class.Term, class.Capacity, class.ClassNumber); err != nil {
		fmt.Printf("Error inserting class to MySQL: %v", err)
	} else {
		fmt.Printf("Class for lesson %s successfully synced to MySQL\n", class.LessonName)
	}
}

func (n Nats) handleClassDeleted(msg *nats.Msg) {
	var class bussinessLogic.Classes
	if err := json.Unmarshal(msg.Data, &class); err != nil {
		fmt.Printf("Error unmarshaling class data: %v", err)
		return
	}

	// Delete from MySQL (read database)
	if err := n.ISubDatabase.DeleteClass(class.Id); err != nil {
		fmt.Printf("Error deleting class from MySQL: %v", err)
	} else {
		fmt.Printf("Class %d successfully deleted from MySQL\n", class.Id)
	}
}

func (n Nats) handleProfessorAdded(msg *nats.Msg) {
	var professor bussinessLogic.Professor
	if err := json.Unmarshal(msg.Data, &professor); err != nil {
		fmt.Printf("Error unmarshaling professor data: %v", err)
		return
	}

	// Add to MySQL (read database)
	if err := n.ISubDatabase.AddProfessorById(professor.Id); err != nil {
		fmt.Printf("Error adding professor to MySQL: %v", err)
	} else {
		fmt.Printf("Professor %d successfully synced to MySQL\n", professor.Id)
	}
}

func (n Nats) handleStudentAdded(msg *nats.Msg) {
	var student struct {
		UserId int `json:"userId"`
	}
	if err := json.Unmarshal(msg.Data, &student); err != nil {
		fmt.Printf("Error unmarshaling student data: %v", err)
		return
	}

	// Add to MySQL (read database)
	if err := n.ISubDatabase.AddStudentById(student.UserId); err != nil {
		fmt.Printf("Error adding student to MySQL: %v", err)
	} else {
		fmt.Printf("Student %d successfully synced to MySQL\n", student.UserId)
	}
}

func (n Nats) handleMarkAdded(msg *nats.Msg) {
	var markData struct {
		UserId  int `json:"userId"`
		ClassId int `json:"classId"`
		Mark    int `json:"mark"`
	}
	if err := json.Unmarshal(msg.Data, &markData); err != nil {
		fmt.Printf("Error unmarshaling mark data: %v", err)
		return
	}

	// Add to MySQL (read database)
	if err := n.ISubDatabase.AddMark(markData.UserId, markData.ClassId, markData.Mark); err != nil {
		fmt.Printf("Error adding mark to MySQL: %v", err)
	} else {
		fmt.Printf("Mark for user %d in class %d successfully synced to MySQL\n", markData.UserId, markData.ClassId)
	}
}

func (n Nats) handleStudentUnitAdded(msg *nats.Msg) {
	var unitData struct {
		UserId  int `json:"userId"`
		ClassId int `json:"classId"`
	}
	if err := json.Unmarshal(msg.Data, &unitData); err != nil {
		fmt.Printf("Error unmarshaling unit data: %v", err)
		return
	}

	// Add to MySQL (read database)
	if err := n.ISubDatabase.InsertUnitForStudent(unitData.UserId, unitData.ClassId); err != nil {
		fmt.Printf("Error adding student unit to MySQL: %v", err)
	} else {
		fmt.Printf("Student unit for user %d in class %d successfully synced to MySQL\n", unitData.UserId, unitData.ClassId)
	}
}

func (n Nats) handleStudentUnitRemoved(msg *nats.Msg) {
	var unitData struct {
		UserId  int `json:"userId"`
		ClassId int `json:"classId"`
	}
	if err := json.Unmarshal(msg.Data, &unitData); err != nil {
		fmt.Printf("Error unmarshaling unit data: %v", err)
		return
	}

	// Remove from MySQL (read database)
	if err := n.ISubDatabase.RemoveStudentUnit(unitData.ClassId, unitData.UserId); err != nil {
		fmt.Printf("Error removing student unit from MySQL: %v", err)
	} else {
		fmt.Printf("Student unit for user %d in class %d successfully removed from MySQL\n", unitData.UserId, unitData.ClassId)
	}
}
