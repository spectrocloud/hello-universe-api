package internal

type Counter struct {
	Id      int    `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Browser string `json:"browser" db:"browser"`
	Os      string `json:"os" db:"os"`
}
