package redisx

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/redis/go-redis/v9"
)

const (
	QueueParseCSV    = "imagedeal:queue:parse_csv"
	QueueDownload    = "imagedeal:queue:download"
	QueueExport      = "imagedeal:queue:export"
	PrefixDiscardDay = "imagedeal:discard:"
)

var (
	client *redis.Client
	once   sync.Once
)

func Client() *redis.Client {
	once.Do(func() {
		addr, _ := beego.AppConfig.String("redis_addr")
		if addr == "" {
			addr = "localhost:6379"
		}
		client = redis.NewClient(&redis.Options{Addr: addr})
	})
	return client
}

func Ping(ctx context.Context) error {
	return Client().Ping(ctx).Err()
}

type TaskMessage struct {
	Type     string          `json:"type"`
	TenantId string          `json:"tenant_id"`
	Payload  json.RawMessage `json:"payload"`
}

func Enqueue(ctx context.Context, queue string, msg TaskMessage) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return Client().LPush(ctx, queue, b).Err()
}

func Dequeue(ctx context.Context, queue string, timeout time.Duration) (*TaskMessage, error) {
	res, err := Client().BRPop(ctx, timeout, queue).Result()
	if err != nil {
		return nil, err
	}
	if len(res) < 2 {
		return nil, redis.Nil
	}
	var msg TaskMessage
	if err := json.Unmarshal([]byte(res[1]), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func discardKey(tenantId, userId string) string {
	return PrefixDiscardDay + tenantId + ":" + userId + ":" + time.Now().Format("20060102")
}

func IncrDiscardCount(ctx context.Context, tenantId, userId string) (int64, error) {
	key := discardKey(tenantId, userId)
	n, err := Client().Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if n == 1 {
		Client().Expire(ctx, key, 48*time.Hour)
	}
	return n, nil
}

func GetDiscardCount(ctx context.Context, tenantId, userId string) (int64, error) {
	n, err := Client().Get(ctx, discardKey(tenantId, userId)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return n, err
}
