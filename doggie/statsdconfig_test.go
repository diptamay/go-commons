package doggie

import (
	Chance "github.com/ZeFort/chance"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

var (
	dogstatsdHost string
	dogstatdPort  int
)

func TestGetValuesFromEnv(t *testing.T) {
	chance := Chance.New()
	dogstatsdHost = chance.Word()
	dogstatdPort = chance.IntBtw(1000, 9999)
	os.Setenv("DOGSTATSD_HOST", dogstatsdHost)
	os.Setenv("DOGSTATSD_PORT", strconv.Itoa(dogstatdPort))

	dataDogHost = nil
	dataDogPort = nil

	env := getValuesFromEnv()

	assert.Equal(t, true, env.hasHost, "Host should be set from env DOGSTATSD_HOST value")
	assert.Equal(t, true, env.hasPort, "Port should be set from env DOGSTATSD_PORT value")

	assert.Equal(t, dogstatsdHost, *dataDogHost, "Host values should match")
	assert.Equal(t, dogstatdPort, *dataDogPort, "Port values should match")
}
