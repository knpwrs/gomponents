package gomponents

import (
	"fmt"
	"html/template"
	"strings"
)

// Node is a DOM node that can Render itself to a string representation.
type Node interface {
	Render() string
}

// NodeFunc is render function that is also a Node.
type NodeFunc func() string

func (n NodeFunc) Render() string {
	return n()
}

// El creates an element DOM Node with a name and child Nodes.
// Use this if no convenience creator exists.
func El(name string, children ...Node) NodeFunc {
	return func() string {
		var b, attrString, childrenString strings.Builder

		b.WriteString("<")
		b.WriteString(name)

		if len(children) == 0 {
			b.WriteString("/>")
			return b.String()
		}

		for _, c := range children {
			s := c.Render()
			if _, ok := c.(attr); ok {
				attrString.WriteString(s)
				continue
			}
			childrenString.WriteString(c.Render())
		}

		b.WriteString(attrString.String())

		if childrenString.Len() == 0 {
			b.WriteString("/>")
			return b.String()
		}

		b.WriteString(">")
		b.WriteString(childrenString.String())
		b.WriteString("</")
		b.WriteString(name)
		b.WriteString(">")
		return b.String()
	}
}

// Attr creates an attr DOM Node.
// If one parameter is passed, it's a name-only attribute (like "required").
// If two parameters are passed, it's a name-value attribute (like `class="header"`).
// More parameter counts make Attr panic.
// Use this if no convenience creator exists.
func Attr(name string, value ...string) Node {
	switch len(value) {
	case 0:
		return attr{name: name}
	case 1:
		return attr{name: name, value: &value[0]}
	default:
		panic("attribute must be just name or name and value pair")
	}
}

type attr struct {
	name  string
	value *string
}

func (a attr) Render() string {
	if a.value == nil {
		return fmt.Sprintf(" %v", a.name)
	}
	return fmt.Sprintf(` %v="%v"`, a.name, *a.value)
}

// Text creates a text DOM Node that Renders the escaped string t.
func Text(t string) NodeFunc {
	return func() string {
		return template.HTMLEscaper(t)
	}
}

// Raw creates a raw Node that just Renders the unescaped string t.
func Raw(t string) NodeFunc {
	return func() string {
		return t
	}
}
