/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/json"
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/shopspring/decimal"
	"os"
	"path"
)

type CountriesGeoJson struct {
	Type     string           `json:"type"`
	Features []CountryFeature `json:"features"`
}
type CountryFeature struct {
	Type       string      `json:"type"`
	Properties CountryInfo `json:"properties"`
	Geometry   Geometry    `json:"geometry"`
}
type CountryInfo struct {
	ADMIN     string `json:"ADMIN"`
	ISOA2     string `json:"ISO_A2"`
	Continent string `json:"continent"`
}

type Geometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

var (
	CountriesGeo CountriesGeoJson
)

func InitCountryLocator() error {
	geoFile := path.Join(conf.GetEnvConf().ConfigPath, "countries_geo.json")
	file, err := os.ReadFile(geoFile)
	if err != nil {
		return err
	}
	//fmt.Println(CountriesGeo.Type)
	err = json.Unmarshal(file, &CountriesGeo)
	if err != nil {
		return err
	}
	return nil
}

func PointValid(latitude, longitude float64) bool {
	lat := decimal.NewFromFloat(latitude)
	lng := decimal.NewFromFloat(longitude)
	minLat := decimal.NewFromFloat(-90)
	maxLat := decimal.NewFromFloat(90)
	minLng := decimal.NewFromFloat(-180)
	maxLng := decimal.NewFromFloat(180)
	if lat.LessThanOrEqual(minLat) || lat.GreaterThanOrEqual(maxLat) {
		return false
	}
	if lng.LessThanOrEqual(minLng) || lng.GreaterThanOrEqual(maxLng) {
		return false
	}

	return true
}

func isPointInCountry(feature CountryFeature, point Point) bool {
	//geometry := feature.Geometry
	//countryGeoType := geometry.Type
	//countryCoordinates := geometry.Coordinates
	//if countryGeoType == "Polygon" {
	//return booleanPointInPolygon(point, polygon(countryCoordinates))
	//} else if countryGeoType == "MultiPolygon" {
	//return booleanPointInPolygon(point, multiPolygon(countryCoordinates))
	//}

	//countryPolygon := &Polygon{}
	//for _, point := range feature.Geometry.Coordinates[0][0] {
	//	countryPolygon.Add(&Point{Lng: point[0], Lat: point[1]})
	//	//fmt.Println(point)
	//}
	flag := false
ok:
	for _, coordinates1 := range feature.Geometry.Coordinates {
		for _, coordinates2 := range coordinates1 {
			countryPolygon := &Polygon{}
			for _, point := range coordinates2 {
				if len(point) >= 2 {
					countryPolygon.Add(&Point{Lng: point[0], Lat: point[1]})
				}
			}
			if countryPolygon.Contains(&point) {
				//fmt.Println(feature.Properties)
				flag = true
				break ok
			}
		}
	}
	return flag
}

func FindCountryByCoordinate(lat, lng float64) CountryInfo {
	point := Point{Lat: lat, Lng: lng}
	var country CountryInfo
	for _, feature := range CountriesGeo.Features {
		//fmt.Println(index)
		//fmt.Println(feature)

		if isPointInCountry(feature, point) {
			country.ADMIN = feature.Properties.ADMIN
			country.ISOA2 = feature.Properties.ISOA2
			country.Continent = feature.Properties.Continent
			break
		}
	}
	if country.ADMIN == "" {
		country.ADMIN = "Global"
		country.ISOA2 = "GL"
		country.Continent = "Global"
	}

	return country
}
