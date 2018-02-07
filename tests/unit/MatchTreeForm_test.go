package unit

import (
	. "../../tests"
	"../../models"
	"../../services"
	"testing"
	"math"
)


func TestMatchTreeHappyPath(t *testing.T) {
	root := CreateBinaryTree(3)
	if root.Validate() != nil {
		t.Error("Symmetric binary tree is not valid")
	}
}

func TestSameTimeInNodes(t *testing.T) {
	root := CreateBinaryTree(3)
	root.LeftChild.StartTime = root.StartTime
	if root.Validate() == nil {
		t.Error("MatchTreeForm.child.Time == MatchTreeForm.Time")
	}
}

func TestTimeInChildMoreThanInParent(t *testing.T) {
	root := CreateBinaryTree(3)
	root.LeftChild.StartTime = root.StartTime.Add(1)
	if root.Validate() == nil {
		t.Error("MatchTreeForm.child.Time > MatchTreeForm.Time")
	}
}

func TestNotBinaryTree(t *testing.T) {
	root := CreateBinaryTree(3)
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
	maxDeep := int(math.Log2(float64(serv.GetConfig().MaxMatchesInTourney + 1)))
	deepTree := CreateBinaryTree(maxDeep)
	if deepTree.Validate() != nil {
		t.Errorf("Tree with depth %d is not a valid", maxDeep)
	}
}

func TestTooDeepRecursion(t *testing.T) {
	maxDeep := int(math.Log2(float64(serv.GetConfig().MaxMatchesInTourney + 1)))
	deepTree := CreateBinaryTree(maxDeep + 1)
	if deepTree.Validate() == nil {
		t.Errorf(
			"Tree with depth %d is valid " +
				"(it's more than max allowable number of nodes = %d)",
			maxDeep, serv.GetConfig().MaxMatchesInTourney,
		)
	}
}

func TestRecursiveTree(t *testing.T) {
	root := CreateBinaryTree(3)
	root.LeftChild.LeftChild = root.LeftChild
	if root.Validate() == nil {
		t.Errorf("Recursive tree is valid")
	}
}

func TestSameChildren(t *testing.T) {
	root := CreateBinaryTree(3)
	root.LeftChild = root.LeftChild.LeftChild
	root.RightChild = root.LeftChild.LeftChild
	if root.Validate() == nil {
		t.Errorf("Tree has the same children")
	}
}

func TestAsymmetricTree(t *testing.T) {
	root := CreateBinaryTree(3)
	root.RightChild.RightChild = nil
	root.RightChild.LeftChild = nil
	if root.Validate() != nil {
		t.Error("Assymetric tree is not valid")
	}
}


func TestCountSimple(t *testing.T) {
	root := CreateBinaryTree(3)
	count := root.GetNodesCount()
	if count != 7 {
		t.Errorf("Incorrect counting nodes in simple symmetric tree with depth 3" +
			"(wrong answer = %d)", count)
	}
}

func TestCountSingleNode(t *testing.T) {
	root := CreateBinaryTree(1)
	count := root.GetNodesCount()
	if count != 1 {
		t.Errorf("Incorrect counting single node (wrong answer = %d)", count)
	}
}

func TestCountAsymmetricTree(t *testing.T) {
	root := CreateBinaryTree(3)
	root.RightChild.RightChild = nil
	root.RightChild.LeftChild = nil
	count := root.GetNodesCount()
	if count != 5 {
		t.Error("Incorrect counting assymetric tree")
	}
}


func TestConvertFromTreeToArray(t *testing.T) {
	tree := CreateBinaryTree(3)

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