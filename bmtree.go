package lightstore

import 
(
	"fmt"
	"sort"
)

//Implementation of B-tree

type BMtree struct{
	count int
	//Максимальный размер  массива в каждом ноде
	degree int
	root[]Node
}

func InitBMTree()(*BMtree){
	bmt := new(BMtree)
	return bmt
}

func (node*Node) GetKeyMaxDegree() int{
	return node.degree * 2 - 1
}

func (node*Node) GetKeeMinDegree(bm*Node) int {
	return bm.degree - 1;
}


func(node*Node) SplitNode()([]Node, []Node){
	//return two nodes
	root:= len(node.item)
	first := make([]Node, root)
	second:= make([]Node, root)
	firststep:= root/2
	for i := 0; i < firststep; i++ {
		dsp:= node.item[i];
		fmt.Println(dsp)
		//first[i] = dsp
	}

	for i := firststep; i < root - firststep; i++ {
		//second[i] = node.item[i];
	}

	return first, second;
}

type Node struct {
	//This is need to be a list
	item []int
	next*Node
	degree int
}

//Ключи должны быть в неубывающем порядке
func SortList(item []int)([]int) {
	sort.Ints(item)
	return item
}

/*
	Append new element
*/
func(node *Node) add(item interface{}){
	//Для добавления ключа нужно разбитие дерева
	if(node.GetKeyMaxDegree() > node.degree){
		//Node is full, split on 2 parts
	}
	result1, result2 := node.SplitNode()
	fmt.Println(result1, result2)
}

func(*Node) remove(item interface{}){
	//Удаление очередного элемента
}


//Probably need to return full list
func (node *Node) Find (item int) bool{
	for i := 0; i < len(node.item); i++ {
		if(node.item[i] == item){
			return true
		}
		if(node.item[i] > item){
			return node.next.Find(item)
		}

		if(node.item[i] < item){
			//Find data
		}
	}
	return false
}

