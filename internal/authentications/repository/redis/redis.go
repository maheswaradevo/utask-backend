package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/maheswaradevo/utask-backend/internal/authentications"
	"github.com/maheswaradevo/utask-backend/internal/models"
	"github.com/maheswaradevo/utask-backend/pkg/common/helpers/json"
)

type googleOauthRedisRepository struct {
	rc *redis.Client
}

func NewGoogleOauthRedisRepository(rc *redis.Client) authentications.GoogleOauthRedisRepository {
	return &googleOauthRedisRepository{rc: rc}
}

func (r googleOauthRedisRepository) Save(ctx context.Context, key string, data *models.OauthRedisData) (*models.OauthRedisData, error) {
	jsonStr, errJson := json.Stringify(data)
	if errJson != nil {
		return nil, errJson
	}
	err := r.rc.Set(ctx, key, jsonStr, 1*time.Hour).Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r googleOauthRedisRepository) FindByKey(ctx context.Context, key string) (*models.OauthRedisData, error) {
	result, err := r.rc.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var oauthResponse = new(models.OauthRedisData)
	if err := json.FromJSON([]byte(result), oauthResponse); err != nil {
		return nil, err
	}
	return oauthResponse, nil
}

func (r googleOauthRedisRepository) GetExpirationTime(ctx context.Context, key string) float64 {
	cmd := r.rc.TTL(ctx, key)
	val := cmd.Val().Seconds()
	return val
}
