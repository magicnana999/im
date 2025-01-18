package enum

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestDeviceType_String(t *testing.T) {
	var v1 Enum = Mobile
	var v2 Enum = PC

	var failValue uint8 = 100
	var v4 DeviceType = DeviceType(failValue)

	assert.Equal(t, v1.Name(), "Mobile")
	assert.Equal(t, v2.Name(), "PC")

	assert.Equal(t, v1.Valid(), true)
	assert.Equal(t, v2.Valid(), true)

	assert.Equal(t, v1.Code(), uint8(Mobile))
	assert.Equal(t, v2.Code(), uint8(PC))

	assert.Equal(t, v4.Name(), "")
	assert.Equal(t, v4.Code(), failValue)
	assert.Equal(t, v4.Valid(), false)

}
