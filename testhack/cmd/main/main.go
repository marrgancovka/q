package main

import (
	"context"
	"fmt"
	"hack/internal/config"
	"hack/internal/pkg"

	_authHandler "hack/internal/auth/handler"
	_sessionRepo "hack/internal/auth/repository/redis"
	_authUCase "hack/internal/auth/usecase"
	_fStorageRepo "hack/internal/file_storage/repository"
	_fStorageUCase "hack/internal/file_storage/usecase"
	_middleware "hack/internal/middlware"
	_userHandler "hack/internal/user/delivery/http"
	_userRepo "hack/internal/user/repository/postgres"
	_userUCase "hack/internal/user/usecase"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
	os.Getenv(config.DBUserEnv),
	os.Getenv(config.DBPasswordEnv),
	os.Getenv(config.DBHostEnv),
	os.Getenv(config.DBPortEnv),
	os.Getenv(config.DBNameEnv),
)

func main() {
	pkg.InitLogger()
	logger := pkg.GetLogger()
	router := mux.NewRouter()
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connStr}), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix:   os.Getenv(config.DBSchemaNameEnv) + ".",
		SingularTable: false,
	}})
	if err != nil {
		logger.Fatalf("db connection error %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv(config.RedisHostEnv) + ":" + os.Getenv(config.RedisPortEnv),
		Password: os.Getenv(config.RedisPasswordEnv),
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed connect to redis : %v", err)
	}

	s3Session, err := session.NewSession(
		&aws.Config{
			Endpoint:         aws.String(config.S3Endpoint),
			Region:           aws.String(config.S3Region),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv(config.S3AccessKeyEnv),
				os.Getenv(config.S3ASecretKeyEnv),
				"",
			),
		},
	)
	if err != nil {
		log.Fatalf("Failed create S3 session : %v", err)
	}

	userRepo := _userRepo.New(db)
	sessionRepo := _sessionRepo.NewSessionRepo(rdb)
	fStorageRepo := _fStorageRepo.New(s3Session)

	fStorageUC := _fStorageUCase.New(fStorageRepo)
	usersUC := _userUCase.New(userRepo, fStorageUC)
	sessionUC := _authUCase.NewSessionUC(sessionRepo, usersUC)
	authUC := _authUCase.NewAuthUC(sessionUC, usersUC)

	middleware := _middleware.New(sessionUC, logger)

	authH := _authHandler.New(authUC)
	userH := _userHandler.New(usersUC)

	_authHandler.RegisterHTTPRoutes(router, authH, middleware)
	_userHandler.RegisterHTTPRoutes(router, userH, middleware)

	router.Use(middleware.HandlerLogger)
	corsRouter := middleware.Cors(router)

	server := http.Server{
		Addr:         ":" + config.Port,
		Handler:      corsRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("server started")
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("server stopped %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Kill, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("failed to gracefully shutdown server")
	}

}
