package domain

type NoWordVM struct {
	tableName struct{} `sql:"kakao_bot_no_word"`
	Idx       int      `json:"idx" sql:"idx,pk"`
	Word      string   `json:"word" sql:"word"`
	PDUserID  string   `json:"pd_user_id" sql:"pd_user_id"`
}
