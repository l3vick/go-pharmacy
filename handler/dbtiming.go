package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/l3vick/go-pharmacy/model"
	"github.com/l3vick/go-pharmacy/util"
)

func GetTiming(idUser string, w http.ResponseWriter, r *http.Request) model.TimingResponse {
	var timingResponse model.TimingResponse

	selDB, err := dbConnector.Query("SELECT morning, afternoon, evening, morning_time, afternoon_time, evening_time FROM timing WHERE id_user=?", idUser)

	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		var morningTime, afternoonTime, eveningTime string
		var morning, afternoon, evening byte
		var morningBoolean, afternoonBoolean, eveningBoolean bool
		err = selDB.Scan(&morning, &afternoon, &evening, &morningTime, &afternoonTime, &eveningTime)

		if err != nil {
			panic(err.Error())
		}

		//TODO: Refactor, make func global
		morningBoolean = util.ByteToBool(morning)
		afternoonBoolean = util.ByteToBool(afternoon)
		eveningBoolean = util.ByteToBool(evening)

		timingResponse.Morning = morningBoolean
		timingResponse.Afternoon = afternoonBoolean
		timingResponse.Evening = eveningBoolean
		timingResponse.Morning_Time = morningTime
		timingResponse.Afternoon_Time = afternoonTime
		timingResponse.Evening_Time = eveningTime
	}

	return timingResponse
}

func GetTimingByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nID := vars["id"]

	var timingResponse model.TimingResponse

	selDB, err := dbConnector.Query("SELECT morning, afternoon, evening, morning_time, afternoon_time, evening_time FROM timing WHERE id_user = " + nID + "")

	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		var morningTime, afternoonTime, eveningTime string
		var morning, afternoon, evening byte
		var morningBoolean, afternoonBoolean, eveningBoolean bool
		err = selDB.Scan(&morning, &afternoon, &evening, &morningTime, &afternoonTime, &eveningTime)

		if err != nil {
			panic(err.Error())
		}

		morningBoolean = util.ByteToBool(morning)
		fmt.Println("morningBoolean %t", morningBoolean)
		afternoonBoolean = util.ByteToBool(afternoon)
		eveningBoolean = util.ByteToBool(evening)

		timingResponse.Morning = morningBoolean
		timingResponse.Afternoon = afternoonBoolean
		timingResponse.Evening = eveningBoolean
		timingResponse.Morning_Time = morningTime
		timingResponse.Afternoon_Time = afternoonTime
		timingResponse.Evening_Time = eveningTime
	}

	output, err := json.Marshal(timingResponse)
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

	query := fmt.Sprintf("INSERT INTO `pharmacy_sh`.`timing` (`id_user`, `morning`, `afternoon`, `evening`, `morning_time`, `afternoon_time`, `evening_time`)  VALUES('%d', '%d', '%d', '%d', '%s', '%s', '%s')", timing.Id_User, timing.Morning, timing.Afternoon, timing.Evening, timing.Morning_Time, timing.Afternoon_Time, timing.Evening_Time)

	fmt.Println(query)
	insert, err := dbConnector.Query(query)

	var timingResponse model.RequestResponse
	if err != nil {
		timingResponse.Code = 500
		timingResponse.Message = err.Error()
	} else {
		timingResponse.Code = 200
		timingResponse.Message = "Timing creado con éxito"
	}

	output, err2 := json.Marshal(timingResponse)
	if err2 != nil {
		http.Error(w, err.Error(), 501)
		return
	}

	w.Write(output)

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

	var timing model.TimingDB
	err = json.Unmarshal(b, &timing)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	timingAux := model.Timing{
		Morning:        util.BoolToByte(timing.Morning),
		Afternoon:      util.BoolToByte(timing.Afternoon),
		Evening:        util.BoolToByte(timing.Evening),
		Morning_Time:   timing.MorningTime,
		Afternoon_Time: timing.AfternoonTime,
		Evening_Time:   timing.EveningTime,
	}

	var query string = fmt.Sprintf("UPDATE `pharmacy_sh`.`timing` SET `morning` = '%d', `afternoon` = '%d', `evening` = '%d', `morning_time` = '%s', `afternoon_time` = '%s', `evening_time` = '%s' WHERE (`id_user` = '%s')", timingAux.Morning, timingAux.Afternoon, timingAux.Evening, timingAux.Morning_Time, timingAux.Afternoon_Time, timingAux.Evening_Time, nID)

	fmt.Println(query)
	update, err := dbConnector.Query(query)

	var timingResponse model.RequestResponse
	if err != nil {
		timingResponse.Code = 500
		timingResponse.Message = err.Error()
	} else {
		timingResponse.Code = 200
		timingResponse.Message = "Timing actualizado con éxito"
	}

	output, err := json.Marshal(timingResponse)
	if err != nil {
		http.Error(w, err.Error(), 501)
		return
	}

	w.Write(output)

	defer update.Close()
}

func DeleteTiming(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nID := vars["id"]

	query := fmt.Sprintf("DELETE FROM `pharmacy_sh`.`timing` WHERE (`id_user` = '%s')", nID)

	fmt.Println(query)
	insert, err := dbConnector.Query(query)

	var timingResponse model.RequestResponse
	if err != nil {
		timingResponse.Code = 500
		timingResponse.Message = err.Error()
	} else {
		timingResponse.Code = 200
		timingResponse.Message = "Timing borrado con éxito"
	}

	output, err := json.Marshal(timingResponse)
	if err != nil {
		http.Error(w, err.Error(), 501)
		return
	}

	w.Write(output)

	defer insert.Close()
}
