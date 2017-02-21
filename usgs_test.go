package usgs_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/odeke-em/usgs"
)

type localFetcher struct {
}

var _ http.RoundTripper = (*localFetcher)(nil)

func (lf *localFetcher) RoundTrip(req *http.Request) (*http.Response, error) {
	splits := strings.Split(req.URL.Path, "/")
	if len(splits) < 1 {
		return nil, fmt.Errorf("no path passed in")
	}

	filename := splits[len(splits)-1]
	localPath := filepath.Join("./testdata", filename)
	f, err := os.Open(localPath)
	if err != nil {
		return nil, err
	}

	resp := &http.Response{Body: f, StatusCode: 200, Status: "200 OK"}
	return resp, nil
}

var customHTTPClient = &http.Client{Transport: &localFetcher{}}

func TestUSGSFetch(t *testing.T) {
	client, err := usgs.NewClient(usgs.WithHTTPClient(customHTTPClient))
	if err != nil {
		t.Fatalf("initializing client: %v", err)
	}

	tests := [...]struct {
		req  *usgs.Request
		want string
	}{
		0: {
			req: &usgs.Request{
				Magnitude: usgs.M2Dot5Plus,
				Period:    usgs.Past30Days,
			},
			want: "./testdata/2.5_month.geojson",
		},
	}

	for i, tt := range tests {
		res, err := client.Request(tt.req)
		if err != nil {
			t.Errorf("#%d err: %v", i, err)
			continue
		}

		gotBlob, err := json.Marshal(res)
		if err != nil {
			t.Errorf("#%d jsonMarshal: %v", i, err)
			continue
		}
		wantBlob, err := ioutil.ReadFile(tt.want)
		if err != nil {
			t.Errorf("#%d readWantBlob: %v", i, err)
			continue
		}
		// For proper comparison we need to unmarshal then remarshal
		want := new(usgs.Response)
		if err := json.Unmarshal(wantBlob, want); err != nil {
			t.Errorf("#%d: re-unmarshal err: %v", i, err)
			continue
		}
		wantBlob, err = json.Marshal(want)
		if err != nil {
			t.Errorf("#%d failed to re-remarshal blob: %v", i, err)
			continue
		}
		if !bytes.Equal(gotBlob, wantBlob) {
			t.Errorf("#%d:\ngot:\n\t:%s\nwant:\n\t:%s", i, gotBlob, wantBlob)
		}
	}
}
