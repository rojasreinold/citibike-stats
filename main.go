package main

/*
StationName, StationSize, month1 trips, month2 trips, etc

*/
import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	csvHeader := []string{
		"station_name",
	}

	stationDockSizes := getStationDockCounts()

	records := make(map[string][]int)
	csvHeader = append(csvHeader, "dock_capacity")

	for stationName, stationDockSize := range stationDockSizes {
		records[stationName] = []int{stationDockSize}
	}

	files, err := ioutil.ReadDir("data/")
	if err != nil {
		log.Fatal(err)
	}
	var baseMonths []int

	for _, file := range files {
		//Skip hidden files and directories
		if file.Name()[0] == '.' || file.IsDir() {
			continue
		}

		csvHeader = append(csvHeader, strings.Split(file.Name(), "-")[0])
		monthUsages := getStationsUsages("data/" + file.Name())

		for stationName, stationUsage := range monthUsages {
			val, ok := records[stationName]
			if ok {
				records[stationName] = append(val, stationUsage)
			} else {
				val = baseMonths[:]
				records[stationName] = append(val, stationUsage)

			}
		}

		// If a station is now closed, set trips to -1
		for stationName, record := range records {
			if len(record) < len(csvHeader)-1 {
				records[stationName] = append(record, -1)
			}

		}
		baseMonths = append(baseMonths, -1)
		fmt.Println("Parsing file data/" + file.Name())
	}

	csvRecords := [][]string{
		csvHeader,
	}

	fmt.Println("Aggregating data...")
	for stationName, stationUsage := range records {
		var stationUsageS []string

		for _, usage := range stationUsage {
			stationUsageS = append(stationUsageS, strconv.Itoa(usage))
		}

		csvRecords = append(csvRecords, append([]string{stationName}, stationUsageS...))

	}

	// Write data to new csv
	outputFileName := "citibike-stats-aggregate-" + strings.Split(files[0].Name(), "-")[0]
	outputFileName = outputFileName + "-" + strings.Split(files[len(files)-1].Name(), "-")[0] + ".csv"
	fmt.Println("Saving to file " + outputFileName)
	outputFile, err := os.Create(outputFileName)

	if err != nil {
		log.Fatalln("Failed to open file", err)
	}

	w := csv.NewWriter(outputFile)
	w.WriteAll(csvRecords)
}

type empData struct {
	rideId           string
	rideType         string
	startTime        string
	endTime          string
	startStationName string
	startStationId   string
	endStationName   string
	endStationId     string
	startLat         string
	startLng         string
	endLat           string
	endLng           string
	memberCasual     string
}

func getStationsUsages(filename string) map[string]int {
	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	stationsUsage := map[string]int{}
	for _, line := range csvLines[1:] {
		emp := empData{
			rideId:           strings.TrimSpace(line[0]),
			rideType:         strings.TrimSpace(line[1]),
			startTime:        strings.TrimSpace(line[2]),
			endTime:          strings.TrimSpace(line[3]),
			startStationName: strings.TrimSpace(line[4]),
			startStationId:   strings.TrimSpace(line[5]),
			endStationName:   strings.TrimSpace(line[6]),
			endStationId:     strings.TrimSpace(line[7]),
			startLat:         strings.TrimSpace(line[8]),
			startLng:         strings.TrimSpace(line[9]),
			endLat:           strings.TrimSpace(line[10]),
			endLng:           strings.TrimSpace(line[11]),
			memberCasual:     strings.TrimSpace(line[12]),
		}

		// Log only the nyc stations
		if !(strings.HasPrefix(emp.startStationId, "JC")) {
			stationsUsage[emp.startStationName] = stationsUsage[emp.startStationName] + 1
		}

		if !(strings.HasPrefix(emp.endStationId, "JC")) {
			stationsUsage[emp.endStationName] = stationsUsage[emp.endStationName] + 1
		}

	}

	return stationsUsage
}

func getStationDockCounts() map[string]int {
	resp, err := http.Get("https://gbfs.citibikenyc.com/gbfs/en/station_information.json")

	if err != nil {
		log.Fatalln("Couldn't get station dock counts ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)

	var jsonSB jsonData
	err = json.Unmarshal([]byte(sb), &jsonSB)

	if err != nil {
		log.Fatalln("JSON Decode error ", err)
	}

	stationDockSizes := map[string]int{}

	for _, station := range jsonSB.Data.Stations {
		stationDockSizes[station.Name] = station.Capacity
	}
	return stationDockSizes
}

type jsonData struct {
	Data struct {
		Stations []struct {
			LegacyID              string `json:"legacy_id"`
			EightdStationServices []any  `json:"eightd_station_services"`
			RentalUris            struct {
				Android string `json:"android"`
				Ios     string `json:"ios"`
			} `json:"rental_uris"`
			RentalMethods               []string `json:"rental_methods"`
			StationID                   string   `json:"station_id"`
			HasKiosk                    bool     `json:"has_kiosk"`
			Lat                         float64  `json:"lat"`
			ShortName                   string   `json:"short_name"`
			RegionID                    string   `json:"region_id,omitempty"`
			Lon                         float64  `json:"lon"`
			ElectricBikeSurchargeWaiver bool     `json:"electric_bike_surcharge_waiver"`
			Name                        string   `json:"name"`
			StationType                 string   `json:"station_type"`
			ExternalID                  string   `json:"external_id"`
			Capacity                    int      `json:"capacity"`
			EightdHasKeyDispenser       bool     `json:"eightd_has_key_dispenser"`
		} `json:"stations"`
	} `json:"data"`
	LastUpdated int `json:"last_updated"`
	TTL         int `json:"ttl"`
}
