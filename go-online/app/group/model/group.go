package model

import (
	"time"
)

type GroupParam struct {
	Page       int      `json:"pn" form:"pn" default:"0"`
	Size       int      `json:"ps" form:"ps" default:"0"`
	Keys       []string `json:"keys" form:"keys,split"`
	GroupIDs   []int64  `json:"group_ids" form:"group_ids,split"`
	UIDs       []int64  `json:"uids" form:"uids,split"`
	GroupNames []string `json:"group_names" form:"group_names,split"`
	Owners     []string `json:"owners" form:"owners,split"`
	Locations  []string `json:"locations" form:"locations,split"`
	Categorys  []string `json:"categorys" form:"categorys,split"`
	Labels     []string `json:"labels" form:"labels,split"`
	IsFree     string   `json:"free" form:"free"`
}

type Group struct {
	GroupID          int64     `json:"group_id" form:"group_id" gorm:"primary_key;column:group_id"`
	UID              int64     `json:"uid" form:"uid" gorm:"column:uid"`
	GroupName        string    `json:"group_name" form:"group_name" gorm:"column:group_name"`
	Owner            string    `json:"owner" form:"owner" gorm:"column:owner"`
	Location         string    `json:"location" form:"location" gorm:"column:location"`
	Category         string    `json:"category" form:"category" gorm:"column:category"`
	Label            string    `json:"label" form:"label" gorm:"column:label"`
	CTime            time.Time `json:"created_time" form:"created_time" gorm:"column:created_time"`
	IsFree           bool      `json:"free" form:"free" gorm:"column:free"`
	GroupScale       int64     `json:"group_scale" form:"group_scale" gorm:"column:group_scale"`
	GroupDescription string    `json:"group_description" form:"group_description" gorm:"column:group_description"`
	GroupActivity    int64     `json:"group_activity" form:"group_activity" gorm:"column:group_activity"`
	GroupPhoto       string    `json:"group_photo" form:"group_photo" gorm:"column:group_photo"`
	OwnerDescription string    `json:"owner_description" form:"owner_description" gorm:"column:owner_description"`
	EntryMethod      string    `json:"entry_method" form:"entry_method" gorm:"column:entry_method"`
	GroupRule        string    `json:"group_rule" form:"group_rule" gorm:"column:group_rule"`
	GroupDetail      string    `json:"group_detail" form:"group_detail" gorm:"column:group_detail"`
	CreatedTime      time.Time `json:"created_at" form:"created_at" gorm:"column:created_at"`
	UpdatedTime      time.Time `json:"updated_at" form:"updated_at" gorm:"column:updated_at"`
}

func (Group) TableName() string {
	return "grouptbl"
}

// type GroupDetail struct {
// 	GroupID        int64     `json:"group_id" form:"group_id" gorm:"column:group_id"`
// 	OwnDescription string    `json:"owner_description" form:"owner_description" gorm:"column:owner_description"`
// 	EntryMethod    string    `json:"entry_method" form:"entry_method" gorm:"column:entry_method"`
// 	GroupRule      string    `json:"group_rule" form:"group_rule" gorm:"column:group_rule"`
// 	GroupDetail    string    `json:"group_detail" form:"group_detail" gorm:"column:group_detail"`
// 	CreatedTime    time.Time `json:"created_at" form:"created_at" gorm:"column:created_at"`
// 	UpdatedTime    time.Time `json:"updated_at" form:"updated_at" gorm:"column:updated_at"`
// }

// func (GroupDetail) TableName() string {
// 	return "group_detail"
// }
