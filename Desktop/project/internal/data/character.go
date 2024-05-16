package data

type Character struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     int64  `json:"age"`
	Version int32  `json:"version"`
}
