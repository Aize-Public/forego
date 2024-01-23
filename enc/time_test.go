package enc_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestTime(t *testing.T) {
	c := test.Context(t)
	h := enc.Handler{
		Debugf: log.Debugf,
	}
	type X struct {
		T time.Time `json:"time"`
	}
	in := X{
		T: time.Now(),
	}
	n, err := h.Marshal(c, in)
	t.Logf("n: %+v", n)
	test.NoError(t, err)
	test.Contains(t, n.GoString(), fmt.Sprint(in.T.Year())) // enc.Pairs
	j := enc.JSON{}.Encode(c, n)
	t.Logf("j: %s", j)
	test.ContainsJSON(c, j, fmt.Sprint(in.T.Year()))
	n2, err := enc.JSON{}.Decode(c, j)
	test.NoError(t, err)
	test.Contains(t, n2.GoString(), fmt.Sprint(in.T.Year())) // enc.Map
	log.Warnf(c, "%#v", n2.(enc.Map)["time"])
	var out X
	err = h.Unmarshal(c, n2, &out)
	test.NoError(t, err)
	test.EqualsGo(t, in, out)
}

func TestTimeConv(t *testing.T) {
	t.Run("time", func(t *testing.T) {
		t0 := enc.Time(time.Now().UTC().Truncate(0))
		s := enc.String(t0.String())
		t1, err := s.AsTime()
		test.NoError(t, err)
		test.EqualsGo(t, t0, t1)
	})

	t.Run("dur", func(t *testing.T) {
		s := enc.String("42s")
		d1, err := s.AsDuration()
		test.NoError(t, err)
		test.EqualsGo(t, 42*time.Second, time.Duration(d1))
	})
	t.Run("dur_f", func(t *testing.T) {
		s := enc.String("4.2s")
		d1, err := s.AsDuration()
		test.NoError(t, err)
		test.EqualsGo(t, 4200*time.Millisecond, time.Duration(d1))
	})
	t.Run("dur_df", func(t *testing.T) {
		s := enc.Digits("4.2")
		test.Assert(t, s.IsFloat())
		d1 := s.Duration(time.Second)
		test.EqualsGo(t, 4200*time.Millisecond, time.Duration(d1))
	})
	t.Run("dur_df", func(t *testing.T) {
		s := enc.Digits("42")
		test.Assert(t, !s.IsFloat())
		d1 := s.Duration(time.Millisecond)
		test.EqualsGo(t, 42*time.Millisecond, time.Duration(d1))
	})

}
