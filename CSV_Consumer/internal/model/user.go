package model

type User struct {
	Id          string `csv:"Id" db:"id"`
	FullName    string `csv:"FullName" db:"full_name"`
	Username    string `csv:"Username" db:"username"`
	Email       string `csv:"Email" db:"email"`
	PhoneNumber string `csv:"PhoneNumber" db:"phone_number"`
}
