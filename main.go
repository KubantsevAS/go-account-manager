package main

import (
	"demo/account-manager/account"
	"demo/account-manager/encryptor"

	// "demo/account-manager/cloud"
	"demo/account-manager/file"
	"demo/account-manager/output"

	"fmt"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("***Account Manager***")
	envErr := godotenv.Load()

	if envErr != nil {
		output.PrintError(envErr)
	}

	var vault = account.InitVault(file.NewJsonDb("vault.vault"), *encryptor.NewEncryptor())
	// var vault = account.GetVault(cloud.NewCloudDb("https://github.com"))

	var menuItems = []string{
		"1. 'create'",
		"2. 'find'",
		"3. 'delete'",
		"4. 'menu'",
		"5. 'exit'",
		"Enter one of commands to start work with accounts",
	}

	printMultiLine(menuItems...)

Menu:
	for {
		var command string

		_, err := fmt.Scanln(&command)

		if err != nil {
			output.PrintError(err)
			continue
		}

		var optionsMap = map[string]func() bool{
			"create": func() bool { createAccount(vault); return false },
			"find":   func() bool { findAccount(vault); return false },
			"delete": func() bool { deleteAccount(vault); return false },
			"menu":   func() bool { printMultiLine(menuItems...); return false },
			"exit":   func() bool { return true },
		}

		menuFunc := optionsMap[command]

		if menuFunc == nil {
			output.PrintError("Invalid input")
			continue
		}

		if menuFunc() {
			break Menu
		}
	}

	fmt.Println("** EXIT Account Manager**")
}

func findAccount(vault *account.VaultWithDb) {
	var menuItems = []string{
		"1. 'Login'",
		"2. 'Url'",
		"Choose witch property to use",
	}
	printMultiLine(menuItems...)

	var property string
	fmt.Scanln(&property)

	if property != "Login" && property != "Url" {
		output.PrintError("Invalid property")
		return
	}

	requiredUrl := promptData("Enter property value to find all accounts: ")

	requiredAccounts := vault.FindAccount(requiredUrl, property, checkPropertyString)

	accountsCount := len(requiredAccounts)

	if accountsCount == 0 {
		output.PrintError("Accounts are not found")
	}

	fmt.Printf("Total found: %d\n", accountsCount)

	for index, acc := range requiredAccounts {
		fmt.Printf("Account %d: \n", index+1)
		acc.OutputData()
	}
}

func deleteAccount(vault *account.VaultWithDb) {
	requiredUrl := promptData("Enter url to delete those accounts: ")
	vault.DeleteAccounts(requiredUrl)
}

func createAccount(vault *account.VaultWithDb) {
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

func printMultiLine(prompt ...string) {
	arrayLength := len(prompt)

	for index, value := range prompt {
		if index+1 == arrayLength {
			fmt.Printf("%v: ", value)
		} else {
			fmt.Println(value)
		}
	}
}

func checkPropertyString(acc account.Account, property string, requiredValue string) bool {
	// Reflection to access struct field by name
	val := reflect.ValueOf(acc)
	field := val.FieldByName(property)

	if !field.IsValid() || field.Kind() != reflect.String {
		return false
	}

	return strings.Contains(field.String(), requiredValue)
}
