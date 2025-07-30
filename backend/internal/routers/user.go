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
	}
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
		ctx.JSON(400, models.ErrorResponse{Error: "Invalid request body"})
		log.Panic(err)
		return
	}

	user, err := r.UserService.Register(ctx, *body)
	if err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: err.Error()})
		log.Panic(err)
		return
	}

	responseBody := models.UserRegisterResponse{
		Username: user.Username,
		Email:    user.Email,
		Age:      user.Age,
		Address:  nil,
	}
	if user.Address != nil {
		responseBody.Address = &models.UserRegisterResponseAddress{
			CityID: user.Address.CityID.String(),
			Street: user.Address.Street,
		}
	}
	ctx.JSON(200, responseBody)
}
