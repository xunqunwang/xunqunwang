package model

import (
	"time"
)

type User struct {
	Id        int64     `json:"id" form:"id" gorm:"primary_key;column:id"`
	NickName  string    `json:"nick_name" form:"nick_name" gorm:"column:nick_name"`
	Email     string    `json:"email" form:"email" gorm:"column:email"`
	LoginName string    `json:"login_name" form:"login_name" gorm:"column:login_name"`
	Password  string    `json:"password" form:"password" gorm:"column:password"`
	RegTime   time.Time `json:"reg_time" form:"reg_time" gorm:"column:reg_time"`
	HeadPhoto string    `json:"head_photo" form:"head_photo" gorm:"column:head_photo"`
	FansNum   int64     `json:"fans_num" form:"fans_num" gorm:"column:fans_num"`
	FollowNum int64     `json:"follow_num" form:"follow_num" gorm:"column:follow_num"`
	CreatedAt time.Time `json:"created_at" form:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" gorm:"column:updated_at"`
	Roles     []*Role   `gorm:"many2many:user_role_relationship;ForeignKey:Id;AssociationForeignKey:Id;"`
}

func (User) TableName() string {
	return "usertbl"
}

type UserRoleR struct {
	UserId int64 `json:"user_id" form:"user_id" gorm:"primary_key;column:user_id"`
	RoleId int64 `json:"role_id" form:"role_id" gorm:"primary_key;column:role_id"`
}

func (UserRoleR) TableName() string {
	return "user_role_relationship"
}

type RoleRightR struct {
	RoleId  int64 `json:"role_id" form:"role_id" gorm:"primary_key;column:role_id"`
	RightId int64 `json:"right_id" form:"right_id" gorm:"primary_key;column:right_id"`
}

func (RoleRightR) TableName() string {
	return "role_right_relationship"
}
