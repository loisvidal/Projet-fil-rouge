package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type NominatimResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Address     struct {
		HouseNumber   string `json:"house_number"`
		Road          string `json:"road"`
		City          string `json:"city"`
		Town          string `json:"town"`
		Village       string `json:"village"`
		County        string `json:"county"`
		State         string `json:"state"`
		Postcode      string `json:"postcode"`
		Country       string `json:"country"`
		CountryCode   string `json:"country_code"`
	} `json:"address"`
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

var (
	lastRequest time.Time
	mu          sync.Mutex
)

func rateLimit() {
	mu.Lock()
	defer mu.Unlock()
	elapsed := time.Since(lastRequest)
	if elapsed < time.Second {
		time.Sleep(time.Second - elapsed)
	}
	lastRequest = time.Now()
}

func GeocodeAddress(city, postalCode, country string) (*VerifiedLocation, error) {
	var queryParts []string

	if postalCode != "" {
		queryParts = append(queryParts, postalCode)
	}
	if city != "" {
		queryParts = append(queryParts, city)
	}
	if country != "" {
		queryParts = append(queryParts, country)
	}

	if len(queryParts) == 0 {
		return &VerifiedLocation{Verified: false, Error: "Aucun critère fourni"}, nil
	}

	q := url.Values{}
	q.Set("q", strings.Join(queryParts, ", "))
	q.Set("format", "json")
	q.Set("addressdetails", "1")
	q.Set("limit", "1")

	rawURL := "https://nominatim.openstreetmap.org/search?" + q.Encode()

	rateLimit()

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return &VerifiedLocation{Verified: false, Error: "Erreur requête"}, err
	}
	req.Header.Set("User-Agent", "YPlaza/1.0 (immobilier)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return &VerifiedLocation{Verified: false, Error: "Erreur appel API"}, err
	}
	defer resp.Body.Close()

	var results []NominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return &VerifiedLocation{Verified: false, Error: "Erreur décodage réponse"}, err
	}

	if len(results) == 0 {
		return &VerifiedLocation{Verified: false, Error: "Aucun résultat trouvé"}, nil
	}

	r := results[0]
	addr := r.Address

	cityName := addr.City
	if cityName == "" {
		cityName = addr.Town
	}
	if cityName == "" {
		cityName = addr.Village
	}
	if cityName == "" {
		cityName = addr.County
	}

	lat, lon := 0.0, 0.0
	fmt.Sscanf(r.Lat, "%f", &lat)
	fmt.Sscanf(r.Lon, "%f", &lon)

	return &VerifiedLocation{
		City:        cityName,
		PostalCode:  addr.Postcode,
		Country:     addr.Country,
		CountryCode: addr.CountryCode,
		Latitude:    lat,
		Longitude:   lon,
		DisplayName: r.DisplayName,
		Verified:    true,
	}, nil
}

func GeocodePostalCode(postalCode, country string) (*VerifiedLocation, error) {
	q := url.Values{}
	q.Set("postalcode", postalCode)
	if country != "" {
		q.Set("country", country)
	}
	q.Set("format", "json")
	q.Set("addressdetails", "1")
	q.Set("limit", "1")

	rawURL := "https://nominatim.openstreetmap.org/search?" + q.Encode()

	rateLimit()

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return &VerifiedLocation{Verified: false, Error: "Erreur requête"}, nil
	}
	req.Header.Set("User-Agent", "YPlaza/1.0 (immobilier)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return &VerifiedLocation{Verified: false, Error: "Erreur appel API"}, nil
	}
	defer resp.Body.Close()

	var results []NominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return &VerifiedLocation{Verified: false, Error: "Erreur décodage"}, nil
	}

	if len(results) == 0 {
		return &VerifiedLocation{Verified: false, Error: "Code postal introuvable"}, nil
	}

	r := results[0]
	addr := r.Address

	cityName := addr.City
	if cityName == "" {
		cityName = addr.Town
	}
	if cityName == "" {
		cityName = addr.Village
	}

	lat, lon := 0.0, 0.0
	fmt.Sscanf(r.Lat, "%f", &lat)
	fmt.Sscanf(r.Lon, "%f", &lon)

	return &VerifiedLocation{
		City:        cityName,
		PostalCode:  addr.Postcode,
		Country:     addr.Country,
		CountryCode: addr.CountryCode,
		Latitude:    lat,
		Longitude:   lon,
		DisplayName: r.DisplayName,
		Verified:    true,
	}, nil
}
