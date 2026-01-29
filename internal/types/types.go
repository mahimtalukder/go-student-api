package types

type Student struct {
	Id    int64
	Email string `validate:"required,email"`
	Name  string `validate:"required"`
	Age   int    `validate:"required,number"`
}
