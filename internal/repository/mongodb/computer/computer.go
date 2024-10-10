package computer

import (
	"context"
	"log/slog"
	"practice/internal/pkg/config"
	"practice/internal/repository/mongodb"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type RepositoryComputer interface {
	Create(ctx context.Context, computer *Computer) (*Computer, error)
	Read(ctx context.Context, compID string) (*Computer, error)
	Update(ctx context.Context, computer *Computer) (string, error)
	Delete(ctx context.Context, compID string) (string, error)
	GetAll(ctx context.Context) ([]*Computer, error)
}

type Repository struct {
	repo       *mongodb.MongoDB
	collection *mongo.Collection
	logger     *slog.Logger
}

type Options struct {
	fx.In
	fx.Lifecycle
	Cfg    *config.Config
	Mongo  *mongodb.MongoDB
	Logger *slog.Logger
}

var _ RepositoryComputer = (*Repository)(nil)

func New(opts Options) RepositoryComputer {
	repo := &Repository{
		logger: opts.Logger,
	}

	opts.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			repo.repo = opts.Mongo
			repo.collection = repo.repo.DB.Collection(opts.Cfg.MongoDB_COLLECTION)
			return nil
		},
		OnStop: func(context.Context) error { return nil },
	})

	return repo
}

func (r *Repository) Create(ctx context.Context, computer *Computer) (*Computer, error) {
	res, err := r.collection.InsertOne(ctx, computer)
	if err != nil {
		return nil, errors.Wrap(err, "error while inserting computer")
	}

	id := res.InsertedID.(primitive.ObjectID)
	computer.ID = &id
	return computer, nil
}

func (r *Repository) Read(ctx context.Context, compID string) (*Computer, error) {
	objID, err := primitive.ObjectIDFromHex(compID)
	if err != nil {
		return nil, errors.Wrap(err, "error while parsing object id")
	}

	var res Computer
	if err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(err, "not found")
		}

		return nil, errors.Wrap(err, "error while finding computer")
	}

	if res.IsDeleted {
		return nil, errors.Wrap(mongo.ErrNoDocuments, "deleted")
	}

	return &res, nil
}

func (r *Repository) Update(ctx context.Context, computer *Computer) (string, error) {
	res, err := r.collection.UpdateByID(ctx, computer.ID, bson.M{"$set": computer})
	if err != nil {
		return "", errors.Wrap(err, "error while updating computer")
	}

	if res.MatchedCount == 0 {
		return "", errors.Wrap(mongo.ErrNoDocuments, "not found")
	}

	if res.ModifiedCount == 0 {
		return "", errors.New("not modified")
	}

	if computer.IsDeleted {
		return "", errors.Wrap(mongo.ErrNoDocuments, "deleted")
	}

	return computer.ID.Hex(), nil
}

func (r *Repository) Delete(ctx context.Context, compID string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(compID)
	if err != nil {
		return "", errors.Wrap(err, "error while parsing object id")
	}

	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID, "isDeleted": false},
		bson.M{"$set": bson.M{"isDeleted": true}},
	).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.Wrap(err, "not found")
		}

		return "", errors.Wrap(err, "error while deleting computer")
	}

	return compID, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*Computer, error) {
	var res []*Computer
	cursor, err := r.collection.Find(ctx, bson.M{"isDeleted": false})
	if err != nil {
		return nil, errors.Wrap(err, "error while finding computers")
	}

	if err = cursor.All(ctx, &res); err != nil {
		return nil, errors.Wrap(err, "error while decoding computers")
	}

	return res, nil
}
