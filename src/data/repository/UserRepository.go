package repository

import (
    "context"
    "Student-Assistant-App/src/data/model"
)

type UserRepository interface {
    Save(ctx context.Context, user *model.User) (*model.User, error)
    FindByID(ctx context.Context, id string) (*model.User, error)
    FindByEmail(ctx context.Context, email string) (*model.User, error)
    FindAll(ctx context.Context) ([]*model.User, error)
    DeleteByID(ctx context.Context, id string) error
    ExistsByEmail(ctx context.Context, email string) (bool, error)
}
