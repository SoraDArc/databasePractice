package materializedpath

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

var DbCon *sql.DB

type AnimalKingdom struct {
	Kingdom_id    int
	Level_kingdom string
	Title         string
	Mpath         string
}

func Start() {
	dbConnectionMPath()
}

func dbConnectionMPath() {
	connStr := "user=postgres password=12345678 dbname=postgres sslmode=disable"
	var err error
	DbCon, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func CorrectInt(i int) bool {
	query := "SELECT * FROM mpath_ak WHERE kingdom_id='" + strconv.Itoa(i) + "'"
	row, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	if !row.Next() {
		fmt.Println("Нет строки с таким id")
		return false
	}
	if i <= 0 {
		fmt.Println("Неверный id")
		return false
	}
	return true
}

func CorrectString(s string) bool {
	if s == "" {
		fmt.Println("Пустая строка")
		return false
	}
	return true
}

// Добавление листа
func AddLeaf(id int, level_kingdom, title string) {
	if !CorrectInt(id) {
		return
	}
	level_kingdom = strings.TrimSpace(level_kingdom)
	title = strings.TrimSpace(title)
	if !CorrectString(level_kingdom) || !CorrectString(title) {
		return
	}
	level_kingdom = strings.TrimSpace(level_kingdom)
	title = strings.TrimSpace(title)
	query := "INSERT INTO mpath_ak (level_kingdom,title) VALUES ('" + level_kingdom + "','" + title + "') RETURNING kingdom_id;"
	var new_id int
	//Вставляем узел
	DbCon.QueryRow(query).Scan(&new_id)
	//Обновляем путь
	query = "UPDATE mpath_ak SET mpath = (SELECT mpath FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(id) + "') || '" + strconv.Itoa(new_id) + "' || '/' WHERE kingdom_id = '" + strconv.Itoa(new_id) + "';"
	_, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Лист %v - %v был вставлен с id - %v\n", level_kingdom, title, new_id)
}

// Удаление листа
func DeleteLeaf(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "DELETE FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(id) + "'"
	res, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		fmt.Println("Такой записи нет")
	} else {
		fmt.Printf("Лист с id - %v был удалён\n", id)
	}
}

// Удаление поддерева
func DeleteSubTree(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "DELETE FROM mpath_ak WHERE kingdom_id IN (SELECT kingdom_id FROM mpath_ak WHERE mpath LIKE (SELECT mpath || '%' FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(id) + "'))"
	res, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		fmt.Println("Такого поддерева нет")
	} else {
		fmt.Printf("Поддерево с id - %v был удалёно\n", id)
	}
}

// Перемещение листа
func MoveLeaf(moveableId, toWhatId int) {
	if !CorrectInt(moveableId) || !CorrectInt(toWhatId) {
		return
	}
	if moveableId == toWhatId {
		fmt.Println("Невозможен перенос в самого себя")
		return
	}
	query := "UPDATE mpath_ak SET mpath = (SELECT mpath FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(toWhatId) + "') || kingdom_id || '/' WHERE kingdom_id = '" + strconv.Itoa(moveableId) + "'"
	res, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		fmt.Println("Такой записи нет")
	} else {
		fmt.Printf("Перемещение листа с id - %v к узлу с id - %v\n", moveableId, toWhatId)
	}
}

// Перемещение поддерева
func MoveSubTree(moveableId, toWhatId int) {
	if !CorrectInt(moveableId) || !CorrectInt(toWhatId) {
		return
	}
	if moveableId == toWhatId {
		fmt.Println("Невозможен перенос в самого себя")
		return
	}
	query := "UPDATE mpath_ak SET mpath = REPLACE(mpath, (SELECT mpath FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(moveableId) + "'), (SELECT mpath FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(toWhatId) + "') || '" + strconv.Itoa(moveableId) + "/')"
	res, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		fmt.Println("Такого поддерева нет")
	} else {
		fmt.Printf("Перемещение поддерева с id - %v к узлу с id - %v\n", moveableId, toWhatId)
	}
}

// Удаление узла без поддерева
func DeleteNodeWithoutSubTree(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "UPDATE mpath_ak SET mpath = REPLACE(mpath, '/" + strconv.Itoa(id) + "/','/')"
	_, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	query = "DELETE FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(id) + "'"
	_, err = DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Узел с id - %v удалён\n", id)
}

// Перемещение узла без поддерева
func MoveNodeWithoutSubTree(moveableId, toWhatId int) {
	if !CorrectInt(moveableId) || !CorrectInt(toWhatId) {
		return
	}
	if moveableId == toWhatId {
		fmt.Println("Невозможен перенос в самого себя")
		return
	}
	query := "UPDATE mpath_ak SET mpath =REPLACE(mpath, '/" + strconv.Itoa(moveableId) + "/','/')"
	_, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	query = "UPDATE mpath_ak SET mpath = (SELECT mpath FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(toWhatId) + "') || kingdom_id || '/' WHERE kingdom_id = '" + strconv.Itoa(moveableId) + "'"
	_, err = DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Перемещение узла с id - %v к узлу id - %v\n", moveableId, toWhatId)
}

// Получения прямых потомков
func GetDirectDescendants(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "SELECT kingdom_id, level_kingdom, title, mpath FROM mpath_ak AS mp WHERE mpath LIKE (SELECT mpath FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(id) + "') || '%' AND (LENGTH(mpath) - LENGTH(REPLACE(mpath, '/', ''))) = ((SELECT LENGTH(mpath) - LENGTH(REPLACE(mpath, '/', '')) FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(id) + "') + 1)"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res := []AnimalKingdom{}
	for rows.Next() {
		temp := AnimalKingdom{}
		err = rows.Scan(&temp.Kingdom_id, &temp.Level_kingdom, &temp.Title, &temp.Mpath)
		if err != nil {
			panic(err)
		}
		res = append(res, temp)
	}
	for _, val := range res {
		fmt.Printf("%v %v %v %v\n", val.Title, val.Kingdom_id, val.Level_kingdom, val.Mpath)
	}
}

// Получение прямого родителя
func GetDirectParent(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "SELECT kingdom_id, level_kingdom, title, mpath FROM mpath_ak WHERE mpath LIKE (SELECT SUBSTRING(mpath, 1, length(SUBSTRING(mpath,1,length(mpath)-1)) - strpos(reverse(SUBSTRING(mpath,1,length(mpath)-1)), '/') + 1) FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(id) + "')"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	ak := AnimalKingdom{}
	if rows.Next() {
		err = rows.Scan(&ak.Kingdom_id, &ak.Level_kingdom, &ak.Title, &ak.Mpath)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Нет прямого родителя")
	}
	fmt.Println(ak)
}

// Получение всех потомков
func GetAllDescendants(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "SELECT lpad(' ',5 * (LENGTH(mpath) - LENGTH(REPLACE(mpath, '/', '')))) || title AS formatted_title, level_kingdom,kingdom_id,mpath FROM mpath_ak WHERE mpath LIKE (SELECT mpath || '%' FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(id) + "') ORDER BY mpath"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res := []AnimalKingdom{}
	for rows.Next() {
		temp := AnimalKingdom{}
		err = rows.Scan(&temp.Title, &temp.Level_kingdom, &temp.Kingdom_id, &temp.Mpath)
		if err != nil {
			panic(err)
		}
		res = append(res, temp)
	}
	for _, val := range res {
		fmt.Printf("%v %v %v %v\n", val.Title, val.Kingdom_id, val.Level_kingdom, val.Mpath)
	}
}

// Получение всех родителей
func GetAllParents(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "SELECT lpad(' ',5 * (LENGTH(mpath) - LENGTH(REPLACE(mpath, '/', '')))) || title AS formatted_title, level_kingdom,kingdom_id,mpath FROM mpath_ak WHERE (SELECT mpath FROM mpath_ak WHERE kingdom_id = '" + strconv.Itoa(id) + "') LIKE mpath || '%'"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res := []AnimalKingdom{}
	for rows.Next() {
		temp := AnimalKingdom{}
		err = rows.Scan(&temp.Title, &temp.Level_kingdom, &temp.Kingdom_id, &temp.Mpath)
		if err != nil {
			panic(err)
		}
		res = append(res, temp)
	}
	for _, val := range res {
		fmt.Printf("%v %v %v %v\n", val.Title, val.Kingdom_id, val.Level_kingdom, val.Mpath)
	}
}

/*
— добавление листа; +
— удаление листа; +
— удаление поддерева; +
— перемещение листа; +
— перемещение поддерева; +
— удаление узла без поддерева; +
— перемещение узла без поддерева; +
— получение прямых потомков; +
— получение прямого родителя; +
— получение всех потомков (с сохранением иерархичности); +
— получение всех родителей (с сохранением иерархичности); +
*/
