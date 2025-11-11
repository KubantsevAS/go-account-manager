package account

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"net/url"
	"time"

	"github.com/fatih/color"
)

type Account struct {
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var letterRunes = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM12346578590-*!")

func NewAccount(login, password, urlString string) (*Account, error) {
	if login == "" {
		return nil, errors.New("empty login")
	}

	_, err := url.ParseRequestURI(urlString)

	if err != nil {
		return nil, errors.New("invalid URL")
	}

	generatedAccount := &Account{
		Login:     login,
		Password:  password,
		Url:       urlString,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if generatedAccount.Password == "" {
		generatedAccount.generatePassword()
	}

	return generatedAccount, nil
}

func (acc *Account) generatePassword() {
	passwordLength := rand.IntN(10)

	for passwordLength == 0 {
		passwordLength = rand.IntN(10)
	}

	var result = make([]rune, passwordLength)

	for integer := range result {
		result[integer] = letterRunes[rand.IntN(len(letterRunes))]
	}

	(*acc).Password = string(result)
}

func (acc Account) OutputData() {
	color.HiBlue(acc.Login)
	color.HiMagenta(acc.Password)
	color.HiGreen(acc.Url)
}

func (acc Account) OutputPassword() {
	fmt.Println("PASSWORD: ", acc.Password)
}

func (acc Account) OutputLogin() {
	color.HiGreen(acc.Login)
}
