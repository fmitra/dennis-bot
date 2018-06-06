package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mocks "github.com/fmitra/dennis-bot/test"
)

func TestGeneratesMessage(t *testing.T) {
	var message BotResponse
	MessageMap = mocks.MessageMapMock

	message = GetMessage("tracking_success", "")
	assert.Equal(t, BotResponse("Roger that!"), message)

	message = GetMessage("period_total_success", "20")
	assert.Equal(t, BotResponse("You spent 20"), message)
}
