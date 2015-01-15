package controllers

//一天的所有時間區段
const TIME_INTERVALS int = 24 * 12
const ROOT_PATH string = "../../../data"

type SpeedTime struct {
	Speed1 int
	Speed2 int
	Time   string
}

type SpeedChart struct {
	LocationID string
	Name       string
	TimeRange  int
	Direction  string
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
	StartMile    string
	EndMile      string
	Time         string
	DirectionId  [2]int
	AverageSpeed [2][TIME_INTERVALS]float64
}

type Interchange struct {
	Name      string   `json:"name"`
	Id        string   `json:"id"`
	FreewayId string   `json:"freeway_id"`
	Locations []string `json:"locations"`
}

type Freeway struct {
	Name      string
	Id        string
	Direction bool
	Locations []string
}

type Road struct {
	Freeways     []Freeway
	Interchanges []Interchange
}
