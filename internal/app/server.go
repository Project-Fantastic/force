package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"tamago/internal/controllers"
	"tamago/internal/dao"
	"tamago/internal/sessions"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/wader/gormstore"
)

// Server contains everything we need to run a Gin web server
type Server struct {
	Config  *viper.Viper
	Engine  *gin.Engine
	DB      *gorm.DB
	Session *gormstore.Store
	Views   *controllers.View
	API     *controllers.API
}

// Init configs and initializes the app server
func (s *Server) Init() {
	appConfig := s.Config.GetStringMap("app")
	dbConfig := s.Config.GetStringMapString("postgres")

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig["host"], dbConfig["port"], dbConfig["user"], dbConfig["password"], dbConfig["db"], dbConfig["ssl_mode"]))

	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	s.DB = db

	s.Session = sessions.NewDBSession(s.DB, appConfig["secret_key"].(string))

	if !appConfig["debug"].(bool) {
		log.Println("Running server in production mode")
		gin.SetMode(gin.ReleaseMode)
	}

	s.Engine = gin.Default()
	s.Engine.Use(cors.Default())
	s.Engine.LoadHTMLGlob(s.Config.GetString("TEMPLATE_PATH"))

	dao := dao.NewDAO(s.DB)

	s.Views = controllers.NewView(s.Engine, s.DB, s.Session, dao)
	s.Views.CreateViewRoutes()

	s.API = controllers.NewAPI(s.Engine, s.DB, s.Session, dao)
	s.API.CreateAPIRoutes()
}

// Start bootstraps and runs the app server
func (s *Server) Start() {
	appConfig := s.Config.GetStringMap("app")

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(appConfig["port"].(int)),
		Handler: s.Engine,
	}

	defer s.DB.Close()

	// TODO: handle quit
	go sessions.CleanDBSession(s.Session, make(chan struct{}))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
}
