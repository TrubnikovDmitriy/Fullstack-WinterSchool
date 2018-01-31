package validation

import (
	"../../models"
	"../../services"
	"testing"
	"time"
	"math"
)

// Функция для создания дерева заданной глубины, возвращает корневую ноду
func createMatchTree(deep int) models.MatchTreeForm {

	nodesCount := int(math.Pow(2, float64(deep))) - 1
	nodes := make([]models.MatchTreeForm, nodesCount)
	times := time.Now()

	for i := 1; i <= nodesCount; i++ {
		if i < nodesCount / 2 {
			nodes[i - 1].LeftChild = &nodes[i * 2 - 1]
			nodes[i - 1].RightChild = &nodes[i * 2]
		}
		nodes[i - 1].StartTime = times.Add(time.Duration(nodesCount - i))
	}
	return nodes[0]
}

func TestHappyPath(t *testing.T) {
	root := createMatchTree(3)
	if !root.Validate() {
		t.Error("HappyPath for validation MatchTreeForm is not happy :(")
	}
}

func TestSameTimeInNodes(t *testing.T) {
	root := createMatchTree(3)
	root.LeftChild.StartTime = root.StartTime
	if root.Validate() {
		t.Error("MatchTreeForm.child.Time == MatchTreeForm.Time")
	}
}

func TestTimeInChildMoreThanInParent(t *testing.T) {
	root := createMatchTree(3)
	root.LeftChild.StartTime = root.StartTime.Add(1)
	if root.Validate() {
		t.Error("MatchTreeForm.child.Time > MatchTreeForm.Time")
	}
}

func TestNotBinaryTree(t *testing.T) {
	root := createMatchTree(4)
	root.LeftChild.LeftChild = nil
	if root.Validate() {
		t.Error("One of node has not both children")
	}
}

func TestOneSingleMatch(t *testing.T) {
	root := models.MatchTreeForm{}
	if !root.Validate() {
		t.Error("Single match without children is not pass test")
	}
}

func TestTheMaximumDepthRecursion(t *testing.T) {
	maxDeep := int(math.Log2(float64(services.MaxMatchesInTournament + 1)))
	deepTree := createMatchTree(maxDeep)
	if !deepTree.Validate() {
		t.Errorf("Tree with depth %d is not a valid", maxDeep)
	}
}

func TestTooDeepRecursion(t *testing.T) {
	maxDeep := int(math.Log2(float64(services.MaxMatchesInTournament + 1)))
	deepTree := createMatchTree(maxDeep + 1)
	if deepTree.Validate() {
		t.Errorf(
			"Tree with depth %d is valid " +
				"(it's more than max allowable number of nodes = %d)",
			maxDeep, services.MaxMatchesInTournament,
		)
	}
}

func TestRecursiveTree(t *testing.T) {
	root := createMatchTree(3)
	root.LeftChild.LeftChild = root.LeftChild
	if root.Validate() {
		t.Errorf("Recursive tree is valid")
	}
}

func TestSameChildren(t *testing.T) {
	root := createMatchTree(3)
	root.LeftChild = root.LeftChild.LeftChild
	root.RightChild = root.LeftChild.LeftChild
	if root.Validate() {
		t.Errorf("Tree has the same children")
	}
}

