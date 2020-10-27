package main

import (
	"sort"
	"time"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"html/template"
	"math/rand"
	"github.com/satori/go.uuid"
	// "github.com/gorilla/sessions"
	_ "github.com/go-sql-driver/mysql"
)
type game struct {
	Root string
	Interval string
	Correct string
	UserAnswer string
	AnswerKey []string
	Message string
	Success bool
	User user
}

type user struct{
	UserName string
	Wins int
	Losses int
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "marcus"
	dbPass := "cifNJ4B4IpQ5nZG9"
	dbName := "music"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
			panic(err.Error())
	}
	return db
}
var tpl *template.Template
var dbUsers = map[string]user{}      
var dbSessions = map[string]string{} 
var win, loss int
var un string
var u user


func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func index(w http.ResponseWriter, r *http.Request){
	roots:=[]string {"C","Db","E","Eb","E","F","Gb","G","Ab","A","Bb","B","C#","D#","F#","G#","A#",}
	intervals:= []string {"Root", "Half", "Whole", "Minor 3rd", "Major 3rd", "Fourth", "Tritone", "Fifth", "Minor 6th", "Major 6th" , "Minor 7th", "Major 7th", "Octave",}
	now := time.Now()
	rand.Seed(now.UnixNano())
	myRootIndex:= rand.Intn(len(roots))
	myIntervalIndex:=rand.Intn(len(intervals))
		
	var g game
	g.Root = roots[myRootIndex]
	g.Interval = intervals[myIntervalIndex]
	g.UserAnswer = r.FormValue("userAnswer")
	
	var root, half, whole, minor3rd, major3rd, fourth, tritone, fifth, minor6th, major6th, minor7th, major7th, octave string
	
	db:= dbConn()
	err := db.QueryRow("SELECT * FROM notes WHERE root =?", g.Root).Scan(&root,	&half, &whole, &minor3rd, &major3rd, &fourth,&tritone,&fifth, &minor6th, &major6th, &minor7th, &major7th, &octave)

	g.AnswerKey = append(g.AnswerKey,half, whole, minor3rd, major3rd,fourth, tritone, fifth, minor6th, major6th, minor7th, major7th, octave)
	if err != nil{
		log.Fatal(err)
	}
	g.Correct = g.AnswerKey[myIntervalIndex-1]

	sort.Sort(sort.StringSlice(g.AnswerKey))

	key := r.FormValue("root")
	uanswer := r.FormValue("userAnswer")
	canswer := r.FormValue("correct")
	interval := r.FormValue("Interval")
//SESSION INFO
// get cookie
	if r.Method == http.MethodPost{
	c, err := r.Cookie("session")
	if err != nil {
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
}
// if the user exists already, get user
if un, ok := dbSessions[c.Value]; ok {
	u = dbUsers[un]
	fmt.Println("*****************************")
	fmt.Println("true")
	fmt.Println("initial sess ", dbSessions)
	fmt.Println("initial users ", u)
} else{
	un = c.Value
	u:= user{un,0,0}
	dbSessions[c.Value] = un
	dbUsers[un] = u
	fmt.Println("*****************************")
		fmt.Println("false", c)
		fmt.Println("false sess ", dbSessions)
		fmt.Println("false users ", dbUsers)
}
}
//SESSION INFO END
	
	if (uanswer == canswer) {
		g.Message = fmt.Sprintf("Correct, %s is the %s from %s", uanswer, interval, key)
		g.Success = true
		u.Wins++
	}else {
		g.Message = fmt.Sprintf("Sorry, %s is NOT the %s from %s",  uanswer, interval, key)
		g.Success = false
		u.Losses++
	}

	g.User = u
	
	err = tpl.ExecuteTemplate(w, "index.gohtml", g)
	if err != nil {
		http.Error(w, err.Error(),500)
		log.Fatalln(err)
	}
}

func main(){
	http.HandleFunc("/", index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8090", nil)
}
