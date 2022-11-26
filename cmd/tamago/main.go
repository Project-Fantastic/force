package main

import (
	tamago "tamago/internal"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	tamago.Run()
}
