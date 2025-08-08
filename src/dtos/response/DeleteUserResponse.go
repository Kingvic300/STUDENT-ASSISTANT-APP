package response

type DeleteUserResponse struct {
	Message string `json:"message"`
}

func (r *DeleteUserResponse) SetMessage(message string) {
	r.Message = message
}
func( r *DeleteUserResponse) GetMessage() string{
	return r.Message
}