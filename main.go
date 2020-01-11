package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Guestbook struct {
	Id      int
	Name    string
	Email   string
	Website string
	Comment string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "db_training"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("views/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDb, err := db.Query("SELECT id_guestbook, name, email FROM tb_guestbook ORDER BY id_guestbook DESC")
	if err != nil {
		panic(err.Error())
	}
	guest := Guestbook{}
	res := []Guestbook{}
	for selDb.Next() {
		var id_guestbook int
		var name, email string
		err = selDb.Scan(&id_guestbook, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		guest.Id = id_guestbook
		guest.Name = name
		guest.Email = email
		res = append(res, guest)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	gId := r.URL.Query().Get("id")
	selDb, err := db.Query("SELECT * FROM tb_guestbook WHERE id_guestbook=?", gId)
	if err != nil {
		panic(err.Error())
	}
	guest := Guestbook{}
	for selDb.Next() {
		var id_guestbook int
		var name, email, website, comment string
		err = selDb.Scan(&id_guestbook, &name, &email, &website, &comment)
		if err != nil {
			panic(err.Error())
		}
		guest.Id = id_guestbook
		guest.Name = name
		guest.Email = email
		guest.Website = website
		guest.Comment = comment
	}
	tmpl.ExecuteTemplate(w, "Show", guest)
	defer db.Close()
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.ListenAndServe(":8080", nil)
}
