package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

const fieldWidth = 30

type Contract struct {
	paymentFrequency string
	amount           int
}

type Client struct {
	company  string
	email    string
	phone    string
	contract Contract
}

var clients []Client

var pages = tview.NewPages()
var clientText = tview.NewTextView()
var app = tview.NewApplication()
var newClientForm = tview.NewForm()
var clientList = tview.NewList().ShowSecondaryText(false)
var flex = tview.NewFlex()
var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("(a) to add a new contact \n(q) to quit")

func setConcatText(client *Client) {
	clientText.Clear()
	text := client.company + "\n" + client.email + "\n" + client.phone
	clientText.SetText(text)
}

func main() {
	// setting the input capture on the flex and not the app (which is the parent of all pages)
	// means that when we go on the form page we can use q and a in our forms without it exiting
	// or opening new forms
	clientList.SetSelectedFunc(func(index int, name string, secondName string, shortcut rune) {
		setConcatText(&clients[index])
	})

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(clientList, 0, 1, true).
			AddItem(clientText, 0, 4, false), 0, 6, false).
		AddItem(text, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.Stop()
		} else if event.Rune() == 97 {
			newClientForm.Clear(true) // to prevent use adding boxes to the same form
			addClientForm()
			pages.SwitchToPage("new client form")
		}
		return event
	})
	pages.AddPage("home", flex, true, true)
	pages.AddPage("new client form", newClientForm, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func addClientForm() {
	client := Client{}
	// we set label of input box appropriately, an empty initial value, a width of 30, no validation function
	// and the last parameter is a function that updates our new client object accordingly
	newClientForm.AddInputField("Company name", "", fieldWidth, nil, func(company string) {
		client.company = company
	})
	newClientForm.AddInputField("Email", "", fieldWidth, nil, func(email string) {
		client.email = email
	})
	newClientForm.AddInputField("Phone", "", fieldWidth, nil, func(phone string) {
		client.phone = phone
	})
	// drop down menu to select contract type (weekly, montly, yearly)
	contractTypes := []string{"Weekly", "Monthly", "Yearly"}
	// initial option set to -1 -> no initial option selected
	newClientForm.AddDropDown("Contract type", contractTypes, -1,
		func(contractType string, index int) {
			client.contract.paymentFrequency = contractType
		})
	newClientForm.AddInputField("Amount", "", fieldWidth, nil, func(amount string) {
		client.contract.amount, _ = strconv.Atoi(amount)
		// TODO: deal with a non numeric input (underscore replace with err)
	})
	newClientForm.AddButton("Add client", func() {
		/* once we click this button, we want to do 3 things:
		1) add the new client to our back-end collection of clients
		2) add the new client to the list displayed
		3) return to the home page
		*/
		clients = append(clients, client)
		addClientList()
		pages.SwitchToPage("home")
	})
}

// this function adds clients to the WIDGET that stores clients
// this doesn't add to our clients array, that is done in addClientForm
func addClientList() {
	clientList.Clear()
	for i, client := range clients {
		// 49 is ascii for 1, rune(49 + i) will number the clients in increasing order
		clientList.AddItem(client.company, "", rune(49+i), nil)
	}
}
