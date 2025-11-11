package file

import (
	"demo/account-manager/output"
	"fmt"
	"os"
)

type JsonDb struct {
	fileName string
}

func NewJsonDb(name string) *JsonDb {
	return &JsonDb{
		fileName: name,
	}
}

func (db *JsonDb) Write(content []byte) {
	file, err := os.Create(db.fileName)

	if err != nil {
		output.PrintError(err)
		return
	}

	defer file.Close()

	_, err = file.Write(content)

	if err != nil {
		output.PrintError(err)
		return
	}

	fmt.Println("Write success")
}

func (db *JsonDb) Read() ([]byte, error) {
	data, err := os.ReadFile(db.fileName)

	if err != nil {
		output.PrintError(err)
		return nil, err
	}

	return data, nil
}
