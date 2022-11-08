package handlers

import (
	docs "cryptocurrency/internal/handlers/docs"
	"cryptocurrency/internal/middleware"
	"cryptocurrency/internal/models"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	ObserverService models.ObserverService
	InfoService     models.InfoService
	AccountService  models.AccountService
	LogService      models.LogService
	TokenService    models.TokenService
	ManagerService  models.ManagerService
}

type Config struct {
	R               *gin.Engine
	ObserverService models.ObserverService
	InfoService     models.InfoService
	AccountService  models.AccountService
	LogService      models.LogService
	TokenService    models.TokenService
	ManagerService  models.ManagerService
}

func NewHandler(c *Config) {
	h := &Handler{
		ObserverService: c.ObserverService,
		InfoService:     c.InfoService,
		AccountService:  c.AccountService,
		LogService:      c.LogService,
		TokenService:    c.TokenService,
		ManagerService:  c.ManagerService,
	}

	docs.SwaggerInfo.BasePath = "/api/"
	c.R.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	g := c.R.Group("/api")

	// Telegram Bot
	// Observer
	g.GET("/observers/price_observer/create",
		middleware.AuthMiddleware(h.TokenService),
		h.CreatePriceObserver)
	g.GET("/observers/percentage_observer/create",
		middleware.AuthMiddleware(h.TokenService),
		h.CreatePercentageObserver)

	// Search
	g.GET("/search/try_find_currency",
		middleware.AuthMiddleware(h.TokenService),
		h.TryFindCurrency)

	// Account
	g.GET("/account/new_account",
		middleware.AuthMiddleware(h.TokenService),
		h.NewAccount)

	g.GET("/account/portfolio/view",
		middleware.AuthMiddleware(h.TokenService),
		h.ViewPortfolio)
	g.GET("/account/portfolio/add",
		middleware.AuthMiddleware(h.TokenService),
		h.AddToPortfolio)
	g.GET("/account/portfolio/update",
		middleware.AuthMiddleware(h.TokenService),
		h.UpdateRecordPortfolio)
	g.GET("/account/portfolio/delete",
		middleware.AuthMiddleware(h.TokenService),
		h.DeleteRecordPortfolio)

	g.GET("/account/price_observers/list",
		middleware.AuthMiddleware(h.TokenService),
		h.GetAccountPriceObservers)
	g.GET("/account/price_observers/delete",
		middleware.AuthMiddleware(h.TokenService),
		h.DeleteAccountPriceObserver)
	g.GET("/account/price_observers/change_status",
		middleware.AuthMiddleware(h.TokenService),
		h.ChangeAccountObserver)

	g.GET("/account/percentage_observers/list",
		middleware.AuthMiddleware(h.TokenService),
		h.GetAccountPercentageObservers)
	g.GET("/account/percentage_observer/delete",
		middleware.AuthMiddleware(h.TokenService),
		h.DeleteAccountPercentageObserver)

	g.GET("/account/subscription/view",
		middleware.AuthMiddleware(h.TokenService),
		h.GetAccountSubscription)
	//g.GET("/account/subscription/extend")

	g.GET("/account/promo_code/check",
		middleware.AuthMiddleware(h.TokenService),
		h.CheckPromoCode)

	// Info
	g.GET("/info/users_languages",
		middleware.AuthMiddleware(h.TokenService),
		h.GetUsersLanguages)
	g.GET("/info/users_admins",
		middleware.AuthMiddleware(h.TokenService),
		h.GetUsersAdmins)

	g.GET("/info/index_rating",
		middleware.AuthMiddleware(h.TokenService),
		h.GetIndexRating)

	g.GET("/info/trending",
		middleware.AuthMiddleware(h.TokenService),
		h.GetTrendingCurrencies)

	g.GET("/info/supported_vs_currencies/view",
		middleware.AuthMiddleware(h.TokenService),
		h.GetSupportedVSCurrencies)

	g.GET("/info/currency/data/full_version",
		middleware.AuthMiddleware(h.TokenService),
		h.GetBasicCurrencyInfo)
	g.GET("/info/currency/data/short_version",
		middleware.AuthMiddleware(h.TokenService),
		h.GetBasicCurrencyInfoShortVersion)

	g.GET("/info/exchange/bestchange",
		middleware.AuthMiddleware(h.TokenService),
		h.ViewExchangeRate)

	// Web Site
	g.GET("/manager/promo_code/create",
		middleware.AuthMiddleware(h.TokenService),
		h.CreatePromoCode)
	g.GET("/manager/promo_code/view",
		middleware.AuthMiddleware(h.TokenService),
		h.ViewPromoCode)

	g.GET("/manager/account/set_premium",
		middleware.AuthMiddleware(h.TokenService),
		h.AddHoursToUserSubscription)

	g.GET("/manager/info/view_online_users",
		middleware.AuthMiddleware(h.TokenService),
		h.ViewOnlineUsers)
	g.GET("/manager/info/view_count_requests",
		middleware.AuthMiddleware(h.TokenService),
		h.ViewCountRequests)
	g.GET("/manager/info/user_actions",
		middleware.AuthMiddleware(h.TokenService),
		h.ViewUserActions)
	g.GET("/manager/info/users_info",
		middleware.AuthMiddleware(h.TokenService),
		h.ViewUsersInfo)
	g.GET("/manager/info/user_info",
		middleware.AuthMiddleware(h.TokenService),
		h.ViewUserInfo)
	g.GET("/manager/info/view_active_observers",
		middleware.AuthMiddleware(h.TokenService),
		h.ViewActiveObservers)

	g.GET("/manager/utils/send_message",
		middleware.AuthMiddleware(h.TokenService),
		h.SendMessage)

	g.GET("/manager/utils/add_user_to_admin_group",
		middleware.AuthMiddleware(h.TokenService),
		h.AddUserToAdmins)

	//g.GET("/manager/utils/change_to_service_mod")
	//g.GET("/manager/utils/change_to_active_mod")
}
