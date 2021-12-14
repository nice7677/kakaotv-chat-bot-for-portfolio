package repository

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"kakaotv-chat-bot/domain"
	"log"
)

type NoWordRepository struct{}

func CreateNoWordTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: false,
	}
	createErr := db.CreateTable(&domain.NoWordVM{}, opts)
	if createErr != nil {
		log.Println("t_log 테이블 생성 실패 : ", createErr)
		return createErr
	}
	log.Println("금칙어 테이블 생성 성공")

	return nil
}

func (noWordRepository *NoWordRepository) SaveNoword(db *pg.DB, noWordVM *domain.NoWordVM) error {

	tx, txErr := db.Begin()

	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	insertErr := db.Insert(noWordVM)

	if insertErr != nil {
		return insertErr
	}

	tx.Commit()
	return nil

}

func (noWordRepository *NoWordRepository) SelectAllNoWord(db *pg.DB) (error, *[]domain.NoWordVM) {

	tx, txErr := db.Begin()

	if txErr != nil {
		tx.Rollback()
		return txErr, nil
	}

	NoWordVMList := &[]domain.NoWordVM{}

	selectErr := db.Model(NoWordVMList).Order("idx desc").Limit(500).Select()

	if selectErr != nil {
		tx.Rollback()
		return selectErr, nil
	}

	tx.Commit()
	return nil, NoWordVMList

}

func (noWordRepository *NoWordRepository) DeleteNoword(db *pg.DB, nowordVM *domain.NoWordVM) error {
	tx, txErr := db.Begin()

	if txErr != nil {
		log.Println("트랜잭션 에러")
		return errors.New("tx error")
	}
	defer tx.Rollback()

	selectErr := db.Model(nowordVM).Where("word = ?0", nowordVM.Word).Select()
	if selectErr != nil {
		tx.Commit()
		return selectErr
	}

	_, deleteErr := db.Model(nowordVM).Where("word = ?0", nowordVM.Word).Delete()
	if deleteErr != nil {
		log.Println("fail delete word", deleteErr)
		return deleteErr
	}

	tx.Commit()
	return nil
}
