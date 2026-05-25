package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"backend/models"
	"backend/pkg/apperr"
	"backend/pkg/redisx"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

type LoginResult struct {
	Token     string       `json:"token"`
	ExpiresIn int64        `json:"expires_in"`
	User      *models.User `json:"user"`
}

type jwtClaims struct {
	UserId   string `json:"sub"`
	TenantId string `json:"tenant_id"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func (s *AuthService) Login(email, password string) (*LoginResult, error) {
	o := orm.NewOrm()
	var user models.User
	err := o.QueryTable(new(models.User)).Filter("email", email).Filter("is_active", 1).One(&user)
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, apperr.New("UNAUTHORIZED", "邮箱或密码错误", apperr.ErrUnauthorized)
		}
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, apperr.New("UNAUTHORIZED", "邮箱或密码错误", apperr.ErrUnauthorized)
	}
	token, exp, err := s.issueToken(&user)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return &LoginResult{Token: token, ExpiresIn: exp, User: &user}, nil
}

func (s *AuthService) issueToken(user *models.User) (string, int64, error) {
	secret, _ := beego.AppConfig.String("jwt_secret")
	if secret == "" {
		secret = "change-me"
	}
	hours, _ := beego.AppConfig.Int("jwt_expire_hours")
	if hours <= 0 {
		hours = 24
	}
	exp := time.Now().Add(time.Duration(hours) * time.Hour)
	claims := jwtClaims{
		UserId:   user.Id,
		TenantId: user.TenantId,
		Role:     user.Role,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return signed, int64(hours * 3600), nil
}

func ParseToken(tokenStr string) (*jwtClaims, error) {
	secret, _ := beego.AppConfig.String("jwt_secret")
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, apperr.ErrUnauthorized
	}
	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, apperr.ErrUnauthorized
	}
	return claims, nil
}

func (s *AuthService) GetUser(userId string) (*models.User, error) {
	o := orm.NewOrm()
	var user models.User
	if err := o.QueryTable(new(models.User)).Filter("id", userId).One(&user); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return &user, nil
}

type UserProfile struct {
	models.User
	DiscardUsedToday int `json:"discard_used_today"`
	MaxDiscardPerDay int `json:"max_discard_per_day"`
}

func (s *AuthService) GetUserProfile(userId, tenantId string) (*UserProfile, error) {
	user, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}
	max := 10
	o := orm.NewOrm()
	var t models.Tenant
	if o.QueryTable(new(models.Tenant)).Filter("id", tenantId).One(&t) == nil && t.Settings != "" {
		var cfg struct {
			MaxDiscardPerDay int `json:"max_discard_per_day"`
		}
		_ = json.Unmarshal([]byte(t.Settings), &cfg)
		if cfg.MaxDiscardPerDay > 0 {
			max = cfg.MaxDiscardPerDay
		}
	}
	used, _ := redisx.GetDiscardCount(context.Background(), tenantId, userId)
	return &UserProfile{User: *user, DiscardUsedToday: int(used), MaxDiscardPerDay: max}, nil
}
