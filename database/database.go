package database

import (
	"os"
	"time"
)

// exported
type Client struct {
	path string
}

// construct a client
func NewClient(path string) Client {
	return Client{path}
}

type databaseSchema struct {
	Users map[string]User `json:"users"`
	Posts map[string]Post `json:"posts"`
}

// User -
type User struct {
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
}

// Post -
type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserEmail string    `json:"userEmail"`
	Text      string    `json:"text"`
}

// create new empty db file (json) at path specified by the client
// overwrite any previous data in file if existed previously
func (c Client) createDB() error {
	err := os.WriteFile(c.path, []byte{}, 0666)
	if err != nil {
		return err
	}
	return nil
}

// EnsureDB -
// check if db exists already, if good do nothing, otherwise create it using createDB
func (c Client) EnsureDB() error {
	_, err := os.ReadFile(c.path)
	// create new db if doesn't exist
	if err != nil {
		return c.createDB()
	}
	// already exists, do nothing
	return nil
}
