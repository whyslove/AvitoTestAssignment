package main

import (
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	avitotest "github.com/whyslove/avito-test"
	"github.com/whyslove/avito-test/core/handler"
	"github.com/whyslove/avito-test/core/repository"
	"github.com/whyslove/avito-test/core/service"
	"github.com/whyslove/avito-test/pkg/converter"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetLevel(logrus.DebugLevel)
	if err := initConfig(); err != nil {
		logrus.Fatalf("error in init config file, %s", err.Error())
	}
	if err := godotenv.Load("./configs/.env"); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDb(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("error in initialize db: %s", err.Error())
	}

	dbRedis, _ := strconv.Atoi(viper.GetString("redis.db"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host") + viper.GetString("redis.port"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbRedis,
	})
	if err := rdb.Ping().Err(); err != nil {
		logrus.Fatalf("error in initialize redis %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	converter := converter.NewCurrencyConverter(rdb)
	cronFunction(converter.UpdateInfo)                  // set function to execute 1m
	handlers := handler.NewHandler(services, converter) //currency converter)
	srv := new(avitotest.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error in running server, error is %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func cronFunction(f func()) {
	f()
	c := cron.New()
	c.AddFunc("@every 1m", f)
	c.Start()
}
