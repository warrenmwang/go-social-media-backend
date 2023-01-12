package main

import (
	"fmt"
	"log"

	"github.com/Warren-Wang-OG/go-social-media-backend/database"
)

func main() {
	c := database.NewClient("db.json")
	err := c.EnsureDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("database ensured!")
}
