package main

import (
	"flamethrower/src/db"
	"fmt"
	"log"

	"github.com/blockloop/scan"
)

func main() {
	db.InitDB("dnd35.db")
	repo := &db.ClassListRepo{BaseRepo: &db.BaseRepo{}}
	rows, err := repo.FindAll().Paginate(1, 6).Query()
	if err != nil {
		log.Fatal(err)
	}
	var data []db.ClassListColumns
	err = scan.Rows(&data, rows)
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range data {
		fmt.Println(row.ClassSkills.String)
	}
}
