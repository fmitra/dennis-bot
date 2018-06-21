package conversation

// GenericResponse is an Intent designed to deliver generic messages
// when we are unable to determine what the user wants.
type GenericResponse struct {
	Context
	actions *Actions
}

// GetResponses returns a list of functions each containing a BotResponse.
func (i *GenericResponse) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.SayDefault,
	}
}

// Respond proccesses a list of response functions.
func (i *GenericResponse) Respond() (BotResponse, *Context) {
	responses := i.GetResponses()
	return i.Process(responses)
}

// SayDefault returns a generic message.
func (i *GenericResponse) SayDefault() (BotResponse, error) {
	return GetMessage(DefaultResponse, ""), nil
}
