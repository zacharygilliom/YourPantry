package models

type UserPOST struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUserPOST struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
