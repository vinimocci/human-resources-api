package server

import (
	"log"
	"fmt"
	"time"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/pelletier/go-toml"
	"human-resources-api/entities/user"
	
	_ 	 "github.com/go-sql-driver/mysql"
	auth "human-resources-api/entities/auth"
)

func Routes() {

	config, err := toml.LoadFile("config.toml")
    if err != nil {
        fmt.Println("Error loading config file:", err)
        return
    }

	APIHost := config.Get("system.host").(string)

	db, err := sql.Open("mysql", "root:An1m@Rum@2000@tcp(localhost:3306)/users?parseTime=true")
	if err != nil {
		log.Fatalf("db: failed to connect./n%s", err)
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)

	routes := gin.Default()
	routes.SetTrustedProxies([]string{})

	routes.Use(cors.New(cors.Config{
		AllowOrigins:     []string{APIHost},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowMethods:     []string{"GET", "POST"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.GET("/signin", func (c*gin.Context)  {
		result, err := auth.Service.SignIn(authService,c)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"message": result,
		})
	})

	routes.GET("/getuserinfo", func(c *gin.Context) {
		result, err := user.Service.GetUserInfoByID(userService, c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error() ,
			})
		}else {
			c.JSON(http.StatusOK, gin.H{
				"message":  result,
			})
		}
	})

	routes.POST("/postUser", func(c *gin.Context) {
		result, err := user.Service.PostUser(userService, c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error() ,
			})
		}else {
			c.JSON(http.StatusOK, gin.H{
				"message":  result,
			})
		}
	})

	routes.Run(":5023")
}
