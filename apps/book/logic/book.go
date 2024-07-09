package logic

import (
	"github.com/yanshicheng/ikube-gin-starter/apps/book"
	"github.com/yanshicheng/ikube-gin-starter/apps/book/model"
	"github.com/yanshicheng/ikube-gin-starter/apps/book/service"
	"github.com/yanshicheng/ikube-gin-starter/apps/book/sql"
	"github.com/yanshicheng/ikube-gin-starter/common/types"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"github.com/yanshicheng/ikube-gin-starter/router"
	"go.uber.org/zap"
)

// 类型检查 Service 是不是等于 LogicService
var _ service.BookService = (*BookService)(nil)
var logic = &BookService{}

type BookService struct {
	l  *zap.Logger
	db *sql.BookSql
}

func (b *BookService) Create(book *model.Book) error {
	if err := b.db.Create(book); err != nil {
		return err
	}
	return nil
}

func (b *BookService) Update(id types.SearchId, book *model.Book) error {
	if err := b.db.Update(id, book); err != nil {
		return err
	}
	return nil
}

func (b *BookService) Delete(id types.SearchId) error {
	if err := b.db.Delete(id); err != nil {
		return err
	}
	return nil
}
func (b *BookService) List(query *model.GetQuery) (*types.QueryResponse, error) {
	if books, err := b.db.GetQuerySlice(query); err != nil {
		return nil, err
	} else {
		return books, nil
	}
}

func (b *BookService) Get(query types.SearchId) (*model.Book, error) {
	if books, err := b.db.GetQuery(query); err != nil {
		return nil, err
	} else {
		return books, nil
	}
}

// 只需要保证 全局对象Config和全局Logger已经加载完成
func (b *BookService) Config() {
	b.l = global.L.Named(book.AppBook).Named(name)
	b.db = sql.NewBookSql()
}

func (b *BookService) Name() string {
	return book.AppBook
}

func init() {
	// 注册
	router.RegistryLogic(logic)
}
