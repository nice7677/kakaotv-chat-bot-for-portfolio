package domain

type PDDomain struct {
	tableName struct{} `sql:"kakao_bot_pd"`
	Idx       int      `json:"idx" sql:"idx,pk"`
	Name      string   `json:"name" sql:"name"`
	UserId    string   `json:"user_id" sql:"user_id"`
}
