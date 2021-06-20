package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/forestgiant/sliceutil"
)

type CowinResp struct {
	Centres []Centre `json:"centers"`
}

type Centre struct {
	Id       int `json:"center_id"`
	Name     string
	Address  string
	State    string `json:"state_name"`
	District string `json:"district_name"`
	Pincode  int64
	Sessions []Session `json:"sessions"`
}

type Session struct {
	AvailableCapacity int `json:"available_capacity"`
	AgeLimit          int `json:"min_age_limit"`
	Vaccine           string
	Dose1Availability int `json:"available_capacity_dose1"`
	Dose2Availability int `json:"available_capacity_dose2"`
	Date              string
}

const vaccineSlotFetchAPI = "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict?district_id=%d&date=%s"

var reader = bufio.NewReader(os.Stdin)

//phc harahua, phc shivpur, aryuvedic college
var centerList = []int{692744, 608946, 596524}

var districtCode = 696
var age, doseIdentifier int64

func main() {
	for {
		fmt.Println("Enter Age limit you want to search for")
		_, err := fmt.Scanf("%d", &age)
		if err == nil {
			break
		}
		log.Println("Insert numeric values for age")
	}

	fmt.Scanf("%v")
	for {
		fmt.Println("Input 1 for Dose 1 search, 2 for Dose 2 Search")
		_, err := fmt.Scanf("%d", &doseIdentifier)
		if err == nil {
			break
		}
		log.Println("Insert numeric values for dose")
	}
	cowinRespBody := sendRequest(districtCode)
	filteredResp := filterResponse(cowinRespBody, int(age), centerList, int(doseIdentifier))
	if len(filteredResp) == 0 {
		log.Println("No Vaccine Session Found")
	}
}

func sendRequest(districtId int) CowinResp {
	dt := time.Now()
	finalAPI := fmt.Sprintf(vaccineSlotFetchAPI, districtId, dt.Format("02-01-2006"))
	client := &http.Client{}
	req, _ := http.NewRequest("GET", finalAPI, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	respBody := CowinResp{}
	json.Unmarshal(body, &respBody)
	return respBody
}

func filterResponse(resp CowinResp, age int, centerCodesList []int, doseIdentifier int) []Centre {
	filterArray := make([]Centre, 0)
	for _, centreObj := range resp.Centres {
		if sliceutil.Contains(centerCodesList, centreObj.Id) {
			for _, session := range centreObj.Sessions {
				if session.AgeLimit <= age {
					if (doseIdentifier == 1 && session.Dose1Availability > 0) || (doseIdentifier == 2 && session.Dose2Availability > 0) {
						log.Println(fmt.Sprintf("Found Session at %s on %s", centreObj.Address, session.Date))
						filterArray = append(filterArray, centreObj)
					}
				}
			}
		}
	}
	return filterArray
}
