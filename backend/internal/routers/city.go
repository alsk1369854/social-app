package routers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

type CityRouter struct {
	CityService *services.CityService
}

var cityOnce sync.Once
var cityRouter *CityRouter

func NewCityRouter() *CityRouter {
	cityOnce.Do(func() {
		cityRouter = &CityRouter{
			CityService: services.NewCityService(),
		}
	})
	return cityRouter
}

func (r *CityRouter) Bind(_router *gin.RouterGroup) {
	router := _router.Group("/city")
	// GET
	{
		router.GET("/all", r.GetAll)
	}
}

// @title City API
// @Summary Get all cities
// @Tags City
// @Accept application/json
// @Produce application/json
// @Success 200 {object} models.CityGetAllResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/city/all [get]
func (r *CityRouter) GetAll(ctx *gin.Context) {
	cities, err := r.CityService.GetAll(ctx)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{Error: "Failed to retrieve cities"})
		log.Panic(err)
		return
	}

	responseBody := make(models.CityGetAllResponse, len(cities))
	for i, city := range cities {
		responseBody[i] = models.CityGetAllResponseItem{
			ID:   city.ID,
			Name: city.Name,
		}
	}
	ctx.JSON(200, responseBody)
}
