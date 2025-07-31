package routers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log"
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

func (r *UserRouter) Login(ctx *gin.Context) {
	// body := &models.UserLoginRequest{}
	// if err := ctx.ShouldBindJSON(body); err != nil {
	// 	ctx.JSON(400, models.ErrorResponse{Error: "invalid request body"})
	// 	log.Panic(err)
	// 	return
	// }

	// user, err := r.UserService.Login(ctx, *body)
	// if err != nil {
	// 	ctx.JSON(400, models.ErrorResponse{Error: err.Error()})
	// 	log.Panic(err)
	// 	return
	// }
	// log.Printf("User logged in: %+v", user)

	// responseBody := models.UserLoginResponse{
	// 	ID:       user.ID,
	// 	Username: user.Username,
	// 	Email:    user.Email,
	// }
	// ctx.JSON(200, responseBody)
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

	user, err := r.UserService.Register(ctx, *body)
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
