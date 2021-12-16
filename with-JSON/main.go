package main

import (
	"sort"
	"os"
	"io/ioutil"
	"time"
	// "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"html/template"
	"math/rand"
	// "github.com/satori/go.uuid"
	//_ "github.com/go-sql-driver/mysql"
)
type game struct {
	Root string
	Interval string
	Correct string
	UserAnswer string
	AnswerKey []string
	sortedScale []string
	Message string
	Success bool
}

type scale struct {
		Root     string `json:"Root"`
		Half     string `json:"Half"`
		Whole    string `json:"Whole"`
		Minor3Rd string `json:"Minor 3rd"`
		Major3Rd string `json:"Major 3rd"`
		Fourth   string `json:"Fourth"`
		Tritone  string `json:"Tritone"`
		Fifth    string `json:"Fifth"`
		Minor6Th string `json:"Minor 6th"`
		Major6Th string `json:"Major 6th"`
		Minor7Th string `json:"Minor 7th"`
		Major7Th string `json:"Major 7th"`
		Octave   string `json:"Octave"`
}

type scales [] scale

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func index(w http.ResponseWriter, r *http.Request){
	rcvd, err := os.Open("scales.json")
	if err != nil {
		fmt.Println(err)
	}
	defer rcvd.Close()

	byteValue, _ := ioutil.ReadAll(rcvd)
	var data scales
	json.Unmarshal(byteValue, &data)
	
	var roots [] string

	for _, v:=range data{
		roots = append(roots, v.Root)
	}

	intervals:= []string {"Root", "Half", "Whole", "Minor 3rd", "Major 3rd", "Fourth", "Tritone", "Fifth", "Minor 6th", "Major 6th" , "Minor 7th", "Major 7th",}
	now := time.Now()
	rand.Seed(now.UnixNano())
	myRootIndex:= rand.Intn(len(roots))
	myIntervalIndex:=rand.Intn(len(intervals))
	var g game

	g.Root = roots[myRootIndex]

	for _, v:=range data{
		if v.Root == g.Root{
		g.AnswerKey = append(g.AnswerKey, v.Root, v.Half, v.Whole, v.Minor3Rd, v.Major3Rd, v.Fourth, v.Tritone, v.Fifth, v.Minor6Th, v.Major6Th, v.Minor7Th, v.Major7Th )
		}
	}
	
	g.Root = roots[myRootIndex]
	g.Interval = intervals[myIntervalIndex]
	g.Correct = g.AnswerKey[myIntervalIndex]
	g.UserAnswer = r.FormValue("userAnswer")
	if err != nil{
		log.Fatal(err)
	}
	sort.Sort(sort.StringSlice(g.AnswerKey))

	key := r.FormValue("root")
	uanswer := r.FormValue("userAnswer")
	canswer := r.FormValue("correct")
	interval := r.FormValue("Interval")
	

	if (uanswer == canswer) {
		g.Message = fmt.Sprintf("Correct, %s is the %s from %s", uanswer, interval, key)
		g.Success = true
		
	}else {
		g.Message = fmt.Sprintf("Sorry, %s is NOT the %s from %s",  uanswer, interval, key)
		g.Success = false
	}

	err = tpl.ExecuteTemplate(w, "index.gohtml", g)
	if err != nil {
		http.Error(w, err.Error(),500)
		log.Fatalln(err)
	}
}

func main(){
	http.HandleFunc("/", index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server is running...")
	http.ListenAndServe(":8070", nil)
}
