package entity // import "github.com/catwashere/microservice/internal/entity"

type Credentials struct {
	Email    *string `json:"email,omitempty" validate:"required"`
	Password *string `json:"password,omitempty" validate:"required"`
}
