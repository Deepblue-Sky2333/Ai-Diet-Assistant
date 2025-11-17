package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/utils"
)

var (
	// ErrAccountLocked 账户被锁定
	ErrAccountLocked = errors.New("account locked due to too many failed login attempts")
	// ErrInvalidCredentials 凭证无效
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrUsernameExists 用户名已存在
	ErrUsernameExists = errors.New("username already exists")
	// ErrRegistrationDisabled 注册已关闭
	ErrRegistrationDisabled = errors.New("registration is currently disabled")
)

// AuthService 认证服务接口
type AuthService interface {
	Login(ctx context.Context, username, password, ipAddress string) (*utils.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	ValidateToken(ctx context.Context, token string) (*utils.Claims, error)
	Logout(ctx context.Context, token string) error
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	Register(ctx context.Context, username, password, email string) (*model.User, error)
}

// authService 认证服务实现
type authService struct {
	userRepo           repository.UserRepository
	loginAttemptRepo   repository.LoginAttemptRepository
	tokenBlacklistRepo repository.TokenBlacklistRepository
	settingsService    SettingsService
	jwtService         *utils.JWTService
	maxLoginAttempts   int
	lockoutDuration    time.Duration
}

// NewAuthService 创建认证服务实例
func NewAuthService(
	userRepo repository.UserRepository,
	loginAttemptRepo repository.LoginAttemptRepository,
	tokenBlacklistRepo repository.TokenBlacklistRepository,
	settingsService SettingsService,
	jwtService *utils.JWTService,
	maxLoginAttempts int,
	lockoutDuration time.Duration,
) AuthService {
	return &authService{
		userRepo:           userRepo,
		loginAttemptRepo:   loginAttemptRepo,
		tokenBlacklistRepo: tokenBlacklistRepo,
		settingsService:    settingsService,
		jwtService:         jwtService,
		maxLoginAttempts:   maxLoginAttempts,
		lockoutDuration:    lockoutDuration,
	}
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, username, password, ipAddress string) (*utils.TokenPair, error) {
	// 检查登录限流
	failedAttempts, err := s.loginAttemptRepo.GetRecentFailedAttempts(ctx, username, s.lockoutDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to check login attempts: %w", err)
	}

	if failedAttempts >= s.maxLoginAttempts {
		// 记录失败的登录尝试
		_ = s.loginAttemptRepo.RecordLoginAttempt(ctx, &model.LoginAttempt{
			Username:    username,
			IPAddress:   ipAddress,
			Success:     false,
			AttemptedAt: time.Now(),
		})
		return nil, ErrAccountLocked
	}

	// 获取用户
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		// 记录失败的登录尝试
		_ = s.loginAttemptRepo.RecordLoginAttempt(ctx, &model.LoginAttempt{
			Username:    username,
			IPAddress:   ipAddress,
			Success:     false,
			AttemptedAt: time.Now(),
		})

		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 验证密码
	if err := utils.VerifyPassword(user.PasswordHash, password); err != nil {
		// 记录失败的登录尝试
		_ = s.loginAttemptRepo.RecordLoginAttempt(ctx, &model.LoginAttempt{
			Username:    username,
			IPAddress:   ipAddress,
			Success:     false,
			AttemptedAt: time.Now(),
		})
		return nil, ErrInvalidCredentials
	}

	// 记录成功的登录尝试
	_ = s.loginAttemptRepo.RecordLoginAttempt(ctx, &model.LoginAttempt{
		Username:    username,
		IPAddress:   ipAddress,
		Success:     true,
		AttemptedAt: time.Now(),
	})

	// 生成 token（包含密码版本）
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Username, user.PasswordVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenPair, nil
}

// RefreshToken 刷新访问令牌
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// 验证 refresh token
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// 验证用户是否存在并获取当前密码版本
	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// 验证密码版本是否匹配
	if err := s.jwtService.ValidatePasswordVersion(claims, user.PasswordVersion); err != nil {
		return "", err
	}

	// 生成新的 access token
	accessToken, err := s.jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	return accessToken, nil
}

// ValidateToken 验证令牌
func (s *authService) ValidateToken(ctx context.Context, token string) (*utils.Claims, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// 验证用户是否存在并获取当前密码版本
	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 验证密码版本是否匹配
	if err := s.jwtService.ValidatePasswordVersion(claims, user.PasswordVersion); err != nil {
		return nil, err
	}

	return claims, nil
}

// Logout 用户登出
func (s *authService) Logout(ctx context.Context, token string) error {
	// 验证令牌并获取声明
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		// 即使令牌无效，也返回成功（幂等性）
		return nil
	}

	// 计算令牌剩余有效期
	expiryDuration := time.Until(claims.ExpiresAt.Time)
	if expiryDuration <= 0 {
		// 令牌已过期，无需加入黑名单
		return nil
	}

	// 将令牌添加到黑名单
	if err := s.tokenBlacklistRepo.Add(ctx, token, expiryDuration); err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}

	return nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrInvalidCredentials
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 验证旧密码
	if err := utils.VerifyPassword(user.PasswordHash, oldPassword); err != nil {
		return repository.ErrInvalidPassword
	}

	// 哈希新密码
	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码版本为当前时间戳
	newPasswordVersion := time.Now().Unix()

	// 更新密码和密码版本
	if err := s.userRepo.UpdatePasswordWithVersion(ctx, userID, newPasswordHash, newPasswordVersion); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// Register 用户注册
func (s *authService) Register(ctx context.Context, username, password, email string) (*model.User, error) {
	// 1. 检查注册开关状态
	registrationEnabled, err := s.settingsService.IsRegistrationEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check registration status: %w", err)
	}

	if !registrationEnabled {
		return nil, ErrRegistrationDisabled
	}

	// 2. 检查用户名唯一性（不区分大小写）
	exists, err := s.userRepo.CheckUsernameExists(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username exists: %w", err)
	}

	if exists {
		return nil, ErrUsernameExists
	}

	// 3. 确定用户角色（第一个用户为管理员）
	userCount, err := s.userRepo.GetUserCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user count: %w", err)
	}

	role := model.RoleUser // 默认为普通用户
	if userCount == 0 {
		role = model.RoleAdmin // 第一个用户为管理员
	}

	// 4. 创建用户
	user := &model.User{
		Username: username,
		Email:    email,
		Role:     role,
	}

	err = s.userRepo.CreateUser(ctx, user, password)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, ErrUsernameExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 5. 返回创建的用户信息
	return user, nil
}
