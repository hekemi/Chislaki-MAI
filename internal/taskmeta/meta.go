package taskmeta

// Meta описывает одно задание, которое нужно показать на фронтенде.
// Этот пакет специально отделен от реестра заданий, чтобы не было циклов
// импортов между общими типами и конкретными модулями.
type Meta struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
