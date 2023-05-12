package hello_world

import (
	context "context"

	"capnproto.org/go/capnp/v3"
)

type GreeterServer struct{}

func (GreeterServer) Greet(ctx context.Context, call Greeter_greet) error {
	res, err := call.AllocResults() // allocate the results struct
	if err != nil {
		return err
	}
	// make a segment
	arena := capnp.SingleSegment(nil)
	_, seg, err := capnp.NewMessage(arena)
	if err != nil {
		panic(err)
	}

	// Make a new response in that segment
	resp, err := NewRootGreeting(seg)
	if err != nil {
		panic(err)
	}

	// get person struct from args and name from person
	person, err := call.Args().Person()
	if err != nil {
		panic(err)
	}
	name, err := person.Name()
	if err != nil {
		panic(err)
	}

	// Set 'text' of response and send back
	resp.SetText("Hello " + name + "!")
	res.SetResponse(resp)
	return nil
}
