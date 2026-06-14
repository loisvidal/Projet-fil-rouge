package models

import "time"

type User struct {
	IdUser         int
	NameUser       string
	Password       string
	Mail           string
	IsConnect      bool
	IsAdmin        bool
	IsConfirmed    bool
	FailedAttempts int
	LockedUntil    *time.Time
	LikedProperty  []Property
	BuyProperty    []Property
	SellProperty   []Property
}

var PropertyTypes = []string{
	"house", "apartment", "studio", "loft", "villa",
	"townhouse", "penthouse", "commercial", "land",
}

var PropertyTypeLabels = map[string]string{
	"house":      "Maison",
	"apartment":  "Appartement",
	"studio":     "Studio",
	"loft":       "Loft",
	"villa":      "Villa",
	"townhouse":  "Maison de ville",
	"penthouse":  "Penthouse",
	"commercial": "Local commercial",
	"land":       "Terrain",
}

var PropertyTypeMultipliers = map[string]float64{
	"house":      1.0,
	"apartment":  0.7,
	"studio":     0.5,
	"loft":       0.85,
	"villa":      1.8,
	"townhouse":  1.1,
	"penthouse":  1.5,
	"commercial": 1.3,
	"land":       0.3,
}

type Property struct {
	NameProperty   string   `json:"NameProperty"`
	DescProprety   string   `json:"DescProprety"`
	IdProperty     int      `json:"IdProperty"`
	PriceProperty  float64  `json:"PriceProperty"`
	IsSellProperty bool     `json:"IsSellProperty"`
	IsLiked        bool     `json:"IsLiked"`
	ImgProperty    []string `json:"ImgProperty"`
	OwnerID        int      `json:"OwnerID"`
	Type           string   `json:"Type"`
	Rooms          int      `json:"Rooms"`
	Location       string   `json:"Location"`
	Surface        float64  `json:"Surface"`
}

func (p Property) TypeLabel() string {
	if label, ok := PropertyTypeLabels[p.Type]; ok {
		return label
	}
	return p.Type
}

type Auction struct {
	ID             int
	PropertyID     int
	PropertyName   string
	SellerID       int
	StartPrice     float64
	CurrentPrice   float64
	MinBidStep     float64
	WinnerID       *int
	WinnerName     string
	StartTime      time.Time
	EndTime        time.Time
	IsActive       bool
	BidCount       int
}

type Bid struct {
	ID         int
	AuctionID  int
	UserID     int
	UserName   string
	Amount     float64
	CreatedAt  time.Time
}

type SalesRecord struct {
	ID         int
	PropertyID int
	Price      float64
	Type       string
	Location   string
	Rooms      int
	Surface    float64
	SoldAt     time.Time
}

type PricePrediction struct {
	PredictedPrice  float64 `json:"PredictedPrice"`
	Confidence      float64 `json:"Confidence"`
	MinRange        float64 `json:"MinRange"`
	MaxRange        float64 `json:"MaxRange"`
	Type            string  `json:"Type"`
	TypeLabel       string  `json:"TypeLabel"`
	Location        string  `json:"Location"`
	PostalCode      string  `json:"PostalCode"`
	Rooms           int     `json:"Rooms"`
	Surface         float64 `json:"Surface"`
	PricePerSqm     float64 `json:"PricePerSqm"`
	DataPoints      int     `json:"DataPoints"`
	Latitude        float64 `json:"Latitude"`
	Longitude       float64 `json:"Longitude"`
}

type AnalyticSummary struct {
	TotalSales       int
	AveragePrice     float64
	AveragePricePerSqm float64
	TotalRevenue     float64
	MostActiveRegion string
	TopPropertyType  string
	PriceTrend       string
	Prediction       *PricePrediction
}

type LocationFactors struct {
	Label     string
	BasePrice float64
	PerSqm    float64
	PerRoom   float64
	Weight    float64
}

type VerifiedLocation struct {
	City        string  `json:"city"`
	PostalCode  string  `json:"postal_code"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	DisplayName string  `json:"display_name"`
	Verified    bool    `json:"verified"`
	Error       string  `json:"error,omitempty"`
}

type GeocodeRequest struct {
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type GeocodeResponse struct {
	Location *VerifiedLocation `json:"location"`
}
