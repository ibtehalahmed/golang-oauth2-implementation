package models
type Oauth2Token struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Access     string    `gorm:"not null" `
	Data     string    `gorm:"not null"`

}
