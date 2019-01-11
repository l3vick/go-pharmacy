package main

import (
	"database/sql"
	"encoding/json"
	_ "errors"
	"fmt"
	"github.com/l3vick/go-pharmacy/model"
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func root(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "App Farmacias" + message
	w.Write([]byte(message))
}

func conectDB() {
	var err error
	db, err = sql.Open("mysql", "rds_pharmacy_00"+":"+"phar00macy"+"@tcp("+"rdspharmacy00.ctiytnyzqbi7.us-east-2.rds.amazonaws.com:3306"+")/"+"rds_pharmacy")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
}

func closeDB() {
	defer db.Close()
}

func GetMeds(w http.ResponseWriter, r *http.Request) {

	pageNumber := r.URL.Query().Get("page")

	intPage, err := strconv.Atoi(pageNumber)

	elementsPage := intPage * 10

	elem := strconv.Itoa(elementsPage) 

	query := fmt.Sprintf("SELECT id, name, pvp, (SELECT COUNT(*)  from rds_pharmacy.med) as count FROM med LIMIT " + elem + ",10")

	fmt.Println(query)

	var meds []*model.Med

	var page model.Page

	selDB, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {
		var id int
		var name string
		var pvp int
		var count int
		err = selDB.Scan(&id, &name, &pvp, &count)
		if err != nil {
			panic(err.Error())
		}
		med := model.Med{
			ID:   id,
			Name: name,
			Pvp:  pvp,
		}
		meds = append(meds, &med)

		var index int
		if (count % 10 == 0){
			index = 1
		}else{
			index = 0
		}
		if intPage == 0 {
			page.First = 0
			page.Previous = 0
			page.Next = intPage+1
			page.Last = (count/10) - index
			page.Count = count
		} else if intPage == (count/10) - index {
			page.First = 0
			page.Previous = intPage -1
			page.Next = intPage
			page.Last = (count/10) - index
			page.Count = count
		} else {
			page.First = 0
			page.Previous = intPage-1
			page.Next = intPage+1
			page.Last = (count/10) - index
			page.Count = count
		}

	}
	response := model.MedResponse{
		Meds: meds,
		Page: page,
	}
	output, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)	
}

func GetMed(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nID := vars["id"]

	selDB, err := db.Query("SELECT * FROM med WHERE id=?", nID)
	if err != nil {
		panic(err.Error())
	}

	med := model.Med{}
	for selDB.Next() {
		var id, pvp int
		var name string
		err = selDB.Scan(&id, &name, &pvp)
		if err != nil {
			panic(err.Error())
		}
		med.ID = id
		med.Name = name
		med.Pvp = pvp
	}

	output, err := json.Marshal(med)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)
}

func CreateMed(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var med model.Med
	err = json.Unmarshal(b, &med)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println(med.Name)
	output, err := json.Marshal(med)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(output)

	query := fmt.Sprintf("INSERT INTO `rds_pharmacy`.`med` (`name`, `pvp`) VALUES('%s','%d')", med.Name, med.Pvp)

	fmt.Println(query)
	insert, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

func UpdateMed(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nID := vars["id"]

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var med model.Med
	err = json.Unmarshal(b, &med)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(med)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(output)

	var query string = fmt.Sprintf("UPDATE `rds_pharmacy`.`med` SET `name` = '%s', `pvp` = '%d' WHERE (`id` = '%s')", med.Name, med.Pvp, nID)

	fmt.Println(query)
	update, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer update.Close()
}

func DeleteMed(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nID := vars["id"]

	query := fmt.Sprintf("DELETE FROM `rds_pharmacy`.`med` WHERE (`id` = '%s')", nID)

	fmt.Println(query)
	insert, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	defer insert.Close()
}


func GetUsers(w http.ResponseWriter, r *http.Request) {

	pageNumber := r.URL.Query().Get("page")

	intPage, err := strconv.Atoi(pageNumber)

	elementsPage := intPage * 10

	elem := strconv.Itoa(elementsPage) 

	query := fmt.Sprintf("SELECT id, name, med_breakfast, med_launch, med_dinner, alarm_breakfast, alarm_launch, alarm_dinner, id_pharmacy, (SELECT COUNT(*)  from rds_pharmacy.users) as count FROM users LIMIT " + elem + ",10 ")

	fmt.Println(query)

	var users []*model.User

	var page model.Page

	selDB, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {
		var id, idPharmacy, count int
		var name, medBreakfast, medLaunch, medDinner, alarmBreakfast, alarmLaunch, alarmDinner, password string
		err = selDB.Scan(&id, &name, &medBreakfast, &medLaunch, &medDinner, &alarmBreakfast, &alarmLaunch, &alarmDinner, &idPharmacy, &count)

		if err != nil {
			panic(err.Error())
		}

		user := model.User{
			ID:             id,
			Name:           name,
			MedBreakfast:   medBreakfast,
			MedLaunch:      medLaunch,
			MedDinner:      medDinner,
			AlarmBreakfast: alarmBreakfast,
			AlarmLaunch:    alarmLaunch,
			AlarmDinner:    alarmDinner,
			Password:       password,
			IDPharmacy:     idPharmacy,
		}

		users = append(users, &user)

		var index int
		if (count % 10 == 0){
			index = 1
		}else{
			index = 0
		}
		if intPage == 0 {
			page.First = 0
			page.Previous = 0
			page.Next = intPage+1
			page.Last = (count/10) - index
			page.Count = count
		} else if intPage == (count/10) - index {
			page.First = 0
			page.Previous = intPage -1
			page.Next = intPage
			page.Last = (count/10) - index
			page.Count = count
		} else {
			page.First = 0
			page.Previous = intPage-1
			page.Next = intPage+1
			page.Last = (count/10) - index
			page.Count = count
		}

	}
	response := model.UserResponse{
		Users: users,
		Page: page,
	}
	output, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)	
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	nID := r.URL.Query().Get("id")
	user := model.User{}

	selDB, err := db.Query("SELECT * FROM users WHERE id=?", nID)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {
		var id, idPharmacy int
		var name, medBreakfast, medLaunch, medDinner, alarmBreakfast, alarmLaunch, alarmDinner, password string
		err = selDB.Scan(&id, &name, &medBreakfast, &medDinner, &medLaunch, &alarmDinner, &alarmLaunch, &alarmBreakfast, &password, &idPharmacy)

		if err != nil {
			panic(err.Error())
		}

		user.ID = id
		user.Name = name
		user.MedBreakfast = medBreakfast
		user.MedLaunch = medLaunch
		user.MedDinner = medDinner
		user.AlarmBreakfast = alarmBreakfast
		user.AlarmLaunch = alarmLaunch
		user.AlarmDinner = alarmDinner
		user.IDPharmacy = idPharmacy
		user.Password = password
	}

	userJSON, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		// handle error
	}

	w.Write([]byte(userJSON))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	var user model.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(output)

	query := fmt.Sprintf("INSERT INTO `rds_pharmacy`.`users` (`name`, `med_breakfast`, `med_launch`, `med_dinner`, `alarm_breakfast`, `alarm_launch`, `alarm_dinner`, `password`, `id_pharmacy`)  VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d')", user.Name, user.MedBreakfast, user.MedLaunch, user.MedDinner, user.AlarmBreakfast, user.AlarmLaunch, user.AlarmDinner, user.Password, user.IDPharmacy)

	fmt.Println(query)
	insert, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nID := vars["id"]

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var user model.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(output)

	var query string = fmt.Sprintf("UPDATE `rds_pharmacy`.`users` SET `name` = '%s', `med_breakfast` = '%s', `med_launch` = '%s', `med_dinner` = '%s', `alarm_breakfast` = '%s', `alarm_launch` = '%s', `alarm_dinner` = '%s', `id_pharmacy` = '%d' WHERE (`id` = '%s')", user.Name, user.MedBreakfast, user.MedLaunch, user.MedDinner, user.AlarmBreakfast, user.AlarmBreakfast, user.AlarmBreakfast, user.IDPharmacy, nID)

	fmt.Println(query)
	update, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer update.Close()
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nID := vars["id"]

	query := fmt.Sprintf("DELETE FROM `rds_pharmacy`.`users` WHERE (`id` = '%s')", nID)

	fmt.Println(query)
	insert, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	defer insert.Close()
}

func GetPharmacies(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	page := vars["page"]

	var pharmacies []*model.Pharmacy
	selDB, err := db.Query("SELECT id, cif, street, number_phone, schedule, `name`, guard FROM med LIMIT" + page + ",10")
	if err != nil {
		panic(err.Error())
	}
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")

	for selDB.Next() {
		var id, numberPhone, guard int
		var name, street, scheduler, cif string
		err = selDB.Scan(&id, &name, &numberPhone, &guard, &street, &scheduler, &cif)
		if err != nil {
			panic(err.Error())
		}
		pharmacy := model.Pharmacy{
			ID:   id,
			Name: name,
			NumberPhone:  numberPhone,
			Guard:	guard,
			Street:	street,
			Schedule:	scheduler,
			Cif:	cif,
		}
		pharmacies = append(pharmacies, &pharmacy)
	}
	output, err := json.Marshal(pharmacies)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	w.Write(output)
}

func GetPharmacy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nID := vars["id"]

	selDB, err := db.Query("SELECT id, cif, street, number_phone, schedule, `name`, guard FROM pharmacy WHERE id=?", nID)
	if err != nil {
		panic(err.Error())
	}

	pharmacy := model.Pharmacy{}
	for selDB.Next() {
		var id, numberPhone, guard int
		var name, street, scheduler, cif string
		err = selDB.Scan(&id, &name, &numberPhone, &guard, &street, &scheduler, &cif)
		if err != nil {
			panic(err.Error())
		}
		pharmacy.ID = id
		pharmacy.Name = name
		pharmacy.NumberPhone = numberPhone
		pharmacy.Guard = guard
		pharmacy.Street = street
		pharmacy.Schedule = scheduler
		pharmacy.Cif = cif
	}

	output, err := json.Marshal(pharmacy)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	w.Write(output)
}

func CreatePharmacy(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Unmarshal
	var pharmacy model.Pharmacy
	err = json.Unmarshal(b, &pharmacy)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(pharmacy)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	w.Write(output)
	query := fmt.Sprintf("INSERT INTO `rds_pharmacy`.`pharmacy` (`name`, `cif`, `street`, `number_phone`, `schedule`, `guard`, `password`)  VALUES('%s', '%s', '%s', '%d', '%s', '%d', '%s')", pharmacy.Name, pharmacy.Cif, pharmacy.Street, pharmacy.NumberPhone, pharmacy.Schedule, pharmacy.Guard, pharmacy.Password)

	fmt.Println(query)
	insert, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

func UpdatePharmacy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nID := vars["id"]

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var pharmacy model.Pharmacy
	err = json.Unmarshal(b, &pharmacy)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(pharmacy)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(output)

	var query string = fmt.Sprintf("UPDATE `rds_pharmacy`.`pharmacy` SET `name` = '%s', `cif` = '%s, `street` = '%s', `schedule` = '%s', `password` = '%s', `phone_number` = '%d', `guard` = '%d' WHERE (`id` = '%s)", pharmacy.Name, pharmacy.Cif, pharmacy.Street, pharmacy.Schedule, pharmacy.Password, pharmacy.NumberPhone, pharmacy.Guard, nID)

	fmt.Println(query)
	update, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer update.Close()
}

func DeletePharmacy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nID := vars["id"]

	query := fmt.Sprintf("DELETE FROM `rds_pharmacy`.`pharmacy` WHERE (`id` = '%s')", nID)

	fmt.Println(query)
	insert, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	defer insert.Close()
}

func Login(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var user model.UserLogin
	err = json.Unmarshal(b, &user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	query := fmt.Sprintf("SELECT * FROM rds_pharmacy.pharmacy WHERE number_phone =  %d  and password = '%s'", user.Phone, user.Password)

	fmt.Println(query)
	selDB, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	pharmacy := model.Pharmacy{}
	for selDB.Next() {
		var id, phoneNumber, guard int
		var cif, street, schedule, name, password string
		err = selDB.Scan(&id, &cif, &street, &phoneNumber, &schedule, &name, &guard, &password)

		pharmacy.ID = id
		pharmacy.Cif = cif
		pharmacy.Street = street
		pharmacy.NumberPhone = phoneNumber
		pharmacy.Schedule = schedule
		pharmacy.Name = name
		pharmacy.Guard = guard
		pharmacy.Password = password

		if err != nil {
			panic(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

	}

	output, err := json.Marshal(pharmacy)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(output)
	defer selDB.Close()
}

func main() {
	conectDB()

	r := mux.NewRouter()
	r.HandleFunc("/", root).Methods("GET")

	r.HandleFunc("/meds", GetMeds).Methods("GET")
	r.HandleFunc("/meds/{id}", GetMed).Methods("GET")
	r.HandleFunc("/meds", CreateMed).Methods("POST")
	r.HandleFunc("/meds/{id}", UpdateMed).Methods("PUT")
	r.HandleFunc("/meds/{id}", DeleteMed).Methods("DELETE")

	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	r.HandleFunc("/pharmacies", GetPharmacies).Methods("GET")
	//r.HandleFunc("/pharmacies/{id}/users", GetUsersByPharmacyID).Methods("GET")
	r.HandleFunc("/pharmacies/{id}", GetPharmacy).Methods("GET")
	r.HandleFunc("/pharmacies", CreatePharmacy).Methods("POST")
	r.HandleFunc("/pharmacies/{id}", UpdatePharmacy).Methods("PUT")
	r.HandleFunc("/pharmacies/{id}", DeletePharmacy).Methods("DELETE")

	r.HandleFunc("/login", Login).Methods("POST")

	http.Handle("/", &MyServer{r})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}

type MyServer struct {
	r *mux.Router
}


func (s* MyServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	if req.Method == "OPTIONS" {
		return
	}

	s.r.ServeHTTP(w, req)
}
