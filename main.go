package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type IngredientModel struct {
	gorm.Model
	Name     string `json:"name"`
	BaseUnit string `json:"baseUnit"`
	Quantity uint64 `json:"quantity"`
	Kcal     uint64 `json:"kcal"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("kc.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&IngredientModel{}); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/ingredient", func(c *gin.Context) {
		ctx := c.Request.Context()

		ingredients, err := gorm.G[IngredientModel](db).Find(ctx)
		if err != nil {
			c.String(http.StatusInternalServerError, "%v", err)
			return
		}

		c.JSON(http.StatusOK, ingredients)
	})

	router.POST("/ingredient", func(c *gin.Context) {
		ctx := c.Request.Context()

		var cmd struct {
			Name     string `json:"name"`
			BaseUnit string `json:"baseUnit"`
			Quantity uint64 `json:"quantity"`
			Kcal     uint64 `json:"kcal"`
		}
		if err := c.BindJSON(&cmd); err != nil {
			c.String(http.StatusInternalServerError, "%v", err)
			return
		}

		ingredient := IngredientModel{
			Name:     cmd.Name,
			BaseUnit: cmd.BaseUnit,
			Quantity: cmd.Quantity,
			Kcal:     cmd.Kcal,
		}

		if err := gorm.G[IngredientModel](db).Create(ctx, &ingredient); err != nil {
			c.String(http.StatusInternalServerError, "%v", err)
			return
		}

		c.JSON(http.StatusOK, ingredient)
	})

	router.Run()
}
