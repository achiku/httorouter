package tiniyrouter

import (
	"log"
	"net/http"
)

const (
	asterisk = byte('*')
	colon    = byte(':')
	slash    = byte('/')
)

func find(src string, target byte, start int, max int) int {
	var i int
	for i = start; i < max && src[i] != target; i++ {
	}
	return i
}

type treeNode struct {
	path      string
	methods   map[string]http.Handler
	indices   string
	maxParams int
	children  []*treeNode
}

func newTreeNode(path string) *treeNode {
	node := &treeNode{
		path: path,
	}
	return node
}

func (tn *treeNode) addMethod(methods []string, handler http.Handler) {
	for _, m := range methods {
		if _, ok := tn.methods[m]; ok {
			panic("handlers conflict")
		}
		tn.methods[m] = handler
	}
}

func (tn *treeNode) getIndexPosition(target byte) int {
	low, high := 0, len(tn.indices)
	for low < high {
		mid := low + ((high - low) >> 1)
		log.Println(mid)
		if tn.indices[mid] < target {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return low
}

func (tn *treeNode) insertChild(index byte, child *treeNode) *treeNode {
	i := tn.getIndexPosition(index)
	tn.indices = tn.indices[:i] + string(index) + tn.indices[i:]
	tn.children = append(tn.children[:i], append([]*treeNode{child}, tn.children[i:]...)...)
	return child
}
