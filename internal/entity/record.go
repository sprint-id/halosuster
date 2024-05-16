package entity

// {
// 	"identityNumber": 123123, // not null, should be 16 digit
// 	"symptoms": "", // not null, minLength 1, maxLength 2000,
// 	"medications" : "" // not null, minLength 1, maxLength 2000
// }

type Record struct {
	ID             string `json:"id"`
	IdentityNumber string `json:"identityNumber"`
	Symptoms       string `json:"symptoms"`
	Medications    string `json:"medications"`

	UserID string `json:"user_id"`
}
