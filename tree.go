package tiniyrouter

import (
	"fmt"
	"net/http"
)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func countParams(path string) uint8 {
	var n uint
	for i := 0; i < len(path); i++ {
		if path[i] != ':' && path[i] != '*' {
			continue
		}
		n++
	}
	if n >= 255 {
		return 255
	}
	return uint8(n)
}

type nodeType uint8

const (
	static nodeType = iota // default
	root
	param
	catchAll
)

type node struct {
	path      string
	wildChild bool
	nType     nodeType
	maxParams uint8
	indices   string
	children  []*node
	handle    http.Handler
	priority  uint32
}

func (n *node) incrementChildPrio(pos int) int {
	n.children[pos].priority++
	prio := n.children[pos].priority

	newPos := pos
	for newPos > 0 && n.children[newPos-1].priority < prio {
		n.children[newPos-1], n.children[newPos] = n.children[newPos], n.children[newPos-1]
		newPos--
	}

	// build new index char string
	if newPos != pos {
		n.indices = n.indices[:newPos] +
			n.indices[pos:pos+1] +
			n.indices[newPos:pos] + n.indices[pos+1:]
	}
	return newPos
}

func (n *node) insertChild(numParams uint8, path, fullPath string, handle http.Handler) {
	var offset int

	// find prefix until first wildcard (begining with ':' or '*')
	for i, max := 0, len(path); numParams > 0; i++ {
		c := path[i]
		if c != ':' && c != '*' {
			continue
		}

		// find wildcard end (either '/' or path end)
		end := i + 1
		for end < max && path[end] != '/' {
			switch path[end] {
			case ':', '*':
				panic(fmt.Sprintf(
					"only one wildcard per path segment is allowed, has: %s in path %s", path[i:], fullPath))
			default:
				end++
			}
		}

		// if this Node has children, existing children would be unreachable
		if len(n.children) > 0 {
			panic(fmt.Sprintf(
				"wildcard route %s conflicts with existing children in path %s", path[i:end], fullPath))
		}

		// check if the wildcard has a name
		if end-i < 2 {
			panic(fmt.Sprintf("wildcards must be named with a non-empty name in path %s", fullPath))
		}

		if c == ':' { // param
			// split path at the beginning of the wildcard
			if i > 0 {
				n.path = path[offset:i]
				offset = i
			}

			child := &node{
				nType:     param,
				maxParams: numParams,
			}
			n.children = []*node{child}
			n.wildChild = true
			n = child
			n.priority++
			numParams--

			// if the path doesn't end with the wildcard, then there
			// will be another non-wildcard subpath starting with '/'
			if end < max {
				n.path = path[offset:end]
				offset = end

				child := &node{
					maxParams: numParams,
					priority:  1,
				}
				n.children = []*node{child}
				n = child
			}
		} else { // catchAll
			// wildcard is placed at the end of path
			if end != max || numParams > 1 {
				panic(fmt.Sprintf(
					"catch-all routers are only allowd at the end of the path: %s", fullPath))
			}
			if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
				panic(fmt.Sprintf(
					"catch-all conflicts with existing handler for the path segmen root: %s", fullPath))
			}
			i--
			if path[i] != '/' {
				panic(fmt.Sprintf(
					"no / before catch-all in path: %s", fullPath))
			}

			n.path = path[offset:i]

			// first node: catchAll node with empty path
			child := &node{
				wildChild: true,
				nType:     catchAll,
				maxParams: 1,
			}
			n.children = []*node{child}
			n.indices = string(path[i])
			n = child
			n.priority++

			// second node: node holding the variable
			child = &node{
				path:      path[i:],
				nType:     catchAll,
				maxParams: 1,
				handle:    handle,
				priority:  1,
			}
			n.children = []*node{child}
			return
		}
	}

	// insert remaining path part and handle to the leaf
	n.path = path[offset:]
	n.handle = handle
}

// addRoute adds a node with the given handle to the path.
// func (n *node) addRoute(path string, handle http.Handler) {
// 	fullPath := path
// 	n.priority++
// 	numParams := countParams(path)
//
// 	if len(n.path) > 0 || len(n.children) > 0 {
// 		for {
// 			if numParams > n.maxParams {
// 				n.maxParams = numParams
// 			}
//
// 			// Find the longest common prefix.
// 			// This alos implies that the common prefix contains no ':' or '*'
// 			// since the existing key can't contain those chars
// 			i := 0
// 			max := min(len(path), len(n.path))
// 			for i < max && path[i] == n.path[i] {
// 				i++
// 			}
//
// 			// Split edge
// 			if i < len(n.path) {
// 				child := node{
// 					path:      n.path[i:],
// 					wildChild: n.wildChild,
// 					nType:     static,
// 				}
// 			}
// 		}
// 	}
// }
