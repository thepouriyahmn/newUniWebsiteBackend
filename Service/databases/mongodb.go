package databases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"UniWebsite/bussinessLogic"
)

type Mongo struct {
	db *mongo.Database
}

func NewMongo(url string) Mongo {
	ctx := context.Background()
	db, _ := mongo.Connect(ctx, options.Client().ApplyURI(url))
	return Mongo{
		db: db.Database("uni"),
	}
}

func (m Mongo) CheackUserByUsernameAndEmail(ClientUsername, ClientEmail string) error {
	ctx := context.Background()
	count, err := m.db.Collection("users").CountDocuments(ctx, bson.M{"name": ClientUsername})
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	if count > 0 {
		return errors.New("username exist")
	}

	count, err = m.db.Collection("users").CountDocuments(ctx, bson.M{"email": ClientEmail})
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	if count > 0 {
		return errors.New("email exist")
	}
	return nil
}

func (m Mongo) InsertUser(username, pass, email string, studentRole, professorRole bool) error {

	err := m.CheackUserByUsernameAndEmail(username, email)
	if err != nil {

		return err
	}

	type User struct {
		Name          string    `bson:"name"`
		Email         string    `bson:"email"`
		Password      string    `bson:"password"`
		StudentRole   bool      `bson:"studentRole"`
		ProfessorRole bool      `bson:"professorRole"`
		CreatedAt     time.Time `bson:"createdAt"`
	}

	var user User
	user.Name = username
	user.Password = pass
	user.Email = email
	user.ProfessorRole = professorRole
	user.StudentRole = studentRole
	user.CreatedAt = time.Now()

	ctx := context.Background()
	_, err = m.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	return nil
}

func (m Mongo) CheackUserByUserNameAndPassword(username, pass string) (int, string, error) {
	fmt.Println("using mongo for login")
	ctx := context.Background()

	type User struct {
		ID       int    `bson:"_id"`
		Name     string `bson:"name"`
		Password string `bson:"password"`
		Email    string `bson:"email"`
	}

	var user User
	err := m.db.Collection("users").FindOne(ctx, bson.M{"name": username, "password": pass}).Decode(&user)
	if err != nil {
		return 0, "", errors.New("username or password is incorrect")
	}

	return user.ID, user.Email, nil
}

func (m Mongo) GetRole(id int) ([]string, string, error) {
	ctx := context.Background()

	type User struct {
		Name          string `bson:"name"`
		StudentRole   bool   `bson:"studentRole"`
		ProfessorRole bool   `bson:"professorRole"`
	}

	var user User
	err := m.db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, "", err
	}

	var roles []string
	if user.StudentRole {
		roles = append(roles, "student")
	}
	if user.ProfessorRole {
		roles = append(roles, "professor")
	}

	return roles, user.Name, nil
}

// Write operations for CQRS
func (m Mongo) InsertLesson(lessonName string, lessonUnit int) error {
	ctx := context.Background()

	type Lesson struct {
		LessonName string    `bson:"lessonName"`
		LessonUnit int       `bson:"lessonUnit"`
		CreatedAt  time.Time `bson:"createdAt"`
	}

	lesson := Lesson{
		LessonName: lessonName,
		LessonUnit: lessonUnit,
		CreatedAt:  time.Now(),
	}

	_, err := m.db.Collection("lessons").InsertOne(ctx, lesson)
	if err != nil {
		fmt.Printf("Error inserting lesson: %v", err)
		return err
	}
	return nil
}

func (m Mongo) DeleteLesson(lessonName string) error {
	ctx := context.Background()

	_, err := m.db.Collection("lessons").DeleteOne(ctx, bson.M{"lessonName": lessonName})
	if err != nil {
		fmt.Printf("Error deleting lesson: %v", err)
		return err
	}
	return nil
}

func (m Mongo) InsertClass(lessonName, professorName, date, term string, capacity, classNumber int) error {
	ctx := context.Background()

	type Class struct {
		LessonName    string    `bson:"lessonName"`
		ProfessorName string    `bson:"professorName"`
		Date          string    `bson:"date"`
		Term          string    `bson:"term"`
		Capacity      int       `bson:"capacity"`
		ClassNumber   int       `bson:"classNumber"`
		CreatedAt     time.Time `bson:"createdAt"`
	}

	class := Class{
		LessonName:    lessonName,
		ProfessorName: professorName,
		Date:          date,
		Term:          term,
		Capacity:      capacity,
		ClassNumber:   classNumber,
		CreatedAt:     time.Now(),
	}

	_, err := m.db.Collection("classes").InsertOne(ctx, class)
	if err != nil {
		fmt.Printf("Error inserting class: %v", err)
		return err
	}
	return nil
}

func (m Mongo) DeleteClass(classID int) error {
	ctx := context.Background()

	_, err := m.db.Collection("classes").DeleteOne(ctx, bson.M{"_id": classID})
	if err != nil {
		fmt.Printf("Error deleting class: %v", err)
		return err
	}
	return nil
}

func (m Mongo) AddProfessorById(userId int) error {
	ctx := context.Background()

	_, err := m.db.Collection("users").UpdateOne(
		ctx,
		bson.M{"_id": userId},
		bson.M{"$set": bson.M{"professorRole": true}},
	)
	if err != nil {
		fmt.Printf("Error adding professor role: %v", err)
		return err
	}
	return nil
}

func (m Mongo) AddStudentById(userId int) error {
	ctx := context.Background()

	_, err := m.db.Collection("users").UpdateOne(
		ctx,
		bson.M{"_id": userId},
		bson.M{"$set": bson.M{"studentRole": true}},
	)
	if err != nil {
		fmt.Printf("Error adding student role: %v", err)
		return err
	}
	return nil
}

func (m Mongo) AddMark(userId, classId int, mark int) error {
	ctx := context.Background()

	type Mark struct {
		UserId    int       `bson:"userId"`
		ClassId   int       `bson:"classId"`
		Mark      int       `bson:"mark"`
		CreatedAt time.Time `bson:"createdAt"`
	}

	markData := Mark{
		UserId:    userId,
		ClassId:   classId,
		Mark:      mark,
		CreatedAt: time.Now(),
	}

	_, err := m.db.Collection("marks").InsertOne(ctx, markData)
	if err != nil {
		fmt.Printf("Error adding mark: %v", err)
		return err
	}
	return nil
}

func (m Mongo) RemoveStudentUnit(classId, userId int) error {
	ctx := context.Background()

	_, err := m.db.Collection("studentUnits").DeleteOne(ctx, bson.M{
		"classId": classId,
		"userId":  userId,
	})
	if err != nil {
		fmt.Printf("Error removing student unit: %v", err)
		return err
	}
	return nil
}

func (m Mongo) InsertUnitForStudent(userId, classId int) error {
	ctx := context.Background()

	type StudentUnit struct {
		UserId    int       `bson:"userId"`
		ClassId   int       `bson:"classId"`
		CreatedAt time.Time `bson:"createdAt"`
	}

	unit := StudentUnit{
		UserId:    userId,
		ClassId:   classId,
		CreatedAt: time.Now(),
	}

	_, err := m.db.Collection("studentUnits").InsertOne(ctx, unit)
	if err != nil {
		fmt.Printf("Error inserting student unit: %v", err)
		return err
	}
	return nil
}

// Read operations for CQRS (available in MongoDB)
func (m Mongo) GetAllProfessors() ([]bussinessLogic.Professor, error) {
	ctx := context.Background()

	cursor, err := m.db.Collection("users").Find(ctx, bson.M{"professorRole": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var professors []bussinessLogic.Professor
	for cursor.Next(ctx) {
		var user map[string]interface{}
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		professor := bussinessLogic.Professor{}
		if id, ok := user["_id"].(int); ok {
			professor.Id = id
		}
		if name, ok := user["name"].(string); ok {
			professor.Name = name
		}
		professors = append(professors, professor)
	}

	return professors, nil
}

func (m Mongo) GetAllUsers(input string) ([]bussinessLogic.User, error) {
	ctx := context.Background()

	var filter bson.M
	if input != "" {
		filter = bson.M{"$or": []bson.M{
			{"name": bson.M{"$regex": input, "$options": "i"}},
			{"email": bson.M{"$regex": input, "$options": "i"}},
		}}
	} else {
		filter = bson.M{}
	}

	cursor, err := m.db.Collection("users").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []bussinessLogic.User
	for cursor.Next(ctx) {
		var user map[string]interface{}
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		userObj := bussinessLogic.User{}
		if id, ok := user["_id"].(int); ok {
			userObj.Id = id
		}
		if username, ok := user["name"].(string); ok {
			userObj.Username = username
		}
		if password, ok := user["password"].(string); ok {
			userObj.Password = password
		}
		if studentRole, ok := user["studentRole"].(bool); ok {
			userObj.StudentRole = studentRole
		}
		if professorRole, ok := user["professorRole"].(bool); ok {
			userObj.ProfessorRole = professorRole
		}
		users = append(users, userObj)
	}

	return users, nil
}

func (m Mongo) GetAllLessons() ([]bussinessLogic.Lesson, error) {
	ctx := context.Background()

	cursor, err := m.db.Collection("lessons").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var lessons []bussinessLogic.Lesson
	for cursor.Next(ctx) {
		var lesson map[string]interface{}
		if err := cursor.Decode(&lesson); err != nil {
			continue
		}

		lessonObj := bussinessLogic.Lesson{}
		if id, ok := lesson["_id"].(int); ok {
			lessonObj.Id = id
		}
		if lessonName, ok := lesson["lessonName"].(string); ok {
			lessonObj.LessonName = lessonName
		}
		if lessonUnit, ok := lesson["lessonUnit"].(int); ok {
			lessonObj.LessonUnit = lessonUnit
		}
		lessons = append(lessons, lessonObj)
	}

	return lessons, nil
}

func (m Mongo) GetAllClassesByTerm(term string) ([]bussinessLogic.Classes, error) {
	ctx := context.Background()

	cursor, err := m.db.Collection("classes").Find(ctx, bson.M{"term": term})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var classes []bussinessLogic.Classes
	for cursor.Next(ctx) {
		var class map[string]interface{}
		if err := cursor.Decode(&class); err != nil {
			continue
		}

		classObj := bussinessLogic.Classes{}
		if id, ok := class["_id"].(int); ok {
			classObj.Id = id
		}
		if lessonName, ok := class["lessonName"].(string); ok {
			classObj.LessonName = lessonName
		}
		if lessonUnit, ok := class["lessonUnit"].(int); ok {
			classObj.LessonUnit = lessonUnit
		}
		if date, ok := class["date"].(string); ok {
			classObj.Date = date
		}
		if capacity, ok := class["capacity"].(int); ok {
			classObj.Capacity = capacity
		}
		if classNumber, ok := class["classNumber"].(int); ok {
			classObj.ClassNumber = classNumber
		}
		if professorName, ok := class["professorName"].(string); ok {
			classObj.ProfessorName = professorName
		}
		if termStr, ok := class["term"].(string); ok {
			classObj.Term = termStr
		}
		classes = append(classes, classObj)
	}

	return classes, nil
}

func (m Mongo) GetAllTerms() ([]string, error) {
	ctx := context.Background()

	cursor, err := m.db.Collection("classes").Distinct(ctx, "term", bson.M{})
	if err != nil {
		return nil, err
	}

	var terms []string
	for _, term := range cursor {
		if termStr, ok := term.(string); ok {
			terms = append(terms, termStr)
		}
	}

	return terms, nil
}

func (m Mongo) GetUsersByRole(roleId int) ([]bussinessLogic.User, error) {
	ctx := context.Background()

	var filter bson.M
	switch roleId {
	case 1: // Student
		filter = bson.M{"studentRole": true}
	case 2: // Professor
		filter = bson.M{"professorRole": true}
	default:
		return nil, fmt.Errorf("invalid role ID: %d", roleId)
	}

	cursor, err := m.db.Collection("users").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []bussinessLogic.User
	for cursor.Next(ctx) {
		var user map[string]interface{}
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		userObj := bussinessLogic.User{}
		if id, ok := user["_id"].(int); ok {
			userObj.Id = id
		}
		if username, ok := user["name"].(string); ok {
			userObj.Username = username
		}
		if password, ok := user["password"].(string); ok {
			userObj.Password = password
		}
		if studentRole, ok := user["studentRole"].(bool); ok {
			userObj.StudentRole = studentRole
		}
		if professorRole, ok := user["professorRole"].(bool); ok {
			userObj.ProfessorRole = professorRole
		}
		users = append(users, userObj)
	}

	return users, nil
}

func (m Mongo) GetStudentsForProfessor(professorId int) ([]bussinessLogic.Student2, error) {
	ctx := context.Background()

	// First get classes taught by this professor
	cursor, err := m.db.Collection("classes").Find(ctx, bson.M{"professorId": professorId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var classIds []int
	for cursor.Next(ctx) {
		var class map[string]interface{}
		if err := cursor.Decode(&class); err != nil {
			continue
		}
		if classId, ok := class["_id"].(int); ok {
			classIds = append(classIds, classId)
		}
	}

	// Then get students enrolled in these classes
	var students []bussinessLogic.Student2
	for _, classId := range classIds {
		cursor, err := m.db.Collection("studentUnits").Find(ctx, bson.M{"classId": classId})
		if err != nil {
			continue
		}

		for cursor.Next(ctx) {
			var unit map[string]interface{}
			if err := cursor.Decode(&unit); err != nil {
				continue
			}
			if userId, ok := unit["userId"].(int); ok {
				// Get user details
				var user map[string]interface{}
				err := m.db.Collection("users").FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
				if err == nil {
					student := bussinessLogic.Student2{}
					if name, ok := user["name"].(string); ok {
						student.Name = name
					}
					student.UserId = userId
					student.ClassId = classId

					// Get class details
					var class map[string]interface{}
					err := m.db.Collection("classes").FindOne(ctx, bson.M{"_id": classId}).Decode(&class)
					if err == nil {
						if classNumber, ok := class["classNumber"].(int); ok {
							student.Class = classNumber
						}
						if lessonName, ok := class["lessonName"].(string); ok {
							student.Lesson = lessonName
						}
						if date, ok := class["date"].(string); ok {
							student.Date = date
						}
					}

					// Get mark if exists
					var mark map[string]interface{}
					err = m.db.Collection("marks").FindOne(ctx, bson.M{"userId": userId, "classId": classId}).Decode(&mark)
					if err == nil {
						if markValue, ok := mark["mark"].(int); ok {
							student.Mark = &markValue
						}
					}

					students = append(students, student)
				}
			}
		}
		cursor.Close(ctx)
	}

	return students, nil
}

func (m Mongo) GetClassesByUserId(userID int) ([]bussinessLogic.StudentClasses, error) {
	ctx := context.Background()

	// Get classes where user is enrolled
	cursor, err := m.db.Collection("studentUnits").Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var classes []bussinessLogic.StudentClasses
	for cursor.Next(ctx) {
		var unit map[string]interface{}
		if err := cursor.Decode(&unit); err != nil {
			continue
		}
		if classId, ok := unit["classId"].(int); ok {
			// Get class details
			var class map[string]interface{}
			err := m.db.Collection("classes").FindOne(ctx, bson.M{"_id": classId}).Decode(&class)
			if err == nil {
				classObj := bussinessLogic.StudentClasses{}
				if id, ok := class["_id"].(int); ok {
					classObj.Id = id
				}
				if lessonName, ok := class["lessonName"].(string); ok {
					classObj.LessonName = lessonName
				}
				if lessonUnit, ok := class["lessonUnit"].(int); ok {
					classObj.LessonUnit = lessonUnit
				}
				if date, ok := class["date"].(string); ok {
					classObj.Date = date
				}
				if capacity, ok := class["capacity"].(int); ok {
					classObj.Capacity = capacity
				}
				if classNumber, ok := class["classNumber"].(int); ok {
					classObj.ClassNumber = classNumber
				}
				if professorName, ok := class["professorName"].(string); ok {
					classObj.ProfessorName = professorName
				}

				// Get mark if exists
				var mark map[string]interface{}
				err := m.db.Collection("marks").FindOne(ctx, bson.M{"userId": userID, "classId": classId}).Decode(&mark)
				if err == nil {
					if markValue, ok := mark["mark"].(int); ok {
						classObj.Mark = markValue
					}
				}

				classes = append(classes, classObj)
			}
		}
	}

	return classes, nil
}

// Sync operations (called by NATS handlers)
// All these methods are already implemented above as Write operations
// They can be used for both MainDatabase and SubDatabase interfaces
