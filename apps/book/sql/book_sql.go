package sql

import (
	"fmt"
	"github.com/yanshicheng/ikube-gin-starter/apps/book"
	"github.com/yanshicheng/ikube-gin-starter/apps/book/model"
	"github.com/yanshicheng/ikube-gin-starter/common/sql"
	"github.com/yanshicheng/ikube-gin-starter/common/types"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BookSql struct {
	l  *zap.Logger
	db *gorm.DB
}

// 只需要保证 全局对象Config和全局Logger已经加载完成
func (b *BookSql) Config() {
	b.l = global.L.Named(book.AppBook).Named(name)
}

func (b *BookSql) Name() string {
	return book.AppBook
}

func (b *BookSql) Create(book *model.Book) error {
	if err := b.db.Create(book).Error; err != nil {
		b.l.Error("create book error", zap.Error(err))
		return fmt.Errorf("创建数据失败")
	}
	return nil
}

func (b *BookSql) Update(id types.SearchId, book *model.Book) error {

	// 修改数据
	if err := b.db.Model(&model.Book{}).Where("id = ?", id.Id).Updates(book).Error; err != nil {
		b.l.Error("update book error", zap.Error(err))
		return fmt.Errorf("修改数据失败")
	}

	return nil
}

func (b *BookSql) Delete(id types.SearchId) error {
	result := b.db.Model(&model.Book{}).Where("id = ?", id.Id).Delete(&model.Book{})
	if result.Error != nil {
		b.l.Error("delete book error", zap.Error(result.Error))
		return fmt.Errorf("删除数据失败")
	}

	if result.RowsAffected == 0 {
		b.l.Error(fmt.Sprintf("book 未找到匹配的记录，id = %d", id.Id))
		// 可以选择返回特定的错误
		return fmt.Errorf("未找到匹配的记录")
	}

	return nil
}

func (b *BookSql) GetQuery(id types.SearchId) (*model.Book, error) {
	var m model.Book
	if err := b.db.Model(&model.Book{}).Where("id = ?", id.Id).First(&m).Error; err != nil {
		b.l.Error("get m error", zap.Error(err))
		return nil, fmt.Errorf("查询数据失败")
	}
	return &m, nil
}

func (b *BookSql) GetQuerySlice(query *model.GetQuery) (*types.QueryResponse, error) {
	var books []model.Book
	db := b.db.Model(&model.Book{})
	// 设置排序
	db = db.Order(fmt.Sprintf("%s %s", "ID", query.Sort))
	// 模糊查询
	if query.Title != "" {
		db = db.Where("title like ?", query.Title+"%")
	}
	// 打印 sql 语句
	queryRes, err := sql.GetQueryResponse(db, query, books)
	if err != nil {
		b.l.Error("get book error", zap.Error(err))
		return nil, fmt.Errorf("查询数据失败")
	}

	return queryRes, nil
}

func NewBookSql() *BookSql {
	return &BookSql{
		l:  global.L.Named(book.AppBook).Named(name),
		db: global.DB.GetDb(),
	}
}
