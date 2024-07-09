package redis_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/yanshicheng/ikube-gin-starter/pkg/redis"
	"testing"
)

func TestIkubeRedis(t *testing.T) {
	// 替换为你的 Redis 服务器连接详情
	addr := "172.16.1.61:6379"
	password := "123456"
	db := 0
	poolSize := 10
	// 初始化 IkubeRedis 实例
	ikube, err := redis.InitIkubeRedis(addr, password, db, poolSize)
	if err != nil {
		t.Fatalf("初始化 IkubeRedis 失败: %v", err)
	}

	// Ping 检查 Redis 服务器是否可达
	err = ikube.Ping()
	assert.NoError(t, err, "Ping Redis 服务器失败")

	// 执行一些操作（例如设置和获取）
	key := "test_key"
	value := "test_value"

	err = ikube.GetClient().Set(context.Background(), key, value, 0).Err()
	assert.NoError(t, err, "设置键值对失败")

	result, err := ikube.GetClient().Get(context.Background(), key).Result()
	assert.NoError(t, err, "从 Redis 获取值失败")
	assert.Equal(t, value, result, "获取的值与预期不符")

	// 清理设置的键
	err = ikube.GetClient().Del(context.Background(), key).Err()
	assert.NoError(t, err, "清理 Redis 键失败")

	// 最后关闭连接
	if err := ikube.Close(); err != nil {
		t.Fatalf("关闭 IkubeRedis 失败: %v", err)
	}
}
