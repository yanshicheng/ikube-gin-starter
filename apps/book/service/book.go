package service

import (
	"github.com/yanshicheng/ikube-gin-starter/apps/book/model"
	"github.com/yanshicheng/ikube-gin-starter/common/types"
)

type BookService interface {
	Create(*model.Book) error
	Update(types.SearchId, *model.Book) error
	Delete(types.SearchId) error
	Get(types.SearchId) (*model.Book, error)
	List(*model.GetQuery) (*types.QueryResponse, error)
}
