package mongo_models

type Cryptocurrency struct {
	Name   string `bson:"name"`
	Slug   string `bson:"slug"`
	Symbol string `bson:"symbol"`
	Rank   int    `bson:"rank"`
}
