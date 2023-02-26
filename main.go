package main

/*
StationName, StationSize, month1 usage, month2 usage, etc

*/
import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
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

	stationsUsage := getStationsUsages()
	sortedStations := sortStationsByUsage(stationsUsage)

	for _, stationStats := range sortedStations {
		fmt.Println(stationStats.Key + "," + strconv.Itoa(stationStats.Value))
	}
}

func getStationsUsages() map[string]int {
	//csvFile, err := os.Open("data/2022d06-citbike-tripdata.csv")
	csvFile, err := os.Open("data/202301-citibike-tripdata.csv")
	//csvFile, err := os.Open("data/202301-citibike-tripdata-reduced.csv")
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
		stationsUsage[emp.startStationName] = stationsUsage[emp.startStationName] + 1
		stationsUsage[emp.endStationName] = stationsUsage[emp.endStationName] + 1

	}

	return stationsUsage
}

type kv struct {
	Key   string
	Value int
}

func sortStationsByUsage(stationsUsage map[string]int) []kv {

	var sortedStations []kv
	for k, v := range stationsUsage {
		sortedStations = append(sortedStations, kv{k, v})
	}

	sort.Slice(sortedStations, func(i, j int) bool {
		return sortedStations[i].Value > sortedStations[j].Value
	})

	return sortedStations
}
