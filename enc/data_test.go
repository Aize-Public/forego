package enc_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

// Make sure the types are marshaled and unmarshaled ok like encoding/json
func TestData(t *testing.T) {
	c := test.Context(t)

	testData(c, t, (any)(nil))
	testData(c, t, (map[string]any)(nil))
	testData(c, t, ([]any)(nil))
	testData(c, t, 3)
	testData(c, t, float64(3.14))
	testData(c, t, float32(0.5))
	testData(c, t, true)
	testData(c, t, "string")
	testData(c, t, "úñí©øðé")
	testData(c, t, "a&b")
	testData(c, t, time.Now().Truncate(0))
	testData(c, t, time.Time{})
	testData(c, t, map[string]any{"1": "one"})
	testData(c, t, map[int]string{1: "one"})
	testData(c, t, map[string]any{"zero": []any{nil, 3}})
	testData(c, t, []float32(nil))
	testData(c, t, []float32{})
	testData(c, t, []uint{1, 2, 3, 4})
	testData(c, t, [4]uint{1, 2, 3, 4})
	testData(c, t, [0]uint{})
	testData(c, t, []any{})
	testData(c, t, []any{nil})
	testData(c, t, []any{map[string]any{"nil": nil}})
}

func testData[T any](c ctx.C, t *testing.T, in T) {
	t.Helper()

	jsonOut, encOut, jsonJson, encJson := testDataHelper(c, t, in)

	test.EqualsJSON(t, in, encOut)
	test.EqualsJSON(t, jsonOut, encOut)
	test.EqualsGo(t, string(jsonJson), string(encJson))
}

func testDataHelper[T any](c ctx.C, t *testing.T, in T) (jsonOut, encOut T, jsonJson, encJson []byte) {
	t.Helper()

	jj, err := json.Marshal(in)
	test.NoError(t, err)
	err = json.Unmarshal(jj, &jsonOut)
	test.NoError(t, err)
	t.Logf("IN: %s (%#v)", jj, jsonOut)

	ej, err := enc.MarshalJSON(c, in)
	test.NoError(t, err)
	err = enc.UnmarshalJSON(c, ej, &encOut)
	test.NoError(t, err)

	return jsonOut, encOut, jj, ej
}

func TestAdvancedData(t *testing.T) {
	c := test.Context(t)

	// Predefine some variables
	tB := true
	tS := "úñí©øðé n23oirj0d r4j0r3k æøå_¤34t3:_;m.`#¤&+"
	tI := int(-42)
	tF32 := float32(-42.5)
	tF64 := float64(42.5)
	tT := time.Now()
	type tStt struct {
		TS *string
		TM map[string]any
	}
	tSt := tStt{
		TS: &tS,
		TM: map[string]any{tS: tB},
	}
	tStP := &struct{}{}

	// Test some advanced cases
	{
		x := [2]tStt{tSt, {}}
		testData(c, t, x)
	}
	{
		x := []any{
			tB,
			tI,
			tS,
			tF32,
			tF64,
		}
		testData(c, t, x)
	}
	{
		testData(c, t, tSt)
		testData(c, t, tStP)
	}
	{
		x := struct {
			TL []int
		}{
			TL: []int{},
		}
		testData(c, t, x)
	}
	{
		x := struct {
			TB   bool
			TS   string
			TI   int
			TF32 float32
			TF64 float64
			TT   time.Time
			TL   []int
			TM   map[int]tStt
		}{
			TB:   tB,
			TS:   tS,
			TI:   tI,
			TF32: tF32,
			TF64: tF64,
			TT:   tT,
			TL:   []int{42, 0, -666},
			TM:   map[int]tStt{42: tSt},
		}
		testData(c, t, x)
	}
	{
		x := struct {
			TB   *bool
			TS   *string
			TI   *int
			TF32 *float32
			TF64 *float64
			TT   *time.Time
			TL   []*int
			TM   map[string]*tStt
		}{
			TB:   &tB,
			TS:   &tS,
			TI:   &tI,
			TF32: &tF32,
			TF64: &tF64,
			TT:   &tT,
			TL:   []*int{&tI, &tI},
			TM:   map[string]*tStt{"42": &tSt},
		}
		testData(c, t, x)
	}
	{
		x := struct {
			I    int64
			UI   uint64
			I32  int32
			UI32 uint32
			I16  int16
			UI16 uint16
			I8   int8
			UI8  uint8
		}{
			I:    42,
			UI:   42,
			I32:  42,
			UI32: 42,
			I16:  42,
			UI16: 42,
			I8:   42,
			UI8:  42,
		}
		testData(c, t, x)
	}
	{
		x := struct {
			X struct {
				X struct {
					S string
				}
			}
			I int
			M map[int]int
		}{
			I: tI,
			M: map[int]int{tI: tI},
		}
		x.X.X.S = tS
		testData(c, t, x)
	}
	{
		x := map[string]any{
			"a": map[string]any{
				"a": map[int]any{
					1: map[int]float64{
						tI: tF64,
					},
				},
			},
		}
		testData(c, t, x)
	}
	{
		x := map[string]any{
			"a": map[string]any{
				"a": map[int]any{
					1: map[int]float64{
						tI: tF64,
					},
				},
			},
			"b": map[string]any{
				"b": map[string]any{
					"b": "bbb",
					"c": "ccc",
					"d": 0xddd,
				},
				"c": tI,
			},
			"c":  tB,
			"-1": -1,
			"+1": struct {
				s string
				i uint16
			}{
				s: "+1",
				i: 1,
			},
			tS: tS,
		}
		testData(c, t, x)
	}
}
