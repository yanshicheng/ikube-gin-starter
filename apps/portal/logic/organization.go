package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-starter/apps/portal"
	"github.com/yanshicheng/ikube-gin-starter/apps/portal/model"
	"github.com/yanshicheng/ikube-gin-starter/apps/portal/service"
	types2 "github.com/yanshicheng/ikube-gin-starter/apps/portal/types"
	"github.com/yanshicheng/ikube-gin-starter/common/types"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"github.com/yanshicheng/ikube-gin-starter/router"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 接口检查
var _ service.OrganizationService = (*OrganizationLogic)(nil)

var logic = &OrganizationLogic{}

type OrganizationLogic struct {
	l  *zap.Logger
	db *gorm.DB
}

func (o *OrganizationLogic) Get(c *gin.Context, search types2.OrganizationSearch) ([]*model.Organization, error) {
	o.l.Info(fmt.Sprintf("查询机构信息, name: %s", search.Name))
	var orgs []model.Organization
	if err := o.db.WithContext(c).Where("name like ?", search.Name+"%").Find(&orgs).Error; err != nil {
		o.l.Error(fmt.Sprintf("查询机构信息失败, id: %s, error: %s", search.Name, err.Error()))
		return nil, err
	}
	resultOrgs := make([]*model.Organization, 0)
	for _, org := range orgs {
		resultOrg, err := o.buildParentHierarchy(&org)
		if err != nil {
			o.l.Error(fmt.Sprintf("构建父级层级结构失败, id: %d, error: %s", org.ID, err.Error()))
			return nil, err
		}
		resultOrgs = append(resultOrgs, resultOrg)
	}
	return resultOrgs, nil
}

// 递归获取父级组织并构建层级结构
func (o *OrganizationLogic) buildParentHierarchy(org *model.Organization) (*model.Organization, error) {
	if org.ParentId == 0 {
		return org, nil // 如果没有父级，返回当前组织
	}

	var parent model.Organization
	if err := o.db.Where("id = ?", org.ParentId).First(&parent).Error; err != nil {
		return nil, err
	}

	// 递归构建父级结构
	parentOrg, err := o.buildParentHierarchy(&parent)
	if err != nil {
		return nil, err
	}

	// 将当前组织添加到父级的 Children 中
	parentOrg.Children = append(parentOrg.Children, org)
	return parentOrg, nil
}
func (o *OrganizationLogic) List(c *gin.Context, search types2.OrganizationSearch) ([]*model.Organization, error) {
	if search.Name == "" {
		var allOrgs []model.Organization
		if err := o.db.WithContext(c).Find(&allOrgs).Error; err != nil {
			o.l.Error(fmt.Sprintf("查询机构信息失败, error: %s", err.Error()))
			return nil, err
		}
		orgMap := make(map[uint]*model.Organization, len(allOrgs))
		for i := range allOrgs {
			orgMap[allOrgs[i].ID] = &allOrgs[i]
		}
		// 构建树形结构
		var orgTree []*model.Organization
		for i := range allOrgs {
			org := &allOrgs[i]
			if org.ParentId != 0 {
				if parent, ok := orgMap[uint(org.ParentId)]; ok {
					if parent.Children == nil {
						parent.Children = make([]*model.Organization, 0)
					}
					parent.Children = append(parent.Children, org)
				}
			} else {
				orgTree = append(orgTree, org)
			}
		}
		return orgTree, nil
	} else {
		var orgs []model.Organization
		if err := o.db.WithContext(c).Where("name like ?", search.Name+"%").Find(&orgs).Error; err != nil {
			o.l.Error(fmt.Sprintf("查询机构信息失败, id: %s, error: %s", search.Name, err.Error()))
			return nil, err
		}
		o.l.Debug(fmt.Sprintf("查询机构信息: %+v, %d", orgs, len(orgs)))
		resultOrgs := make([]*model.Organization, 0)
		for _, org := range orgs {
			o.l.Debug(fmt.Sprintf("查询机构信息org: %+v", org))
			resultOrg, err := o.buildParentHierarchy(&org)
			if err != nil {
				o.l.Error(fmt.Sprintf("构建父级层级结构失败, id: %d, error: %s", org.ID, err.Error()))
				return nil, err
			}
			o.l.Debug(fmt.Sprintf("查询机构信息resultOrg: %+v", resultOrg))
			resultOrgs = append(resultOrgs, resultOrg)
		}
		return resultOrgs, nil
	}
}

func (o *OrganizationLogic) Put(c *gin.Context, id types.SearchId, org *model.Organization) error {
	// 查询
	if err := o.db.WithContext(c).Where("id = ?", id.Id).First(&model.Organization{}).Updates(org).Error; err != nil {
		o.l.Error(fmt.Sprintf("更新机构信息失败, id: %d, error: %s", id.Id, err.Error()))
		return err
	}
	return nil
}

func (o *OrganizationLogic) Create(c *gin.Context, org *model.Organization) error {
	if err := o.db.WithContext(c).Create(org).Error; err != nil {
		o.l.Error(fmt.Sprintf("创建机构信息失败, error: %s", err.Error()))
		return err
	}
	return nil
}

func (o *OrganizationLogic) Delete(c *gin.Context, id types.SearchId) error {
	if err := o.db.WithContext(c).Where("id = ?", id.Id).Delete(&model.Organization{}).Error; err != nil {
		o.l.Error(fmt.Sprintf("删除机构信息失败, id: %d, error: %s", id.Id, err.Error()))
		return err
	}
	return nil
}

// 只需要保证 全局对象Config和全局Logger已经加载完成
func (o *OrganizationLogic) Config() {
	o.l = global.L.Named(portal.AppName).Named(portal.AppOrganization).Named("logic")
	o.db = global.DB.GetDb()
}

func (o *OrganizationLogic) Name() string {
	return fmt.Sprintf("%s.%s", portal.AppName, portal.AppOrganization)
}

func init() {
	// 注册
	router.RegistryLogic(logic)
}
