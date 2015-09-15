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

func (self *Document) All(filters ...Filter) []*Node {
	filter := all(filters...)
	var found []*Node
	search(self.root, func(n *html.Node) {
		if filter(n) {
			found = append(found, newNode(n))
		}
	})
	return found
}

func (self *Document) First(filters ...Filter) (*Node, bool) {
	var found *html.Node
	filter := all(filters...)
	search(self.root, func(n *html.Node) {
		if filter(n) {
			found = n
			return
		}
	})
	return newNode(found), (found != nil)
}

func (self *Document) Tag(name string) *FilterChain {
	return newFilterChain(self).Tag(name)
}

func (self *Document) Class(name string) *FilterChain {
	return newFilterChain(self).Class(name)
}

func (self *Document) Attribute(key, value string) *FilterChain {
	return newFilterChain(self).Attribute(key, value)
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
	search(self.node, func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
	})
	return strings.TrimSpace(buf.String())
}

type Filter func(*html.Node) bool

func newFilterChain(doc *Document) *FilterChain {
	return &FilterChain{doc: doc}
}

type FilterChain struct {
	chain []Filter
	doc   *Document
}

func (self *FilterChain) Tag(name string) *FilterChain {
	return self.add(Tag(name))
}

func (self *FilterChain) Class(name string) *FilterChain {
	return self.add(Class(name))
}

func (self *FilterChain) Attribute(key, value string) *FilterChain {
	return self.add(Attribute(key, value))
}

func (self *FilterChain) add(f Filter) *FilterChain {
	self.chain = append(self.chain, f)
	return self
}

func (self *FilterChain) First() (*Node, bool) {
	return self.doc.First(all(self.chain...))
}

func (self *FilterChain) All() []*Node {
	return self.doc.All(all(self.chain...))
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

// search applies visit using Depth First Search
func search(root *html.Node, visit func(n *html.Node)) {
	var dfs func(*html.Node)
	dfs = func(n *html.Node) {
		visit(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}
	dfs(root)
}
