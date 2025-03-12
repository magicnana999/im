package error

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError(t *testing.T) {
	err := New(1001, "Test error")
	assert.Equal(t, 1001, err.GetCode())
	assert.Equal(t, "Test error", err.GetMessage())

	err = err.SetMessage("Updated error")
	assert.Equal(t, "Updated error", err.GetMessage())

	err = err.SetDetail("Error occurred during processing")
	assert.Equal(t, "Error occurred during processing", err.GetDetail())

	err = err.SetDetail(42)
	assert.Equal(t, "42", err.GetDetail())

	jsonStr := err.JsonString()
	assert.NotEqual(t, "{}", jsonStr)

	formattedError := Format(err)
	assert.NotNil(t, formattedError)
	assert.Equal(t, err.GetCode(), formattedError.GetCode())
	assert.Equal(t, err.GetMessage(), formattedError.GetMessage())

	nilError := Format(nil)
	assert.Nil(t, nilError)
}

func TestSetDetailErrorHandling(t *testing.T) {
	err := New(2001, "Invalid type")
	err = err.SetDetail(map[string]int{"key": 1})
	assert.NotEmpty(t, err.GetDetail())
}
