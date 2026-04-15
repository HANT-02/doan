package main

import (
	httpConfig "doan/cmd/http/config"
	"doan/cmd/http/controllers/academic"
	"doan/cmd/http/controllers/account"
	"doan/cmd/http/controllers/class"
	"doan/cmd/http/controllers/course"
	"doan/cmd/http/controllers/leave"
	lessons "doan/cmd/http/controllers/lesson"
	"doan/cmd/http/controllers/material"
	"doan/cmd/http/controllers/predictive"
	"doan/cmd/http/controllers/program"
	"doan/cmd/http/controllers/room"
	"doan/cmd/http/controllers/scheduling"
	"doan/cmd/http/controllers/shift"
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
	Name                   string
	Version                string
	ConfigFilePath         string
	ConfigFile             string
	router                 *gin.Engine
	restConfig             httpConfig.RestServer
	userControllerV1       user.Controller
	userControllerV2       user.Controller
	accountControllerV1    account.Controller
	academicControllerV1   academic.Controller
	classControllerV1      class.Controller
	roomControllerV1       room.Controller
	shiftControllerV1      shift.Controller
	teacherControllerV1    teacher.Controller
	studentControllerV1    student.Controller
	courseControllerV1     course.Controller
	programControllerV1    program.Controller
	schedulingControllerV1 scheduling.Controller
	materialControllerV1   material.Controller
	predictiveControllerV1 predictive.Controller
	lessonControllerV1     lessons.Controller
	leaveControllerV1      leave.Controller
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
	account.RegisterRoutesV1(api, a.accountControllerV1, config.GetManager())
	academic.RegisterRoutesV1(api, a.academicControllerV1, config.GetManager())
	class.RegisterRoutesV1(api, a.classControllerV1, config.GetManager())
	room.RegisterRoutesV1(api, a.roomControllerV1, config.GetManager())
	shift.RegisterRoutesV1(api, a.shiftControllerV1, config.GetManager())
	teacher.RegisterRoutesV1(api, a.teacherControllerV1, config.GetManager())
	student.RegisterRoutesV1(api, a.studentControllerV1, config.GetManager())
	course.RegisterRoutesV1(api, a.courseControllerV1, config.GetManager())
	program.RegisterRoutesV1(api, a.programControllerV1, config.GetManager())
	scheduling.RegisterRoutesV1(api, a.schedulingControllerV1, config.GetManager())
	material.RegisterRoutesV1(api, a.materialControllerV1, config.GetManager())
	predictive.RegisterRoutesV1(api, a.predictiveControllerV1, config.GetManager())
	lessons.RegisterRoutesV1(api, a.lessonControllerV1, config.GetManager())
	leave.RegisterRoutesV1(api, a.leaveControllerV1, config.GetManager())

}

func inject(
	app *App,
	userControllerV1 *user.ControllerV1,
	userControllerV2 *user.ControllerV2,
	accountControllerV1 account.Controller,
	academicControllerV1 academic.Controller,
	classControllerV1 *class.ControllerV1,
	roomControllerV1 *room.ControllerV1,
	shiftControllerV1 shift.Controller,
	teacherControllerV1 teacher.Controller,
	studentControllerV1 student.Controller,
	courseControllerV1 course.Controller,
	programControllerV1 program.Controller,
	schedulingControllerV1 scheduling.Controller,
	materialControllerV1 material.Controller,
	predictiveControllerV1 predictive.Controller,
	lessonControllerV1 lessons.Controller,
	leaveControllerV1 leave.Controller,
) error {
	app.userControllerV1 = userControllerV1
	app.userControllerV2 = userControllerV2
	app.accountControllerV1 = accountControllerV1
	app.academicControllerV1 = academicControllerV1
	app.classControllerV1 = classControllerV1
	app.roomControllerV1 = roomControllerV1
	app.shiftControllerV1 = shiftControllerV1
	app.teacherControllerV1 = teacherControllerV1
	app.studentControllerV1 = studentControllerV1
	app.courseControllerV1 = courseControllerV1
	app.programControllerV1 = programControllerV1
	app.schedulingControllerV1 = schedulingControllerV1
	app.materialControllerV1 = materialControllerV1
	app.predictiveControllerV1 = predictiveControllerV1
	app.lessonControllerV1 = lessonControllerV1
	app.leaveControllerV1 = leaveControllerV1
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
