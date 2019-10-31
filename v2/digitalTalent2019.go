package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//Object
type responseObject struct {
	Response string
}

type inputData struct {
	People string
}

type updateDataObject struct {
	Name        string
	Temperature string
	Humidity    string
	OldName     string
}

type readDataObject struct {
	Name        string
	Temperature string
	Humidity    string
	LED         string
}

var ledHolder = ""
var tmpl = template.Must(template.ParseFiles("forms.html"))

//Function Helper
func initDatabase(database *sql.DB) *sql.Tx {
	tx, err2 := database.Begin()
	if err2 != nil {
		log.Println(err2)
	}

	stmt, err3 := tx.Prepare("CREATE TABLE IF NOT EXISTS equipmentList (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, temperature TEXT, humidity TEXT)")
	if err3 != nil {
		log.Println(err3)
	}
	stmt.Exec()
	defer stmt.Close()

	return tx

}

func updateResponseParser(request *http.Request) *updateDataObject {
	body, err0 := ioutil.ReadAll(request.Body)
	if err0 != nil {
		log.Println(err0)
	}
	var m updateDataObject
	err1 := json.Unmarshal(body, &m)
	if err1 != nil {
		log.Println(err1)
	}

	return &m
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/createData", createDataHandler)
	mux.HandleFunc("/readData", readDataHandler)
	mux.HandleFunc("/updateData", updateDataHandler)
	mux.HandleFunc("/updateData2", updateDataHandler2)
	mux.HandleFunc("/updateData3", updateDataHandler3)

	http.ListenAndServe(":8080", mux)
}

func createDataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	mName := r.FormValue("name")
	mTemperature := r.FormValue("temperature")
	mHumidity := r.FormValue("humidity")
	log.Println(mName)
	log.Println(mTemperature)
	log.Println(mHumidity)

	database, err0 := sql.Open("sqlite3", "./digitalTalent2019.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	stmt, err1 := tx.Prepare("INSERT INTO equipmentList (name, temperature, humidity) VALUES (?, ?, ?)")
	if err1 != nil {
		log.Println(err1)
	}
	stmt.Exec(mName, mTemperature, mHumidity)
	defer stmt.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	m2 := responseObject{"Create data success"}
	b, err2 := json.Marshal(m2)
	if err2 != nil {
		log.Println(err2)
	}
	w.Write(b)

}

func readDataHandler(w http.ResponseWriter, r *http.Request) {
	database, err0 := sql.Open("sqlite3", "./digitalTalent2019.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	mName := ""
	mTemperature := ""
	mHumidity := ""
	var mDeviceDataList []readDataObject
	rows, err1 := tx.Query("SELECT name, temperature, humidity FROM equipmentList")
	if err1 != nil {
		log.Println(err1)
	}
	for rows.Next() {
		rows.Scan(&mName, &mTemperature, &mHumidity)
		mCondition := "Mati"
		ledHolder2, _ := strconv.ParseInt(ledHolder, 10, 32)
		mTemperature2, _ := strconv.ParseFloat(mTemperature, 32)
		if (mTemperature2 >= 30.00) && (ledHolder2 >= 50) {
			mCondition = "Hidup"

		} else {

			mCondition = "Mati"

		}
		mDeviceDataList = append(mDeviceDataList, readDataObject{mName, mTemperature, mHumidity, mCondition})

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	b, err2 := json.Marshal(mDeviceDataList)
	if err2 != nil {
		log.Println(err2)
	}
	w.Write(b)

}

func updateDataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	mName := r.FormValue("name")
	mTemperature := r.FormValue("temperature")
	mHumidity := r.FormValue("humidity")
	mOldName := r.FormValue("oldName")

	database, err0 := sql.Open("sqlite3", "./digitalTalent2019.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	stmt, err0 := tx.Prepare("UPDATE equipmentList SET name=?, temperature=?, humidity=? WHERE name=?")
	if err0 != nil {
		log.Println(err0)
	}
	stmt.Exec(mName, mTemperature, mHumidity, mOldName)
	defer stmt.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	m2 := responseObject{"Update data success"}
	b, err1 := json.Marshal(m2)
	if err1 != nil {
		log.Println(err1)
	}
	w.Write(b)

}

func updateDataHandler2(w http.ResponseWriter, r *http.Request) {
	m := updateResponseParser(r)

	database, err0 := sql.Open("sqlite3", "./digitalTalent2019.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	stmt, err0 := tx.Prepare("UPDATE equipmentList SET name=?, temperature=?, humidity=? WHERE name=?")
	if err0 != nil {
		log.Println(err0)
	}
	stmt.Exec(m.Name, m.Temperature, m.Humidity, m.OldName)
	defer stmt.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	mValue := m.Temperature
	mValue2, _ := strconv.ParseFloat(mValue, 64)

	mCond := "Mati"
	ledHolder2, _ := strconv.ParseInt(ledHolder, 10, 32)
	if (mValue2 >= 30.00) && (ledHolder2 >= 50) {
		mCond = "Hidup"

	} else {

		mCond = "Mati"
	}

	w.Write([]byte(mCond))

}

func updateDataHandler3(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := inputData{
		People: r.FormValue("people"),
	}

	// do something with details
	inputHolder2 := details.People
	ledHolder = inputHolder2

	log.Println(inputHolder2)

	tmpl.Execute(w, struct{ Success bool }{true})

}
