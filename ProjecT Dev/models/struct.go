package models

type User struct {
	NameUser      string
	Password      string
	Mail          string
	IdUser        int
	IsConnect     bool
	LikedProperty []Property
	BuyProperty   []Property
	SellProperty  []Property
}

type Property struct {
	NameProperty   string   `json:"NameProperty"`
	DescProprety   string   `json:"DescProprety"`
	IdProperty     int      `json:"IdProperty"`
	PriceProperty  float32  `json:"PriceProperty"`
	IsSellProperty bool     `json:"IsSellProperty"`
	IsLiked        bool     `json:"IsLiked"`
	ImgProperty    []string `json:"ImgProperty"`
}


