package deviceprocess

import (
	"encoding/json"
	"os"
)

// Configuration define a struct for configuration
type Configuration struct {
	Interval  int        // interval to send in seconds
	IoTHubs   []string   // slice of strings to hold IoTHubs to send to
	SasTokens []string   // slice of strings to hold SaS tokens
	DevGroups []DevGroup // slice of devgroup
}

// Conf is global configuration variable
var Conf Configuration


// GetConf read the configuration from configuration file
func GetConf(configFile string) (Configuration, error) {
	file, _ := os.Open(configFile)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	decoder := json.NewDecoder(file)
	Conf := Configuration{}
	err := decoder.Decode(&Conf)

	// in case of err Conf will be empty struct
	return Conf, err

}