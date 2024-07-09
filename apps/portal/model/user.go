package model

import (
	"fmt"
	"github.com/yanshicheng/ikube-gin-starter/common/model"
	"gorm.io/gorm"
	"time"
)

// 用户表 和 机构表

func init() {
	model.Register(&Account{}, &Organization{})
}

const OrganizationLevel = 5

type Account struct {
	model.Model
	UserName       string    `json:"userName" binding:"required,max=32" gorm:"type:varchar(32);not null;comment:姓名"`
	Account        string    `json:"account" binding:"required,max=32" gorm:"type:varchar(32);unique_index;not null;comment:账号"`
	Password       string    `json:"password" binding:"max=24" gorm:"type:varchar(256);not null;comment:密码"`
	Icon           string    `json:"icon"  gorm:"type:varchar(256);not null;comment:头像"`
	Mobile         string    `json:"mobile" binding:"required,max=11" gorm:"type:char(11);unique_index;not null;comment:手机号"`
	Email          string    `json:"email" binding:"required,max=36,email" gorm:"type:varchar(36);unique_index;not null;comment:邮箱"`
	WorkNumber     string    `json:"workNumber" binding:"required,max=24" gorm:"type:varchar(24);unique_index;not null;comment:工号"`
	HireDate       time.Time `json:"hireDate" binding:"required" gorm:"type:date;not null;comment:入职时间"`
	IsFrozen       bool      `json:"isFrozen" binding:"boolean" gorm:"type:tinyint(1);not null;default:false;comment:是否冻结"`
	IsDisabled     bool      `json:"isDisabled" binding:"boolean" gorm:"type:tinyint(1);not null;default:false;comment:是否禁用"`
	IsLeave        bool      `json:"isLeave" binding:"boolean" gorm:"type:tinyint(1);not null;default:false;comment:是否离职"`
	Position       int       `json:"position" binding:"required,number" gorm:"type:int;not null;comment:职位"` // 对应职位表
	OrganizationId uint      `json:"organizationId" binding:"required,number" gorm:"type:int;not null;comment:组织"`
}

// 定义表名
func (u *Account) TableName() string {
	return "ikubeops_portal_account"
}

type Organization struct {
	model.Model
	Name     string          `json:"name" binding:"required,max=32" gorm:"type:varchar(32);ngit ot null;comment:团队"`
	ParentId uint            `json:"parentId" binding:"number" gorm:"type:int;not null;comment:父级"`
	Level    int             `json:"level" gorm:"type:int;not null;comment:层级"`
	Desc     string          `json:"Desc" gorm:"type:varchar(56);not null;comment:描述"`
	Children []*Organization `gorm:"-"` // 使用指针类型存储子组织
}

func (u *Organization) TableName() string {
	return "ikubeops_portal_organization"
}

// 机构表 创建钩子函数
func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	// 检查是否有父节点，如果没有父节点，则为根节点
	if o.ParentId != 0 {
		// 如果 ParentID 不为0，说明此节点有父节点

		var parent Organization
		// 查询父节点的详细信息
		// 这里使用 tx.First 来查询具有指定 ID 的父节点
		// o.ParentID 是父节点的 ID，将结果存储在 parent 变量中
		if err := tx.First(&parent, o.ParentId).Error; err != nil {
			// 如果查询过程中出现错误，例如数据库连接错误或找不到指定的父节点
			return err // 返回错误，中断创建操作
		}

		// 如果父节点查询成功，设置当前节点的层级为父节点层级 + 1
		o.Level = parent.Level + 1

		// 检查层级是否超过5
		if o.Level > OrganizationLevel {
			// 如果层级超过5，返回错误
			return fmt.Errorf("cannot add beyond level %d", OrganizationLevel)
		}
	} else {
		// 如果 ParentID 为0，说明此节点没有父节点，即它是一个根节点
		o.Level = 1 // 设置根节点的层级为1
	}
	// 如果所有检查都通过，没有错误，则返回 nil，允许创建操作继续进行
	return nil
}

func (org *Organization) GetAllDescendants(db *gorm.DB) ([]*Organization, error) {
	var allOrgs []*Organization
	err := db.Find(&allOrgs).Error
	if err != nil {
		return nil, err
	}

	orgMap := make(map[uint]*Organization)
	var rootOrgs []*Organization

	// 将所有组织存入 map，便于快速查找
	for _, o := range allOrgs {
		orgMap[o.ID] = o
	}

	// 构建树形结构
	for _, o := range allOrgs {
		if o.ParentId == 0 {
			rootOrgs = append(rootOrgs, o)
		} else {
			parent := orgMap[uint(o.ParentId)]
			parent.Children = append(parent.Children, o)
		}
	}

	if org.ID == 0 {
		return rootOrgs, nil
	} else {
		return orgMap[uint(org.ID)].Children, nil
	}
}
