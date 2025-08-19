package main

import (
	nestedsetstree "Lab4/nestedSetsTree"
	"bufio"
	"fmt"
	"os"
)

func main() {
	nestedsetstree.Start()
	defer nestedsetstree.DbCon.Close()
	var cmd int
	for {
		fmt.Printf("1 — добавление листа\n2 — удаление листа\n3 — удаление поддерева\n4 — перемещение листа\n5 — перемещение поддерева\n6 — удаление узла\n7 — перемещение узла\n8 — получение прямых потомков\n9 — получение прямого родителя\n10 — получение всех потомков (с сохранением иерархичности)\n11 — получение всех родителей (с сохранением иерархичности)\n")
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
			nestedsetstree.AddLeaf(id, typ, title)
		case 2:
			fmt.Println("Id удаляемого")
			var id int
			fmt.Scanln(&id)
			nestedsetstree.DeleteLeafOrNode(id)
		case 3:
			fmt.Println("Id удаляемого")
			var id int
			fmt.Scanln(&id)
			nestedsetstree.DeleteSubTree(id)
		case 4:
			fmt.Println("Id перемещаемого - к чему через пробел")
			var moveableId, toWhatId int
			fmt.Scanln(&moveableId, &toWhatId)
			nestedsetstree.MoveLeafOrNode(moveableId, toWhatId)
		case 5:
			fmt.Println("Id перемещаемого - к чему через пробел")
			var moveableId, toWhatId int
			fmt.Scanln(&moveableId, &toWhatId)
			nestedsetstree.MoveSubTree(moveableId, toWhatId)
		case 6:
			fmt.Println("Id удаляемого")
			var id int
			fmt.Scanln(&id)
			nestedsetstree.DeleteLeafOrNode(id)
		case 7:
			fmt.Println("Id перемещаемого - к чему через пробел")
			var moveableId, toWhatId int
			fmt.Scanln(&moveableId, &toWhatId)
			nestedsetstree.MoveLeafOrNode(moveableId, toWhatId)
		case 8:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			nestedsetstree.GetDirectDescendants(id)
		case 9:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			nestedsetstree.GetDirectParent(id)
		case 10:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			nestedsetstree.GetAllDescendants(id)
		case 11:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			nestedsetstree.GetAllParents(id)
		default:
			fmt.Println("Неверная команда")
		}
	}
}
