package neighbortreeDB

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type AnimalKingdom struct {
	Kingdom_id    int
	Parent_id     int
	Level_kingdom string
	Title         string
}

var DbCon *sql.DB

func Start() {
	dbConnectionNTRee()
}

func dbConnectionNTRee() {
	connStr := "user=postgres password=12345678 dbname=postgres sslmode=disable"
	var err error
	DbCon, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func ValidStrings(a ...*string) bool {
	for _, val := range a {
		(*val) = strings.TrimSpace(*val)
		if (*val) == "" {
			return false
		}
	}
	return true
}

// + 9
func GetDirectParent(id int) {
	query := "SELECT kingdom_id,parent_id,level_kingdom,title FROM animal_kingdom WHERE kingdom_id = (SELECT parent_id FROM animal_kingdom WHERE kingdom_id = ('" + strconv.Itoa(id) + "'))"
	row, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer row.Close()
	ak := AnimalKingdom{}
	if row.Next() {
		var help interface{}
		err = row.Scan(&ak.Kingdom_id, &help, &ak.Level_kingdom, &ak.Title)
		if err != nil {
			panic(err)
		}
		if i, ok := help.(int64); ok {
			ak.Parent_id = int(i)
		}
	} else {
		fmt.Println("Нет прямого родителя")
		return
	}
	fmt.Printf("Прямой родитель %v\n", ak)
}

// + 11
func GetAllParents(id int) {
	query := "WITH RECURSIVE parents AS (SELECT kingdom_id, parent_id, level_kingdom, title, 1 AS level FROM animal_kingdom WHERE kingdom_id = '" + strconv.Itoa(id) + "' UNION SELECT ak.kingdom_id, ak.parent_id, ak.level_kingdom, ak.title, level+1 AS level FROM animal_kingdom ak JOIN parents a ON ak.kingdom_id = a.parent_id) SELECT kingdom_id, parent_id, level_kingdom, lpad(' ', 5* level) || title FROM parents;"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res := []AnimalKingdom{}
	for rows.Next() {
		ak := AnimalKingdom{}
		var help interface{}
		err = rows.Scan(&ak.Kingdom_id, &help, &ak.Level_kingdom, &ak.Title)
		if err != nil {
			panic(err)
		}
		if i, ok := help.(int64); ok {
			ak.Parent_id = int(i)
		}
		res = append(res, ak)
	}
	if len(res) == 0 {
		fmt.Println("Такой записи нет")
		return
	}
	for _, val := range res {
		fmt.Printf("%v %v %v \"%v\"\n", val.Title, val.Kingdom_id, val.Parent_id, val.Level_kingdom)
	}
}

// + 8
func GetDirectDescendants(id int) {
	query := "SELECT kingdom_id,parent_id,level_kingdom,title FROM animal_kingdom WHERE parent_id = '" + strconv.Itoa(id) + "'"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res := []AnimalKingdom{}
	for rows.Next() {
		ak := AnimalKingdom{}
		var help interface{}
		err = rows.Scan(&ak.Kingdom_id, &help, &ak.Level_kingdom, &ak.Title)
		if err != nil {
			panic(err)
		}
		if i, ok := help.(int64); ok {
			ak.Parent_id = int(i)
		}
		res = append(res, ak)
	}
	if len(res) == 0 {
		fmt.Println("Нет прямого потомка")
		return
	}
	for _, val := range res {
		fmt.Printf("%v %v %v \"%v\"\n", val.Title, val.Kingdom_id, val.Parent_id, val.Level_kingdom)
	}
}

// + 10
func GetAllDescendants(id int) {
	query := "WITH RECURSIVE t AS (SELECT ARRAY[kingdom_id] AS hierarchy, kingdom_id, parent_id,level_kingdom, title, 1 AS level FROM animal_kingdom WHERE  kingdom_id = '" + strconv.Itoa(id) + "' UNION ALL SELECT t.hierarchy || a.kingdom_id, a.kingdom_id, a.parent_id,a.level_kingdom, a.title, t.level+1 AS level FROM   animal_kingdom a JOIN t ON a.parent_id = t.kingdom_id) SELECT kingdom_id, parent_id, level_kingdom, (lpad(' ', 5 * level) || title) AS title FROM t ORDER BY hierarchy;"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	res := []AnimalKingdom{}
	for rows.Next() {
		ak := AnimalKingdom{}
		var help interface{}
		err = rows.Scan(&ak.Kingdom_id, &help, &ak.Level_kingdom, &ak.Title)
		if err != nil {
			panic(err)
		}
		if i, ok := help.(int64); ok {
			ak.Parent_id = int(i)
		}
		res = append(res, ak)
	}
	if len(res) == 0 {
		fmt.Println("Такой записи нет")
		return
	}
	for _, val := range res {
		fmt.Printf("%v %v %v \"%v\"\n", val.Title, val.Kingdom_id, val.Parent_id, val.Level_kingdom)
	}
}

// Добавление листа + 1
func AddLeaf(parent_id int, level_kingdom, cur_title string) {
	if !ValidStrings(&level_kingdom, &cur_title) {
		fmt.Println("Поля level_kingdom, cur_title не должны быть пустыми")
		return
	}
	query := "INSERT INTO animal_kingdom (parent_id, level_kingdom, title) VALUES ('" + strconv.Itoa(parent_id) + "','" + level_kingdom + "','" + cur_title + "')"
	_, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Лист %v - %v добавлен\n", level_kingdom, cur_title)
}

// Удаление листа + 2
func DeleteLeaf(id int) {
	query := "DELETE FROM animal_kingdom WHERE kingdom_id = '" + strconv.Itoa(id) + "'"
	res, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		fmt.Println("Такой записи нет")
	} else {
		fmt.Printf("Лист с id - %v удалён\n", id)
	}
}

// Удаление поддерева + 3
func DeleteSubTree(id int) {
	query := "DELETE FROM animal_kingdom WHERE kingdom_id = '" + strconv.Itoa(id) + "'"
	res, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		fmt.Println("Такого поддерева нет")
	} else {
		fmt.Printf("Поддерево с id - %v удалёно\n", id)
	}
}

// Удаления узла без поддерева + 6
func DeleteNodeWithoutSubTree(id int) {
	//Запрос предка
	selectQuery := "SELECT parent_id FROM animal_kingdom WHERE kingdom_id ='" + strconv.Itoa(id) + "'"
	row, err := DbCon.Query(selectQuery)
	if err != nil {
		panic(err)
	}
	defer row.Close()
	var parent_id int
	if row.Next() {
		err = row.Scan(&parent_id)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Такой записи нет")
		return
	}
	//Переподчинение поддерева
	updateQuery := "UPDATE animal_kingdom SET parent_id = '" + strconv.Itoa(parent_id) + "' WHERE parent_id = '" + strconv.Itoa(id) + "'"
	_, err = DbCon.Exec(updateQuery)
	if err != nil {
		panic(err)
	}
	//Удаление узла
	deleteQuery := "DELETE FROM animal_kingdom WHERE kingdom_id = '" + strconv.Itoa(id) + "'"
	_, err = DbCon.Exec(deleteQuery)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Узел с id - %v удалён\n", id)
}

// Перемещение листа + 4
func MoveLeaf(id_Movable, id_toWhat int) {
	query := "UPDATE animal_kingdom SET parent_id = '" + strconv.Itoa(id_toWhat) + "' WHERE kingdom_id = '" + strconv.Itoa(id_Movable) + "'"
	result, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		fmt.Println("Такой записи нет")
		return
	}
	fmt.Printf("Перемещение листа с id - %v к узлу с id - %v\n", id_Movable, id_toWhat)
}

// Перемещение узла без поддерева + 7
func MoveNodeWithoutSubTree(id_moveable, id_toWhat int) {
	//Запросим предка
	selectQuery := "SELECT parent_id FROM animal_kingdom WHERE kingdom_id = '" + strconv.Itoa(id_moveable) + "'"
	row, err := DbCon.Query(selectQuery)
	if err != nil {
		panic(err)
	}
	defer row.Close()
	var parent_id int
	if row.Next() {
		err = row.Scan(&parent_id)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Такой записи нет")
		return
	}
	//Переподчиняем поддерево
	updateQuery := "UPDATE animal_kingdom SET parent_id = '" + strconv.Itoa(parent_id) + "' WHERE parent_id = '" + strconv.Itoa(id_moveable) + "'"
	_, err = DbCon.Exec(updateQuery)
	if err != nil {
		panic(err)
	}
	//Перемещаем узел
	update2Query := "UPDATE animal_kingdom SET parent_id = '" + strconv.Itoa(id_toWhat) + "' WHERE kingdom_id = '" + strconv.Itoa(id_moveable) + "'"
	_, err = DbCon.Exec(update2Query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Узел с id - %v перемещён к id - %v\n", id_moveable, id_toWhat)
}

// Перемешение поддерева + 5
func MoveSubTree(id_moveableSubTree, id_toWhat int) {
	query := "UPDATE animal_kingdom SET parent_id = '" + strconv.Itoa(id_toWhat) + "' WHERE kingdom_id = '" + strconv.Itoa(id_moveableSubTree) + "'"
	result, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		fmt.Println("Такой записи нет")
		return
	}
	fmt.Printf("Перемещение поддерева с id - %v к узлу  с id - %v\n", id_moveableSubTree, id_toWhat)
}

/*
CREATE TABLE animal_kingdom (
	kingdom_id SERIAL PRIMARY KEY,
	parent_id BIGINT,
	level_kingdom TEXT NOT NULL,
	title TEXT NOT NULL,
	FOREIGN KEY(parent_id) REFERENCES animal_kingdom(kingdom_id) ON DELETE CASCADE
);
*/

/*
— добавление листа;
— удаление листа;
— удаление поддерева;
— перемещение листа;
— перемещение поддерева;
— удаление узла без поддерева;
— перемещение узла без поддерева;
— получение прямых потомков;
— получение прямого родителя;
— получение всех потомков (с сохранением иерархичности);
— получение всех родителей (с сохранением иерархичности);
*/
