package slack_bot

import (
	"cryptocurrency/internal/models/mongo_models"
	"time"
)

type UserInfo struct {
	TelegramUserId         int64                               `json:"telegramUserId"         bson:"telegram_user_id"`
	Status                 string                              `json:"status"                 bson:"status"`
	PremiumDuration        float64                             `json:"premiumDuration"        bson:"premium_duration"`
	CreatedAt              time.Time                           `json:"createdAt"              bson:"created_at"`
	PromoCodeActivations   []mongo_models.PromoCodeActivations `json:"promoCodeActivations"   bson:"promo_code_activations"`
	PriceObservers         []mongo_models.Observer             `json:"priceObservers"         bson:"price_observers"`
	PercentageObserver     mongo_models.PercentageObserver     `json:"percentageObserver"     bson:"percentage_observer"`
	Portfolio              []mongo_models.Portfolio            `json:"portfolio"              bson:"portfolio"`
	CountRequestsLast24H   int                                 `json:"countRequestsLast24H"   bson:"count_requests_last_24_h"`
	CountRequestsLast7D    int                                 `json:"countRequestsLast7D"    bson:"count_requests_last_7_d"`
	CountRequestsLastMonth int                                 `json:"countRequestsLastMonth" bson:"count_requests_last_month"`
	CountRequests          int                                 `json:"countRequests"          bson:"count_requests"`
}
