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
    <ul class="nav">
      <li><a href="/log_in"<button>Log In</button></a></li>
      <li><a href=" /sign_up " class="active "><button>Sign Up</button></a></li>
      <li><a href="/home" class="one two three">Home</a></li>
      <li><a href="/about">About</a></li>
    </ul>
  </nav>
  <h1>Register</h1>
  <form method="post" action="/sign_up">
    <label><input type="text" name="name" /></label>
    <label><input type="email" name="email" /></label>
    <label><input type="password" name="password" /></label>
    <label><textarea name="bio">Biography</textarea></label>

    <input type="submit" name="Register" />
  </form>
</body>
</html>
`

const TagNotFound = "tag not found"

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
