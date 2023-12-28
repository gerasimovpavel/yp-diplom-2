package domain

type Club struct {
	ClubId   string `json:"clubId" db:"clubId"`
	Name     string `json:"name" db:"name"`
	Password string `json:"password" db:"password"`
}
