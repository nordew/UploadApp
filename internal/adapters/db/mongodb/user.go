package mongodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrFailedToMarshal = errors.New("failed to marshal user")
	ErrFailedToInsert  = errors.New("failed to insert user")
	ErrFailedToDecode  = errors.New("failed to decode user")
	ErrDuplicateKey    = errors.New("user with this email already exists")
)

type UserStorage interface {
	// Create creates a new user in the database.
	Create(ctx context.Context, user entity.User) error
	// GetByCredentials retrieves a user from the database by email and password.
	// It returns an error if the operation fails or the user is not found.
	GetByCredentials(ctx context.Context, email, password string) (*entity.User, error)
}

type userStorage struct {
	db *mongo.Collection
}

func NewUserStorage(db *mongo.Collection) *userStorage {
	return &userStorage{db}
}

func IsDuplicateKeyError(err error) bool {
	var writeException mongo.WriteException
	return errors.As(err, &writeException) && writeException.HasErrorCode(11000)
}

func (s *userStorage) Create(ctx context.Context, user entity.User) error {
	user.ID = uuid.NewString()

	marshalledUser, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedToMarshal, err)
	}

	_, err = s.db.InsertOne(ctx, marshalledUser)
	if err != nil {
		if IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateKey, err)
		}

		return fmt.Errorf("%w: %v", ErrFailedToInsert, err)
	}

	return nil
}

func (s *userStorage) GetByCredentials(ctx context.Context, email, password string) (*entity.User, error) {
	filter := bson.M{"email": email, "password": password}

	var user entity.User

	err := s.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w: user not found for email %s", ErrUserNotFound, email)
		}
		return nil, fmt.Errorf("%w: %v", ErrFailedToDecode, err)
	}

	return &user, nil
}
