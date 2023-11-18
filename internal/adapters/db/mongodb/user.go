package mongodb

import (
	"context"

	uuid "github.com/google/uuid"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) error
	GetByCredentials(ctx context.Context, email, password string) (entity.User, error)
}

type UserStorage struct {
	db *mongo.Collection
}

func NewUserStorage(db *mongo.Collection) UserStorage {
	return UserStorage{db}
}

func (s *UserStorage) Create(ctx context.Context, user *entity.User) error {
	id := uuid.New()
	user.ID = id.String()

	marshalledUser, err := bson.Marshal(user)
	if err != nil {
		return err
	}

	_, err = s.db.InsertOne(context.TODO(), marshalledUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) GetByCredentials(ctx context.Context, email, password string) (*entity.User, error) {
	filter := bson.M{"email": email, "password": password}

	var user entity.User

	err := s.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
