package repository

import (
    "Student-Assistant-App/src/data/model"
    "context"
    "errors"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
    collection *mongo.Collection
}

func NewMongoUserRepository(database *mongo.Database) UserRepository {
    return &MongoUserRepository{
        collection: database.Collection("users"),
    }
}

func (r *MongoUserRepository) Save(ctx context.Context, user *model.User) (*model.User, error) {
    if user.ID.IsZero() {
        result, err := r.collection.InsertOne(ctx, user)
        if err != nil {
            return nil, err
        }
        user.ID = result.InsertedID.(primitive.ObjectID)
    } else {
        filter := bson.M{"_id": user.ID}
        _, err := r.collection.ReplaceOne(ctx, filter, user)
        if err != nil {
            return nil, err
        }
    }
    return user, nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    var user model.User
    err = r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *MongoUserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var users []*model.User
    for cursor.Next(ctx) {
        var user model.User
        if err := cursor.Decode(&user); err != nil {
            return nil, err
        }
        users = append(users, &user)
    }
    if err := cursor.Err(); err != nil {
        return nil, err
    }
    return users, nil
}

func (r *MongoUserRepository) DeleteByID(ctx context.Context, id string) error {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
    return err
}

func (r *MongoUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
    count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
