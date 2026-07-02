package services

import (
	"testing"

	"pcms/store"
	"pcms/utils"
)

func TestAuthService_Register(t *testing.T) {
	utils.InitJWT("test-secret")
	s := NewAuthService(store.NewMemStore())

	t.Run("注册成功", func(t *testing.T) {
		result, err := s.Register(RegisterInput{
			Username: "testuser",
			Password: "pass1234",
			Email:    "a@b.com",
		})
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
		if result.Token == "" {
			t.Fatal("token 不应为空")
		}
		if result.User.Username != "testuser" {
			t.Fatalf("期望 testuser, 得到 %s", result.User.Username)
		}
		if result.User.ID != 1 {
			t.Fatalf("期望 ID=1, 得到 %d", result.User.ID)
		}
	})

	t.Run("用户名已存在", func(t *testing.T) {
		_, err := s.Register(RegisterInput{
			Username: "testuser",
			Password: "pass1234",
		})
		if err == nil {
			t.Fatal("期望报错")
		}
		if err.Error() != "用户名已存在" {
			t.Fatalf("期望 '用户名已存在', 得到 '%s'", err.Error())
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	utils.InitJWT("test-secret")
	s := NewAuthService(store.NewMemStore())

	s.Register(RegisterInput{Username: "user1", Password: "secret456"})

	t.Run("登录成功", func(t *testing.T) {
		result, err := s.Login(LoginInput{
			Username: "user1",
			Password: "secret456",
		})
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
		if result.Token == "" {
			t.Fatal("token 不应为空")
		}
	})

	t.Run("密码错误", func(t *testing.T) {
		_, err := s.Login(LoginInput{
			Username: "user1",
			Password: "wrong",
		})
		if err == nil {
			t.Fatal("期望报错")
		}
	})

	t.Run("用户不存在", func(t *testing.T) {
		_, err := s.Login(LoginInput{
			Username: "noone",
			Password: "pass",
		})
		if err == nil {
			t.Fatal("期望报错")
		}
	})
}

func TestAuthService_GetUserByID(t *testing.T) {
	utils.InitJWT("test-secret")
	s := NewAuthService(store.NewMemStore())
	s.Register(RegisterInput{Username: "u1", Password: "p1"})

	t.Run("存在的用户", func(t *testing.T) {
		u, err := s.GetUserByID(1)
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
		if u.Username != "u1" {
			t.Fatalf("期望 u1, 得到 %s", u.Username)
		}
	})

	t.Run("不存在的用户", func(t *testing.T) {
		_, err := s.GetUserByID(999)
		if err == nil {
			t.Fatal("期望报错")
		}
	})
}
