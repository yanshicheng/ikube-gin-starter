package model

import (
	"database/sql/driver"
	"fmt"
	"github.com/yanshicheng/ikube-gin-starter/common/model"
)

// 角色表，角色菜单关联表，角色账户关联表， 权限表，
func init() {
	model.Register(&Role{}, &RoleMenu{}, &RoleAccount{}, &Upms{})
}

type Role struct {
	model.Model
	Name          string `json:"name" binding:"required,alphanum,max=32" gorm:"type:varchar(32);not null;unique;comment:角色"`
	ApplicationId uint   `json:"applicationId" binding:"required,number" gorm:"type:int;not null;comment:应用" `
}

func (r *Role) TableName() string {
	return "ikubeops_portal_role"
}

type RoleMenu struct {
	model.Model
	RoleId uint `json:"roleId" binding:"required,number" gorm:"type:int;not null;uniqueIndex:idx_role_menu;comment:角色" `
	MenuId uint `json:"menuId" binding:"required,number" gorm:"type:int;not null;uniqueIndex:idx_role_menu;comment:菜单" `
}

func (r *RoleMenu) TableName() string {
	return "ikubeops_portal_role_menu"
}

type RoleAccount struct {
	model.Model
	RoleId    uint `json:"roleId" binding:"required,number" gorm:"type:int;not null;uniqueIndex:idx_role_account;comment:角色"`
	AccountId uint `json:"accountId" binding:"required,number" gorm:"type:int;not null;uniqueIndex:idx_role_account;comment:用户"`
}

func (r *RoleAccount) TableName() string {
	return "ikubeops_portal_role_account"
}

type Upms struct {
	model.Model
	Name     string     `json:"name" binding:"required,max=32" gorm:"type:varchar(32);not null;unique;comment:权限名称"`
	RoleId   uint       `json:"roleId" binding:"required,number" gorm:"type:int;not null;comment:角色"`
	Resource string     `json:"resource" binding:"required,max=255" gorm:"type:varchar(255);not null;comment:资源"`
	Type     ActionType `json:"type"  binding:"required,oneof=0 1" gorm:"type:tinyint;not null;comment:操作类型"`
}

func (u *Upms) TableName() string {
	return "ikubeops_portal_upms"
}

// ActionType  定义 ActionType 类型
type ActionType uint

const (
	ReadAction  ActionType = 0
	WriteAction ActionType = 1
)

// ActionTypeToString 将 ActionType 转换为字符串
func (a ActionType) String() string {
	switch a {
	case ReadAction:
		return "read"
	case WriteAction:
		return "write"
	default:
		return "unknown"
	}
}

// Scan 实现  接口
func (a *ActionType) Scan(value interface{}) error {
	v, ok := value.(uint)
	if !ok {
		return fmt.Errorf("invalid value for ActionType: %v", value)
	}
	*a = ActionType(v)
	return nil
}

// Value 实现 Valuer 接口
func (a ActionType) Value() (driver.Value, error) {
	return uint(a), nil
}
