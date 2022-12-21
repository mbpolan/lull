package ui

import "github.com/rivo/tview"

// URLBox is a view that contains an HTTP method, URL and other input components.
type URLBox struct {
	flex           *tview.Flex
	form           *tview.Form
	method         *tview.DropDown
	url            *tview.InputField
	allowedMethods []string
	vm             *viewModel
}

type viewModel struct {
	method int
	url    string
}

// NewURLBox returns a new instance of URLBox.
func NewURLBox() *URLBox {
	u := new(URLBox)
	u.allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	u.vm = new(viewModel)
	u.vm.method = 0
	u.build()

	return u
}

// Widget returns a primitive widget containing this component.
func (u *URLBox) Widget() *tview.Form {
	return u.form
}

func (u *URLBox) build() {
	u.form = tview.NewForm()
	u.form.SetBorder(true)
	u.form.SetHorizontal(true)

	u.method = tview.NewDropDown()
	u.method.SetCurrentOption(u.vm.method)
	u.method.SetOptions(u.allowedMethods, u.handleMethodChanged)

	u.url = tview.NewInputField()

	u.form.AddDropDown("", u.allowedMethods, u.vm.method, u.handleMethodChanged)
	u.form.AddInputField("", u.vm.url, 500, nil, nil)
}

func (u *URLBox) handleMethodChanged(text string, index int) {
	u.vm.method = index
}
