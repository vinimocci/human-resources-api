package server

import (
	"os"
	"log"
	"time"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"
	"web-solutions-api/utils"
	"github.com/gin-contrib/cors"
	"web-solutions-api/entities/user"
	"web-solutions-api/entities/notifications"

	_ 	 "github.com/go-sql-driver/mysql"
	auth "web-solutions-api/entities/auth"
)

func Routes() {
	config, tomlErr := utils.GetCurrentEnvironment(os.Getenv("ENVIRONMENT"))
    if tomlErr != nil {
		panic("error loading config file")
    }

	APIHost := config.Get("system.host").(string)
	APIPort := config.Get("system.apiPort").(string)
	databaseHost := config.Get("database.host").(string)
	baseTable := config.Get("database.baseTable").(string)
	databasePass := config.Get("database.password").(string)

	db, err := sql.Open("mysql", "root:"+databasePass+"@tcp("+databaseHost+")/"+baseTable+"?parseTime=true")
	if err != nil {
		log.Fatalf("db: failed to connect./n%s", err)
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)

	notificationRepo := notifications.NewRepository(db)
	notificationService := notifications.NewService(notificationRepo)

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

	routes.POST("/signin", func (c*gin.Context)  {
		result, err := auth.Service.SignIn(authService,c)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		}else{
			c.JSON(http.StatusOK, gin.H{
				"message": result,
			})
		}
	})

	routes.GET("/getuserinfo/:id", func(c *gin.Context) {
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

	routes.GET("/getnotifications", func(c *gin.Context) {
		 notifications.Service.GetNotificationMessages(notificationService, c)
	})

	routes.POST("/changeNotification", func (c *gin.Context)  {
		producer, prdErr := utils.CreateKafkaProducer(config.Get("kafka.host").(string))
		if prdErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": prdErr,
			})
		}

		result, err := notifications.Service.RegisterNotificationMessage(notificationService, c, producer)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		}else{
			c.JSON(http.StatusOK, gin.H{
				"message": result,
			})
		}
	})

	//testing purposes
	routes.GET("/ping", func(c *gin.Context) {		
		c.JSON(http.StatusOK, gin.H{
			"message":  "pong! =D",
		})
	})

	routes.Run(APIPort)
}