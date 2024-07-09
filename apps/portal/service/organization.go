package service

import (
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-starter/apps/portal/model"
	otypes "github.com/yanshicheng/ikube-gin-starter/apps/portal/types"
	"github.com/yanshicheng/ikube-gin-starter/common/types"
)

type OrganizationService interface {
	Get(*gin.Context, otypes.OrganizationSearch) ([]*model.Organization, error)
	List(*gin.Context, otypes.OrganizationSearch) ([]*model.Organization, error)
	Create(*gin.Context, *model.Organization) error
	Put(*gin.Context, types.SearchId, *model.Organization) error
	Delete(*gin.Context, types.SearchId) error
}
