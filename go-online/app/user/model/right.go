package model

type Right struct {
	Id        int64  `json:"id" form:"id" gorm:"primary_key;column:id"`
	RightName string `json:"right_name" form:"right_name" gorm:"column:right_name"`
}

func (Right) TableName() string {
	return "right_tbl"
}
