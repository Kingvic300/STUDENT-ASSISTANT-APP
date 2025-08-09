package repository

import (
	"Student-Assistant-App/src/data/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OTPRepository interface {
	Save(ctx context.Context, otp *model.OTP) (*model.OTP, error)
	FindByEmailAndCode(ctx context.Context, email, code string) (*model.OTP, error)
	FindLatestByEmailAndPurpose(ctx context.Context, email, purpose string) (*model.OTP, error)
	MarkAsUsed(ctx context.Context, id primitive.ObjectID) error
	DeleteExpired(ctx context.Context) error
	DeleteByEmail(ctx context.Context, email string) error
}

type OTPRepositoryImpl struct {
	collection *mongo.Collection
}

func NewOTPRepositoryImpl(database *mongo.Database) OTPRepository {
	collection := database.Collection("otps")
	
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	collection.Indexes().CreateOne(context.Background(), indexModel)
	
	return &OTPRepositoryImpl{
		collection: collection,
	}
}

func (r *OTPRepositoryImpl) Save(ctx context.Context, otp *model.OTP) (*model.OTP, error) {
	if otp.ID.IsZero() {
		otp.CreatedAt = time.Now()
		result, err := r.collection.InsertOne(ctx, otp)
		if err != nil {
			return nil, err
		}
		otp.ID = result.InsertedID.(primitive.ObjectID)
	} else {
		filter := bson.M{"_id": otp.ID}
		_, err := r.collection.ReplaceOne(ctx, filter, otp)
		if err != nil {
			return nil, err
		}
	}
	return otp, nil
}

func (r *OTPRepositoryImpl) FindByEmailAndCode(ctx context.Context, email, code string) (*model.OTP, error) {
	var otp model.OTP
	filter := bson.M{
		"email": email,
		"code":  code,
		"used":  false,
		"expires_at": bson.M{"$gt": time.Now()},
	}
	
	err := r.collection.FindOne(ctx, filter).Decode(&otp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

func (r *OTPRepositoryImpl) FindLatestByEmailAndPurpose(ctx context.Context, email, purpose string) (*model.OTP, error) {
	var otp model.OTP
	filter := bson.M{
		"email":   email,
		"purpose": purpose,
	}
	
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	err := r.collection.FindOne(ctx, filter, opts).Decode(&otp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

func (r *OTPRepositoryImpl) MarkAsUsed(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"used": true}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *OTPRepositoryImpl) DeleteExpired(ctx context.Context) error {
	filter := bson.M{"expires_at": bson.M{"$lt": time.Now()}}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

func (r *OTPRepositoryImpl) DeleteByEmail(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}