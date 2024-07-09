package pagination

import (
	"github.com/yanshicheng/ikube-gin-starter/common/types"
	"gorm.io/gorm"
)

func PaginateQuery(page types.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page.PageNumber - 1) * page.PageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(page.PageSize)
	}
}
