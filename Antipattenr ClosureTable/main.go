package main

import (
	closuretable "Lab5/closureTable"
	"bufio"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	closuretable.Start()
	defer closuretable.DbCon.Close()
	var cmd int
	for {
		fmt.Printf("1 — добавление листа\n2 — удаление листа\n3 — удаление поддерева\n4 — перемещение листа, поддерева, узла\n5 — удаление узла\n6 — получение прямых потомков\n7 — получение прямого родителя\n8 — получение всех потомков (с сохранением иерархичности)\n9 — получение всех родителей (с сохранением иерархичности)\n")
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
			closuretable.AddLeaf(id, typ, title)
		case 2:
			fmt.Println("Id удаляемого")
			var id int
			fmt.Scanln(&id)
			closuretable.DeleteLeafOrNode(id)
		case 3:
			fmt.Println("Id удаляемого")
			var id int
			fmt.Scanln(&id)
			closuretable.DeleteSubTree(id)
		case 4:
			fmt.Println("Id перемещаемого - к чему через пробел")
			var moveableId, toWhatId int
			fmt.Scanln(&moveableId, &toWhatId)
			closuretable.MoveSubTree(moveableId, toWhatId)
		case 5:
			fmt.Println("Id удаляемого")
			var id int
			fmt.Scanln(&id)
			closuretable.DeleteLeafOrNode(id)
		case 6:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			closuretable.GetDirectDescendants(id)
		case 7:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			closuretable.GetDirectParent(id)
		case 8:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			closuretable.ViewTree(id)
		case 9:
			fmt.Println("Id?")
			var id int
			fmt.Scanln(&id)
			closuretable.ViewParentsTree(id)
		default:
			fmt.Println("Неверная команда")
		}
	}
}
