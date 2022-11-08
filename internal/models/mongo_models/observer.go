package mongo_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Observer struct {
	Id              primitive.ObjectID `json:"id"              bson:"_id"`
	CryptoID        string             `json:"cryptoID"        bson:"crypto_id"`
	CryptoName      string             `json:"cryptoName"      bson:"crypto_name"`
	CryptoSymbol    string             `json:"cryptoSymbol"    bson:"crypto_symbol"`
	TelegramUserID  int64              `json:"telegramUserID"  bson:"telegram_user_id"`
	CurrencyOfValue string             `json:"currencyOfValue" bson:"currency_of_value"`
	ExpectedValue   float64            `json:"expectedValue"   bson:"expected_value"`
	IsUpDirection   bool               `json:"isUpDirection"   bson:"is_up_direction"`
	IsActive        bool               `json:"isActive"        bson:"is_active"`
	Tier            int                `json:"tier"            bson:"tier"`
	CreatedAt       time.Time          `json:"createdAt"       bson:"created_at"`
}
