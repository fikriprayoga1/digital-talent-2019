func deleteDataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	mName := r.FormValue("name")

	database, err0 := sql.Open("sqlite3", "./digitalTalent2019.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	stmt, err0 := tx.Prepare("DELETE FROM equipmentList WHERE name=?")
	if err0 != nil {
		log.Println(err0)
	}
	stmt.Exec(mName)
	defer stmt.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

    m2 := responseObject{"Delete data success"}
    b, err1 := json.Marshal(m2)
    if err1 != nil {
		log.Println(err1)
	}
	w.Write(b)

}