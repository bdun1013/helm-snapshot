package valueutils_test

import (
	"testing"

	"github.com/bpdunni/helm-unittest/unittest/common"
	. "github.com/bpdunni/helm-unittest/unittest/valueutils"
	"github.com/stretchr/testify/assert"
)

func TestGetValueOfSetPath(t *testing.T) {
	a := assert.New(t)
	data := common.K8sManifest{
		"a": map[string]interface{}{
			"b": []interface{}{"_", map[string]interface{}{"c": "yes"}},
		},
	}

	var expectionsMapping = map[string]interface{}{
		"a.b[1].c": "yes",
		"a.b[0]":   "_",
		"a.b":      []interface{}{"_", map[string]interface{}{"c": "yes"}},
	}

	for path, expect := range expectionsMapping {
		actual, err := GetValueOfSetPath(data, path)
		a.Equal(actual, expect)
		a.Nil(err)
	}
}

func TestBuildValueOfSetPath(t *testing.T) {
	a := assert.New(t)
	data := map[string]interface{}{"foo": "bar"}

	var expectionsMapping = map[string]interface{}{
		"a.b":    map[string]interface{}{"a": map[string]interface{}{"b": data}},
		"a[1]":   map[string]interface{}{"a": []interface{}{nil, data}},
		"a[1].b": map[string]interface{}{"a": []interface{}{nil, map[string]interface{}{"b": data}}},
	}

	for path, expected := range expectionsMapping {
		actual, err := BuildValueOfSetPath(data, path)
		a.Equal(actual, expected)
		a.Nil(err)
	}
}
