package api_test

import (
	"testing"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/test"
)

func TestDataJSON(t *testing.T) {
	c := test.C(t)

	data := api.JSON{}
	data.Marshal(c, "one", "foo")
	t.Logf("JSON: %s", data)

	var x string
	data.Unmarshal(c, "one", &x)
	test.EqualsJSON(t, "foo", x)
}
