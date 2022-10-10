package main

import (
	"fmt"
	"net/http"
	"database/sql"
	"html/template"
	"io/ioutil"

	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Trener struct{
	Id string
	Surname string
	Name string
	Patronymic string
}


var database *sql.DB


func CheckError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func CheckError2(err error) {
	if err != nil {
		log.Fatal("Error2:", err)
	}
}

func CheckError3(err error) {
	if err != nil {
		log.Fatal("Error3:", err)
	}
}

func DeleteTrener(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.Exec("DELETE FROM tbtrener WHERE id_passport = $1", id)
	CheckError3(err)

	http.Redirect(w, r, "/", 301)
}

func EditSelectIdTrener(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	rowEditTrener := database.QueryRow("SELECT * FROM tbtrener WHERE id_passport = $1;", id)

	structTrener := Trener{}
	err := rowEditTrener.Scan(&structTrener.Id, &structTrener.Surname, &structTrener.Name, &structTrener.Patronymic)
	CheckError(err)

	if err == nil {
		tmpl, err := template.ParseFiles("templates/viewEditTrener.html")
		CheckError(err)
		tmpl.Execute(w, structTrener)
	} 
	
}

func EditTrener(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	CheckError(err)

	id := r.FormValue("id")
	surname := r.FormValue("surname")
	name := r.FormValue("name")
	patronymic := r.FormValue("patronymic")

	_, err = database.Exec("UPDATE tbtrener SET surname_by_trener=$1, name_by_trener=$2, patronumic_by_trener=$3 WHERE id_passport = $4;", surname, name, patronymic, id)
	CheckError(err)

	http.Redirect(w, r, "/", 301)
}

func CreateTrener(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err :=r.ParseForm()
		CheckError(err)
		id := r.FormValue("id")
		surname := r.FormValue("surname")
		name := r.FormValue("name")
		patronymic := r.FormValue("patronymic")

		_, err = database.Exec("INSERT INTO tbtrener (id_passport, surname_by_trener, name_by_trener, patronumic_by_trener) VALUES ($1, $2, $3, $4);", id, surname, name, patronymic)
		CheckError(err)
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "templates/viewCreateTrener.html")
	}
}

func getTrener(w http.ResponseWriter, r *http.Request) {

	rowsTrener, err := database.Query("SELECT * FROM tbtrener;")
	CheckError(err)
	defer rowsTrener.Close()

	structTrener := []Trener{}

	for rowsTrener.Next(){
		trener := Trener{}
		err := rowsTrener.Scan(&trener.Id, &trener.Surname, &trener.Name, &trener.Patronymic)
		CheckError(err)
		structTrener = append(structTrener, trener)
	}

	tmpl, err := template.ParseFiles("templates/viewTrener.html")
	CheckError(err)
	tmpl.Execute(w, structTrener)
}

func main(){

	connectString, err := ioutil.ReadFile("connect_db.txt")
	CheckError(err)

	db, err := sql.Open("postgres", string(connectString))
	CheckError(err)
	database = db
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", getTrener)
	router.HandleFunc("/create", CreateTrener)
	router.HandleFunc("/edit/{id:[0-9]+}", EditSelectIdTrener).Methods("GET")
	router.HandleFunc("/edit/{id[0-9]+}", EditTrener).Methods("POST")
	router.HandleFunc("/delete/{id:[0-9]+}", DeleteTrener)

	http.Handle("/", router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
