package request

import "go.mongodb.org/mongo-driver/bson/primitive"

type DeleteUserRequest struct {
    Id primitive.ObjectID `json:"id" bson:"_id"`
}
func (req *DeleteUserRequest) SetId(ID primitive.ObjectID){
    req.Id = ID
}
func (req *DeleteUserRequest) GetId() primitive.ObjectID{
    return req.Id
}