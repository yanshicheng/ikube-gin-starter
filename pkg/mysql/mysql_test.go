package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yanshicheng/ikube-gin-starter/pkg/mysql"
)

func TestIkubeGorm(t *testing.T) {
	// 假设你的数据库连接信息如下
	dsn := "root:123456@tcp(172.16.1.61:3307)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	maxIdleConns := 10
	maxOpenConns := 100
	logToFile := false // 这里可以根据需要修改

	// 初始化一个新的 IkubeGorm 实例
	ikube, err := mysql.NewIkubeGorm(dsn, maxIdleConns, maxOpenConns, logToFile, "error")
	assert.NoError(t, err, "初始化 IkubeGorm 实例应该成功")

	// 确保数据库实例不为空
	db := ikube.GetDb()
	assert.NotNil(t, db, "获取数据库实例不应为空")

	// 测试数据库连接是否存活
	err = ikube.Ping()
	assert.NoError(t, err, "数据库连接应该存活")

	// 在这里可以添加更多的测试逻辑，例如执行数据库查询等
	// 测试一个简单的数据库操作（例如创建一个模型）
	type User struct {
		ID   uint
		Name string
	}

	// 自动迁移模式以确保表结构存在
	err = db.AutoMigrate(&User{})
	assert.NoError(t, err, "自动迁移模式时出错")

	// 创建一个新用户
	newUser := User{Name: "John Doe"}
	result := db.Create(&newUser)
	assert.NoError(t, result.Error, "创建新用户时出错")

	// 查询刚才创建的用户
	var fetchedUser User
	result = db.First(&fetchedUser, newUser.ID)
	assert.NoError(t, result.Error, "查询用户时出错")
	assert.Equal(t, newUser.Name, fetchedUser.Name, "查询到的用户名与预期不符")

	// 删除表格（示例中删除User表）
	err = db.Exec("DROP TABLE users").Error
	assert.NoError(t, err, "删除表时出错")

	// 关闭数据库连接
	err = ikube.Close()
	assert.NoError(t, err, "关闭数据库连接应该成功")
}
