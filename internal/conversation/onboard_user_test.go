package conversation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/pkg/telegram"
	mocks "github.com/fmitra/dennis-bot/test"
)

func TestOnboardUser(t *testing.T) {
	t.Run("Should return a list of possible responses", func(t *testing.T) {
		onboardUser := &OnboardUser{}
		assert.Equal(t, 4, len(onboardUser.GetResponses()))
	})

	t.Run("Should ask for password", func(t *testing.T) {
		onboardUser := &OnboardUser{
			Context{
				Step: 0,
			},
			&Actions{},
		}
		MessageMap = mocks.MessageMapMock

		response, _ := onboardUser.AskForPassword()
		assert.Equal(t, BotResponse("What's your password?"), response)
	})

	t.Run("Should ask to confirm password", func(t *testing.T) {
		var incMessage telegram.IncomingMessage
		message := mocks.GetMockMessage("")
		json.Unmarshal(message, &incMessage)

		onboardUser := &OnboardUser{
			Context{
				Step:       1,
				IncMessage: incMessage,
			},
			&Actions{},
		}
		MessageMap = mocks.MessageMapMock

		response, _ := onboardUser.ConfirmPassword()
		assert.Equal(t, BotResponse("Your password is Hello world"), response)
	})

	t.Run("Should validate password", func(t *testing.T) {
		MessageMap = mocks.MessageMapMock
		var incMessage telegram.IncomingMessage
		message := mocks.GetMockMessage("No")
		json.Unmarshal(message, &incMessage)

		onboardUser := &OnboardUser{
			Context{
				Step:       2,
				IncMessage: incMessage,
			},
			&Actions{},
		}

		response, err := onboardUser.ValidatePassword()
		assert.Equal(t, BotResponse("Okay try again later"), response)
		assert.EqualError(t, err, "Password rejected")

		message = mocks.GetMockMessage("YES")
		json.Unmarshal(message, &incMessage)

		onboardUser = &OnboardUser{
			Context{
				Step:       2,
				IncMessage: incMessage,
			},
			&Actions{},
		}

		response, err = onboardUser.ValidatePassword()
		assert.Equal(t, BotResponse(""), response)
		assert.NoError(t, err)

		message = mocks.GetMockMessage("Invalid")
		json.Unmarshal(message, &incMessage)

		onboardUser = &OnboardUser{
			Context{
				Step:       2,
				IncMessage: incMessage,
			},
			&Actions{},
		}

		response, err = onboardUser.ValidatePassword()
		assert.Equal(t, BotResponse("I didn't understand that"), response)
		assert.EqualError(t, err, "Response invalid")
	})

	t.Run("Should say outro", func(t *testing.T) {
		onboardUser := &OnboardUser{
			Context{
				Step: 3,
			},
			&Actions{},
		}
		MessageMap = mocks.MessageMapMock

		response, _ := onboardUser.SayOutro()
		assert.Equal(t, BotResponse("Outro message"), response)
	})
}
