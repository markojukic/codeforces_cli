package parser

import (
	"io"

	"golang.org/x/net/html"
)

// Node - wrapper for html.Node
type Node html.Node

// Nodes - wrapper for []*Node
type Nodes []*Node

func Parse(r io.Reader) (*Node, error) {
	n, err := html.Parse(r)
	return (*Node)(n), err
}

func (n *Node) GetAttr(key string) string {
	if n == nil {
		return ""
	}
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

type elementNodeTester func(n *Node) bool

func (n *Node) GetNode(tester elementNodeTester) *Node {
	if n == nil {
		return nil
	}
	for c := (*Node)(n.FirstChild); c != nil; c = (*Node)(c.NextSibling) {
		if c.Type == html.ElementNode && tester(c) {
			return c
		} else if res := c.GetNode(tester); res != nil {
			return res
		}
	}
	return nil
}

func (n *Node) GetNodes(tester elementNodeTester) Nodes {
	if n == nil {
		return nil
	}
	var result Nodes
	for c := (*Node)(n.FirstChild); c != nil; c = (*Node)(c.NextSibling) {
		if c.Type == html.ElementNode && tester(c) {
			result = append(result, c)
		}
		result = append(result, c.GetNodes(tester)...)
	}
	return result
}

func (nodes Nodes) GetNodes(tester elementNodeTester) Nodes {
	if len(nodes) == 0 {
		return nil
	}
	var result Nodes
	for _, n := range nodes {
		result = append(result, n.GetNodes(tester)...)
	}
	return result
}

func (node *Node) GetNodeByTag(tagName string) *Node {
	return node.GetNode(func(n *Node) bool {
		return n.Data == tagName
	})
}

func (node *Node) GetNodesByTag(tagName string) Nodes {
	return node.GetNodes(func(n *Node) bool {
		return n.Data == tagName
	})
}

func (node *Node) GetNodeByTagAttr(tagName, key, val string) *Node {
	return node.GetNode(func(n *Node) bool {
		return n.Data == tagName && n.GetAttr(key) == val
	})
}

func (node *Node) GetNodesByTagAttr(tagName, key, val string) Nodes {
	return node.GetNodes(func(n *Node) bool {
		return n.Data == tagName && n.GetAttr(key) == val
	})
}

// Possibly returns a node multiple times, e.g. if [parent, child] is given
func (nodes Nodes) GetNodesByTag(tagName string) Nodes {
	return nodes.GetNodes(func(n *Node) bool {
		return n.Data == tagName
	})
}

func (nodes Nodes) GetNodesByTagAttr(tagName, key, val string) Nodes {
	return nodes.GetNodes(func(n *Node) bool {
		return n.Data == tagName && n.GetAttr(key) == val
	})
}

func (n *Node) GetText() string {
	if n == nil {
		return ""
	}
	text := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
		} else if c.Type == html.ElementNode && c.Data == "br" {
			text += "\n"
		}
	}
	return text
}

func (nodes Nodes) GetText() []string {
	if len(nodes) == 0 {
		return nil
	}
	var result []string
	for _, n := range nodes {
		result = append(result, n.GetText())
	}
	return result
}
