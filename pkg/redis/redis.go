package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// IkubeRedis 结构体用于管理 Redis 连接
type IkubeRedis struct {
	client   *redis.Client // Redis 客户端实例
	addr     string        // Redis 服务器地址
	password string        // Redis 密码
	db       int           // Redis 数据库编号
	poolSize int           // 连接池大小

}

// NewIkubeRedis 初始化一个新的 IkubeRedis 实例
func InitIkubeRedis(addr, password string, db, poolSize int) (*IkubeRedis, error) {
	ikube := &IkubeRedis{
		addr:     addr,
		password: password,
		db:       db,
		poolSize: poolSize,
	}

	// 初始化并加载 Redis 连接
	if err := ikube.load(); err != nil {
		return nil, err
	}

	return ikube, nil
}

// load 初始化 Redis 连接并设置客户端
func (ikube *IkubeRedis) load() error {
	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     ikube.addr,
		Password: ikube.password,
		DB:       ikube.db,
		PoolSize: ikube.poolSize,
	})

	// Ping Redis 服务器
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		zap.L().Error("Redis 启动异常", zap.Error(err))
		return err
	}

	ikube.client = client
	return nil
}

// GetClient 返回 Redis 客户端实例
func (ikube *IkubeRedis) GetClient() *redis.Client {
	return ikube.client
}

// Close 关闭 Redis 客户端连接
func (ikube *IkubeRedis) Close() error {
	return ikube.client.Close()
}

// Ping 检查 Redis 服务器连接是否存活
func (ikube *IkubeRedis) Ping() error {
	ctx := context.Background()
	return ikube.client.Ping(ctx).Err()
}
