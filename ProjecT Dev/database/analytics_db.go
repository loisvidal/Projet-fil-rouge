package database

import (
	"RedProject/models"
	"RedProject/services"
	"math"
	"strings"
)

func RecordSale(p models.Property, price float64) error {
	_, err := DB.Exec(
		`INSERT INTO sales_history (property_id, price, property_type, rooms, location, surface)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		p.IdProperty, price, p.Type, p.Rooms, p.Location, p.Surface)
	return err
}

func GetAnalyticSummary() (*models.AnalyticSummary, error) {
	s := &models.AnalyticSummary{}

	DB.QueryRow("SELECT COUNT(*), COALESCE(AVG(price),0), COALESCE(AVG(price/surface),0), COALESCE(SUM(price),0) FROM sales_history WHERE surface > 0").
		Scan(&s.TotalSales, &s.AveragePrice, &s.AveragePricePerSqm, &s.TotalRevenue)

	DB.QueryRow(`SELECT location FROM sales_history GROUP BY location ORDER BY COUNT(*) DESC LIMIT 1`).
		Scan(&s.MostActiveRegion)

	DB.QueryRow(`SELECT property_type FROM sales_history GROUP BY property_type ORDER BY COUNT(*) DESC LIMIT 1`).
		Scan(&s.TopPropertyType)

	DB.QueryRow(`SELECT
		CASE WHEN AVG(price) > (SELECT AVG(price) FROM sales_history WHERE sold_at < DATE_SUB(NOW(), INTERVAL 30 DAY))
		THEN 'hausse' ELSE 'baisse' END FROM sales_history`).Scan(&s.PriceTrend)

	return s, nil
}

func GetSalesHistory() ([]models.SalesRecord, error) {
	rows, err := DB.Query(`SELECT id, property_id, price, property_type, rooms, location, surface, sold_at
		FROM sales_history ORDER BY sold_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.SalesRecord
	for rows.Next() {
		var r models.SalesRecord
		if err := rows.Scan(&r.ID, &r.PropertyID, &r.Price, &r.Type, &r.Rooms, &r.Location, &r.Surface, &r.SoldAt); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}

var locationFactors = []models.LocationFactors{
	{Label: "Paris", BasePrice: 500000, PerSqm: 10000, PerRoom: 80000, Weight: 1.5},
	{Label: "Lyon", BasePrice: 300000, PerSqm: 5500, PerRoom: 50000, Weight: 1.2},
	{Label: "Marseille", BasePrice: 250000, PerSqm: 4500, PerRoom: 40000, Weight: 1.1},
	{Label: "Bordeaux", BasePrice: 280000, PerSqm: 5000, PerRoom: 45000, Weight: 1.15},
	{Label: "Toulouse", BasePrice: 260000, PerSqm: 4200, PerRoom: 38000, Weight: 1.1},
	{Label: "Lille", BasePrice: 240000, PerSqm: 4000, PerRoom: 35000, Weight: 1.0},
	{Label: "Nice", BasePrice: 350000, PerSqm: 6500, PerRoom: 55000, Weight: 1.3},
	{Label: "Nantes", BasePrice: 250000, PerSqm: 4000, PerRoom: 35000, Weight: 1.0},
	{Label: "Strasbourg", BasePrice: 230000, PerSqm: 3800, PerRoom: 32000, Weight: 1.0},
	{Label: "Montpellier", BasePrice: 270000, PerSqm: 4500, PerRoom: 40000, Weight: 1.15},
	{Label: "Rennes", BasePrice: 220000, PerSqm: 3800, PerRoom: 30000, Weight: 1.0},
	{Label: "Sud", BasePrice: 350000, PerSqm: 5800, PerRoom: 50000, Weight: 1.25},
	{Label: "Côte d'Azur", BasePrice: 450000, PerSqm: 7500, PerRoom: 65000, Weight: 1.4},
}

func getLocationFactors(location string) models.LocationFactors {
	loc := strings.ToLower(location)
	for _, lf := range locationFactors {
		if strings.Contains(loc, strings.ToLower(lf.Label)) {
			return lf
		}
	}
	return models.LocationFactors{Label: "France", BasePrice: 200000, PerSqm: 3500, PerRoom: 30000, Weight: 1.0}
}

func getTypeMultiplier(propType string) float64 {
	if m, ok := models.PropertyTypeMultipliers[propType]; ok {
		return m
	}
	return 1.0
}

func getTypeLabel(propType string) string {
	if l, ok := models.PropertyTypeLabels[propType]; ok {
		return l
	}
	return propType
}

func PredictPrice(propType string, rooms int, surface float64, location string) *models.PricePrediction {
	geoLoc, _ := services.GeocodeAddress(location, "", "France")

	resolvedLocation := "France"
	resolvedPostalCode := ""
	lat, lon := 0.0, 0.0

	if geoLoc != nil && geoLoc.Verified {
		resolvedLocation = geoLoc.City
		resolvedPostalCode = geoLoc.PostalCode
		lat = geoLoc.Latitude
		lon = geoLoc.Longitude
		if location != "" {
			resolvedLocation = location
		}
	} else if location != "" {
		resolvedLocation = location
	}

	locFilter := "%"
	if location != "" {
		locFilter = "%" + location + "%"
	}

	rows, err := DB.Query(
		`SELECT price, rooms, surface, location, property_type FROM sales_history
		 WHERE property_type = ? AND location LIKE ?
		 ORDER BY sold_at DESC LIMIT 50`,
		propType, locFilter)
	if err != nil || rows == nil {
		return fallbackPrediction(propType, rooms, surface, location, resolvedLocation, resolvedPostalCode, lat, lon)
	}
	defer rows.Close()

	var prices []float64
	var roomsList []float64
	var surfaces []float64
	var locations []string
	var types []string

	for rows.Next() {
		var price, surf float64
		var r int
		var loc, ptype string
		if err := rows.Scan(&price, &r, &surf, &loc, &ptype); err == nil {
			prices = append(prices, price)
			roomsList = append(roomsList, float64(r))
			surfaces = append(surfaces, surf)
			locations = append(locations, loc)
			types = append(types, ptype)
		}
	}

	if len(prices) < 3 {
		if location != "" {
			rows2, err2 := DB.Query(
				`SELECT price, rooms, surface, property_type FROM sales_history
				 WHERE property_type = ? ORDER BY sold_at DESC LIMIT 30`,
				propType)
			if err2 == nil {
				defer rows2.Close()
				for rows2.Next() {
					var price, surf float64
					var r int
					var ptype string
					if err2 := rows2.Scan(&price, &r, &surf, &ptype); err2 == nil {
						prices = append(prices, price)
						roomsList = append(roomsList, float64(r))
						surfaces = append(surfaces, surf)
						types = append(types, ptype)
					}
				}
			}
		}
		if len(prices) < 3 {
			return fallbackPrediction(propType, rooms, surface, location, resolvedLocation, resolvedPostalCode, lat, lon)
		}
	}

	avgPrice := average(prices)
	avgRooms := average(roomsList)

	var pricePerSqmList []float64
	for i, p := range prices {
		if surfaces[i] > 0 {
			pricePerSqmList = append(pricePerSqmList, p/surfaces[i])
		}
	}
	avgPricePerSqm := 0.0
	if len(pricePerSqmList) > 0 {
		avgPricePerSqm = average(pricePerSqmList)
	}

	pricePerRoom := 0.0
	if avgRooms > 0 {
		pricePerRoom = avgPrice / avgRooms
	}

	predictedBySqm := avgPricePerSqm * surface
	predictedByRoom := pricePerRoom * float64(rooms)
	predicted := (predictedBySqm + predictedByRoom) / 2

	if predicted <= 0 {
		predicted = avgPrice
	}

	if location != "" {
		lf := getLocationFactors(location)
		predicted *= lf.Weight
	}

	variance := 0.0
	for _, p := range prices {
		variance += (p - avgPrice) * (p - avgPrice)
	}
	variance /= float64(len(prices))
	stdDev := math.Sqrt(variance)

	confidence := 0.0
	if avgPrice > 0 {
		confidence = 1.0 - (stdDev / avgPrice)
	}
	if confidence < 0.3 {
		confidence = 0.3
	}
	if confidence > 0.95 {
		confidence = 0.95
	}

	return &models.PricePrediction{
		PredictedPrice: math.Round(predicted*100) / 100,
		Confidence:     math.Round(confidence*100) / 100,
		MinRange:       math.Round((predicted-stdDev)*100) / 100,
		MaxRange:       math.Round((predicted+stdDev)*100) / 100,
		Type:           propType,
		TypeLabel:      getTypeLabel(propType),
		Location:       resolvedLocation,
		PostalCode:     resolvedPostalCode,
		Rooms:          rooms,
		Surface:        surface,
		PricePerSqm:    math.Round(avgPricePerSqm*100) / 100,
		DataPoints:     len(prices),
		Latitude:       lat,
		Longitude:      lon,
	}
}

func fallbackPrediction(propType string, rooms int, surface float64, location, resolvedLocation, postalCode string, lat, lon float64) *models.PricePrediction {
	lf := getLocationFactors(location)
	mult := getTypeMultiplier(propType)

	basePrice := lf.BasePrice * mult
	perRoom := lf.PerRoom * mult
	perSqm := lf.PerSqm * mult

	predicted := basePrice + float64(rooms)*perRoom + surface*perSqm
	predicted *= lf.Weight

	stdDev := predicted * 0.15

	confidence := 0.60
	if lf.Weight > 1.0 {
		confidence = 0.65
	}
	if mult > 1.3 {
		confidence += 0.05
	}

	return &models.PricePrediction{
		PredictedPrice: predicted,
		Confidence:     confidence,
		MinRange:       predicted - stdDev,
		MaxRange:       predicted + stdDev,
		Type:           propType,
		TypeLabel:      getTypeLabel(propType),
		Location:       resolvedLocation,
		PostalCode:     postalCode,
		Rooms:          rooms,
		Surface:        surface,
		PricePerSqm:    perSqm,
		DataPoints:     0,
		Latitude:       lat,
		Longitude:      lon,
	}
}

func average(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}
