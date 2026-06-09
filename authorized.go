package esaccount

type AuthorizedUser struct {
	Id         int64  `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	PictureURL string `json:"picture_url"`
	Email      string `json:"email"`
	Access     string `json:"access"`

	TelegramId string `json:"telegram_id"`
	GoogleId   string `json:"google_id"`
}
