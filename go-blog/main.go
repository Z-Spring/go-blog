package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Post struct {
	Id        int
	Content   string
	Author    string `sql: "not null"`
	CreatedAt time.Time
}

type Student struct {
	Id        int
	Name      string `sql: "not null"`
	Age       int
	CreatedAt *time.Time
}

func main() {
	openDb()
	mux := httprouter.New()
	mux.GET("/student/:name", getStudentInfo)
	mux.DELETE("/student/:name", deleteStudentByName)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func getStudentInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	result, _ := SearchByAuthor(p.ByName("name"))
	fmt.Fprintln(w, result)
}
func deleteStudentByName(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var s *Student
	s.DeleteStudent(p.ByName("name"))
	fmt.Fprint(w,"delete sucess!")
}

func (s *Student) Create() (err error) {
	result := DB.Create(&s)
	affected := result.RowsAffected
	fmt.Println(affected)
	return
}

func SearchByAuthor(name string) (student Student, err error) {
	DB.Where("name=?", name).Find(&student)
	return
}

func (s *Student) DeleteStudent(name string) (err error) {
	DB.Where("name=?", name).Delete(&s)
	return
}

func openDb() {
	var err error
	dbname := "root:admin@tcp(localhost:3306)/chat?parseTime=true"
	DB, err = gorm.Open(mysql.Open(dbname), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&Student{})
	if err != nil {
		log.Fatal(err)
	}
}
