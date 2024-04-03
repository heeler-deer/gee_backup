package gee

import (

	"strings"
)

type node struct {
	pattern  string
	part     string
	childern []*node
	isWild   bool
}

func (n *node) MatchChild(part string) *node {
	for _, child := range n.childern {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) MatchChildern(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.childern {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
			return nodes
		}
	}
	return nil
}

func (n *node) Insert(pattern string, parts []string, height int) {
	
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	
	part := parts[height]
	child := n.MatchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == '*' || part[0] == ':'}
		n.childern = append(n.childern, child)
	}
	child.Insert(pattern, parts, height+1)
}

func (n *node) Search(parts []string, height int) *node {
	
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	
	children := n.MatchChildern(part)
	for _, child := range children {
		result := child.Search(parts, height+1)
		if result != nil {
			return result
		}

	}
	return nil
}


