package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/lab5e/go-spanapi/v4"
)

func main() {
	// It's always a good idea to leave authentication tokens out of the source
	// code so we use a command line parameter here.
	token := ""
	flag.StringVar(&token, "token", "", "API token for the Span service")
	flag.Parse()

	if token == "" {
		fmt.Println("Missing token parameter")
		flag.PrintDefaults()
		return
	}

	config := spanapi.NewConfiguration()

	// Set this to true to list the requests and responses in the client. It can
	// be useful if you are wondering what is happening.
	config.Debug = false

	client := spanapi.NewAPIClient(config)

	// In the Real World this context should also have a context.WithTimeout
	// call to ensure we time out if there's no response.
	var reqCtx context.Context
	reqCtx, done := context.WithTimeout(context.Background(), 30*time.Second)
	defer done()

	keys := make(map[string]spanapi.APIKey)
	keys["APIToken"] = spanapi.APIKey{Key: token, Prefix: ""}
	ctx := context.WithValue(reqCtx, spanapi.ContextAPIKeys, keys)

	collections, _, err := client.CollectionsApi.ListCollections(ctx).Execute()
	if err != nil {
		fmt.Println("Error listing collections: ", err.Error())
		return
	}

	fmt.Println("Collections and devices")
	fmt.Println("=======================")

	for _, collection := range collections.Collections {
		fmt.Printf("Collection ID = %s\n", *collection.CollectionId)

		devices, _, err := client.DevicesApi.ListDevices(ctx, *collection.CollectionId).Execute()
		if err != nil {
			fmt.Println("Error listing devices: ", err.Error())
			continue
		}
		if devices.Devices != nil {
			for _, device := range devices.Devices {
				fmt.Printf("   Device ID = %s,  IMSI = %s,  IMEI = %s\n", device.GetDeviceId(), device.GetConfig().Ciot.GetImsi(), device.GetConfig().Ciot.GetImei())
			}
			fmt.Println(len(devices.Devices), " devices in collection")
		}
		fmt.Println()
	}
}
