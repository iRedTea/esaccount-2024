package main

import (
	"esaccount"
	"esaccount/pkg/handler"
	"esaccount/pkg/repository"
	"esaccount/pkg/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	logrus.SetFormatter(new(logrus.TextFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error occured while init config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env vars: %s", err.Error())
	}

	db, err := repository.NewDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("failed to init database: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService()

	handlers := handler.NewHandler(
		services,
		repos,
	)

	srv := new(esaccount.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
