package services

import (
	"errors"

	"pcms/models"
	"pcms/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Email    string `json:"email"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResult struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *AuthService) Register(input RegisterInput) (*AuthResult, error) {
	// 检查用户名是否存在
	var count int64
	s.DB.Model(&models.User{}).Where("username = ?", input.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Email:    input.Email,
	}

	if err := s.DB.Create(user).Error; err != nil {
		return nil, errors.New("创建用户失败")
	}

	// 生成 token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &AuthResult{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(input LoginInput) (*AuthResult, error) {
	var user models.User
	if err := s.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, errors.New("查询用户失败")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &AuthResult{
		Token: token,
		User:  &user,
	}, nil
}

func (s *AuthService) GetUserByID(id uint64) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}
