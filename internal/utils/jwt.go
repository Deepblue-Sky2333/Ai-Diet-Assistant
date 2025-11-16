package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken 无效的 token
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken token 已过期
	ErrExpiredToken = errors.New("token has expired")
	// ErrTokenGenerationFailed token 生成失败
	ErrTokenGenerationFailed = errors.New("token generation failed")
	// ErrPasswordChanged 密码已修改
	ErrPasswordChanged = errors.New("password has been changed")
)

// Claims JWT 声明结构
type Claims struct {
	UserID          int64  `json:"user_id"`
	Username        string `json:"username"`
	PasswordVersion int64  `json:"pwd_ver"` // 密码版本（密码修改时间戳）
	jwt.RegisteredClaims
}

// TokenPair token 对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// JWTService JWT 服务
type JWTService struct {
	secret               []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewJWTService 创建 JWT 服务实例
func NewJWTService(secret string, accessTokenHours, refreshTokenHours int) *JWTService {
	return &JWTService{
		secret:               []byte(secret),
		accessTokenDuration:  time.Duration(accessTokenHours) * time.Hour,
		refreshTokenDuration: time.Duration(refreshTokenHours) * time.Hour,
	}
}

// GenerateTokenPair 生成 token 对（Access Token 和 Refresh Token）
func (j *JWTService) GenerateTokenPair(userID int64, username string, passwordVersion int64) (*TokenPair, error) {
	// 生成 Access Token
	accessToken, err := j.generateToken(userID, username, passwordVersion, j.accessTokenDuration)
	if err != nil {
		return nil, err
	}

	// 生成 Refresh Token
	refreshToken, err := j.generateToken(userID, username, passwordVersion, j.refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.accessTokenDuration.Seconds()),
	}, nil
}

// generateToken 生成 JWT token
func (j *JWTService) generateToken(userID int64, username string, passwordVersion int64, duration time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:          userID,
		Username:        username,
		PasswordVersion: passwordVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", ErrTokenGenerationFailed
	}

	return tokenString, nil
}

// ValidateToken 验证并解析 token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshAccessToken 使用 Refresh Token 生成新的 Access Token
func (j *JWTService) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 生成新的 Access Token（保持相同的密码版本）
	accessToken, err := j.generateToken(claims.UserID, claims.Username, claims.PasswordVersion, j.accessTokenDuration)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// ValidatePasswordVersion 验证令牌中的密码版本是否匹配
func (j *JWTService) ValidatePasswordVersion(claims *Claims, currentPasswordVersion int64) error {
	if claims.PasswordVersion != currentPasswordVersion {
		return ErrPasswordChanged
	}
	return nil
}
