package main

import (
	"cryptocurrency/internal/handlers"
	"cryptocurrency/internal/repository"
	"cryptocurrency/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
)

func inject(d *dataSource, cronObservers *cron.Cron) (*gin.Engine, error) {
	log.Printf("Injecting data sources")

	/*
	*  repository level
	 */
	observerRepository := repository.NewObserverRepository(d.MongoDBClient, d.BadgerClient, cronObservers)
	infoRepository := repository.NewInfoRepository(d.MongoDBClient)
	accountRepository := repository.NewAccountRepository(d.MongoDBClient, d.BadgerClient, cronObservers)
	logRepository := repository.NewLogRepository(d.MongoDBClient)
	managerRepository := repository.NewManagerRepository(d.MongoDBClient)

	/*
	*  repository level
	 */
	observerService := service.NewObserverService(&service.OSConfig{
		ObserverRepository: observerRepository,
		InfoRepository:     infoRepository,
	})
	infoService := service.NewInfoService(&service.ISConfig{InfoRepository: infoRepository})
	accountService := service.NewAccountService(&service.ASConfig{AccountRepository: accountRepository})
	logService := service.NewLogService(&service.LSConfig{LogRepository: logRepository})
	tokenService := service.NewTokenService(&service.TSConfig{})
	managerService := service.NewManagerService(&service.MSConfig{
		ManagerRepository: managerRepository,
		AccountRepository: accountRepository,
	})

	// initialize gin.Engine
	router := gin.Default()

	handlers.NewHandler(&handlers.Config{
		R:               router,
		ObserverService: observerService,
		InfoService:     infoService,
		AccountService:  accountService,
		LogService:      logService,
		TokenService:    tokenService,
		ManagerService:  managerService,
	})

	return router, nil
}
