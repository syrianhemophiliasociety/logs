package actions

import "time"

type Diagnosis struct {
	Id        uint   `json:"id"`
	GroupName string `json:"group_name"`
	Title     string `json:"title"`

	CreatedAt time.Time `json:"created_at"`
}
