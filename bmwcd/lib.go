package bmwcd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	}

	return float64(fullTime.Unix())
}

func getOAuthToken(username, password string) string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var url_encoded_data = "scope=authenticate_user+vehicle_data+remote_services&username=" + url.QueryEscape(username) + "&password=" + url.QueryEscape(password) + "&client_id=dbf0a542-ebd1-4ff0-a9a7-55172fbfce35&response_type=token&redirect_uri=https%3A%2F%2Fwww.bmw-connecteddrive.com%2Fapp%2Fstatic%2Fexternal-dispatch.html"

	var data = strings.NewReader(url_encoded_data)

	req, err := http.NewRequest("POST", "https://customer.bmwgroup.com/gcdm/oauth/authenticate", data)

	if err != nil {
		log.Errorln(err)
	}

	req.Header.Set("Authorization", "Basic blF2NkNxdHhKdVhXUDc0eGYzQ0p3VUVQOjF6REh4NnVuNGNEanliTEVOTjNreWZ1bVgya0VZaWdXUGNRcGR2RFJwSUJrN3JPSg==")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Credentials", "nQv6CqtxJuXWP74xf3CJwUEP:1zDHx6un4cDjybLENN3kyfumX2kEYigWPcQpdvDRpIBk7rOJ")

	resp, err := client.Do(req)

	if err != nil {
		log.Errorln(err)
	}

	headerMap := resp.Header
	for key, value := range headerMap {
		if key == "Location" {
			token_url := strings.Join(value, "")
			u, err := url.Parse(token_url)

			if err != nil {
				log.Errorln(err)
			}

			fragments, _ := url.ParseQuery(u.Fragment)
			if fragments["access_token"] != nil {
				access_token := strings.Join(fragments["access_token"], "")
				return access_token
			} else {
				log.Errorln(err)
			}
		}
	}
	return ""
}

func getVehicleStatus(token, vin string) string {
	client := &http.Client{}
	url := fmt.Sprintf("https://b2vapi.bmwgroup.com/webapi/v1/user/vehicles/%s/status", vin)
	req, err := http.NewRequest("GET", url, strings.NewReader(""))

	if err != nil {
		log.Errorln(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("referer", "https://www.bmw-connecteddrive.de/app/index.html")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)

	if err != nil {
		log.Errorln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorln(err)
	}

	status := gjson.Get(string(body), "vehicleStatus")
	return status.String()
}

func getVehicleVin(token string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://b2vapi.bmwgroup.com/webapi/v1/user/vehicles", strings.NewReader(""))

	if err != nil {
		log.Errorln(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("referer", "https://www.bmw-connecteddrive.de/app/index.html")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)

	if err != nil {
		log.Errorln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorln(err)
	}

	vin := gjson.Get(string(body), "vehicles.0.vin")
	return vin.String()
}

func StartPolling(username, password string) {
	vin := getVehicleVin(getOAuthToken(username, password))

	for {
		token := string(getOAuthToken(username, password))
		status := getVehicleStatus(token, vin)
		go jsonToProm(status)
		time.Sleep(300 * time.Second)
	}
}
