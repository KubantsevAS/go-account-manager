package main

import (
	"demo/account-manager/account"
	// "demo/account-manager/cloud"
	"demo/account-manager/output"

	"demo/account-manager/file"
	"fmt"
	"reflect"
	"strings"
)

// var vault = account.GetVault(cloud.NewCloudDb("https://github.com"))
var vault = account.GetVault(file.NewJsonDb("vault.json"))

func main() {
	fmt.Println("***Account Manager***")
	genericCallMenu(&[]string{
		"1. 'create'",
		"2. 'find'",
		"3. 'delete'",
		"4. 'exit'",
		"Enter one of commands to start work with accounts",
	})
	callMenu()

Menu:
	for {
		var command string

		_, err := fmt.Scanln(&command)

		if err != nil {
			output.PrintError(err)
			continue
		}

		var optionsMap = map[string]func() bool{
			"create": func() bool { createAccount(); return false },
			"find":   func() bool { findAccount(); return false },
			"delete": func() bool { deleteAccount(); return false },
			"menu":   func() bool { callMenu(); return false },
			"exit":   func() bool { return true },
		}

		menuFunc := optionsMap[command]

		if menuFunc == nil {
			output.PrintError("Invalid input")
		}

		if menuFunc() {
			break Menu
		}
	}

	fmt.Println("** EXIT Account Manager**")
}

func callMenu() {
	fmt.Println(`Enter one of commands to start work with accounts:
1. 'create'
2. 'find'
3. 'delete'
4. 'exit'`)
}

func findAccount() {
	requiredUrl := promptData("Enter url to find all accounts: ")
	requiredAccounts := vault.FindAccount(requiredUrl, "Url", checkString)

	accountsCount := len(requiredAccounts)

	if accountsCount == 0 {
		output.PrintError("Accounts are not found")
	}

	fmt.Printf("Total found: %d\n", accountsCount)

	for index, acc := range requiredAccounts {
		fmt.Printf("Account %d:\n", index+1)
		acc.OutputData()
	}
}

func deleteAccount() {
	requiredUrl := promptData("Enter url to delete those accounts: ")
	vault.DeleteAccounts(requiredUrl)
}

func createAccount() {
	login := promptData("Enter login: ")
	password := promptData("Enter password: ")
	url := promptData("Enter url: ")

	myAccount, err := account.NewAccount(login, password, url)

	if err != nil {
		output.PrintError(err)
		return
	}

	vault.AddAccount(*myAccount)
}

func promptData(prompt string) string {
	fmt.Print(prompt)
	var res string
	fmt.Scanln(&res)
	return res
}

func genericCallMenu[T any](prompt *[]T) {
	arrayLength := len(*prompt)

	for index, value := range *prompt {
		if index+1 == arrayLength {
			fmt.Printf("%v:", value)
		} else {
			fmt.Println(value)
		}
	}
}

func checkString(acc account.Account, property string, requiredValue string) bool {
	// Use reflection to access struct field by name
	val := reflect.ValueOf(acc)
	field := val.FieldByName(property)

	if !field.IsValid() || field.Kind() != reflect.String {
		return false
	}

	return strings.Contains(field.String(), requiredValue)
}
