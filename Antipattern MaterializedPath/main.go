package main

import (
	"bufio"
	"fmt"
	materializedpath "lab3/materializedPath"
	"os"
)

func main() {
	materializedpath.Start()
	defer materializedpath.DbCon.Close()
	var cmd int
	for {
		fmt.Printf("1 — добавление листа\n2 — удаление листа\n3 — удаление поддерева\n4 — перемещение листа\n5 — перемещение поддерева\n6 — удаление узла без поддерева\n7 — перемещение узла без поддерева\n8 — получение прямых потомков\n9 — получение прямого родителя\n10 — получение всех потомков (с сохранением иерархичности)\n11 — получение всех родителей (с сохранением иерархичности)\n")
		fmt.Scanln(&cmd)
		switch cmd {
		case 1:
			fmt.Println("Id к чему добавляем, level_kingdom и title через переход на новую строку")
			var id int
			var typ, title string
			fmt.Scanln(&id)
			text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			typ = text
			text, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			title = text
			materializedpath.AddLeaf(id, typ, title)
		case 2:
			fmt.Println("Id удаляемого")
			var id int
			fmt.Scanln(&id)
			materializedpath.DeleteLeaf(id)
		case 3:
			fmt.Println("Id удаляемого")
			var id int
			fmt.Scanln(&id)
			materializedpath.DeleteSubTree(id)
		case 4:
			fmt.Println("Id перемещаемого - к чему через пробел")
			var moveableId, toWhatId int
			fmt.Scanln(&moveableId, &toWhatId)
			materializedpath.MoveLeaf(moveableId, toWhatId)
		case 5:
			fmt.Println("Id перемещаемого - к чему через пробел")
			var moveableId, toWhatId int
			fmt.Scanln(&moveableId, &toWhatId)
			materializedpath.MoveSubTree(moveableId, toWhatId)
		case 6:
			fmt.Println("Id удаляемого")
			var id int
			fmt.Scanln(&id)
			materializedpath.DeleteNodeWithoutSubTree(id)
		case 7:
			fmt.Println("Id перемещаемого - к чему через пробел")
			var moveableId, toWhatId int
			fmt.Scanln(&moveableId, &toWhatId)
			materializedpath.MoveNodeWithoutSubTree(moveableId, toWhatId)
		case 8:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			materializedpath.GetDirectDescendants(id)
		case 9:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			materializedpath.GetDirectParent(id)
		case 10:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			materializedpath.GetAllDescendants(id)
		case 11:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			materializedpath.GetAllParents(id)
		default:
			fmt.Println("Неверная команда")
		}
	}
}
