package repository

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"kakaotv-chat-bot/domain"
	"log"
)

type InstructionRepository struct{}

func CreateInstructionTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: false,
	}
	createErr := db.CreateTable(&domain.InstructionVM{}, opts)
	if createErr != nil {
		log.Println("t_log 테이블 생성 실패 : ", createErr)
		return createErr
	}
	log.Println("명령어 테이블 생성 성공")

	return nil
}

func (instructionRepository *InstructionRepository) SaveInstruction(db *pg.DB, instructionVM *domain.InstructionVM) error {

	tx, txErr := db.Begin()

	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	insertErr := db.Insert(instructionVM)

	if insertErr != nil {
		return insertErr
	}

	tx.Commit()
	return nil

}

func (instructionRepository *InstructionRepository) SelectAllInstructionVMList(db *pg.DB) (error, *[]domain.InstructionVM) {

	tx, txErr := db.Begin()

	if txErr != nil {
		return txErr, nil
	}
	defer tx.Rollback()

	instructionVMList := &[]domain.InstructionVM{}

	selectErr := db.Model(instructionVMList).Order("idx desc").Limit(500).Select()

	if selectErr != nil {
		return selectErr, nil
	}

	tx.Commit()
	return nil, instructionVMList

}

func (instructionRepository *InstructionRepository) DeleteInstruction(db *pg.DB, instructionVM *domain.InstructionVM) error {
	tx, txErr := db.Begin()

	if txErr != nil {
		log.Println("트랜잭션 에러")
		return errors.New("tx error")
	}
	defer tx.Rollback()

	selectErr := db.Model(instructionVM).Where("function = ?0", instructionVM.Function).Select()
	if selectErr != nil {
		tx.Commit()
		return selectErr
	}

	_, deleteErr := db.Model(instructionVM).Where("function = ?", instructionVM.Function).Delete()
	if deleteErr != nil {
		log.Println("fail delete function", deleteErr)
		return deleteErr
	}

	tx.Commit()
	return nil
}
