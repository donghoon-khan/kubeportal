package common

type ResourceStatus struct {
	Running     int `json:"running"`
	Pending     int `json:"pending"`
	Failed      int `json:"failed"`
	Succeeded   int `json:"succeeded"`
	Unknows     int `json:"unknown`
	Terminating int `json:"terminating"`
}
