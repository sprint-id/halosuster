package entity

type User struct {
	ID        string `json:"id"`
	NIP       string `json:"nip"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}
