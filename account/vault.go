package account

import (
	"demo/account-manager/encrypter"
	"demo/account-manager/output"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Embedded interface
type Database interface {
	ByteReader
	ByteWriter
}

type ByteReader interface {
	Read() ([]byte, error)
}

type ByteWriter interface {
	Write([]byte)
}

type Vault struct {
	Accounts  []Account `json:"accounts"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type VaultWithDb struct {
	Vault
	db  Database
	enc encrypter.Encrypter
}

func GetVault(db Database, enc encrypter.Encrypter) *VaultWithDb {
	file, err := db.Read()

	if err != nil {
		return &VaultWithDb{
			Vault: Vault{
				Accounts:  []Account{},
				UpdatedAt: time.Now(),
			},
			db:  db,
			enc: enc,
		}
	}

	var vault Vault

	parseErr := json.Unmarshal(file, &vault)

	if parseErr != nil {
		output.PrintError(parseErr)
		return &VaultWithDb{
			Vault: Vault{
				Accounts:  []Account{},
				UpdatedAt: time.Now(),
			},
			db:  db,
			enc: enc,
		}
	}

	return &VaultWithDb{
		Vault: vault,
		db:    db,
		enc:   enc,
	}
}

func (vault *VaultWithDb) AddAccount(acc Account) {
	vault.Accounts = append(vault.Accounts, acc)
	vault.save()
}

func (vault *VaultWithDb) FindAccount(requiredUrl string, property string, checker func(Account, string, string) bool) []Account {
	var requiredAccounts []Account

	for _, vaultAccount := range vault.Accounts {
		hasAccount := checker(vaultAccount, property, requiredUrl)

		if hasAccount {
			requiredAccounts = append(requiredAccounts, vaultAccount)
		}
	}

	return requiredAccounts
}

func (vault *VaultWithDb) DeleteAccounts(requiredUrl string) {
	var requiredAccounts []Account

	for _, vaultAccount := range vault.Accounts {
		if !strings.Contains(vaultAccount.Url, requiredUrl) {
			requiredAccounts = append(requiredAccounts, vaultAccount)
		}
	}

	vault.Accounts = requiredAccounts
	vault.save()
	fmt.Println("Delete Success")
}

func (vault *Vault) ToBytes() ([]byte, error) {
	file, err := json.Marshal(vault)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (vault *VaultWithDb) save() {
	vault.UpdatedAt = time.Now()

	data, err := vault.Vault.ToBytes()

	if err != nil {
		output.PrintError(err)
		return
	}

	vault.db.Write(data)
}
