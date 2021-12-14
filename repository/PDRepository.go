package repository

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"kakaotv-chat-bot/domain"
	"log"
)

type PDRepository struct {
}

func CreatePDTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: false,
	}
	createErr := db.CreateTable(&domain.PDDomain{}, opts)
	if createErr != nil {
		log.Println("t_log 테이블 생성 실패 : ", createErr)
		return createErr
	}
	log.Println("PD 테이블 생성 성공")

	return nil
}

func (pDRepository *PDRepository) SelectAllPD(db *pg.DB) (error, *[]domain.PDDomain) {

	tx, txErr := db.Begin()

	if txErr != nil {
		tx.Rollback()
		return txErr, nil
	}

	pdVMList := &[]domain.PDDomain{}

	selectErr := db.Model(pdVMList).Order("idx asc").Limit(500).Select()

	if selectErr != nil {
		tx.Rollback()
		return selectErr, nil
	}

	tx.Commit()
	return nil, pdVMList

}
