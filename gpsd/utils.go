package gpsd

import (
	"github.com/fatih/structs"
)

func init() {
	structs.DefaultTagName = "json"
}

//TPVToGeoJSON adds a location field which is in geojson format
func TPVToGeoJSON(report *TPVReport) (jsonReport map[string]interface{}) {
	jsonReport = structs.Map(report)
	jsonReport["location"] = map[string]interface{}{
		"latitude":  report.Lat,
		"longitude": report.Lon,
	}
	delete(jsonReport, "lat")
	delete(jsonReport, "lon")
	return
}
