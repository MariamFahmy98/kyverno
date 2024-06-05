package v2

import (
	"encoding/json"
	"testing"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/stretchr/testify/assert"
)

func Test_Convert_v2_AnyAllConditions_To_v1_AnyAllCondition(t *testing.T) {
	testcases := []struct {
		name             string
		anyAllConditions []byte
		expected         []byte
	}{
		{
			name: "convert v2 any condition to v1 any condition",
			anyAllConditions: []byte(`
			{
				"any": [
				  {
					"key": "{{ request.object.data.team }}",
					"operator": "Equals",
					"value": "eng",
					"message": "The expression team = eng failed."
				  },
				  {
					"key": "{{ request.object.data.unit }}",
					"operator": "Equals",
					"value": "green",
					"message": "The expression unit = green failed."
				  }
				]
			}`),
			expected: []byte(`
			{
				"any": [
				  {
					"key": "{{ request.object.data.team }}",
					"operator": "Equals",
					"value": "eng",
					"message": "The expression team = eng failed."
				  },
				  {
					"key": "{{ request.object.data.unit }}",
					"operator": "Equals",
					"value": "green",
					"message": "The expression unit = green failed."
				  }
				]
			}`),
		},
		{
			name: "convert v2 all condition to v1 all condition",
			anyAllConditions: []byte(`
			{
				"all": [
				  {
					"key": "{{request.operation}}",
					"operator": "Equals",
					"value": "DELETE"
				  }
				]
			}`),
			expected: []byte(`
			{
				"all": [
				  {
					"key": "{{request.operation}}",
					"operator": "Equals",
					"value": "DELETE"
				  }
				]
			}`),
		},
	}
	for _, testcase := range testcases {
		var anyAllConditions AnyAllConditions
		err := json.Unmarshal(testcase.anyAllConditions, &anyAllConditions)
		assert.Nil(t, err)

		var v1anyAllConditions kyvernov1.AnyAllConditions
		err = Convert_v2_AnyAllConditions_To_v1_AnyAllConditions(&anyAllConditions, &v1anyAllConditions, nil)
		assert.Nil(t, err)

		converted, err := json.Marshal(v1anyAllConditions)
		assert.Nil(t, err)

		assert.JSONEq(t, string(testcase.expected), string(converted))
	}
}

func Test_Convert_v1_AnyAllConditions_To_v2_AnyAllCondition(t *testing.T) {
	testcases := []struct {
		name             string
		anyAllConditions []byte
		expected         []byte
	}{
		{
			name: "convert v1 any condition to v2 any condition",
			anyAllConditions: []byte(`
			{
				"any": [
				  {
					"key": "{{ request.object.data.team }}",
					"operator": "Equals",
					"value": "eng",
					"message": "The expression team = eng failed."
				  },
				  {
					"key": "{{ request.object.data.unit }}",
					"operator": "Equals",
					"value": "green",
					"message": "The expression unit = green failed."
				  }
				]
			}`),
			expected: []byte(`
			{
				"any": [
				  {
					"key": "{{ request.object.data.team }}",
					"operator": "Equals",
					"value": "eng",
					"message": "The expression team = eng failed."
				  },
				  {
					"key": "{{ request.object.data.unit }}",
					"operator": "Equals",
					"value": "green",
					"message": "The expression unit = green failed."
				  }
				]
			}`),
		},
		{
			name: "convert v1 all condition to v2 all condition",
			anyAllConditions: []byte(`
			{
				"all": [
				  {
					"key": "{{request.operation}}",
					"operator": "Equals",
					"value": "DELETE"
				  }
				]
			}`),
			expected: []byte(`
			{
				"all": [
				  {
					"key": "{{request.operation}}",
					"operator": "Equals",
					"value": "DELETE"
				  }
				]
			}`),
		},
	}
	for _, testcase := range testcases {
		var anyAllConditions kyvernov1.AnyAllConditions
		err := json.Unmarshal(testcase.anyAllConditions, &anyAllConditions)
		assert.Nil(t, err)

		var v1anyAllConditions AnyAllConditions
		err = Convert_v1_AnyAllConditions_To_v2_AnyAllConditions(&anyAllConditions, &v1anyAllConditions, nil)
		assert.Nil(t, err)

		converted, err := json.Marshal(v1anyAllConditions)
		assert.Nil(t, err)

		assert.JSONEq(t, string(testcase.expected), string(converted))
	}
}
