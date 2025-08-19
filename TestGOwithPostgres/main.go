package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func connectionToDb() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=12345678 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	handleError(err)
	return db
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isValidString(attrName string, attr *interface{}) string {
	var res string = ""
	v, ok := (*attr).(string)
	if !ok {
		res += "Поле " + attrName + " должно быть string!"
		return res
	}
	(*attr) = strings.TrimSpace(v)
	if (*attr) == "" {
		res += "Поле " + attrName + " не может быть пустым!"
		return res
	}
	return res
}

func isValidInt(attrName string, attr *interface{}) string {
	var res string = ""
	if v, ok := (*attr).(string); ok {
		v = strings.TrimSpace(v)
		if v == "" {
			return "Поле " + attrName + " не может быть пустым!"
		}
		v = strings.ReplaceAll(v, ",", ".")
		temp, err := strconv.ParseFloat(v, 64)
		if err != nil || temp <= 0 || temp != float64(int(temp)) {
			res += "Поле " + attrName + " должно быть целым и неотрицательным!"
		} else {
			(*attr) = int(temp)
		}
	} else {
		if i, ok := (*attr).(float64); ok && i > 0 {
			if i-float64(int(i)) != 0 {
				res += "Поле " + attrName + " должно быть целым и не отрицательным!"
				return res
			} else {
				if i > 0 {
					(*attr) = int(i)
				} else {
					res += "Поле " + attrName + " должно быть целым и не отрицательным!"
					return res
				}
			}
		} else {
			res += "Поле " + attrName + " должно быть целым и не отрицательным!"
			return res
		}
	}
	return res
}

func isValidFloat(attrName string, attr *interface{}) string {
	var res string = ""
	i, ok := (*attr).(float64)
	if !ok {
		res += "Поле " + attrName + " должнть быть не отрицательным и целым или вещественным числом!"
	}
	help := fmt.Sprintf("%.2f", i)
	temp, err := strconv.ParseFloat(help, 64)
	handleError(err)
	(*attr) = temp
	if temp < 0 {
		res += "Поле " + attrName + " не может быть отрицательным!"
	}
	return res
}

func getComputers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := "SELECT id, model, videocard, videocard_memory, cpu, cpu_hz, type_of_storage, storage_capacity, ram, price FROM computersdb"
	rows, err := db.Query(query)
	handleError(err)

	var computersData []Computer
	for rows.Next() {
		temp := Computer{}
		err := rows.Scan(&temp.Id, &temp.Model, &temp.Videocard, &temp.Videocard_memory, &temp.Cpu, &temp.Cpu_hz, &temp.Type_of_storage, &temp.Storage_capacity, &temp.Ram, &temp.Price)
		handleError(err)
		computersData = append(computersData, temp)
	}

	jsonData, err := json.Marshal(computersData)
	handleError(err)

	w.Write(jsonData)
}

func getComputer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := "SELECT  id, model, videocard, videocard_memory, cpu, cpu_hz, type_of_storage, storage_capacity, ram, price FROM computersdb WHERE id = $1"
	params := mux.Vars(r)
	if _, err := strconv.Atoi(params["id"]); err != nil {
		jsonData, err1 := json.Marshal("404 Not Found - Ресурс не найден")
		handleError(err1)
		w.Write(jsonData)
		return
	}
	rows, err := db.Query(query, params["id"])
	handleError(err)

	var computersData []Computer
	for rows.Next() {
		temp := Computer{}
		err := rows.Scan(&temp.Id, &temp.Model, &temp.Videocard, &temp.Videocard_memory, &temp.Cpu, &temp.Cpu_hz, &temp.Type_of_storage, &temp.Storage_capacity, &temp.Ram, &temp.Price)
		handleError(err)
		computersData = append(computersData, temp)
	}

	if computersData == nil {
		jsonData, err1 := json.Marshal("404 Not Found - Ресурс не найден")
		handleError(err1)
		w.Write(jsonData)
		return
	}

	jsonData, err := json.Marshal(computersData)
	handleError(err)

	w.Write(jsonData)
}

func createComputer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	temp := Computer{}
	json.NewDecoder(r.Body).Decode(&temp)
	var help []string = []string{isValidString("Model", &temp.Model), isValidString("Videocard", &temp.Videocard), isValidInt("Videocard_memory", &temp.Videocard_memory), isValidString("Cpu", &temp.Cpu), isValidString("Cpu_hz", &temp.Cpu_hz), isValidString("Type_of_storage", &temp.Type_of_storage), isValidInt("Storage_capacity", &temp.Storage_capacity), isValidInt("Ram", &temp.Ram), isValidFloat("Price", &temp.Price)}
	var resErr []string
	for _, str := range help {
		if str != "" {
			resErr = append(resErr, str)
		}
	}
	if len(resErr) != 0 {
		jsonData, err1 := json.Marshal(resErr)
		handleError(err1)
		w.Write(jsonData)
		return
	}
	query := "INSERT INTO computersdb(model, videocard, videocard_memory, cpu, cpu_hz, type_of_storage, storage_capacity, ram, price) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id"
	stmt, err := db.Prepare(query)
	handleError(err)

	var id string
	err = stmt.QueryRow(temp.Model, temp.Videocard, temp.Videocard_memory, temp.Cpu, temp.Cpu_hz, temp.Type_of_storage, temp.Storage_capacity, temp.Ram, temp.Price).Scan(&id)
	handleError(err)
	temp.Id, err = strconv.Atoi(id)
	handleError(err)

	jsonData, err := json.Marshal(&temp)
	handleError(err)

	w.Write(jsonData)
}

func deleteComputer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	if _, err := strconv.Atoi(params["id"]); err != nil {
		jsonData, err1 := json.Marshal("404 Not Found - Ресурс не найден")
		handleError(err1)
		w.Write(jsonData)
		return
	}
	query := "DELETE FROM computersdb WHERE id = $1"
	stmt, err := db.Prepare(query)
	handleError(err)
	res, err := stmt.Exec(params["id"])
	handleError(err)
	affected, err := res.RowsAffected()
	handleError(err)
	message := ""
	if affected == 0 {
		message = fmt.Sprintf("Удалено %v строк", affected)
	} else {
		message = fmt.Sprintf("Удалена %v строка", affected)
	}
	json.NewEncoder(w).Encode(message)
}

func updateComputer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	if _, err := strconv.Atoi(params["id"]); err != nil {
		jsonData, err1 := json.Marshal("404 Not Found - Ресурс не найден")
		handleError(err1)
		w.Write(jsonData)
		return
	}

	comp := Computer{}
	json.NewDecoder(r.Body).Decode(&comp)
	var help []string = []string{isValidString("Model", &comp.Model), isValidString("Videocard", &comp.Videocard), isValidInt("Videocard_memory", &comp.Videocard_memory), isValidString("Cpu", &comp.Cpu), isValidString("Cpu_hz", &comp.Cpu_hz), isValidString("Type_of_storage", &comp.Type_of_storage), isValidInt("Storage_capacity", &comp.Storage_capacity), isValidInt("Ram", &comp.Ram), isValidFloat("Price", &comp.Price)}
	var resErr []string
	for _, str := range help {
		if str != "" {
			resErr = append(resErr, str)
		}
	}
	if len(resErr) != 0 {
		jsonData, err1 := json.Marshal(resErr)
		handleError(err1)
		w.Write(jsonData)
		return
	}
	query := "UPDATE computersdb SET model = $1, videocard = $2, videocard_memory = $3, cpu = $4, cpu_hz = $5, type_of_storage = $6, storage_capacity = $7, ram = $8, price = $9 WHERE id = $10"
	stmt, err := db.Prepare(query)
	handleError(err)

	RES, err := stmt.Exec(comp.Model, comp.Videocard, comp.Videocard_memory, comp.Cpu, comp.Cpu_hz, comp.Type_of_storage, comp.Storage_capacity, comp.Ram, comp.Price, params["id"])
	handleError(err)
	affected, _ := RES.RowsAffected()
	if affected == 0 {
		jsonData, err1 := json.Marshal("404 Not Found - Ресурс не найден")
		handleError(err1)
		w.Write(jsonData)
		return
	}

	comp.Id, err = strconv.Atoi(params["id"])
	handleError(err)

	json.NewEncoder(w).Encode(comp)
}

func main() {
	db = connectionToDb()
	defer db.Close()
	r := mux.NewRouter()

	r.HandleFunc("/computersData", getComputers).Methods("GET")
	r.HandleFunc("/computersData/{id}", getComputer).Methods("GET")
	r.HandleFunc("/computersData/create", createComputer).Methods("POST")
	r.HandleFunc("/computersData/delete/{id}", deleteComputer).Methods("DELETE")
	r.HandleFunc("/computerData/update/{id}", updateComputer).Methods("PUT")

	http.ListenAndServe(":8000", r)
}

//бД- место, где хранятся данные и её структура.
//СУбд облегчает контроль и манипуляцию данными
//клиент-серверная архитектура это выделенный сервер, к-рый хранит в
//себе бизнес логику и хранится общие данные (клиент - сторонняя программа)
/*
pgadmin - клиент для постгрес(если есть выбор к кому коннектится тот это клиент)
протокол - свод правил(пример: пожатие рук)  - цифры - протокольные мерояприятия
протоколы в нашей области:
TCP/IP
UDP
HTTP-80
HTTPS - 443,81
FTP - 21
SFTP - 22
SSH - 22
TELNET - 23
представиться -> доверительное согласие ->
Разница Query и Exec в том что они возвращают. Query возвращает строки , а Exec результат выполнения запроса,
например сколько удалено или обновлено.
*/
