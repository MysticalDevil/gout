package gout

import (
	"iter"
	"strings"
)

// Tree node
type node struct {
	pattern  string  // route to be matched, e.g. /p/:lang
	part     string  // part of the route, e.g. :lang
	children []*node // children node, e.g. [doc, tutorial, intro]
	isWild   bool    // exact match or not, true when containing : and *
}

// matchChild finds the first child node that matches the given part.
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren finds all child nodes that match the given part.
func (n *node) matchChildren(part string) iter.Seq[*node] {
	return func(yield func(*node) bool) {
		for _, child := range n.children {
			if child.part == part || child.isWild {
				if !yield(child) {
					return
				}
			}
		}
	}
}

// insert adds a new route pattern to the trie.
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// search finds the node matching the given path parts.
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]

	for child := range n.matchChildren(part) {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
