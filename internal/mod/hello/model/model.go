package model

type Hello struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}
