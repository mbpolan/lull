package ui

import (
	"github.com/mbpolan/lull/internal/state"
	"github.com/rivo/tview"
)

// URLBox is a view that contains an HTTP method, URL and other input components.
type URLBox struct {
	flex           *tview.Flex
	form           *tview.Form
	allowedMethods []string
	state          *state.AppState
}

// NewURLBox returns a new instance of URLBox.
func NewURLBox(state *state.AppState) *URLBox {
	u := new(URLBox)
	u.state = state
	u.allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	u.build()

	return u
}

// Widget returns a primitive widget containing this component.
func (u *URLBox) Widget() *tview.Form {
	return u.form
}

func (u *URLBox) build() {
	curMethod := u.currentMethod()

	u.form = tview.NewForm()
	u.form.SetBorder(true)
	u.form.SetHorizontal(true)

	u.form.AddDropDown("", u.allowedMethods, curMethod, u.handleMethodChanged)
	u.form.AddInputField("", u.state.URL, 500, nil, u.handleURLChanged)
}

func (u *URLBox) currentMethod() int {
	for i, method := range u.allowedMethods {
		if method == u.state.Method {
			return i
		}
	}

	return -1
}

func (u *URLBox) handleMethodChanged(text string, index int) {
	u.state.Method = text
}

func (u *URLBox) handleURLChanged(text string) {
	u.state.URL = text
}
