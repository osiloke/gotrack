package gpsd

import (
	"github.com/fatih/structs"
)

//TPVToGeoJSON adds a location field which is in geojson format
func TPVToGeoJSON(report *TPVReport) (jsonReport map[string]interface{}) {
	jsonReport = structs.Map(report)
	jsonReport["location"] = map[string]interface{}{
		"latitude":  jsonReport["lat"],
		"longitude": jsonReport["lon"],
	}
	delete(jsonReport, "lat")
	delete(jsonReport, "lon")
	return
}
