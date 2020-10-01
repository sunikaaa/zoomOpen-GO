package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ZoomData zoomデータデコード用
type ZoomData struct {
	ID   string `json:"id"`
	Pass string `json:"pass"`
}

func getTodayLecture(day int, hourMinute int, timeTable WeekdayData, schedule []int) (result ZoomData) {
	var zoomData []ZoomData
	switch day {
	case 0:
		zoomData = timeTable.Sun
	case 1:
		zoomData = timeTable.Mon
	case 2:
		zoomData = timeTable.Tue

	case 3:
		zoomData = timeTable.Wed

	case 4:
		zoomData = timeTable.Thu
	case 5:
		zoomData = timeTable.Fri
	case 6:
		zoomData = timeTable.Sat
	}

	result = lectureTime(hourMinute, zoomData, schedule)
	return result
}

func lectureTime(hourMinute int, timeTable []ZoomData, schedule []int) (result ZoomData) {
	fmt.Println(timeTable)
	for i, time := range schedule {
		if time == hourMinute {
			result = timeTable[i]
		}
	}
	if result.ID == "" {
		log.Println("Zoom自動起動実行中です")
	}

	return result
}

// func formatTxt(zoomData [][]ZoomData) {
// 	for toDayLecture, youbi := range zoomData {
// 		for lecture,lectureNumber := range toDayLecture {
// 			switch youbi {
// 			case 0:

// 			}
// 		}
// 	}
// }

func nowTimeLoop(wg *sync.WaitGroup, zoomData WeekdayData, schedule []int) {
	for {

		var t = time.Now()
		day := int(t.Weekday())
		minute := t.Minute()
		hour := t.Hour()
		getData := getTodayLecture(day, hour*60+minute, zoomData, schedule)
		if getData.ID != "" {
			maincmd := fmt.Sprintf(`zoommtg:"//zoom.us/join?confno=%v&pwd=%v"`, getData.ID, getData.Pass)
			cmd := exec.Command(`rundll32.exe`, `url.dll,FileProtocolHandler`, maincmd)
			fmt.Print(cmd)
			err := cmd.Run()
			if err != nil {
				fmt.Println("err Occur", err.Error())
				time.Sleep(10 * time.Second)
				break
			}
		}
		time.Sleep(60 * time.Second)
	}
	wg.Done()
}

func main() {
	zoomData, schedule := loadCSV()

	// var zoomData = make([][]ZoomData, 5)
	// if err := json.Unmarshal(bytes, &zoomData); err != nil {
	// 	log.Fatal(err)
	// }
	var wg sync.WaitGroup
	wg.Add(1)
	go nowTimeLoop(&wg, zoomData, schedule)

	wg.Wait()

}

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

// WeekdayData は曜日ごとのズームを管理する構造体です。
type WeekdayData struct {
	Sun []ZoomData
	Mon []ZoomData
	Tue []ZoomData
	Wed []ZoomData
	Thu []ZoomData
	Fri []ZoomData
	Sat []ZoomData
}

func appendWeekdayData(zoomdatas []ZoomData, isID bool, IDorPASS string) []ZoomData {
	if isID {
		zoomdatas = append(zoomdatas, ZoomData{ID: IDorPASS})
	} else {
		zoomdatas[len(zoomdatas)-1].Pass = IDorPASS
	}
	return zoomdatas
}

func loadCSV() (WeekdayData, []int) {
	file1, err := os.Open("./zoomData.csv")
	// failOnError(err)
	if err != nil {
		log.Print("zoomData.csvが同じディレクトリに存在しません。")
		time.Sleep(time.Second * 3)
		log.Fatal(err)
	}
	reader := csv.NewReader(file1)
	reader.LazyQuotes = true

	log.Printf("zoomdata.csv Loading")
	weekdayData := WeekdayData{}
	var schedule []int
	for row := 0; ; row++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if row == 0 {
			continue
		} else {
			failOnError(err)
		}

		isID := false

		for i, v := range record {
			// newRecord = append(newRecord, v)
			switch i {
			case 1:
				if v == "" {
					continue
				}
				timeSliceStr := strings.Split(v, ":")
				timeSlice := make([]int, 2)
				for i, v := range timeSliceStr {

					startTime, err := strconv.Atoi(v)
					failOnError(err)
					timeSlice[i] = startTime
				}
				schedule = append(schedule, timeSlice[0]*60+timeSlice[1])
			case 2:
				if v == "id" {
					isID = true
				}
			case 3:
				weekdayData.Sun = appendWeekdayData(weekdayData.Sun, isID, v)
			case 4:
				weekdayData.Mon = appendWeekdayData(weekdayData.Mon, isID, v)
			case 5:
				weekdayData.Tue = appendWeekdayData(weekdayData.Tue, isID, v)
			case 6:
				weekdayData.Wed = appendWeekdayData(weekdayData.Wed, isID, v)
			case 7:
				weekdayData.Thu = appendWeekdayData(weekdayData.Thu, isID, v)
			case 8:
				weekdayData.Fri = appendWeekdayData(weekdayData.Fri, isID, v)
			case 9:
				weekdayData.Sat = appendWeekdayData(weekdayData.Sat, isID, v)
			}
		}

	}
	if len(weekdayData.Sun) != len(schedule) {
		log.Fatal("不正なCSVデータです。")
	}
	log.Println("finish Load")
	fmt.Println("日", weekdayData.Sun)
	fmt.Println("月", weekdayData.Mon)
	fmt.Println("火", weekdayData.Tue)
	fmt.Println("水", weekdayData.Wed)
	fmt.Println("木", weekdayData.Thu)
	fmt.Println("金", weekdayData.Fri)
	fmt.Println("土", weekdayData.Sat)
	fmt.Println("開始時刻")
	for i, v := range schedule {
		fmt.Printf("%d限%d時%d分\n", i+1, v/60, v%60)
	}
	return weekdayData, schedule
}
