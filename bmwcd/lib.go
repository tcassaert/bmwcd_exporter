package bmwcd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
)

func convertToEpoch(date string) float64 {
	dateString := fmt.Sprintf("%s-01T00:00:00+02:00", date)
	fullTime, err := time.Parse(time.RFC3339, dateString)

	if err != nil {
		log.Warnln("Failed to parse time format")
		return float64(0)
	}

	return float64(fullTime.Unix())
}

func getOAuthToken(username, password, region string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var url_encoded_data = "scope=authenticate_user+vehicle_data+remote_services&username=" + url.QueryEscape(username) + "&password=" + url.QueryEscape(password) + "&client_id=dbf0a542-ebd1-4ff0-a9a7-55172fbfce35&response_type=token&redirect_uri=https%3A%2F%2Fwww.bmw-connecteddrive.com%2Fapp%2Fstatic%2Fexternal-dispatch.html"

	var data = strings.NewReader(url_encoded_data)

	authUrl := getRegionUrls(region)[0]

	req, err := http.NewRequest("POST", authUrl, data)

	if err != nil {
		log.Errorln(err)
		return "", err
	}

	req.Header.Set("Authorization", "Basic blF2NkNxdHhKdVhXUDc0eGYzQ0p3VUVQOjF6REh4NnVuNGNEanliTEVOTjNreWZ1bVgya0VZaWdXUGNRcGR2RFJwSUJrN3JPSg==")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Credentials", "nQv6CqtxJuXWP74xf3CJwUEP:1zDHx6un4cDjybLENN3kyfumX2kEYigWPcQpdvDRpIBk7rOJ")

	resp, err := client.Do(req)

	if err != nil {
		log.Errorln(err)
		return "", err
	}

	headerMap := resp.Header
	for key, value := range headerMap {
		if key == "Location" {
			token_url := strings.Join(value, "")
			u, err := url.Parse(token_url)

			if err != nil {
				log.Errorln(err)
				return "", err
			}

			fragments, _ := url.ParseQuery(u.Fragment)
			if fragments["access_token"] != nil {
				access_token := strings.Join(fragments["access_token"], "")
				return access_token, nil
			} else {
				log.Errorln(err)
				return "", err
			}
		}
	}
	return "", nil
}

func getRegionUrls(region string) []string {
	if region == "rest_of_world" {
		urls := []string{"https://customer.bmwgroup.com/gcdm/oauth/authenticate", "b2vapi.bmwgroup.com"}
		return urls
	} else if region == "us" {
		urls := []string{"https://customer.bmwgroup.com/gcdm/usa/oauth/authenticate", "b2vapi.bmwgroup.us"}
		return urls
	} else if region == "cn" {
		urls := []string{"https://customer.bmwgroup.com/gcdm/oauth/authenticate", "b2vapi.bmwgroup.cn:8592"}
		return urls
	} else {
		log.Errorln("Unsupported region")
		os.Exit(1)
	}
	return nil
}

func getVehicleStatus(token, vin, region string) (string, error) {
	client := &http.Client{}
	apiUrl := getRegionUrls(region)[1]

	url := fmt.Sprintf("https://%s/webapi/v1/user/vehicles/%s/status", apiUrl, vin)
	req, err := http.NewRequest("GET", url, strings.NewReader(""))

	if err != nil {
		log.Errorln(err)
		return "", err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("referer", "https://www.bmw-connecteddrive.de/app/index.html")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)

	if err != nil {
		log.Errorln(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorln(err)
		return "", err
	}

	status := gjson.Get(string(body), "vehicleStatus")
	return status.String(), nil
}

func getVehicleVin(token, region string) (string, error) {
	client := &http.Client{}
	apiUrl := getRegionUrls(region)[1]

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/webapi/v1/user/vehicles", apiUrl), strings.NewReader(""))

	if err != nil {
		log.Errorln(err)
		return "", err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("referer", "https://www.bmw-connecteddrive.de/app/index.html")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)

	if err != nil {
		log.Errorln(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorln(err)
		return "", err
	}

	vin := gjson.Get(string(body), "vehicles.0.vin")
	return vin.String(), nil
}
