package slack_bot

import "cryptocurrency/internal/models/mongo_models"

type PromoCodesView struct {
	PromoCode            mongo_models.PromoCode
	PromoCodeActivations []mongo_models.PromoCodeActivations
}
