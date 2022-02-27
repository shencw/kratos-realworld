package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/shencw/kratos-realworld/internal/biz"
	"gorm.io/gorm"
	"time"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

var _ biz.UserRepo = (*userRepo)(nil)

type Users struct {
	gorm.Model
	ID           uint64    `gorm:"primaryKey"`
	Email        string    `gorm:"size:not_null"`
	Username     string    `gorm:"size:not_null"`
	Bio          string    `gorm:"size:not_null,1000"`
	Image        string    `gorm:"size:not_null,1000"`
	PasswordHash string    `gorm:"size:not_null,500"`
	CreatedAt    time.Time `gorm:"<-:create"`
	UpdatedAt    time.Time
}

func (Users) TableName() string { return "users" }

func (r *userRepo) Login(ctx context.Context, user biz.User) error {
	return nil
}

// CreateUser 创建成功
func (r *userRepo) CreateUser(ctx context.Context, user *biz.User) (uint64, error) {
	u := &Users{
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.Password,
	}
	result := r.data.realWorldDB.WithContext(ctx).Create(u)
	return u.ID, result.Error
}

// GetUserByUsername 根据用户名获取用户信息
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (bizUser *biz.User, err error) {
	u := new(Users)
	bizUser = new(biz.User)
	cacheKey := fmt.Sprintf("GetUserByUsername:%s", username)
	jsonCode := encoding.GetCodec("json")

	if redisResult, redisErr := r.data.redisCli.Get(ctx, cacheKey).Result(); redisErr == nil {
		if redisErr = jsonCode.Unmarshal([]byte(redisResult), bizUser); redisErr != nil && bizUser.ID > 0 {
			return
		}
	}

	if err = r.data.realWorldDB.WithContext(ctx).Where("username = ?", username).First(u).Error; err != nil {
		return
	}

	bizUser = &biz.User{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		Image:        u.Image,
		PasswordHash: u.PasswordHash,
	}

	bytes, _ := jsonCode.Marshal(bizUser)
	_, _ = r.data.redisCli.Set(ctx, cacheKey, bytes, time.Hour).Result()

	return
}
