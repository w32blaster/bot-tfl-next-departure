package commands

import (
	"testing"
	"time"

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

// test parse date
func TestParseDate(t *testing.T) {

	// Given:
	strDate := "2017-12-08T08:58:00"

	// When:
	date, err := parseTflDate(strDate)

	// Then:
	assert.Nil(t, err)

	// and:
	if assert.NotNil(t, date) {
		assert.Equal(t, 8, date.Hour())
		assert.Equal(t, 58, date.Minute())
		assert.Equal(t, time.Month(12), date.Month())
		assert.Equal(t, 8, date.Day())
	}
}
