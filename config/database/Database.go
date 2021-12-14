package database

import (
	"github.com/go-pg/pg"
	"log"
	"os"
)

func Connect() *pg.DB {

	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "",
		Database: "postgres",
		Addr:     ":5432",
	})

	if db == nil {
		log.Println("DB 연결 실패")
	}
	else {
		log.Println("DB 연결 성공")
	}

	///////////////////////////////////////
	//repository.CreateInstructionTable(db)
	//repository.CreateNoWordTable(db)
	//repository.CreatePDTable(db)
	//////////////////////////////////////
	return db

}

func Close(db *pg.DB) {
	closeErr := db.Close()
	if closeErr != nil {
		log.Println("연결 종료 실패 : ", closeErr)
		os.Exit(100)
	}
	else {
		log.Println("DB 연결 종료")
	}
}
