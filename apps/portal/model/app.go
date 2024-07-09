package model

import (
	"fmt"
	"github.com/yanshicheng/ikube-gin-starter/common/model"
	"gorm.io/gorm"
)

// 应用表，用户应用关联表，菜单表
func init() {
	model.Register(&Application{}, &Menu{}, &AccountApplication{})
}

const MenuLevel = 5

type Application struct {
	model.Model
	Name string `json:"name" binding:"required,max=32" gorm:"type:varchar(32);not null;uniqueIndex;comment:应用"`
	Path string `json:"path" binding:"required,max=10" gorm:"type:varchar(10);not null;unique;comment:路径"`
	Icon string `json:"icon"  gorm:"type:varchar(32);not null;unique;comment:图标"`
	Desc string `json:"desc" gorm:"type:varchar(56);not null;comment:描述"`
}

func (a *Application) TableName() string {
	return "ikubeops_portal_application"
}

type AccountApplication struct {
	model.Model
	AccountId     uint `json:"accountId" binding:"required,number" gorm:"type:int;not null;uniqueIndex:idx_account_app;comment:用户"`
	ApplicationId uint `json:"applicationId" binding:"required,number" gorm:"type:int;not null;uniqueIndex:idx_account_app;comment:应用"`
}

func (a *AccountApplication) TableName() string {
	return "ikubeops_portal_account_application"
}

type Menu struct {
	model.Model
	Path             string  `json:"path" binding:"required,max=32" gorm:"type:varchar(32);not null;comment:路由路径"`
	Name             string  `json:"name" binding:"required,max=32" gorm:"type:varchar(32);not null;unique;comment:唯一标识名称" `
	Component        string  `json:"component" binding:"required,max=255" gorm:"type:varchar(255);not null;comment:组件路径" `
	Redirect         string  `json:"redirect" binding:"max=255" gorm:"type:varchar(255);comment:重定向路径" `
	Title            string  `json:"title" binding:"max=26" gorm:"type:varchar(26);not null;comment:菜单标题" `
	Icon             string  `json:"icon"  binding:"max=32" gorm:"type:varchar(32);comment:菜单图标" `
	Expanded         bool    `json:"expanded"  binding:"boolean" gorm:"type:tinyint(1);default:false;comment:是否默认展开" `
	OrderNo          int     `json:"orderNo" binding:"required,number" gorm:"type:tinyint;not null;comment:菜单顺序编号" `
	Hidden           bool    `json:"hidden" binding:"required" gorm:"default:false;comment:是否隐藏菜单"`
	HiddenBreadcrumb bool    `json:"hiddenBreadcrumb" binding:"boolean" gorm:"type:tinyint(1);default:false;comment:是否隐藏面包屑"`
	Single           bool    `json:"single" binding:"boolean" gorm:"type:tinyint(1);default:false;comment:是否单级菜单显示"`
	FrameSrc         string  `json:"frameSrc" gorm:"type:varchar(255);comment:内嵌iframe的地址"`
	FrameBlank       bool    `json:"frameBlank" binding:"boolean" gorm:"type:tinyint(1);default:false;comment:内嵌iframe是否新窗口打开" `
	KeepAlive        bool    `json:"keepAlive" binding:"boolean" gorm:"type:tinyint(1);default:true;comment:开启keep-alive"`
	ParentId         uint    `json:"parentId"  binding:"required,number" gorm:"type:int;not null;comment:父级"` // 关联父级路由
	ApplicationId    uint    `json:"applicationId" binding:"required,number"  gorm:"type:int;not null"`
	Level            int     `json:"level" gorm:"type:int;not null;comment:层级"`
	Children         []*Menu `gorm:"-" json:"children"` // 子路由，不存储在数据库中，只用于加载和显示

}

func (m *Menu) TableName() string {
	return "ikubeops_portal_menu"
}

// 机构表 创建钩子函数
func (o *Menu) BeforeCreate(tx *gorm.DB) error {
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
		if o.Level > MenuLevel {
			// 如果层级超过5，返回错误
			return fmt.Errorf("cannot add beyond level %d", MenuLevel)
		}
	} else {
		// 如果 ParentID 为0，说明此节点没有父节点，即它是一个根节点
		o.Level = 1 // 设置根节点的层级为1
	}
	// 如果所有检查都通过，没有错误，则返回 nil，允许创建操作继续进行
	return nil
}
