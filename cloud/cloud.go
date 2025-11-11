package cloud

type CloudDb struct {
	url string
}

func NewCloudDb(url string) *CloudDb {
	return &CloudDb{
		url: url,
	}
}

func (db *CloudDb) Write(content []byte) {
}

func (db *CloudDb) Read() ([]byte, error) {
	return nil, nil
}
