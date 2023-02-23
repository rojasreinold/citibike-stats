package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func main() {
	csvFile, err := os.Open("data/202301-citibike-tripdata.csv")
	//csvFile, err := os.Open("data/202301-citibike-tripdata-reduced.csv")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("opened csv")
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
		stationsUsage[emp.startStationName+"-"+emp.startStationId] = stationsUsage[emp.startStationName+"-"+emp.startStationId] + 1
		stationsUsage[emp.endStationName+"-"+emp.endStationId] = stationsUsage[emp.endStationName+"-"+emp.endStationId] + 1

	}
	// fmt.Println(stationsUsage)
	for stationName, stationUsage := range stationsUsage {
		fmt.Println(stationName + ":" + strconv.Itoa(stationUsage))
	}
}

//func getStationsUsage()  {}
