package usgs

type Period uint

const (
	PastHour Period = 1 + iota
	PastDay
	Past7Days
	Past30Days
)

func (p Period) String() string {
	str, _ := p.summaryString()
	return str
}

func (p Period) summaryString() (string, bool) {
	switch p {
	default:
		return "", false
	case PastHour:
		return "hour", true
	case PastDay:
		return "day", true
	case Past7Days:
		return "week", true
	case Past30Days:
		return "month", true
	}
}

type Magnitude uint

const (
	MAll Magnitude = iota
	MSignificant
	M4Dot5Plus
	M2Dot5Plus
	M1Dot0Plus
)

func (m Magnitude) summaryString() (string, bool) {
	switch m {
	default:
		return "", false
	case MAll:
		return "all", true
	case MSignificant:
		return "significant", true
	case M4Dot5Plus:
		return "4.5", true
	case M2Dot5Plus:
		return "2.5", true
	case M1Dot0Plus:
		return "1.0", true
	}
}

/*
Within each period we have:
* Significant Earthquakes:
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/significant_month.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/significant_week.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/significant_day.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/significant_hour.geojson
* M4.5+ Earthquakes:
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/4.5_month.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/4.5_week.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/4.5_day.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/4.5_hour.geojson
* M2.5+ Earthquakes:
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/2.5_month.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/2.5_week.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/2.5_day.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/2.5_hour.geojson
* M1.0+ Earthquakes:
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/1.0_month.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/1.0_week.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/1.0_day.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/1.0_hour.geojson
* All Earthquakes:
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_month.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_week.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_day.geojson
  + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_hour.geojson
*/
