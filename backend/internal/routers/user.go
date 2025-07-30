package routers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	CityService *services.CityService
}

var userOnce sync.Once
var userRouter *UserRouter

func NewUserRouter() *UserRouter {
	userOnce.Do(func() {
		userRouter = &UserRouter{
			CityService: services.NewCityService(),
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
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/user/register [post]
func (r *UserRouter) Register(ctx *gin.Context) {
	body := &models.UserRegisterRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "Invalid request body"})
		log.Panic(err)
		return
	}
	if body.Address != nil {
		cityID := body.Address.CityID
		_, err := r.CityService.GetByID(ctx, cityID)
		if err != nil {
			ctx.JSON(400, models.ErrorResponse{Error: "City not found"})
			log.Panic(err)
			return
		}
	}
	ctx.JSON(200, models.SuccessResponse{Success: true})
}
