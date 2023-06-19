package authentications

import (
	"context"

	"github.com/maheswaradevo/utask-backend/internal/models"
)

type GoogleOauthRedisRepository interface {
	Save(ctx context.Context, key string, data *models.OauthRedisData) (*models.OauthRedisData, error)
	FindByKey(ctx context.Context, key string) (*models.OauthRedisData, error)
	GetExpirationTime(ctx context.Context, key string) float64
	DeleteKey(ctx context.Context, key string) string
}
