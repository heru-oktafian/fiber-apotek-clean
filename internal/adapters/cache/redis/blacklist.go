package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Blacklist struct {
	Client *goredis.Client
}

func (b Blacklist) Blacklist(ctx context.Context, token string, ttl time.Duration) error {
	if ttl <= 0 {
		return nil
	}
	return b.Client.Set(ctx, "blacklist:"+token, "blacklisted", ttl).Err()
}

func (b Blacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	res, err := b.Client.Exists(ctx, "blacklist:"+token).Result()
	return res > 0, err
}
