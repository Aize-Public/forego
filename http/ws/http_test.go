package ws_test

import (
	"io"
	"testing"

	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/http"
	"github.com/Aize-Public/forego/http/ws"
	"github.com/Aize-Public/forego/test"
	"golang.org/x/net/websocket"
)

// test websocket implementation
func TestHttp(t *testing.T) {
	c := test.Context(t)
	s := http.NewServer(c)
	addr, err := s.Listen(c, "127.0.0.1:0")
	test.NoError(t, err)
	test.NotEqualsGo(t, 0, addr.Port)

	h := &ws.RpcHandler{}
	h.MustRegister(c, &Echo{})
	s.Mux().Handle("/ws", h.Server())

	conf, err := websocket.NewConfig("ws://"+addr.String()+"/ws", "http://"+addr.String()+"/")
	test.NoError(t, err)
	t.Logf("connecting to %+v", conf.Location)
	conn, err := websocket.DialConfig(conf)
	test.NoError(t, err)

	_, err = conn.Write(enc.MustMarshalJSON(c, ws.Frame{
		Channel: "c0",
		Path:    "echo",
		Type:    "open",
		Data:    enc.String("ping"),
	}))
	test.NoError(t, err)

	buf := make([]byte, 1024)
LOOP:
	for {
		ct, err := conn.Read(buf)
		switch err {
		case io.EOF:
			break LOOP
		case nil:
		default:
			test.NoError(t, err)
		}
		var f ws.Frame
		test.NoError(t, enc.UnmarshalJSON(c, buf[0:ct], &f))
		switch f.Type {
		default:
			test.Fail(t, "unexpected %s", buf[0:ct])
		case "":
			t.Logf("recv %v", f.Data)
		}
	}
}

type Echo struct {
}

func (this *Echo) Init(c ws.C, in string) error {
	defer c.Close()
	return c.Reply("echo", in)
}
