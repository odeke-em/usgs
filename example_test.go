package usgs_test

import (
	"fmt"
	"log"

	"github.com/odeke-em/usgs"
)

func Example() {
	client, err := usgs.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	req := &usgs.Request{
		Period:    usgs.Past30Days,
		Magnitude: usgs.M4Dot5Plus,
	}

	details, err := client.Request(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("BoundingBox: %#v\n", details.BoundingBox)
	fmt.Printf("Metadata: %#v\n", details.Metadata)
	for i, feature := range details.Features {
		fmt.Printf("#%d: Id: %s Geometry: %#v\n", i, feature.Id, feature.Geometry)
	}
}
