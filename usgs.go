package usgs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type Response struct {
	Type        string       `json:"type"`
	Metadata    *Metadata    `json:"metadata"`
	Features    []*Feature   `json:"features"`
	BoundingBox *BoundingBox `json:"bbox"`
}

type Metadata struct {
	GeneratedTimeMs float64 `json:"generated"`
	URL             string  `json:"url"`
	Title           string  `json:"title"`
	Status          uint    `json:"status"`
	APIVersion      string  `json:"api"`
	Count           uint    `json:"count"`
}

type Feature struct {
	Id   string `json:"id"`
	Type string `json:"type"`

	Properties *Property `json:"properties"`
	Geometry   *Geometry `json:"geometry"`
}

type Timezone int
type Property struct {
	Magnitude     float32  `json:"mag"`
	Place         string   `json:"place"`
	Time          uint64   `json:"time"`
	Timezone      Timezone `json:"tz"`
	URL           string   `json:"url"`
	Detail        string   `json:"detail"`
	Felt          int      `json:"felt,omitempty"`
	CDI           float32  `json:"cdi,omitempty"`
	MMI           float32  `json:"mmi,omitempty"`
	Status        string   `json:"status"`
	Tsunami       int      `json:"tsunami"`
	Significance  int      `json:"sig"`
	Net           string   `json:"net"`
	Code          string   `json:"code"`
	IDS           string   `json:"ids"`
	Sources       string   `json:"sources"`
	Types         string   `json:"types"`
	NST           int      `json:"nst,omitempty"`
	DMin          float32  `json:"dmin"`
	RMS           float32  `json:"rms"`
	GAP           float32  `json:"gap"`
	MagnitudeType string   `json:"magType"`
	Type          string   `json:"type"`
}

type Geometry struct {
	Type        string     `json:"type"`
	Coordinates Coordinate `json:"coordinates"`
}

type Coordinate struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Depth     float32 `json:"depth"`
}

var (
	lBrace = []byte("[")
	rBrace = []byte("]")
	comma  = []byte(",")
)

var _ json.Unmarshaler = (*Coordinate)(nil)

type Client struct {
	sync.RWMutex
	httpClient *http.Client
	apiVersion string
}

type Request struct {
	Period    Period
	Magnitude Magnitude
}

func NewClient(options ...Option) (*Client, error) {
	client := new(Client)
	for _, opt := range options {
		opt.apply(client)
	}

	return client, nil
}

const defaultPeriod = PastDay
const defaultMagnitude = MAll

func (c *Client) Request(req *Request) (*Response, error) {
	if req == nil {
		req = new(Request)
	}

	urlToRequest := fullURL(req.Magnitude, req.Period)
	httpClient := c._httpClient()
	res, err := httpClient.Get(urlToRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if !statusOK(res.StatusCode) {
		return nil, fmt.Errorf("%s", res.Status)
	}

	blob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	recv := new(Response)
	if err := json.Unmarshal(blob, recv); err != nil {
		return nil, err
	}

	return recv, nil
}

func statusOK(code int) bool {
	return code >= 200 && code <= 299
}

var defaultHTTPClient = http.DefaultClient

func (c *Client) _httpClient() *http.Client {
	c.RLock()
	defer c.RUnlock()

	if c.httpClient != nil {
		return c.httpClient
	}
	return defaultHTTPClient
}

func fullURL(mag Magnitude, period Period) string {
	// First vet the period
	periodSummaryStr, known := period.summaryString()
	if !known {
		periodSummaryStr, _ = defaultPeriod.summaryString()
	}

	magnitudeSummaryStr, known := mag.summaryString()
	if !known {
		magnitudeSummaryStr, _ = defaultMagnitude.summaryString()
	}

	// Samples:
	// + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_hour.geojson
	// + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/1.0_hour.geojson
	// + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/2.5_week.geojson
	// + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/4.5_week.geojson
	// + https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/significant_month.geojson
	return fmt.Sprintf("https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/%s_%s.geojson",
		magnitudeSummaryStr, periodSummaryStr)
}

type BoundingBox struct {
	MinLatitude  float32 `json:"min_lat"`
	MinLongitude float32 `json:"min_lon"`
	MinDepth     float32 `json:"min_depth"`
	MaxLatitude  float32 `json:"max_lat"`
	MaxLongitude float32 `json:"max_lon"`
	MaxDepth     float32 `json:"max_depth"`
}

var _ json.Unmarshaler = (*BoundingBox)(nil)

func (bb *BoundingBox) UnmarshalJSON(b []byte) error {
	ptrsToSet := []*float32{
		0: &bb.MinLongitude,
		1: &bb.MinLatitude,
		2: &bb.MinDepth,
		3: &bb.MaxLongitude,
		4: &bb.MaxLatitude,
		5: &bb.MaxDepth,
	}
	return parseFloat32Ptrs(b, ptrsToSet...)
}

func (c *Coordinate) UnmarshalJSON(b []byte) error {
	ptrsToSet := []*float32{
		0: &c.Longitude,
		1: &c.Latitude,
		2: &c.Depth,
	}
	return parseFloat32Ptrs(b, ptrsToSet...)
}

func parseFloat32Ptrs(b []byte, ptrsToSet ...*float32) error {
	b = bytes.TrimSpace(b)
	b = bytes.TrimPrefix(b, lBrace)
	b = bytes.TrimSuffix(b, rBrace)

	splits := bytes.Split(b, comma)
	var cleaned [][]byte
	for _, split := range splits {
		cleaned = append(cleaned, bytes.TrimSpace(split))
	}

	for i, ptr := range ptrsToSet {
		if i >= len(cleaned) { // Done receiving values
			break
		}

		f64, err := strconv.ParseFloat(string(cleaned[i]), 32)
		if err != nil {
			return err
		}
		*ptr = float32(f64)
	}

	return nil
}
