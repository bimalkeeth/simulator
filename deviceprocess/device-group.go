package deviceprocess

// DevGroup define a struct for a device group (used in configuration struct)
type DevGroup struct {
	Prefix    string // prefix for device name
	DeviceNum int    // number of devices with this prefix
	Firmware  string // firmware version like 1.20
	IoTHub    int    // index into IoTHubs slice in configuration struct

}

type Device struct {
	Name     string // name of the device (deviceprefix + number)
	Firmware string // firmware of device (results in other data to be sent)
	InHub    bool   // is the device in IoT Hub
	IoTHub   int    // index into IoTHubs slice in configuration object
}

type DeviceBody struct {
	DeviceID string
}

type DeviceMessage struct {
	Temperature float64
	Humidity    float64
}
