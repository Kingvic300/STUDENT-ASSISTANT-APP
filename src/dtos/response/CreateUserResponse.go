package response

import "Student-Assistant-App/src/data/model"

type CreateUserResponse struct {
	Message 	string      	`json:"message"`
	User    	*model.User 	`json:"user"`
}
func (req *CreateUserResponse) SetMessage(Message string){
	req.Message = Message
}
func (req *CreateUserResponse) GetMessage() string {
	return req.Message
}
func (req *CreateUserResponse) SetUser(User *model.User){
	req.User = User
}
func (req *CreateUserResponse) GetUser() *model.User{
	return req.User
}