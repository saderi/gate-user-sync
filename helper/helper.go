package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	stdLog = log.New(os.Stdout, "", 0)
)

func Log(v ...interface{}) {
	stdLog.Println(green(time.Now().Format("02 01 2006 15:04:05")), fmt.Sprint(v...))
}

func SendAPIRequest(method string, url string, bearerToken string, data interface{}) ([]byte, error) {
	var req *http.Request
	var err error

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API %s request failed with url %s status code %d: %s", method, url, resp.StatusCode, string(body))
	}

	return body, nil
}

func ToMap(slice []string) map[string]bool {
	result := make(map[string]bool)

	for _, item := range slice {
		result[item] = true
	}

	return result
}

func FindPhones(phones []string, alternativePhones map[string]bool, shouldBeIn bool) []string {
	var result []string

	for _, phone := range phones {
		_, exists := alternativePhones[phone]
		if exists == shouldBeIn {
			result = append(result, phone)
		}
	}

	return result
}

func FindToBeRemoved(passivePhones []string, alternativePhones []string) []string {
	alternativePhoneMap := ToMap(alternativePhones)

	return FindPhones(passivePhones, alternativePhoneMap, true)
}

func FindToBeAdded(activePhones []string, alternativePhones []string) []string {
	alternativePhoneMap := ToMap(alternativePhones)

	return FindPhones(activePhones, alternativePhoneMap, false)
}
