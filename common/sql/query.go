package sql

import (
	"fmt"
	"github.com/yanshicheng/ikube-gin-starter/apps/book/model"
	"github.com/yanshicheng/ikube-gin-starter/common/pagination"
	"github.com/yanshicheng/ikube-gin-starter/common/types"
	"gorm.io/gorm"
	"math"
)

func GetQueryResponse(db *gorm.DB, query *model.GetQuery, modelSlice interface{}) (*types.QueryResponse, error) {
	// 获取总记录数
	var totalRecords int64
	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, fmt.Errorf("数据统计查询失败: %s", err.Error())
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(totalRecords) / float64(query.Pagination.PageSize)))

	// 应用分页查询条件
	if err := db.Scopes(
		pagination.PaginateQuery(query.Pagination),
	).Find(&modelSlice).Error; err != nil {
		return nil, fmt.Errorf("数据查询失败: %s", err)
	}

	resp := &types.QueryResponse{
		Page:       query.Pagination.PageSize,
		PageNumber: query.Pagination.PageNumber,
		TotalPage:  totalPages,
		Total:      int(totalRecords),
		Data:       modelSlice,
	}

	return resp, nil
}
