package auth

type UserAuthInfo struct {
	User  User
	Token string
}

type User struct {
	ID       uint
	GuestID  string
	Username string
	Password string
	Nama     string
	IsAdmin  bool
	UmkmID   uint
}
