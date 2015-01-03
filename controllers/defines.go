package controllers

//一天的所有時間區段
const TimeIntervals int = 24 * 12

type SpeedTime struct {
	Speed1     int
	Speed2     int
	Direction1 int
	Direction2 int
	Time       string
}

type SpeedChart struct {
	LocationID string
	TimeRange  int
	Data       []SpeedTime
}

type Location struct {
	TimeStamp  string
	FreewayID  string
	LocationID string
	Direction1 int
	Direction2 int
	Speed1     int
	Speed2     int
}

type LocationInfo struct {
	FreewayId    string
	LocationId   string
	LocationName string
	StartMile    string
	EndMile      string
	Time         string
	DirectionId  [2]int
	AverageSpeed [2][TimeIntervals]float64
}

type Freeway struct {
	Name      string
	Id        string
	Direction bool
	Locations []string
}

type Road struct {
	Freeways []Freeway
}
