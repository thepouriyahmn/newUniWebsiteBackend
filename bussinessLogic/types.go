package bussinessLogic

import (
	"sync"
	"time"
)

type Professor struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}
type User struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	StudentRole bool   `json:"studentRole"`

	ProfessorRole bool `json:"professorRole"`
}
type Lesson struct {
	LessonName string `json:"lessonName"`
	LessonUnit int    `json:"lessonUnit"`
	Id         int    `json:"id"`
}
type Classes struct {
	LessonName    string `json:"lessonName"`
	LessonUnit    int    `json:"lessonUnit"`
	Date          string `json:"date"`
	Capacity      int    `json:"capacity"`
	ClassNumber   int    `json:"classNumber"`
	ProfessorName string `json:"professorName"`
	Id            int    `json:"id"`
	Term          string `json:"term"`
}
type StudentClasses struct {
	LessonName    string `json:"lessonName"`
	LessonUnit    int    `json:"lessonUnit"`
	Date          string `json:"date"`
	Capacity      int    `json:"capacity"`
	ClassNumber   int    `json:"classNumber"`
	ProfessorName string `json:"professorName"`
	Id            int    `json:"id"`
	Mark          int    `json:"mark"`
}

type Student struct {
	Name   string `json:"name"`
	Id     int    `json:"id"`
	Mark   *int   `json:"mark"`
	Class  int    `json:"class"`
	Lesson string `json:"lesson"`
}
type Student2 struct {
	Name    string `json:"studentName"`
	UserId  int    `json:"userId"`
	Mark    *int   `json:"mark"`
	Class   int    `json:"classNumber"`
	Lesson  string `json:"lessonName"`
	Date    string `json:"date"`
	ClassId int    `json:"classId"`
}
type CodeInfo struct {
	Code      string
	CreatedAt time.Time
}

var mu sync.Mutex
var verificationCodes = make(map[int]CodeInfo)
