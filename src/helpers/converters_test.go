package helpers

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseJSONData(t *testing.T) {
	t.Run("ValidJSON", func(t *testing.T) {
		input := []byte(`{"key": "value"}`)
		expectedOutput := map[string]interface{}{"key": "value"}

		output, err := ParseJSONData(input)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(output, expectedOutput) {
			t.Errorf("Expected %v, but got %v", expectedOutput, output)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		input := []byte(`invalid-json`)
		expectedError := &json.SyntaxError{}

		_, err := ParseJSONData(input)
		if reflect.TypeOf(err) != reflect.TypeOf(expectedError) {
			t.Errorf("Expected error of type %T, but got %T", expectedError, err)
		}
	})

}
