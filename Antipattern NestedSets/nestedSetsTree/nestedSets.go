package nestedsetstree

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type AnimalKingdom struct {
	Kingdom_id    int
	Level_kingdom string
	Title         string
	lft           int
	rgt           int
}

var DbCon *sql.DB

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
	query := "SELECT * FROM nested_sets WHERE kingdom_id='" + strconv.Itoa(i) + "'"
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
func AddLeaf(toWhatId int, level_kingdom, title string) {
	if !CorrectInt(toWhatId) {
		return
	}
	level_kingdom = strings.TrimSpace(level_kingdom)
	title = strings.TrimSpace(title)
	if !CorrectString(level_kingdom) || !CorrectString(title) {
		return
	}
	level_kingdom = strings.TrimSpace(level_kingdom)
	title = strings.TrimSpace(title)
	query := "SELECT lft FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(toWhatId) + "'"
	var lft int
	row, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	if row.Next() {
		err = row.Scan(&lft)
		if err != nil {
			panic(err)
		}
	}

	query = "UPDATE nested_sets SET lft = CASE WHEN lft>='" + strconv.Itoa(lft+1) + "' THEN lft + 2 ELSE lft END, rgt= rgt + 2 WHERE rgt>='" + strconv.Itoa(lft+1) + "'"
	_, err = DbCon.Exec(query)
	if err != nil {
		panic(err)
	}

	query = "INSERT INTO nested_sets (level_kingdom, title, lft, rgt) VALUES ('" + level_kingdom + "','" + title + "', '" + strconv.Itoa(lft+1) + "', '" + strconv.Itoa(lft+2) + "') RETURNING kingdom_id"
	var new_id int
	DbCon.QueryRow(query).Scan(&new_id)
	fmt.Printf("Лист %v - %v был вставлен с id - %v\n", level_kingdom, title, new_id)
}

// Удаление листа или узла
func DeleteLeafOrNode(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "DELETE FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(id) + "'"
	res, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	aff, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if aff == 0 {
		fmt.Println("Такой записи нет")
	} else {
		fmt.Printf("Запись с id - %v удалена\n", id)
	}
}

// Удаление поддерева
func DeleteSubTree(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "DELETE FROM nested_sets WHERE lft BETWEEN (SELECT lft FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(id) + "') AND (SELECT rgt FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(id) + "')"
	res, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	aff, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if aff == 0 {
		fmt.Println("Такого поддерева нет")
	} else {
		fmt.Printf("Поддерево с id - %v было удалено\n", id)
	}
}

// Перемещение листа или узла
func MoveLeafOrNode(moveableId, toWhatId int) {
	if !CorrectInt(moveableId) || !CorrectInt(toWhatId) {
		return
	}
	if moveableId == toWhatId {
		fmt.Println("Невозможен перенос в самого себя")
		return
	}
	query := "SELECT lft FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(toWhatId) + "'"
	var lft int
	row, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	if row.Next() {
		err = row.Scan(&lft)
		if err != nil {
			panic(err)
		}
	}

	query = "UPDATE nested_sets SET lft = CASE WHEN lft>='" + strconv.Itoa(lft+1) + "' THEN lft + 2 ELSE lft END, rgt= rgt + 2 WHERE rgt>='" + strconv.Itoa(lft+1) + "'"
	_, err = DbCon.Exec(query)
	if err != nil {
		panic(err)
	}

	query = "UPDATE nested_sets SET lft = '" + strconv.Itoa(lft+1) + "', rgt = '" + strconv.Itoa(lft+2) + "' WHERE kingdom_id = '" + strconv.Itoa(moveableId) + "'"
	_, err = DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Узел с id %v перенесён к узлу id - %v\n", moveableId, toWhatId)
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
	query := "SELECT lft, rgt FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(moveableId) + "'"
	var lft, rgt int
	row, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	if row.Next() {
		err = row.Scan(&lft, &rgt)
		if err != nil {
			panic(err)
		}
	}
	space := rgt - lft + 1

	query = "SELECT lft FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(toWhatId) + "'"
	var TWIlft int
	row, err = DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	if row.Next() {
		err = row.Scan(&TWIlft)
		if err != nil {
			panic(err)
		}
	}

	query = "UPDATE nested_sets SET lft = CASE WHEN lft>='" + strconv.Itoa(TWIlft+1) + "' THEN lft + '" + strconv.Itoa(space) + "' ELSE lft END, rgt = rgt + '" + strconv.Itoa(space) + "' WHERE rgt>='" + strconv.Itoa(TWIlft+1) + "'"
	_, err = DbCon.Exec(query)
	if err != nil {
		panic(err)
	}

	query = "SELECT lft, rgt FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(moveableId) + "'"
	row, err = DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	if row.Next() {
		err = row.Scan(&lft, &rgt)
		if err != nil {
			panic(err)
		}
	}

	var newLft int = TWIlft + 1 - lft
	query = "UPDATE nested_sets SET lft = lft + '" + strconv.Itoa(newLft) + "', rgt = rgt + '" + strconv.Itoa(newLft) + "' WHERE lft BETWEEN '" + strconv.Itoa(lft) + "' AND '" + strconv.Itoa(rgt) + "'"
	_, err = DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	//fmt.Println(query)
	fmt.Printf("Перемещение поддерева с id - %v к узлу с id - %v\n", moveableId, toWhatId)
}

// Получение прямых потомков
func GetDirectDescendants(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "WITH Dep AS (SELECT n1.kingdom_id, COUNT(n2.kingdom_id) AS depth FROM nested_sets AS n1 JOIN nested_sets AS n2 ON n1.lft BETWEEN n2.lft AND n2.rgt GROUP BY n1.kingdom_id) SELECT NS.kingdom_id, NS.level_kingdom, NS.title, NS.lft, NS.rgt FROM nested_sets AS NS JOIN Dep AS D ON NS.kingdom_id = D.kingdom_id WHERE lft > (SELECT lft FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(id) + "') AND rgt < (SELECT rgt FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(id) + "') AND d.depth = (SELECT depth + 1 FROM Dep WHERE kingdom_id = '" + strconv.Itoa(id) + "') ORDER BY lft"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res := []AnimalKingdom{}
	for rows.Next() {
		temp := AnimalKingdom{}
		err = rows.Scan(&temp.Kingdom_id, &temp.Level_kingdom, &temp.Title, &temp.lft, &temp.rgt)
		if err != nil {
			panic(err)
		}
		res = append(res, temp)
	}
	for _, val := range res {
		fmt.Printf("%v %v %v %v %v\n", val.Title, val.Kingdom_id, val.Level_kingdom, val.lft, val.rgt)
	}
}

// Получение прямого родителя
func GetDirectParent(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "SELECT kingdom_id, level_kingdom, title, lft, rgt FROM nested_sets WHERE lft < (SELECT lft FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(id) + "') AND rgt > (SELECT rgt FROM nested_sets WHERE kingdom_id = '" + strconv.Itoa(id) + "') ORDER BY lft DESC LIMIT 1;"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	ak := AnimalKingdom{}
	if rows.Next() {
		err = rows.Scan(&ak.Kingdom_id, &ak.Level_kingdom, &ak.Title, &ak.lft, &ak.rgt)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Нет прямого родителя")
	}
	fmt.Println(ak)
}

// Получение всех родителей
func GetAllParents(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "WITH Dep AS (SELECT n1.kingdom_id, COUNT(n2.kingdom_id) AS depth FROM nested_sets AS n1 JOIN nested_sets AS n2 ON n1.lft BETWEEN n2.lft AND n2.rgt GROUP BY n1.kingdom_id) SELECT lpad('  ', 5 * CAST(D.depth AS int)) || h2.title as title, h2.level_kingdom, h2.kingdom_id, h2.lft, h2.rgt FROM nested_sets AS h1 INNER JOIN nested_sets AS h2 ON h1.lft BETWEEN h2.lft AND h2.rgt JOIN Dep AS D ON h2.kingdom_id = D.kingdom_id WHERE h1.kingdom_id = '" + strconv.Itoa(id) + "' ORDER BY lft"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res := []AnimalKingdom{}
	for rows.Next() {
		temp := AnimalKingdom{}
		err = rows.Scan(&temp.Title, &temp.Level_kingdom, &temp.Kingdom_id, &temp.lft, &temp.rgt)
		if err != nil {
			panic(err)
		}
		res = append(res, temp)
	}
	for _, val := range res {
		fmt.Printf("%v %v %v %v %v\n", val.Title, val.Kingdom_id, val.Level_kingdom, val.lft, val.rgt)
	}
}

// Получение всех потомков
func GetAllDescendants(id int) {
	if !CorrectInt(id) {
		return
	}
	query := "WITH Dep AS (SELECT n1.kingdom_id, COUNT(n2.kingdom_id) AS depth FROM nested_sets AS n1 JOIN nested_sets AS n2 ON n1.lft BETWEEN n2.lft AND n2.rgt GROUP BY n1.kingdom_id), AllDescen as (SELECT h2.kingdom_id, h2.level_kingdom, h2.title, h2.lft, h2.rgt FROM nested_sets AS h1 INNER JOIN nested_sets AS h2 ON h2.lft BETWEEN h1.lft AND h1.rgt WHERE h1.kingdom_id = '" + strconv.Itoa(id) + "') SELECT lpad('  ', 2 * CAST(D.depth AS int)) || AD.title as title, AD.level_kingdom, AD.kingdom_id, AD.lft, AD.rgt FROM AllDescen AS AD JOIN Dep AS D ON AD.kingdom_id=D.kingdom_id ORDER BY AD.lft"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res := []AnimalKingdom{}
	for rows.Next() {
		temp := AnimalKingdom{}
		err = rows.Scan(&temp.Title, &temp.Level_kingdom, &temp.Kingdom_id, &temp.lft, &temp.rgt)
		if err != nil {
			panic(err)
		}
		res = append(res, temp)
	}
	for _, val := range res {
		fmt.Printf("%v %v %v %v %v\n", val.Title, val.Kingdom_id, val.Level_kingdom, val.lft, val.rgt)
	}
}

/*
— добавление листа; +
— удаление листа; +
— удаление поддерева; +
— перемещение листа; +
— перемещение поддерева;
— удаление узла; +
— перемещение узла; +
— получение прямых потомков; +
— получение прямого родителя; +
— получение всех потомков (с сохранением иерархичности);
— получение всех родителей (с сохранением иерархичности); +
*/
