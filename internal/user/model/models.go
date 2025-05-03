package model

type PersonFromRequest struct {
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
}

// Person to DB
type Person struct {
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         int     `json:"age"`
	Gender      string  `json:"gender"`
	Nationality string  `json:"nationality"`
	IsDeleted   bool
}
