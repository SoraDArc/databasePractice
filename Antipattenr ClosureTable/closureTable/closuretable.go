package closuretable

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

var DbCon *sql.DB

func Start() {
	dbConnectionMPath()
}

func dbConnectionMPath() {
	connStr := "user=postgres password=DP123456 dbname=postgres sslmode=disable"
	var err error
	DbCon, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

type AnimalKingdom struct {
	Kingdom_id    int
	Level_kingdom string
	Title         string
	lvl           int
}

type ClosureTable struct {
	Ascenstor  int
	Descendant int
	lvl        int
}

// Получение прямых потомков
func GetDirectDescendants(id int) {
	query := "SELECT p1.ancestor, p1.descendant FROM closure_table as p1 LEFT JOIN (closure_table as p2 INNER JOIN closure_table as p3 ON p2.descendant = p3.ancestor) ON p2.ancestor = p1.ancestor AND p3.descendant = p1.descendant AND p2.ancestor <> p2.descendant AND p3.ancestor <> p3.descendant WHERE p1.ancestor = '" + strconv.Itoa(id) + "' and p2.ancestor is NULL;"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	CTable := []ClosureTable{}
	for rows.Next() {
		temp := ClosureTable{}
		rows.Scan(&temp.Ascenstor, &temp.Descendant)
		CTable = append(CTable, temp)
	}

	mas := []AnimalKingdom{}
	for _, val := range CTable {
		if val.Ascenstor == val.Descendant {
			continue
		}
		query := "SELECT kingdom_id, level_kingdom, title FROM anim_kingdom WHERE kingdom_id = '" + strconv.Itoa(val.Descendant) + "'"
		row, err := DbCon.Query(query)
		if err != nil {
			panic(err)
		}
		if row.Next() {
			temp := AnimalKingdom{}
			row.Scan(&temp.Kingdom_id, &temp.Level_kingdom, &temp.Title)
			temp.lvl = val.lvl
			mas = append(mas, temp)
		}
	}
	for _, v := range mas {
		fmt.Printf("%v %v %v\n", v.Title, v.Level_kingdom, v.Kingdom_id)
	}
}

// Получение прямого родителя
func GetDirectParent(id int) {
	query := "SELECT p1.ancestor, p1.descendant FROM closure_table as p1 LEFT JOIN (closure_table as p2 INNER JOIN closure_table as p3 ON p2.ancestor = p3.descendant) ON p2.descendant = p1.descendant AND p3.ancestor = p1.ancestor AND p2.ancestor <> p2.descendant AND p3.ancestor <> p3.descendant WHERE p1.descendant = '" + strconv.Itoa(id) + "' and p2.descendant is NULL;"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	CTable := []ClosureTable{}
	for rows.Next() {
		temp := ClosureTable{}
		rows.Scan(&temp.Ascenstor, &temp.Descendant)
		CTable = append(CTable, temp)
	}

	mas := []AnimalKingdom{}
	for _, val := range CTable {
		if val.Ascenstor == val.Descendant {
			continue
		}
		query := "SELECT kingdom_id, level_kingdom, title FROM anim_kingdom WHERE kingdom_id = '" + strconv.Itoa(val.Descendant) + "'"
		row, err := DbCon.Query(query)
		if err != nil {
			panic(err)
		}
		if row.Next() {
			temp := AnimalKingdom{}
			row.Scan(&temp.Kingdom_id, &temp.Level_kingdom, &temp.Title)
			temp.lvl = val.lvl
			mas = append(mas, temp)
		}
	}
	for _, v := range mas {
		fmt.Printf("%v %v %v\n", v.Title, v.Level_kingdom, v.Kingdom_id)
	}
}

// Получение всех потомков через closure_table
var CT []ClosureTable

func getAllDescendant(id, lvl int) {
	query := "SELECT p1.ancestor, p1.descendant FROM closure_table as p1 LEFT JOIN (closure_table as p2 INNER JOIN closure_table as p3 ON p2.descendant = p3.ancestor) ON p2.ancestor = p1.ancestor AND p3.descendant = p1.descendant AND p2.ancestor <> p2.descendant AND p3.ancestor <> p3.descendant WHERE p1.ancestor = '" + strconv.Itoa(id) + "' and p2.ancestor is NULL;"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		temp := ClosureTable{}
		rows.Scan(&temp.Ascenstor, &temp.Descendant)
		temp.lvl = lvl
		CT = append(CT, temp)
		if temp.Ascenstor != temp.Descendant {
			getAllDescendant(temp.Descendant, lvl+1)
		}
	}
}

// Вывод в иерархии потомков из closure_table
func ViewTree(id int) {
	CT = []ClosureTable{}
	getAllDescendant(id, 1)
	mas := []AnimalKingdom{}
	CT[0].lvl -= 1
	for ind, val := range CT {
		if val.Ascenstor == val.Descendant && ind != 0 {
			continue
		}
		query := "SELECT kingdom_id, level_kingdom, title FROM anim_kingdom WHERE kingdom_id = '" + strconv.Itoa(val.Descendant) + "'"
		row, err := DbCon.Query(query)
		if err != nil {
			panic(err)
		}
		if row.Next() {
			temp := AnimalKingdom{}
			row.Scan(&temp.Kingdom_id, &temp.Level_kingdom, &temp.Title)
			temp.lvl = val.lvl
			mas = append(mas, temp)
		}
	}
	for _, v := range mas {
		fmt.Printf("%v%v %v %v\n", strings.Repeat(" ", v.lvl*3), v.Title, v.Level_kingdom, v.Kingdom_id)
	}
}

// Получение всех родителей
var CTParents []ClosureTable

func getAllParents(id int) {
	query := "SELECT p1.ancestor, p1.descendant FROM closure_table as p1 LEFT JOIN (closure_table as p2 INNER JOIN closure_table as p3 ON p2.ancestor = p3.descendant) ON p2.descendant = p1.descendant AND p3.ancestor = p1.ancestor AND p2.ancestor <> p2.descendant AND p3.ancestor <> p3.descendant WHERE p1.descendant = '" + strconv.Itoa(id) + "' and p2.descendant is NULL;"
	rows, err := DbCon.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		temp := ClosureTable{}
		rows.Scan(&temp.Ascenstor, &temp.Descendant)
		if temp.Ascenstor != temp.Descendant {
			getAllParents(temp.Ascenstor)
		}
		CTParents = append(CTParents, temp)
	}
}

func ViewParentsTree(id int) {
	CTParents = []ClosureTable{}
	getAllParents(id)
	mas := []AnimalKingdom{}
	for idx, val := range CTParents {
		if val.Ascenstor == val.Descendant && idx != len(CTParents)-1 {
			continue
		}
		query := "SELECT kingdom_id, level_kingdom, title FROM anim_kingdom WHERE kingdom_id = '" + strconv.Itoa(val.Ascenstor) + "'"
		row, err := DbCon.Query(query)
		if err != nil {
			panic(err)
		}
		if row.Next() {
			temp := AnimalKingdom{}
			row.Scan(&temp.Kingdom_id, &temp.Level_kingdom, &temp.Title)
			temp.lvl = val.lvl
			mas = append(mas, temp)
		}
	}

	for i := 0; i < len(mas); i++ {
		fmt.Printf("%v%v %v %v\n", strings.Repeat(" ", (i)*3), mas[i].Title, mas[i].Level_kingdom, mas[i].Kingdom_id)
	}
}

// Удаление узла или листа
func DeleteLeafOrNode(id int) {
	query := "DELETE FROM closure_table WHERE descendant = '" + strconv.Itoa(id) + "'"
	res, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	aff, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if aff == 0 {
		fmt.Println("Ничего не удалилось")
	} else {
		fmt.Println("Узел или лист был удалён")
	}
}

// Удаление поддерева
func DeleteSubTree(id int) {
	query := "DELETE FROM closure_table WHERE descendant IN (SELECT descendant FROM closure_table WHERE ancestor = '" + strconv.Itoa(id) + "')"
	res, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	aff, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if aff == 0 {
		fmt.Println("Ничего не удалилось")
	} else {
		fmt.Println("Поддерево было удалёно")
	}
}

// Добавление листа
func AddLeaf(toWhatId int, level_kingdom, title string) {
	level_kingdom = strings.TrimSpace(level_kingdom)
	title = strings.TrimSpace(title)
	query := "INSERT INTO anim_kingdom (level_kingdom, title) VALUES ('" + level_kingdom + "', '" + title + "') RETURNING kingdom_id"
	var new_id int
	err := DbCon.QueryRow(query).Scan(&new_id)
	if err != nil {
		panic(err)
	}

	query = "INSERT INTO closure_table (ancestor, descendant) SELECT closure_table.ancestor, " + strconv.Itoa(new_id) + " FROM closure_table WHERE closure_table.descendant = '" + strconv.Itoa(toWhatId) + "' UNION ALL SELECT '" + strconv.Itoa(new_id) + "','" + strconv.Itoa(new_id) + "'"
	fmt.Println(query)
	_, err = DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Лист %v - %v был вставлен с id - %v\n", level_kingdom, title, new_id)
}

// Перемещение поддерева, узла, листа
func MoveSubTree(moveableId, toWhatId int) {
	query := "DELETE FROM closure_table WHERE descendant IN (SELECT descendant FROM closure_table WHERE ancestor = '" + strconv.Itoa(moveableId) + "') AND ancestor IN (SELECT ancestor FROM closure_table WHERE descendant = '" + strconv.Itoa(moveableId) + "' AND ancestor <> descendant)"
	_, err := DbCon.Exec(query)
	if err != nil {
		panic(err)
	}

	query = "INSERT INTO closure_table (ancestor, descendant) SELECT super.ancestor, sub.descendant FROM closure_table AS super CROSS JOIN closure_table AS sub WHERE super.descendant = '" + strconv.Itoa(toWhatId) + "' AND sub.ancestor = '" + strconv.Itoa(moveableId) + "'"
	_, err = DbCon.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Перенесено")
}

/*
— добавление листа; +
— удаление листа; +
— удаление поддерева; +
— перемещение листа; +
— перемещение поддерева; +
— удаление узла; +
— перемещение узла; +
— получение прямых потомков; +
— получение прямого родителя; +
— получение всех потомков (с сохранением иерархичности); +
— получение всех родителей (с сохранением иерархичности); +
*/
