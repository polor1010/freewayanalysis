package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ByAverageSpeed []Location

func (locations ByAverageSpeed) Len() int { return len(locations) }

func (locations ByAverageSpeed) Swap(i, j int) {
	locations[i], locations[j] = locations[j], locations[i]
}

func (locations ByAverageSpeed) Less(i, j int) bool { return locations[i].Speed1 < locations[j].Speed1 }

//
//回傳所有路段的車速
func GetAll(date string) []Location {

	var locations []Location

	fileName := GetFileName(date)

	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Println("Error:", err)
	} else {

		lines := strings.Split(string(content), "\r\n")

		for i := 1; i < len(lines)-1; i++ {

			var location Location
			location.TimeStamp = strings.Split(lines[i], ",")[0]
			location.FreewayID = strings.Split(lines[i], ",")[1]
			location.LocationID = strings.Split(lines[i], ",")[2]
			location.Direction1, _ = strconv.Atoi(strings.Split(lines[i], ",")[3])
			location.Direction2, _ = strconv.Atoi(strings.Split(lines[i], ",")[5])
			speed1, _ := strconv.ParseFloat(strings.Split(lines[i], ",")[4], 64)
			speed2, _ := strconv.ParseFloat(strings.Split(lines[i], ",")[6], 64)
			location.Speed1 = int(speed1 + 0.5)
			location.Speed2 = int(speed2 + 0.5)

			locations = append(locations, location)

		}
	}

	//SpeedMapJson, _ := json.Marshal(locations)
	//fmt.Println(string(SpeedMapJson))

	return locations

}

//
//回傳目前某個路段的過去一個月的速度歷史資料
func GetMonthByLocationID(date string, location Location) SpeedChart {

	var t time.Time
	if date != "" {
		t, _ = time.Parse("200601021504", date)
		t = t.UTC()
	} else {
		t = time.Now().UTC()
	}

	var results []string

	var index = t.Minute() / 5
	var hour = t.Hour()

	var speedDays []SpeedTime

	for i := 0; i < 30; i++ {
		var speedDay SpeedTime
		t = t.Add(-time.Hour * 24)

		fileName := fmt.Sprintf("%s%d%.2d%.2d%.2d%.2d.csv", ROOT_PATH, t.Year(), t.Month(), t.Day(), hour, index*5)

		content, err := ioutil.ReadFile(fileName)

		if err != nil {
			fmt.Println("Error:", err)
		} else {

			lines := strings.Split(string(content), "\r\n")

			for i := 1; i < len(lines)-1; i++ {

				//freewayID := strings.Split(lines[i], ",")[1]
				locationID := strings.Split(lines[i], ",")[2]

				//(freewayID == location.FreewayID) &&
				if locationID == location.LocationID {

					//speedDay.Direction1, _ = strconv.Atoi(strings.Split(lines[i], ",")[3])
					//speedDay.Direction2, _ = strconv.Atoi(strings.Split(lines[i], ",")[5])
					speed1, _ := strconv.ParseFloat(strings.Split(lines[i], ",")[4], 64)
					speed2, _ := strconv.ParseFloat(strings.Split(lines[i], ",")[6], 64)
					speedDay.Speed1 = int(speed1 + 0.5)
					speedDay.Speed2 = int(speed2 + 0.5)
					results = append(results, lines[i])

				}
			}
		}
		speedDay.Time = strconv.Itoa(int(t.Year())) + "/" + strconv.Itoa(int(t.Month())) + "/" + strconv.Itoa(t.Day())
		speedDays = append(speedDays, speedDay)
	}

	//fmt.Println(results)

	speedChartData := SpeedChart{
		LocationID: location.LocationID,
		TimeRange:  30,
		Data:       speedDays}

	//SpeedChartJson, _ := json.Marshal(speedChartData)
	//fmt.Println(string(SpeedChartJson))

	return speedChartData
}

//
//回傳某一路段一整天的車速資料,會包含(過去|預測）
func GetDetailByLocationID(date string, locations Location) []string {

	var t time.Time

	//UTC time
	if date != "" {
		t, _ = time.Parse("200601021504", date)
	} else {
		t = time.Now().UTC()
	}

	day := fmt.Sprintf("%d%.2d%.2d%.2d%.2d", t.Year(), t.Month(), t.Day(), 0, 0)
	t2, _ := time.Parse("200601021504", day)
	t2 = t2.Add(time.Hour*time.Duration(-8) + time.Minute*time.Duration(+5)) //限定台灣時區
	fmt.Println(t2)
	var results []string

	for h := 0; h < 24; h++ {
		for m := 0; m < 60; m += 5 {

			var fileName string

			sum1 := t.Hour()*60 + t.Minute()/5*5
			sum2 := h*60 + m

			if sum1 > sum2 {
				fileName = fmt.Sprintf("%s%d%.2d%.2d%.2d%.2d.csv", ROOT_PATH, t2.Year(), t2.Month(), t2.Day(), t2.Hour(), m)
				//fmt.Println(fileName)
			} else {
				fileName = fmt.Sprintf("%spredict/%d%.2d%.2d%.2d%.2d_.csv", ROOT_PATH, t2.Year(), t2.Month(), t2.Day(), t2.Hour(), m)
				//fmt.Println(fileName)
			}

			content, err := ioutil.ReadFile(fileName)

			if err != nil {
				fmt.Println("Error:", err)
				s := fmt.Sprintf("%2d:%2d %s s1:%d s2:%d tm:%d\r\n", h, m, fileName, sum1, sum2, t.Minute())
				results = append(results, s)
			} else {

				lines := strings.Split(string(content), "\r\n")

				locationID := locations.LocationID
				//freewayID := locations.FreewayID

				for k := 1; k < len(lines)-1; k++ {

					if locationID == strings.Split(lines[k], ",")[2] { //freewayID == strings.Split(lines[k], ",")[1] { //&& {
						s := fmt.Sprintf("%s,%2d:%2d,%s\r\n", lines[k], h, m, fileName)
						results = append(results, s)
					}
				}
			}
		}

		t2 = t2.Add(time.Hour)
		//fmt.Println(t2, count, speed1Sum, speed2Sum)
	}
	return results
}

//
//回傳某一路段一整天的車速資料,會包含(過去|預測）
func GetDayByLocationID(date string, locations Location) SpeedChart {

	var t time.Time

	//UTC time
	if date != "" {
		t, _ = time.Parse("200601021504", date)
	} else {
		t = time.Now().UTC()
	}

	day := fmt.Sprintf("%d%.2d%.2d%.2d%.2d", t.Year(), t.Month(), t.Day(), 0, 0)
	t2, _ := time.Parse("200601021504", day)
	t2 = t2.Add(time.Hour*time.Duration(-8) + time.Minute*time.Duration(+5)) //限定台灣時區
	//fmt.Println(t2)
	var searchResults []string
	var speedHours []SpeedTime
	var freewayID string
	var direction string

	for h := 0; h < 24; h++ {

		var speedHour SpeedTime
		var speed1Sum float64 = 0.0
		var speed2Sum float64 = 0.0
		var count int = 0
		for m := 0; m < 60; m += 5 {

			var fileName string

			sum1 := t.Hour()*60 + t.Minute()/5*5
			sum2 := h*60 + m

			if sum1 > sum2 {
				fileName = fmt.Sprintf("%s%d%.2d%.2d%.2d%.2d.csv", ROOT_PATH, t2.Year(), t2.Month(), t2.Day(), t2.Hour(), m)
				//fmt.Println(fileName)
			} else {
				fileName = fmt.Sprintf("%spredict/%d%.2d%.2d%.2d%.2d_.csv", ROOT_PATH, t2.Year(), t2.Month(), t2.Day(), t2.Hour(), m)
				//fmt.Println(fileName)
			}

			content, err := ioutil.ReadFile(fileName)

			if err != nil {
				fmt.Println("Error:", err)
			} else {

				lines := strings.Split(string(content), "\r\n")

				locationID := locations.LocationID
				//freewayID := locations.FreewayID

				for k := 1; k < len(lines)-1; k++ {

					if locationID == strings.Split(lines[k], ",")[2] { //freewayID == strings.Split(lines[k], ",")[1] { //&& {
						freewayID = strings.Split(lines[k], ",")[1]
						direction = strings.Split(lines[k], ",")[3]
						fmt.Println(lines[k])
						speed1, _ := strconv.ParseFloat(strings.Split(lines[k], ",")[4], 64)
						speed2, _ := strconv.ParseFloat(strings.Split(lines[k], ",")[6], 64)
						speed1Sum += speed1
						speed2Sum += speed2
						searchResults = append(searchResults, lines[k])
						count++
					}
				}
			}
			speedHour.Time = strconv.Itoa(h) + ":00"

		}

		t2 = t2.Add(time.Hour)
		//fmt.Println(t2, count, speed1Sum, speed2Sum)

		if count > 0 {
			speedHour.Speed1 = int(speed1Sum/float64(count) + 0.5)
			speedHour.Speed2 = int(speed2Sum/float64(count) + 0.5)
		} else {
			speedHour.Speed1 = 0
			speedHour.Speed2 = 0
		}
		fmt.Println("average:", speedHour.Speed1, speedHour.Speed2)

		speedHours = append(speedHours, speedHour)
	}

	speedChartData := SpeedChart{
		LocationID: locations.LocationID,
		Name:       GetInterchangeName(freewayID, locations.LocationID),
		TimeRange:  1,
		Direction:  direction,
		Data:       speedHours}

	return speedChartData
}

func GetLocationsByRegion(regionID string) [][]string {

	var results [][]string

	return results

}

func GetSmoothData() {

	locationList := GetLocationList()
	//t2, _ := time.Parse("200601021504", "201501210100")
	//t2 = t2.UTC()
	t2 := time.Now()
	t2 = t2.Add(time.Hour * (time.Duration(-8)))

	timeSting := fmt.Sprintf("%d%.2d%.2d%.2d%.2d", t2.Year(), t2.Month(), t2.Day(), 0, 0)
	t, _ := time.Parse("200601021504", timeSting)
	t = t.Add(time.Hour * time.Duration(-8)) //utc

	fmt.Println("GetSmoothData", "time is", t)

	for h := 0; h < 24; h++ {
		for m := 0; m < 60; m += 5 {
			for f := 0; f < 2; f++ {
				var filePath string

				if f == 0 {
					filePath = fmt.Sprintf("%spredict/%d%.2d%.2d%.2d%.2d_.csv", ROOT_PATH, t.Year(), t.Month(), t.Day(), t.Hour(), m)
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						//第一個星期的data, 所以沒有預測資料
						t3 := t.Add(time.Hour * time.Duration(-24*7))
						filePath = fmt.Sprintf("%s%d%.2d%.2d%.2d%.2d.csv", ROOT_PATH, t3.Year(), t3.Month(), t3.Day(), t3.Hour(), m)

						//fmt.Println("err", filePath)
					}
				} else {
					filePath = fmt.Sprintf("%s%d%.2d%.2d%.2d%.2d.csv", ROOT_PATH, t.Year(), t.Month(), t.Day(), t.Hour(), m)
				}

				//fmt.Println(filePath)

				data := GetCSVData(filePath)

				if data == nil {
					//fmt.Println("data is nil..", filePath)
					//log
				}

				//row0 is column name
				for k := 1; k < len(data); k++ {

					freewayId := data[k][1]
					locationId := data[k][2]

					for i := 0; i < len(locationList); i++ {
						if (locationList[i].FreewayId == freewayId) && (locationList[i].LocationId == locationId) {

							index := h*12 + m/5

							locationList[i].DirectionId[0], _ = strconv.Atoi(data[k][3])
							locationList[i].DirectionId[1], _ = strconv.Atoi(data[k][5])

							speed1, _ := strconv.ParseFloat(data[k][4], 32)
							speed2, _ := strconv.ParseFloat(data[k][6], 32)

							if locationList[i].AverageSpeed[0][index] == 0 {
								locationList[i].AverageSpeed[0][index] = speed1

							} else {
								locationList[i].AverageSpeed[0][index] = speed1*0.3 + locationList[i].AverageSpeed[0][index]*0.7
							}

							if locationList[i].AverageSpeed[1][index] == 0 {
								locationList[i].AverageSpeed[1][index] = speed2

							} else {
								locationList[i].AverageSpeed[1][index] = speed2*0.3 + locationList[i].AverageSpeed[1][index]*0.7
							}

							continue
						}
					}

				}
			}

		}
		t = t.Add(time.Hour)
	}

	GetError(t)

	SaveCSVData(t, locationList)

}

func GetError(t time.Time) {

	for h := 0; h < 24; h++ {
		for m := 0; m < 60; m += 5 {

			filePath := fmt.Sprintf("%s%d%.2d%.2d%.2d%.2d.csv", ROOT_PATH, t.Year(), t.Month(), t.Day(), t.Hour(), m)

			data := GetCSVData(filePath)

			filePathP := fmt.Sprintf("%spredict/%d%.2d%.2d%.2d%.2d_.csv", ROOT_PATH, t.Year(), t.Month(), t.Day(), t.Hour(), m)

			dataP := GetCSVData(filePathP)

			filePathErr := fmt.Sprintf("%s/predict/%.4d%.2d%.2d%.2d%.2d_err.csv", ROOT_PATH, t.Year(), t.Month(), t.Day(), t.Hour(), m)

			csvfile, err := os.Create(filePathErr)

			if err != nil {
				fmt.Println("Error:", err)
			}

			defer csvfile.Close()

			fmt.Fprintf(csvfile, "freeway_id,location_id,err1,err2\r\n")

			//row0 is column name
			for k := 1; k < len(data); k++ {

				for i := 1; i < len(dataP); i++ {
					if (data[k][1] == dataP[i][1]) && (data[k][2] == dataP[i][2]) {
						freewayId := data[k][1]
						locationId := data[k][2]

						s1, _ := strconv.ParseFloat(data[k][4], 32)
						s2, _ := strconv.ParseFloat(dataP[i][4], 32)
						s3, _ := strconv.ParseFloat(data[k][6], 32)
						s4, _ := strconv.ParseFloat(dataP[i][6], 32)
						err1 := math.Abs(s1 - s2)
						err2 := math.Abs(s3 - s4)

						//fmt.Println(h, m, dataP[i][1], dataP[i][2], ":", s1, s2, err1, " - ", s3, s4, err2)

						fmt.Fprintf(csvfile, "%s,%s,%.0f,%.0f\r\n", freewayId, locationId, err1, err2)
					}
				}
			}

		}
		t = t.Add(time.Hour)
	}
	/*
		if f == 0 {
			filePath =  fmt.Sprintf("%spredict/%d%.2d%.2d%.2d%.2d", ROOT_PATH, t.Year(), t.Month(), t.Day(), t.Hour(), m)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				//第一個星期的data, 所以沒有預測資料
				t3 := t.Add(time.Hour * time.Duration(-24*7))
				filePath = fmt.Sprintf("%s%d%.2d%.2d%.2d%.2d.csv", ROOT_PATH, t3.Year(), t3.Month(), t3.Day(), t3.Hour(), m)

				//fmt.Println("err", filePath)
			}
		} else {
			filePath = fmt.Sprintf("%s%d%.2d%.2d%.2d%.2d.csv", ROOT_PATH, t.Year(), t.Month(), t.Day(), t.Hour(), m)
		}
	*/

}

func GetInterchangeName(freewayID string, locationID string) string {

	interchangeList := GetInterchangeList()

	var name string
	count := 0
	for i := 0; i < len(interchangeList); i++ {

		if interchangeList[i].FreewayId == freewayID {

			for j := 0; j < len(interchangeList[i].Locations); j++ {

				if interchangeList[i].Locations[j] == locationID {
					if count != 0 {
						name += " - "
					}
					name += interchangeList[i].Name
					count++
				}
			}
		}
	}

	//fmt.Println(name)
	return name
}

func GetInterchangeList() []Interchange {

	var interchangeList []Interchange

	content, err := ioutil.ReadFile("highway.json")
	if err != nil {
		fmt.Println("Error:", err)
	}

	var roads Road
	err = json.Unmarshal(content, &roads)
	if err != nil {
		fmt.Println("Error:", err)
	}

	for i := 0; i < len(roads.Interchanges); i++ {

		interchangeList = append(interchangeList, roads.Interchanges[i])
		fmt.Println(interchangeList[i].Name)
	}

	return interchangeList
}

func GetLocationList() []LocationInfo {

	var locationList []LocationInfo

	content, err := ioutil.ReadFile("highway.json")
	if err != nil {
		fmt.Println("Error:", err)
	}

	var roads Road
	err = json.Unmarshal(content, &roads)
	if err != nil {
		fmt.Println("Error:", err)
	}

	for i := 0; i < len(roads.Freeways); i++ {

		freewayID := roads.Freeways[i].Id
		for j := 0; j < len(roads.Freeways[i].Locations); j++ {

			var LI LocationInfo
			LI.FreewayId = freewayID
			LI.LocationId = roads.Freeways[i].Locations[j]

			locationList = append(locationList, LI)

		}
	}

	return locationList
}

//應該修改輸入某一天, 回傳同一個星期的所有data
func GetFileList() []string {

	var fileList []string

	_ = filepath.Walk(ROOT_PATH, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "csv") {
			fileList = append(fileList, path)
		}
		return nil
	})

	return fileList
}

func GetCSVData(filePath string) [][]string {

	csvfile, err := os.Open(filePath)

	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	defer csvfile.Close()

	content := csv.NewReader(csvfile)
	data, err := content.ReadAll()

	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return data
}

/*
func GetData(selectDate string) {

	fileList := GetFileList()

	for _, filePath := range fileList {

		fileName := strings.Split(filePath, "/")[1]
		fileName = strings.Split(fileName, ".")[0]

		//if fileName == selectDate {
		//	fmt.Println(filePath)
		//}

	}
}

func SaveDataByLocation(locationList []LocationInfo) {

	for i := 0; i < len(locationList); i++ {

	}

}
*/

func RenameFiles() {

	fileList := GetFileList()

	for _, filePath := range fileList {

		fmt.Println(filePath)
		t, _ := time.Parse("../../../data/060102_1504.csv", filePath)
		newFileName := fmt.Sprintf("%.2d%.2d%.2d_%.2d%.2d.csv", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()/10*10)

		fmt.Println(newFileName)
		os.Rename(filePath, ROOT_PATH+newFileName)
	}
}

func SaveCSVData(t time.Time, locationList []LocationInfo) {

	t2 := t.Add(time.Hour * (time.Duration(24 * 6)))

	for j := 0; j < TIME_INTERVALS; j++ {

		minute := j % 12

		filePath := fmt.Sprintf("%s/predict/%.4d%.2d%.2d%.2d%.2d_.csv", ROOT_PATH, t2.Year(), t2.Month(), t2.Day(), t2.Hour(), minute*5)
		t2 = t2.Add(time.Duration(5) * time.Minute)
		//fmt.Println("save ", filePath)

		csvfile, err := os.Create(filePath)

		if err != nil {
			fmt.Println("Error:", err)
		}

		defer csvfile.Close()

		fmt.Fprintf(csvfile, "timestamp,freeway_id,location_id,direction_1,speed_1,direction_2,speed_2\r\n")

		for i := 0; i < len(locationList); i++ {
			timestamp := t2.Unix()
			freewayId := locationList[i].FreewayId
			locationId := locationList[i].LocationId

			direction1 := locationList[i].DirectionId[0]
			direction2 := locationList[i].DirectionId[1]
			speed1 := locationList[i].AverageSpeed[0][j]
			speed2 := locationList[i].AverageSpeed[1][j]

			fmt.Fprintf(csvfile, "%d,%s,%s,%d,%.0f,%d,%.0f\r\n", timestamp, freewayId, locationId, direction1, speed1, direction2, speed2)

		}
	}
}

func GetFileName(date string) string {

	var t time.Time
	var fileName string

	if date != "" {

		t, err := time.Parse("200601021504", date)
		t = t.Add(time.Duration(-8)*time.Hour + time.Duration(-5)*time.Minute)

		if err != nil {
			fmt.Println("Error:", err)
		}

		index := t.Minute() / 5

		fileName = fmt.Sprintf("%s/%d%.2d%.2d%.2d%.2d.csv", ROOT_PATH, t.Year(), t.Month(), t.Day(), t.Hour(), index*5)

	} else {
		t = time.Now().UTC()

		index := t.Minute() / 5

		fileName = fmt.Sprintf("%s/predict/%d%.2d%.2d%.2d%.2d_.csv", ROOT_PATH, t.Year(), t.Month(), t.Day(), t.Hour(), index*5)

	}

	return fileName
}
