package models

const GifPath = "./app/gifs"

type Gif struct {
	ID         int64  `json:"id" db:"id"`
	Name       string `json:"path" db:"path"`
	EventClass int32  `json:"class_name" db:"class_name"`
}

func GetEvent(id int) string {
	Events := map[int]string{
		1: "Вброс",
	}
	ans, ok := Events[id]
	if !ok {
		return ""
	}
	return ans
}
