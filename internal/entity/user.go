package entity // import "github.com/catwashere/microservice/internal/entity"

type User struct {
	ID                *string `bson:"_id,omitempty" json:"id,omitempty"`
	Email             *string `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
	Name              *string `bson:"name,omitempty" json:"name" validate:"required"`
	Password          *string `bson:"-" json:"password,omitempty" validate:"required"`
	BPassword         []byte  `bson:"-" json:"bpassword,omitempty"`
	EncryptedPassword *string `bson:"password,omitempty" json:"-"`
}
