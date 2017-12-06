package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test the command extractor
func TestCommandExtractorStripSlash(t *testing.T) {

	// when:
	command := extractCommand("/start")

	// then:
	assert.Equal(t, "start", command)
}

func TestCommandExtractor(t *testing.T) {

	// when:
	command := extractCommand("start")

	// then:
	assert.Equal(t, "start", command)
}

// the full command when someone calls the bot from another (not private) chat
func TestCommandExtractorWithBotName(t *testing.T) {

	// when:
	command := extractCommand("/help@nextTrainLondonBot")

	// then:
	assert.Equal(t, "help", command)
}
