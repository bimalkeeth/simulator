package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"simulator/deviceprocess"
)

var RemoveDevices = flag.Bool("r", false, "Remove devices from IoT Hub")

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// create a channel to allow goroutines to communicate
	// in this case goroutines will communicate with main() which is also a goroutine
	ch := make(chan string)

	// http client
	deviceprocess.Client = &http.Client{
		Timeout: time.Second * 10,
	}

	// parse flags
	flag.Parse()

	// read configuration into conf struct
	config, err := deviceprocess.GetConf("config.json")

	// store config in global Conf variable
	deviceprocess.Conf = config

	if err != nil {
		// could be due to file not found or some JSON parsing error
		fmt.Println("Error occurred: ", err)
		os.Exit(1)
	}

	// get device list (returns pointer) created from the configuration
	devList := deviceprocess.GetDeviceList()

	if *RemoveDevices {
		for _, device := range *devList {
			deleteResp, err := device.DeleteDevice()
			if err != nil {
				fmt.Printf("Could not delete device %s. Error: %v\n", device.Name, err)
			} else {
				fmt.Printf("Deleted device %s with reponse code %d\n", device.Name, deleteResp.StatusCode)
			}
		}
		os.Exit(0)
	}

	for _, device := range *devList {
		// for each device, run the deviceSend method as a goroutine
		go device.DeviceSend(deviceprocess.Conf.Interval, ch)
	}

	// not sure if this is good practice, but we keep reading from
	// the channel forever to pick up messages from the goroutines
	for {
		fmt.Println(<-ch)
	}

}
