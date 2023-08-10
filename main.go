package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	getSerialURLTemplate = "http://api.morphix.online:8080/api/v1/device_by_serial/%s"
	postURLTemplate      = "http://api.morphix:8080/api/v1/stat/%s"
	serialFilePath       = "/etc/thoth/serial.id"
	defaultLoopInterval  = 300 * time.Second // 300 seconds (5 minutes)
)

func main() {
	var loopInterval time.Duration
	flag.DurationVar(&loopInterval, "interval", defaultLoopInterval, "Loop interval duration")
	flag.Parse()
	for {
		err := processLoop()
		if err != nil {
			fmt.Println("Error in loop:", err)
		}

		time.Sleep(loopInterval)
	}
}

func processLoop() error {
	serial, err := readSerial(serialFilePath)
	if err != nil {
		return fmt.Errorf("error reading serial: %v", err)
	}

	deviceID, err := getDeviceID(serial)
	if err != nil {
		return fmt.Errorf("error getting device ID: %v", err)
	}

	postData := map[string]interface{}{}

	err = sendPOSTRequest(deviceID, postData)
	if err != nil {
		return fmt.Errorf("error sending POST request: %v", err)
	}

	return nil
}

func readSerial(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getDeviceID(serial string) (string, error) {
	serial = strings.TrimSpace(serial)
	getURL := fmt.Sprintf(getSerialURLTemplate, serial)
	resp, err := http.Get(getURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(respBody, &responseData)
	if err != nil {
		return "", err
	}

	result, ok := responseData["result"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid result in GET response")
	}

	deviceID, ok := result["device_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid device_id in GET response")
	}

	return deviceID, nil
}

func sendPOSTRequest(deviceID string, data interface{}) error {
	postURL := fmt.Sprintf(postURLTemplate, deviceID)
	postPayload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	postResp, err := http.Post(postURL, "application/json", bytes.NewBuffer(postPayload))
	if err != nil {
		return err
	}
	defer postResp.Body.Close()

	postRespBody, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		return err
	}

	fmt.Println("POST Response:", string(postRespBody))
	return nil
}
