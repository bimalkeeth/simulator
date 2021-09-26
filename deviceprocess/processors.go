package deviceprocess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"
)

var Client *http.Client

var NatsConn *nats.Conn

func init(){

	var natsURL = nats.DefaultURL
	if len(os.Args) == 2 {
		natsURL = os.Args[1]
	}
	// Connect to the NATS server.
	NatsConn, _ = nats.Connect(natsURL, nats.Timeout(5*time.Second))
}

func GetDeviceList() *[]Device {
	devices := make([]Device, 0)

	for _, group := range Conf.DevGroups {

		// for current group, append devices to the list
		for i := 1; i <= group.DeviceNum; i++ {
			devices = append(devices, Device{group.Prefix + strconv.Itoa(i), group.Firmware, false, group.IoTHub})
		}
	}
	return &devices
}

// DeviceSend function for use as a goroutine
func (d Device) DeviceSend(interval int, ch chan<- string) {
	// check if device exists
	if getResp, err := d.GetDevice(); err == nil {
		if getResp.StatusCode == 404 {
			// create the device because it does not exist
			if createResp, err := d.CreateDevice(); err == nil {
				// device was created
				ch <- fmt.Sprint("Device created with response code ", createResp.StatusCode)
			} else {
				// there was an error creating the device
				ch <- fmt.Sprintf("Could not create device %s. Error %v", d.Name, err)
			}

		}
	} else {
		// there was an error calling getDevice
		ch <- fmt.Sprintf("Could not get device %s. Error %v", d.Name, err)
	}

	for {
		temperature := 20.0 + rand.Float64()*10
		humidity := 40.0 + rand.Float64()*10
		message := DeviceMessage{Temperature: temperature, Humidity: humidity}

		err := d.sendData(message)
		if err == nil {
			ch <- fmt.Sprintf("Sent message from %s", d.Name)
		} else {
			// there was a send error
			ch <- fmt.Sprintf("Error sending message from %s. Error %v", d.Name, err)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (d Device) CreateDevice() (*http.Response, error) {

	reqBody, _ := json.Marshal(DeviceBody{DeviceID: d.Name})

	req, _ := http.NewRequest("PUT", "https://"+Conf.IoTHubs[d.IoTHub]+"/devices/"+d.Name+"?api-version=2016-02-03", bytes.NewBuffer(reqBody))

	// add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", Conf.SasTokens[d.IoTHub])

	// do the request
	resp, err := Client.Do(req)

	return resp, err
}

// GetDevice get device in IoT Hub
// return 200 if device is found, 404 if not
func (d Device) GetDevice() (*http.Response, error) {

	req, _ := http.NewRequest("GET", "https://"+Conf.IoTHubs[d.IoTHub]+"/devices/"+d.Name+"?api-version=2016-02-03", nil)

	// add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", Conf.SasTokens[d.IoTHub])

	// do the request
	resp, err := Client.Do(req)

	return resp, err
}

func (d Device) sendData(message DeviceMessage)  error {

	type SendingMessage struct {
		DeviceId string
		Message DeviceMessage
		DeviceName string
		Firmware string
	}

	sendMessage :=SendingMessage{
		Message: message,
		DeviceId:uuid.NewV4().String(),
		DeviceName: d.Name,
		Firmware: d.Firmware,
	}
	messageVal, _ := json.Marshal(sendMessage)

	fmt.Println(sendMessage)

	err:=NatsConn.Publish("message",messageVal)

    return err
}

func (d Device) DeleteDevice() (*http.Response, error) {

	//fmt.Println("https://" + Conf.IoTHubs[d.IoTHub] + "/devices/" + d.Name + "?api-version=2016-11-14")

	req, _ := http.NewRequest("DELETE", "https://"+ Conf.IoTHubs[d.IoTHub]+"/devices/"+d.Name+"?api-version=2016-11-14", nil)

	// add headers; DELETE requires If-Match with * for unconditional removal
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", Conf.SasTokens[d.IoTHub])
	req.Header.Add("If-Match", "*")

	// do the request
	resp, err := Client.Do(req)

	return resp, err
}
