package domain

type InstructionVM struct {
	tableName struct{} `sql:"kakao_bot_instruction"`
	Idx       int      `json:"idx" sql:"idx,pk"`
	Word      string   `json:"word" sql:"word"`
	Function  string   `json:"function" sql:"function"`
	PDUserID  string   `json:"pd_user_id" sql:"pd_user_id"`
}
