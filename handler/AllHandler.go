package handler

import (
	"kakaotv-chat-bot/config/database"
	"kakaotv-chat-bot/domain"
	"kakaotv-chat-bot/repository"
	"log"
)

type AllHandler struct{}

func (allHandler *AllHandler) GetInstruction() *[]domain.InstructionVM {
	instructRepo := &repository.InstructionRepository{}
	db := database.Connect()
	_, instList := instructRepo.SelectAllInstructionVMList(db)
	database.Close(db)
	return instList
}

func (allHandler *AllHandler) SaveInstruction(instructVM *domain.InstructionVM) error {
	instructRepo := &repository.InstructionRepository{}
	db := database.Connect()
	error := instructRepo.SaveInstruction(db, instructVM)
	database.Close(db)
	if error != nil {
		log.Println(error)
	}
	return error

}

func (allHandler *AllHandler) GetNoword() *[]domain.NoWordVM {
	nowordRepo := &repository.NoWordRepository{}
	db := database.Connect()
	_, nowordList := nowordRepo.SelectAllNoWord(db)
	database.Close(db)
	return nowordList
}

func (allHandler *AllHandler) SaveNoword(noWordVM *domain.NoWordVM) error {
	nowordRepo := &repository.NoWordRepository{}
	db := database.Connect()
	error := nowordRepo.SaveNoword(db, noWordVM)
	database.Close(db)
	if error != nil {
		log.Println(error)
	}
	return error
}

func (allHandler *AllHandler) GetPD() *[]domain.PDDomain {
	pdRepo := &repository.PDRepository{}
	db := database.Connect()
	_, pdList := pdRepo.SelectAllPD(db)
	database.Close(db)
	return pdList
}

func (allHandler *AllHandler) DeleteNoword(noWordVM *domain.NoWordVM) error {
	nowordRepo := &repository.NoWordRepository{}
	db := database.Connect()
	error := nowordRepo.DeleteNoword(db, noWordVM)
	database.Close(db)
	if error != nil {
		log.Println(error)
	}
	return error
}

func (allHandler *AllHandler) DeleteInstruct(instructVM *domain.InstructionVM) error {
	instructRepo := &repository.InstructionRepository{}
	db := database.Connect()
	error := instructRepo.DeleteInstruction(db, instructVM)
	database.Close(db)
	if error != nil {
		log.Println(error)
	}
	return error
}
