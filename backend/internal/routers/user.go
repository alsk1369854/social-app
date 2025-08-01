package routers

import (
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/services"
	"errors"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRouter struct {
	UserService    *services.UserService
	CityService    *services.CityService
	AddressService *services.AddressService

	CryptoUtils *pkg.CryptoUtils
}

var userOnce sync.Once
var userRouter *UserRouter

func NewUserRouter() *UserRouter {
	userOnce.Do(func() {
		userRouter = &UserRouter{
			UserService:    services.NewUserService(),
			CityService:    services.NewCityService(),
			AddressService: services.NewAddressService(),

			CryptoUtils:    pkg.NewCryptoUtils(),
		}
	})
	return userRouter
}

func (r *UserRouter) Bind(_router *gin.RouterGroup) {
	router := _router.Group("/user")
	// POST
	{
		router.POST("/register", r.Register)
		router.POST("/login", r.Login)
	}
}

// @title User API
// @Summary Login a user
// @Tags User
// @Accept application/json
// @Produce application/json
// @Param user body models.UserLoginRequest true "User login request"
// @Success 200 {object} models.UserLoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/user/login [post]
func (r *UserRouter) Login(ctx *gin.Context) {
	body := &models.UserLoginRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "invalid request body"})
		log.Panic(err)
		return
	}

	// 檢查使用者是否存在
	user, err := r.UserService.GetByEmail(ctx, body.Email)
	if err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "email not found"})
		log.Panic(err)
		return
	}

	// 驗證密碼
	cryptoUtils := pkg.NewCryptoUtils()
	isValid := cryptoUtils.VerifyPasswordHash(user.HashedPassword, &pkg.CryptoUtilsPasswordHashInput{
		Email:    user.Email,
		Username: user.Username,
		Password: body.Password,
	})
	if !isValid {
		ctx.JSON(400, models.ErrorResponse{Error: "incorrect email or password"})
		return
	}

	// 生成 JWT Token
	jwtUtils := pkg.NewJWTUtils()
	accessToken, err := jwtUtils.GenerateToken(&models.JWTClaimsData{UserID: user.ID}, nil)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{Error: "failed to generate access token"})
		log.Panic(err)
		return
	}

	// 構建回應
	response := &models.UserLoginResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		AccessToken: accessToken,
	}
	ctx.JSON(200, response)
}

// @title User API
// @Summary Register a new user
// @Tags User
// @Accept application/json
// @Produce application/json
// @Param user body models.UserRegisterRequest true "User registration request"
// @Success 200 {object} models.UserRegisterResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/user/register [post]
func (r *UserRouter) Register(ctx *gin.Context) {
	reqBody := &models.UserRegisterRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Validate password length
	passwordLen := len(reqBody.Password)
	if passwordLen < 6 || passwordLen > 12 {
		ctx.JSON(400, models.ErrorResponse{Error: "password length must be between 6 and 12 characters"})
		return
	}

	// 檢查 email 是否存在
	if _, err := r.UserService.GetByEmail(ctx, reqBody.Email); err == nil {
		ctx.JSON(400, models.ErrorResponse{Error: "email already exists"})
		return
	}

	// 檢查是否有填寫 address
	var addressID *uuid.UUID
	if reqBody.Address != nil {
		if _, err := r.CityService.GetByID(ctx, reqBody.Address.CityID); err != nil {
			ctx.JSON(400, models.ErrorResponse{Error: "city not found"})
			return
		}
		addressSlice, err := r.AddressService.Create(ctx, []models.AddressBase{{
			CityID: reqBody.Address.CityID,
			Street: reqBody.Address.Street,
		}})
		if err != nil {
			ctx.JSON(400, models.ErrorResponse{Error: "invalid address data"})
			return
		}
		}
		addressID = &(addressSlice[0].ID)
	}
	
	// 建立 User 資料
	hashedPassword := r.CryptoUtils.GeneratePasswordHash(&pkg.CryptoUtilsPasswordHashInput{
		Email:    reqBody.Email,
		Username: reqBody.Username,
		Password: reqBody.Password,
	})
	userBase := models.UserBase{
		Username:       reqBody.Username,
		Email:          reqBody.Email,
		HashedPassword: hashedPassword,
		Age:            reqBody.Age,
		AddressID:      addressID,
	}
	users, err := r.UserRepository.Create(ctx, []models.UserBase{userBase})
	if err != nil {
		return nil, err
	}

	// 載入地址資訊
	result := &users[0]
	if addressID != nil {
		address, err := s.AddressService.GetByID(ctx, *addressID)
		if err != nil {
			return nil, err
		}
		result.Address = address
	}
	return result, nil

	// user, err := r.UserService.Register(ctx, body)
	// if err != nil {
	// 	ctx.JSON(400, models.ErrorResponse{Error: err.Error()})
	// 	log.Panic(err)
	// 	return
	// }

	responseBody := models.UserRegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Age:      user.Age,
		Address:  nil,
	}
	if user.Address != nil {
		responseBody.Address = &models.UserRegisterResponseAddress{
			CityID: user.Address.CityID,
			Street: user.Address.Street,
		}
	}
	ctx.JSON(200, responseBody)
}
