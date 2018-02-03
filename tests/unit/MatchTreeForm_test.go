package unit

import (
	"../../models"
	"../../services"
	"testing"
	"time"
	"math"
)

// Функция для создания дерева заданной глубины, возвращает корневую ноду
func createBinaryTree(deep int) models.MatchesTreeForm {

	nodesCount := int(math.Pow(2, float64(deep))) - 1
	nodes := make([]models.MatchesTreeForm, nodesCount)
	times := time.Now()

	for i := 1; i <= nodesCount; i++ {
		if i <= nodesCount / 2 {
			nodes[i - 1].LeftChild = &nodes[i * 2 - 1]
			nodes[i - 1].RightChild = &nodes[i * 2]
		}
		nodes[i - 1].StartTime = times.Add(time.Duration(nodesCount - i))
	}
	return nodes[0]
}

func TestHappyPath(t *testing.T) {
	root := createBinaryTree(3)
	if root.Validate() != nil {
		t.Error("Symmetrical binary tree is not valid")
	}
}

func TestSameTimeInNodes(t *testing.T) {
	root := createBinaryTree(3)
	root.LeftChild.StartTime = root.StartTime
	if root.Validate() == nil {
		t.Error("MatchTreeForm.child.Time == MatchTreeForm.Time")
	}
}

func TestTimeInChildMoreThanInParent(t *testing.T) {
	root := createBinaryTree(3)
	root.LeftChild.StartTime = root.StartTime.Add(1)
	if root.Validate() == nil {
		t.Error("MatchTreeForm.child.Time > MatchTreeForm.Time")
	}
}

func TestNotBinaryTree(t *testing.T) {
	root := createBinaryTree(3)
	root.LeftChild.LeftChild = nil
	if root.Validate() == nil {
		t.Error("One of node has not both children")
	}
}

func TestOneSingleMatch(t *testing.T) {
	root := models.MatchesTreeForm{}
	if root.Validate() != nil {
		t.Error("Single match without children is not pass test")
	}
}

func TestTheMaximumDepthRecursion(t *testing.T) {
	maxDeep := int(math.Log2(float64(serv.MaxMatchesInTournament + 1)))
	deepTree := createBinaryTree(maxDeep)
	if deepTree.Validate() != nil {
		t.Errorf("Tree with depth %d is not a valid", maxDeep)
	}
}

func TestTooDeepRecursion(t *testing.T) {
	maxDeep := int(math.Log2(float64(serv.MaxMatchesInTournament + 1)))
	deepTree := createBinaryTree(maxDeep + 1)
	if deepTree.Validate() == nil {
		t.Errorf(
			"Tree with depth %d is valid " +
				"(it's more than max allowable number of nodes = %d)",
			maxDeep, serv.MaxMatchesInTournament,
		)
	}
}

func TestRecursiveTree(t *testing.T) {
	root := createBinaryTree(3)
	root.LeftChild.LeftChild = root.LeftChild
	if root.Validate() == nil {
		t.Errorf("Recursive tree is valid")
	}
}

func TestSameChildren(t *testing.T) {
	root := createBinaryTree(3)
	root.LeftChild = root.LeftChild.LeftChild
	root.RightChild = root.LeftChild.LeftChild
	if root.Validate() == nil {
		t.Errorf("Tree has the same children")
	}
}

func TestAsymmetricalTree(t *testing.T) {
	root := createBinaryTree(3)
	root.RightChild.RightChild = nil
	root.RightChild.LeftChild = nil
	if root.Validate() != nil {
		t.Error("Assymetrical tree is not valid")
	}
}


func TestCountSimple(t *testing.T) {
	root := createBinaryTree(3)
	count := root.GetNodesCount()
	if count != 7 {
		t.Errorf("Incorrect counting nodes in simple symmetrical tree with depth 3" +
			"(wrong answer = %d)", count)
	}
}

func TestCountSingleNode(t *testing.T) {
	root := createBinaryTree(1)
	count := root.GetNodesCount()
	if count != 1 {
		t.Errorf("Incorrect counting single node (wrong answer = %d)", count)
	}
}

func TestCountAsymmetricalTree(t *testing.T) {
	root := createBinaryTree(3)
	root.RightChild.RightChild = nil
	root.RightChild.LeftChild = nil
	count := root.GetNodesCount()
	if count != 5 {
		t.Error("Incorrect counting assymetrical tree")
	}
}


func TestConvertFromTreeToArray(t *testing.T) {
	tree := createBinaryTree(3)

	array := tree.CreateArrayMatch()

	// Считаем, что массив создавался прямым обходом дерева
	if array[0].ID != *array[1].NextMatch {
		t.Error("Left son has lost link to parent")
	}
	if array[0].ID != *array[4].NextMatch {
		t.Error("Right son has lost link to parent")
	}
	if array[1].ID != *array[2].NextMatch {
		t.Error("Left-left grandson has lost link to parent")
	}
	if array[1].ID != *array[3].NextMatch {
		t.Error("Left-right granson has lost link to parent")
	}
	if array[4].ID != *array[5].NextMatch {
		t.Error("Right-left grandson has lost link to parent")
	}
	if array[4].ID != *array[6].NextMatch {
		t.Error("Right-right granson has lost link to parent")
	}
}