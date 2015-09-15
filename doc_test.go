package htmldoc

import (
	"reflect"
	"strings"
	"testing"
)

const TestPage = `
<!doctype html>
<html>
<head><title>Test Page</title></head>
<body>
  <nav>
    <ul id="nav" class="fancy">
      <li><a href="/log_in"><button>Log In</button></a></li>
      <li><a href=" /sign_up " class="active "><button>Sign Up</button></a></li>
      <li><a href="/home" class="one two three">Home</a></li>
      <li><a href="/about">About</a></li>
    </ul>
  </nav>
  <h1>Register</h1>
  <form method="post" action="/sign_up">
    <fieldset>
      <label><input type="text" name="Name" /></label>
      <label><input type="email" name="Email" /></label>
      <label><input type="password" name="Password" /></label>
      <label><textarea name="Bio">Biography</textarea></label>
    </fieldset>

    <input type="submit" name="Register" />
  </form>
</body>
</html>
`

const TagNotFound = "tag not found"

func TestDocument_Filter(t *testing.T) {
	writeable := Or(
		Tag("textarea"),
		And(
			Tag("input"),
			Or(
				Attribute("type", "text"),
				Attribute("type", "email"),
				Attribute("type", "password"),
			),
		),
	)
	all := MustNew(TestPage).Filter(writeable).All()

	if len(all) != -1 {
		t.Error("IMPLEMENT ME")
	}
}

func TestBySelector_AllTag(t *testing.T) {
	if len(MustNew(TestPage).All("a")) != 4 {
		t.Error("want", 4)
	}
}

func TestBySelector_Tag(t *testing.T) {
	node, _ := MustNew(TestPage).First("a")
	Expect(t, "Found By Selector").Equal(node.Text(), "About")
}

func TestBySelector_TagAndClass(t *testing.T) {
	node, _ := MustNew(TestPage).First("a.active")

	Expect(t, "Found By Selector").Equal(node.Text(), "Sign Up")
}

func TestBySelector_TagAndID(t *testing.T) {
	node, found := MustNew(TestPage).First("ul#nav")
	Expect(t, "Found ul#nav").Equal(found, true)
	Expect(t, "ul#nav class name").Equal(node.Attribute("class"), "fancy")
}

func TestParent(t *testing.T) {
	input, _ := MustNew(TestPage).
		Tag("input").
		Attribute("type", "email").
		First()

	form, found := input.Parent().Tag("form").First()

	Expect(t, "Found").Equal(found, true)
	Expect(t, "Form Action").Equal(form.Attribute("action"), "/sign_up")
}

func TestAll_Simple(t *testing.T) {
	nodes := MustNew(TestPage).Tag("a").All()

	Expect(t, "Number Found").Equal(len(nodes), 4)
}

func TestFirst_Simple(t *testing.T) {
	node, found := MustNew(TestPage).Tag("a").First()

	Expect(t, "Found").Equal(found, true)
	Expect(t, "Node Text").Equal(node.Text(), "About")
}

func TestFirst_Multi(t *testing.T) {
	node, found := MustNew(TestPage).
		Tag("a").
		Class("active").
		Attribute("href", "/sign_up").
		First()

	Expect(t, "Found").Equal(found, true)
	Expect(t, "Node Text").Equal(node.Text(), "Sign Up")
}

func TestFirst_MultiNotFound(t *testing.T) {
	_, found := MustNew(TestPage).Tag("a").
		Class("bork").
		Attribute("href", "/bork").
		First()

	Expect(t, "Not Found").Equal(found, false)
}

func Expect(t *testing.T, title ...string) *expectation {
	return &expectation{
		t:     t,
		title: strings.Join(title, " "),
	}
}

type expectation struct {
	t     *testing.T
	title string
}

func (self *expectation) Equal(a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		self.error(a, "does not equal", b)
	}
}

func (self *expectation) error(message ...interface{}) {
	args := []interface{}{self.title}
	args = append(args, message...)
	self.t.Error(args...)
}
