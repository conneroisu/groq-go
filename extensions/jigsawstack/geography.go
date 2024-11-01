package jigsawstack

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	geographyEndpoint Endpoint = "/v1/geo/search"
)

type (
	// GeographyRequest represents a request structure for geography API.
	GeographyRequest struct {
		Query        string  `json:"query"`
		Country      string  `json:"country,omitempty"`
		Latitude     float64 `json:"latitude,omitempty"`
		ProximityLat float64 `json:"proximity_lat,omitempty"`
		Longitude    float64 `json:"longitude,omitempty"`
		ProximityLng float64 `json:"proximity_lng,omitempty"`
		Types        string  `json:"types,omitempty"`
	}
	// GeographyResponse represents a response structure for geography API.
	GeographyResponse struct {
		Success bool `json:"success"`
		Data    []struct {
			Type           string `json:"type"`
			FullAddress    string `json:"full_address"`
			Name           string `json:"name"`
			PlaceFormatted string `json:"place_formatted"`
			Postcode       string `json:"postcode"`
			Place          string `json:"place"`
			Region         struct {
				Name           string `json:"name"`
				RegionCode     string `json:"region_code"`
				RegionCodeFull string `json:"region_code_full"`
			} `json:"region"`
			Country struct {
				Name              string `json:"name"`
				CountryCode       string `json:"country_code"`
				CountryCodeAlpha3 string `json:"country_code_alpha_3"`
			} `json:"country"`
			Language string `json:"language"`
			Geoloc   struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geoloc"`
			PoiCategory         []string `json:"poi_category"`
			AddtionalProperties struct {
				Phone     string `json:"phone"`
				Website   string `json:"website"`
				OpenHours struct {
				} `json:"open_hours"`
			} `json:"addtional_properties"`
		} `json:"data"`
	}
)

// URLQuery converts the params into params on the given url.
func (r *GeographyRequest) URLQuery(url *url.URL) {
	values := url.Query()
	if r.Query != "" {
		values.Set("query", r.Query)
	}
	if r.Country != "" {
		values.Set("country", r.Country)
	}
	var strLat, strLng string
	if r.Latitude != 0 {
		strLat = strconv.FormatFloat(r.Latitude, 'f', -1, 64)
		values.Set("latitude", strLat)
	}
	if r.ProximityLat != 0 {
		strLat = strconv.FormatFloat(r.ProximityLat, 'f', -1, 64)
		values.Set("proximity_lat", strLat)
	}
	if r.Longitude != 0 {
		strLng = strconv.FormatFloat(r.Longitude, 'f', -1, 64)
		values.Set("longitude", strLng)
	}
	if r.ProximityLng != 0 {
		strLng = strconv.FormatFloat(r.ProximityLng, 'f', -1, 64)
		values.Set("proximity_lng", strLng)
	}
	if r.Types != "" {
		values.Set("types", r.Types)
	}
	url.RawQuery = values.Encode()
}

// GeographySearch performs a geography search api call over a query string.
//
// https://api.jigsawstack.com/v1/geo/search
//
// https://docs.jigsawstack.com/api-reference/geo/search
func (j *JigsawStack) GeographySearch(
	ctx context.Context,
	request GeographyRequest,
) (response GeographyResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(geographyEndpoint),
		builders.WithQuerier(&request),
	)
	if err != nil {
		return
	}
	var resp GeographyResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}

// GeographyGeocode performs a geography geocode api call over a query string.
//
// GET https://api.jigsawstack.com/v1/geo/geocode
//
// https://docs.jigsawstack.com/api-reference/geo/geocode
func (j *JigsawStack) GeographyGeocode(
	ctx context.Context,
	request GeographyRequest,
) (response GeographyResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodGet,
		j.baseURL+string(geographyEndpoint),
		builders.WithQuerier(&request),
	)
	if err != nil {
		return
	}
	var resp GeographyResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
