package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/l3vick/go-pharmacy/model"
	"io/ioutil"
	"net/http"
)

func GetTiming(w http.ResponseWriter, r *http.Request, idUser int) {
	timing := model.Timing{}

	selDB, err := dbConnector.Query("SELECT * FROM timing WHERE id=?", idUser)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {
		var idUser int
		var morningTime, afternoonTime, eveningTime string
		var morning, afternoon, evening bool
		err = selDB.Scan(&idUser, &morning, &afternoon, &evening, &morningTime, &afternoonTime, &eveningTime)

		if err != nil {
			panic(err.Error())
		}

		timing.Id_User = idUser
		timing.Morning = morning
		timing.Afternoon = afternoon
		timing.Evening = evening
		timing.Morning_Time = morningTime
		timing.Afternoon_Time = afternoonTime
		timing.Evening_Time = eveningTime
	}

	output, err := json.Marshal(timing)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)
}

func CreateTiming(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var timing model.Timing
	err = json.Unmarshal(b, &timing)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(timing)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(output)

	query := fmt.Sprintf("INSERT INTO `pharmacy_sh`.`timing` (`id_user`, `morning`, `afternoon`, `evening`, `morning_time`, `afternoon_time`, `evening_time`)  VALUES('%d', '%t', '%t', '%t', '%s', '%s', '%s',)", timing.Id_User, timing.Morning, timing.Afternoon, timing.Evening, timing.Morning_Time, timing.Afternoon_Time, timing.Evening_Time)

	fmt.Println(query)
	insert, err := dbConnector.Query(query)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

func UpdateTiming(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nID := vars["id"]

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var timing model.Timing
	err = json.Unmarshal(b, &timing)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(timing)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(output)

	var query string = fmt.Sprintf("UPDATE `pharmacy_sh`.`timing` SET `id_user` = '%d', `morning` = '%t', `afternoon` = '%t', `evening` = '%t', `morning_time` = '%s', `afternoon_time` = '%s', `evening_time` = '%s' WHERE (`id_user` = '%s')", timing.Id_User, timing.Morning, timing.Afternoon, timing.Evening, timing.Morning_Time, timing.Afternoon_Time, timing.Evening_Time, nID)

	fmt.Println(query)
	update, err := dbConnector.Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer update.Close()
}

func DeleteTiming(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nID := vars["id"]

	query := fmt.Sprintf("DELETE FROM `pharmacy_sh`.`timing` WHERE (`id_user` = '%s')", nID)

	fmt.Println(query)
	insert, err := dbConnector.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	defer insert.Close()
}
