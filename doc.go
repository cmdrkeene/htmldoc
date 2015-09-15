package htmldoc

import (
	"strings"

	"golang.org/x/net/html"
)
import "bytes"

func New(s string) (*Document, error) {
	root, err := parseString(s)
	if err != nil {
		return nil, err
	}
	return newDocument(root), nil
}

func MustNew(s string) *Document {
	root, err := parseString(s)
	if err != nil {
		panic(err)
	}
	return newDocument(root)
}

func newDocument(root *html.Node) *Document {
	return &Document{root: root}
}

func parseString(s string) (*html.Node, error) {
	return html.Parse(bytes.NewBufferString(s))
}

type Document struct {
	root *html.Node
}

type Node struct {
	node *html.Node
}

func newNode(n *html.Node) *Node {
	return &Node{
		node: n,
	}
}

func (self *Node) Text() string {
	buf := bytes.NewBufferString("")
	find(self.node, func(n *html.Node) bool {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		return false
	})
	return strings.TrimSpace(buf.String())
}

type Filter func(*html.Node) bool

func (self *Document) First(filters ...Filter) (*Node, bool) {
	node := find(self.root, all(filters...))
	return newNode(node), (node != nil)
}

func Tag(name string) Filter {
	return func(node *html.Node) bool {
		return (node.Type == html.ElementNode &&
			node.Data == name)
	}
}

func Class(name string) Filter {
	return func(node *html.Node) bool {
		return strings.TrimSpace(attribute(node, "class")) == name
	}
}

func Attribute(key, value string) Filter {
	return func(node *html.Node) bool {
		return strings.TrimSpace(attribute(node, key)) == value
	}
}

// attribute returns value for key
func attribute(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func all(filters ...Filter) Filter {
	return func(n *html.Node) bool {
		for _, f := range filters {
			if !f(n) {
				return false
			}
		}
		return true
	}
}

func find(root *html.Node, filter Filter) *html.Node {
	var found *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if filter(n) {
			found = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(root)
	return found
}
