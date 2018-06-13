package conversation

// An Intent designed to deliver generic responses when we
// are unable to determine what the user wants.
type GenericResponse struct {
	Context
	actions *Actions
}

func (i *GenericResponse) GetResponses() []func() (BotResponse, error) {
	return []func() (BotResponse, error){
		i.SayDefault,
	}
}

func (i *GenericResponse) Respond() (BotResponse, int) {
	responses := i.GetResponses()
	return i.Process(responses)
}

func (i *GenericResponse) SayDefault() (BotResponse, error) {
	return GetMessage(DEFAULT, ""), nil
}
