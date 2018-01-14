package main

const trackExpenses string = "track"

type Intent interface {
	respond() string
}

type IntentSession struct {
	Keyword string
	UserLocation string
}

type IntentTrack struct {
	userId int
	lastMessageRecieved string
}

func (i IntentTrack) respond() (string) {
	intentSession, err := getIntentSession(i.userId)
	if err != nil {
		intentSession = IntentSession{
			trackExpenses,
			"askForAmount",
		}
	}

	userLocation := intentSession.UserLocation
	switch userLocation {
	case "askForAmount":
		return i.askForAmount()
	case "askForDetails":
		return i.askForDetails()
	case "askForDate":
		return i.askForDate()
	case "endConversation":
		return i.endConversation()
	default:
		return i.askForAmount()
	}
}

func (i IntentTrack) endConversation() (string) {
	clearIntentSession(i.userId)
	return getMessage("track.askForDate")
}

func (i IntentTrack) askForDate() (string) {
	isInvalid := true

	if isInvalid {
		return getMessage("track.error.askForDate")
	}

	return getMessage("track.endConversation")
}

func (i IntentTrack) askForDetails() (string) {
	isInvalid := true

	if isInvalid {
		return getMessage("track.error.askForAmount")
	}

	intentSession := IntentSession{
		trackExpenses,
		"askForDate",
	}
	updateIntentSession(i.userId, intentSession)

	return getMessage("track.askForDetails")
}

func (i IntentTrack) askForAmount() (string) {
	intentSession := IntentSession{
		trackExpenses,
		"askForDetails",
	}
	updateIntentSession(i.userId, intentSession)
	return getMessage("track.askForAmount")
}

func getIntentResponse(keyword string, incMessage IncomingMessage) (string) {
	var i Intent
	user := incMessage.getUser()

	if keyword == trackExpenses {
		i = IntentTrack{
			user.Id,
			incMessage.getMessage(),
		}
	}

	return i.respond()
}
