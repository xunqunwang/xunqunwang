package model

type Role struct {
	Id       int64    `json:"id" form:"id" gorm:"primary_key;column:id"`
	RoleName string   `json:"role_name" form:"role_name" gorm:"column:role_name"`
	Rights   []*Right `gorm:"many2many:role_right_relationship;ForeignKey:Id;AssociationForeignKey:Id;"`
}

func (Role) TableName() string {
	return "role"
}
