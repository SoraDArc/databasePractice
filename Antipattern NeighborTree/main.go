package main

import (
	"bufio"
	"fmt"
	"lab2/neighbortreeDB"
	"os"
)

// Тематика - классификация животных
func main() {
	neighbortreeDB.Start()
	defer neighbortreeDB.DbCon.Close()
	var cmd int
	for {
		fmt.Printf("1 — добавление листа\n2 — удаление листа\n3 — удаление поддерева\n4 — перемещение листа\n5 — перемещение поддерева\n6 — удаление узла без поддерева\n7 — перемещение узла без поддерева\n8 — получение прямых потомков\n9 — получение прямого родителя\n10 — получение всех потомков (с сохранением иерархичности)\n11 — получение всех родителей (с сохранением иерархичности)\n")
		fmt.Scanln(&cmd)
		switch cmd {
		case 1:
			fmt.Println("id родителя, тип, имя через переход на новую строку")
			var p_id int
			var typ, title string
			fmt.Scanln(&p_id)
			text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			typ = text
			text, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			title = text
			neighbortreeDB.AddLeaf(p_id, typ, title)
		case 2:
			fmt.Println("id")
			var id int
			fmt.Scanln(&id)
			neighbortreeDB.DeleteLeaf(id)
		case 3:
			fmt.Println("id")
			var id int
			fmt.Scanln(&id)
			neighbortreeDB.DeleteSubTree(id)
		case 4:
			fmt.Println("id Перемещаемого, id к чему через переход на новую строку")
			var moveable, toWhat int
			fmt.Scanln(&moveable)
			fmt.Scanln(&toWhat)
			neighbortreeDB.MoveLeaf(moveable, toWhat)
		case 5:
			fmt.Println("id Перемещаемого, id к чему через переход на новую строку")
			var moveable, toWhat int
			fmt.Scanln(&moveable)
			fmt.Scanln(&toWhat)
			neighbortreeDB.MoveSubTree(moveable, toWhat)
		case 6:
			fmt.Println("id")
			var id int
			fmt.Scanln(&id)
			neighbortreeDB.DeleteNodeWithoutSubTree(id)
		case 7:
			fmt.Println("Перемещаемое, к чему через переход на новую строку")
			var moveable, toWhat int
			fmt.Scanln(&moveable)
			fmt.Scanln(&toWhat)
			neighbortreeDB.MoveNodeWithoutSubTree(moveable, toWhat)
		case 8:
			fmt.Println("id")
			var id int
			fmt.Scanln(&id)
			neighbortreeDB.GetDirectDescendants(id)
		case 9:
			fmt.Println("id")
			var id int
			fmt.Scanln(&id)
			neighbortreeDB.GetDirectParent(id)
		case 10:
			fmt.Println("id")
			var id int
			fmt.Scanln(&id)
			neighbortreeDB.GetAllDescendants(id)
		case 11:
			fmt.Println("id")
			var id int
			fmt.Scanln(&id)
			neighbortreeDB.GetAllParents(id)
		default:
			fmt.Println("Неверная команда")
		}
	}
}
