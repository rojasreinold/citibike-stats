package main

/*
StationName, StationSize, month1 usage, month2 usage, etc

*/
import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
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
	csvHeader := []string{
		"station_name",
	}

	records := make(map[string][]int)

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
				fmt.Println(baseMonths)
				val = baseMonths[:]
				records[stationName] = append(val, stationUsage)

			}
		}
		baseMonths = append(baseMonths, -1)
		fmt.Println(records)
	}

	fmt.Println(csvHeader)

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

		if !(strings.HasPrefix(emp.startStationId, "JC")) {
			stationsUsage[emp.startStationName] = stationsUsage[emp.startStationName] + 1
		}

		if !(strings.HasPrefix(emp.endStationId, "JC")) {
			stationsUsage[emp.endStationName] = stationsUsage[emp.endStationName] + 1
		}

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
