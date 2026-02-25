package main

import (
	httpConfig "doan/cmd/http/config"
	"doan/cmd/http/controllers/class"
	"doan/cmd/http/controllers/room"
	"doan/cmd/http/controllers/student"
	"doan/cmd/http/controllers/teacher"
	"doan/cmd/http/controllers/user"
	_ "doan/cmd/http/docs"
	"doan/cmd/http/middleware"
	"doan/pkg/config"
	"doan/pkg/constants"
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	Name                string
	Version             string
	ConfigFilePath      string
	ConfigFile          string
	router              *gin.Engine
	restConfig          httpConfig.RestServer
	userControllerV1    user.Controller
	userControllerV2    user.Controller
	classControllerV1   class.Controller
	roomControllerV1    room.Controller
	teacherControllerV1 teacher.Controller
	studentControllerV1 student.Controller
}

func (a *App) initFlag() {
	flag.StringVar(&a.Name, "name", "service-name", "")
	flag.StringVar(&a.Version, "version", "1.0.0", "")
	flag.StringVar(&a.ConfigFilePath, "config-file-path", "./configs", "Config file path: path to config dir")
	flag.StringVar(&a.ConfigFile, "config-file", "config", "Config file path: path to config dir")
	flag.Parse()
}

func (a *App) initConfig() {
	configSource := &config.Viper{
		ConfigType: constants.ConfigTypeFile,
		FilePath:   a.ConfigFilePath,
		ConfigFile: a.ConfigFile,
	}
	err := configSource.InitConfig()
	if err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	a.registerRoute()
	err := a.router.Run(fmt.Sprintf("%s:%s", a.restConfig.Path, a.restConfig.Port))
	if err != nil {
		return err
	}
	return nil
}

func (a *App) registerRoute() {
	// Base API group
	api := a.router.Group("/api")

	// Routes
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Swagger route under /api
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/swagger/doc.json")))

	user.RegisterRoutesV1(api, a.userControllerV1)
	user.RegisterRoutesV2(api, a.userControllerV2)
	class.RegisterRoutesV1(api, a.classControllerV1, config.GetManager())
	room.RegisterRoutesV1(api, a.roomControllerV1, config.GetManager())
	teacher.RegisterRoutesV1(api, a.teacherControllerV1, config.GetManager())
	student.RegisterRoutesV1(api, a.studentControllerV1, config.GetManager())

}

func inject(
	app *App,
	userControllerV1 *user.ControllerV1,
	userControllerV2 *user.ControllerV2,
	classControllerV1 *class.ControllerV1,
	roomControllerV1 *room.ControllerV1,
	teacherControllerV1 teacher.Controller,
	studentControllerV1 student.Controller,
) error {
	app.userControllerV1 = userControllerV1
	app.userControllerV2 = userControllerV2
	app.classControllerV1 = classControllerV1
	app.roomControllerV1 = roomControllerV1
	app.teacherControllerV1 = teacherControllerV1
	app.studentControllerV1 = studentControllerV1
	return nil
}

// @title Devices manager API
// @version 1.0
// @description Restfull API Application for web devices management
// @scheme https
// @host localhost:9000
// @BasePath /api
func main() {
	app := &App{}
	app.initFlag()
	app.initConfig()
	restConfig := httpConfig.RestServer{}
	err := viper.UnmarshalKey("http", &restConfig)
	if err != nil {
		panic(err)
	}
	app.restConfig = restConfig

	// Init gin router
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(middleware.JsonLogMiddleware), gin.Recovery(), middleware.CorsMiddleware())
	router.HandleMethodNotAllowed = true
	router.NoMethod(
		func(context *gin.Context) {
			context.AbortWithStatus(http.StatusMethodNotAllowed)
		},
	)
	app.router = router

	err = wireApp(app)
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}

}
