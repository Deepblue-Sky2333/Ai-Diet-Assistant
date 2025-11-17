package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/config"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/database"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/model"
	"github.com/Deepblue-Sky2333/Ai-Diet-Assistant/internal/repository"
)

var (
	username   = flag.String("username", "", "用户名 (3-50个字符，仅字母和数字)")
	password   = flag.String("password", "", "密码 (至少8个字符)")
	email      = flag.String("email", "", "电子邮件 (可选)")
	role       = flag.String("role", "", "用户角色: admin 或 user (可选，默认根据是否为第一个用户自动判断)")
	configPath = flag.String("config", "", "配置文件路径 (默认: ./configs/config.yaml)")
)

func main() {
	flag.Parse()

	// 验证必填参数
	if *username == "" || *password == "" {
		fmt.Println("错误: 用户名和密码是必填项")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  create-user -username <用户名> -password <密码> [-email <邮箱>] [-role admin|user] [-config <配置文件路径>]")
		fmt.Println()
		fmt.Println("参数:")
		fmt.Println("  -username  用户名 (3-50个字符，仅字母和数字)")
		fmt.Println("  -password  密码 (至少8个字符)")
		fmt.Println("  -email     电子邮件 (可选)")
		fmt.Println("  -role      用户角色: admin 或 user (可选，默认自动判断)")
		fmt.Println("  -config    配置文件路径 (可选)")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  # 创建第一个用户（自动成为管理员）")
		fmt.Println("  create-user -username admin -password adminpass123 -email admin@example.com")
		fmt.Println()
		fmt.Println("  # 创建普通用户")
		fmt.Println("  create-user -username testuser -password userpass123")
		fmt.Println()
		fmt.Println("  # 显式指定角色创建管理员")
		fmt.Println("  create-user -username admin2 -password admin2pass -role admin")
		os.Exit(1)
	}

	// 验证用户名格式
	if err := validateUsername(*username); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 验证密码格式
	if err := validatePassword(*password); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 验证邮箱格式（如果提供）
	if *email != "" {
		if err := validateEmail(*email); err != nil {
			fmt.Printf("错误: %v\n", err)
			os.Exit(1)
		}
	}

	// 验证角色参数
	if *role != "" && *role != model.RoleAdmin && *role != model.RoleUser {
		fmt.Printf("错误: 角色必须是 'admin' 或 'user'，当前值: %s\n", *role)
		os.Exit(1)
	}

	// 加载配置
	fmt.Println("正在加载配置...")
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 连接数据库
	fmt.Println("正在连接数据库...")
	err = database.Init(&cfg.Database)
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// 创建用户仓储
	db := database.GetDB()
	userRepo := repository.NewUserRepository(db)

	// 确定用户角色
	userRole := *role
	if userRole == "" {
		// 如果未指定角色，检查是否为第一个用户
		count, err := userRepo.GetUserCount(context.Background())
		if err != nil {
			fmt.Printf("检查用户数量失败: %v\n", err)
			os.Exit(1)
		}

		if count == 0 {
			userRole = model.RoleAdmin
			fmt.Println("✓ 检测到这是第一个用户，将设置为管理员")
		} else {
			userRole = model.RoleUser
		}
	}

	// 检查用户名是否已存在
	exists, err := userRepo.CheckUsernameExists(context.Background(), *username)
	if err != nil {
		fmt.Printf("检查用户名失败: %v\n", err)
		os.Exit(1)
	}

	if exists {
		fmt.Printf("错误: 用户名 '%s' 已存在\n", *username)
		os.Exit(1)
	}

	// 创建用户
	fmt.Println("正在创建用户...")
	user := &model.User{
		Username: *username,
		Email:    *email,
		Role:     userRole,
	}

	err = userRepo.CreateUser(context.Background(), user, *password)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		os.Exit(1)
	}

	// 显示创建结果
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("✓ 用户创建成功！")
	fmt.Println("========================================")
	fmt.Printf("ID:         %d\n", user.ID)
	fmt.Printf("用户名:     %s\n", user.Username)
	fmt.Printf("角色:       %s %s\n", user.Role, getRoleEmoji(user.Role))
	if user.Email != "" {
		fmt.Printf("邮箱:       %s\n", user.Email)
	}
	fmt.Printf("创建时间:   %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("========================================")
}

// validateUsername 验证用户名格式
func validateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("用户名长度不能少于3个字符")
	}
	if len(username) > 50 {
		return fmt.Errorf("用户名长度不能超过50个字符")
	}

	// 仅允许字母和数字
	matched, err := regexp.MatchString("^[a-zA-Z0-9]+$", username)
	if err != nil {
		return fmt.Errorf("验证用户名格式失败: %v", err)
	}
	if !matched {
		return fmt.Errorf("用户名只能包含字母和数字")
	}

	return nil
}

// validatePassword 验证密码格式
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("密码长度不能少于8个字符")
	}
	if len(password) > 128 {
		return fmt.Errorf("密码长度不能超过128个字符")
	}
	return nil
}

// validateEmail 验证邮箱格式
func validateEmail(email string) error {
	if len(email) > 100 {
		return fmt.Errorf("邮箱长度不能超过100个字符")
	}

	// 简单的邮箱格式验证
	matched, err := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	if err != nil {
		return fmt.Errorf("验证邮箱格式失败: %v", err)
	}
	if !matched {
		return fmt.Errorf("邮箱格式无效")
	}

	return nil
}

// getRoleEmoji 获取角色对应的 emoji
func getRoleEmoji(role string) string {
	if role == model.RoleAdmin {
		return "👑"
	}
	return "👤"
}
