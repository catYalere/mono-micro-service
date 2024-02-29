package entity // import "github.com/catwashere/microservice/internal/entity"

type Session struct {
	ID     *string `bson:"_id,omitempty" json:"id,omitempty"`
	UserID *string `bson:"user_id,omitempty" json:"-"`
	User   User    `bson:"-" json:"user" validate:"required"`
}
