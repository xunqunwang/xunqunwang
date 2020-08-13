package model

type Ecode struct {
	ID       int64  `json:"id" gorm:"primary_key;column:id"`
	Category string `json:"category" gorm:"column:category"`
	Code     int    `json:"code" gorm:"column:code"`
	Message  string `json:"message" gorm:"column:message"`
	Remark   string `json:"remark" gorm:"column:remark"`
}

func (Ecode) TableName() string {
	return "ecode"
}

type EcodeData struct {
	Ver  int64
	MD5  string
	Code map[int]string
}
