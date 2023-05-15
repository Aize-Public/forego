package http

import (
	"encoding/json"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
)

// TODO
func TestAPI[T Doable, UID any](c ctx.C, uid UID, obj T) error {
	cli, err := api.NewClient(c, obj)
	if err != nil {
		return err
	}
	ser, err := api.NewServer(c, obj)
	if err != nil {
		return err
	}
	var data api.JSON
	err = cli.Send(c, obj, &data)
	if err != nil {
		return err
	}

	{
		data.UID, err = json.Marshal(uid)
		if err != nil {
			return ctx.NewErrorf(c, "can't marshal UID: %w", err)
		}
		op, err := ser.Recv(c, &data)
		if err != nil {
			return ctx.NewErrorf(c, "remote: %w", err)
		}
		err = op.Do(c)
		if err != nil {
			return ctx.NewErrorf(c, "remote: %w", err)
		}

		data.Data = nil // clear the buffer
		err = ser.Send(c, op, &data)
		if err != nil {
			return ctx.NewErrorf(c, "remote: %w", err)
		}
	}
	err = cli.Recv(c, &data, obj)
	if err != nil {
		return err
	}
	return nil
}
