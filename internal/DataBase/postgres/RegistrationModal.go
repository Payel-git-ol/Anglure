package postgres

type UserRegister struct {
	Email    string
	Password string
	Name     string
	Region   string
	ID       uint
}
