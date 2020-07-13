package model

type TokenObj struct {
	Id     int64  `json:"id" form:"id" gorm:"primary_key;column:id"`
	Token  string `json:"token" form:"token" gorm:"column:token"`
	UserId int64  `json:"user_id" form:"user_id" gorm:"column:user_id"`
}

func (TokenObj) TableName() string {
	return "token_tbl"
}
