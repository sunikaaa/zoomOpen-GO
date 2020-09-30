package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"sync"
	"time"
)

// ZoomData zoomデータデコード用
type ZoomData struct {
	ID   string `json:"id"`
	Pass string `json:"pass"`
}

func getTodayLecture(day int, hourMinute int, timeTable [][]ZoomData) (result ZoomData) {
	if day > 0 && day < 6 {
		zoom := timeTable[day-1]
		result = lectureTime(hourMinute, setLen5Array(zoom))
	}
	return result
}

func setLen5Array(timeTable []ZoomData) []ZoomData {
	for len(timeTable) < 5 {
		var initData ZoomData
		timeTable = append(timeTable, initData)
	}
	return timeTable
}

func lectureTime(hourMinute int, timeTable []ZoomData) (result ZoomData) {
	fmt.Print(hourMinute, timeTable)
	switch hourMinute {
	case 570:
		result = timeTable[0]
	case 670:
		result = timeTable[1]
	case 810:
		result = timeTable[2]
	case 930:
		result = timeTable[3]
	}

	return result
}

func nowTimeLoop(wg *sync.WaitGroup, zoomData [][]ZoomData) {
	for {

		var t = time.Now()
		day := int(t.Weekday())
		minute := t.Minute()
		hour := t.Hour()
		getData := getTodayLecture(day, hour*60+minute, zoomData)
		if getData.ID != "" {
			maincmd := fmt.Sprintf(`start zoommtg:"//zoom.us/join?confno=%v&pwd=%v"`, getData.ID, getData.Pass)
			cmd := exec.Command(`bash`, `-c`, maincmd)
			fmt.Print(cmd)
			err := cmd.Run()
			if err != nil {
				fmt.Println("err Occur", err.Error())
				break
			}
		}
		time.Sleep(60 * time.Second)
	}
	wg.Done()
}

func main() {

	bytes, err := ioutil.ReadFile("./zoomData.json")
	if err != nil {
		log.Fatal(err)
	}
	var zoomData = make([][]ZoomData, 5)
	if err := json.Unmarshal(bytes, &zoomData); err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go nowTimeLoop(&wg, zoomData)

	wg.Wait()

}
