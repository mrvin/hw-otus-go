//go:generate codecgen -o user_generated.go user.go

package user

type User struct {
	ID       int    `codec:",omitempty"`
	Name     string `codec:",omitempty"`
	Username string `codec:",omitempty"`
	Email    string `codec:",omitempty"`
	Phone    string `codec:",omitempty"`
	Password string `codec:",omitempty"`
	Address  string `codec:",omitempty"`
}
