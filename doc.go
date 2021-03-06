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
	root, err := New(s)
	if err != nil {
		panic(err)
	}
	return root
}

func parseString(s string) (*html.Node, error) {
	return html.Parse(bytes.NewBufferString(s))
}

type Node struct {
	node *html.Node
}

func (self *Node) Attribute(key string) string {
	return attribute(self.node, key)
}

// Parent returns a Document that searches upward through parents
func (self *Node) Parent() *Document {
	return &Document{
		root:       self.node,
		searchFunc: searchParent,
	}
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

func newDocument(root *html.Node) *Document {
	return &Document{
		root:       root,
		searchFunc: search,
	}
}

type Document struct {
	root       *html.Node
	chain      []Filter
	match      Filter
	searchFunc func(*html.Node, func(*html.Node))
}

func (self *Document) Tag(name string) *Document {
	return self.add(Tag(name))
}

func (self *Document) Class(name string) *Document {
	return self.add(Class(name))
}

func (self *Document) Attribute(key, value string) *Document {
	return self.add(Attribute(key, value))
}

func (self *Document) add(f Filter) *Document {
	if f == nil {
		return self
	}

	self.chain = append(self.chain, f)
	self.match = all(self.chain...)
	return self
}

func (self *Document) First(selectors ...string) (*Node, bool) {
	self.addSelectors(selectors)
	var found *html.Node
	self.searchFunc(self.root, func(n *html.Node) {
		if self.match(n) {
			found = n
			return
		}
	})
	return newNode(found), (found != nil)
}

func (self *Document) addSelectors(selectors []string) {
	for _, s := range selectors {
		for _, f := range self.newSelectorFilter(s) {
			self.add(f)
		}
	}
}

func (self *Document) newSelectorFilter(s string) []Filter {
	chain := []Filter{}
	buf := bytes.NewBuffer(nil)

	var tagSet bool
	var settingClass bool
	var settingID bool

	appendFilter := func(f Filter) {
		if buf.Len() > 0 {
			chain = append(chain, f)
			buf.Truncate(0)
		}
	}

	appendTag := func() {
		appendFilter(Tag(buf.String()))
		tagSet = true
	}

	for _, r := range s {
		if r == '.' {
			settingClass = true
			appendTag()
			continue
		}

		if r == '#' {
			settingID = true
			appendTag()
			continue
		}

		buf.WriteRune(r)
	}

	if tagSet {
		if settingClass {
			appendFilter(Class(buf.String()))
			settingClass = false
		} else if settingID {
			appendFilter(Attribute("id", buf.String()))
			settingID = false
		} else {
			// ...
		}
	} else {
		appendTag()
	}

	return chain
}

func (self *Document) All(selectors ...string) []*Node {
	self.addSelectors(selectors)
	var found []*Node
	self.searchFunc(self.root, func(n *html.Node) {
		if self.match(n) {
			found = append(found, newNode(n))
		}
	})
	return found
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
	if node == nil {
		return ""
	}

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

// search applies visit func using Depth First Search over children
func search(root *html.Node, visit func(n *html.Node)) {
	if root == nil {
		return
	}

	var dfs func(*html.Node)
	dfs = func(n *html.Node) {
		visit(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}
	dfs(root)
}

// searchParent applies visit func using Depth First Search over parents
func searchParent(root *html.Node, visit func(n *html.Node)) {
	if root == nil {
		return
	}

	var dfs func(*html.Node)
	dfs = func(n *html.Node) {
		visit(n)
		for c := n.Parent; c != nil; c = c.Parent {
			dfs(c)
		}
	}
	dfs(root.Parent)
}
