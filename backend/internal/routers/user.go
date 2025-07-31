package routers

import (
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/services"
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	UserService    *services.UserService
	CityService    *services.CityService
	AddressService *services.AddressService
}

var userOnce sync.Once
var userRouter *UserRouter

func NewUserRouter() *UserRouter {
	userOnce.Do(func() {
		userRouter = &UserRouter{
			UserService:    services.NewUserService(),
			CityService:    services.NewCityService(),
			AddressService: services.NewAddressService(),
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
	accessToken, err := jwtUtils.GenerateToken(
		&models.JWTClaimsData{UserID: user.ID},
		os.Getenv(jwtUtils.DefaultEnvKey),
	)
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
	body := &models.UserRegisterRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "invalid request body"})
		log.Panic(err)
		return
	}

	// Validate password length
	passwordLen := len(body.Password)
	if passwordLen < 6 || passwordLen > 12 {
		ctx.JSON(400, models.ErrorResponse{Error: "password length must be between 6 and 12 characters"})
		log.Panic("Password length is invalid")
		return
	}

	user, err := r.UserService.Register(ctx, body)
	if err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: err.Error()})
		log.Panic(err)
		return
	}
	log.Printf("User registered: %+v", user)

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
