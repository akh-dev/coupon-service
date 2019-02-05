package api

import (
	"encoding/json"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Request struct {
	ApiKey string          `json:"apiKey"`
	Data   json.RawMessage `json:"data"`
}

type Response struct {
	Error  []string    `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

type CouponCollection struct {
	Coupons []Coupon `json:"coupons"`
}

type Coupon struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Brand     string             `json:"brand" bson:"brand"`
	Value     float64            `json:"value" bson:"value"`
	Expiry    time.Time          `json:"expiry" bson:"expiry"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

type CouponFilter struct {
	IdIn          []string  `json:"idIn,omitempty"`
	NameContains  string    `json:"nameContains,omitempty"`
	BrandEqual    string    `json:"brandEqual,omitempty"`
	ValueFrom     *int      `json:"valueFrom,omitempty"`
	ValueTo       *int      `json:"valueTo,omitempty"`
	ExpiryFrom    time.Time `json:"expiryFrom"`
	ExpiryTo      time.Time `json:"expiryTo"`
	CreatedAtFrom time.Time `json:"createdAtFrom"`
	CreatedAtTo   time.Time `json:"createdAtTo"`
}
