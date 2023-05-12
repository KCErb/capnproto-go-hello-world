#using Go = import "../../go/src/capnproto.org/go/capnp/std/go.capnp";
using Go = import "/go.capnp";

@0xfbc2c91ea16d009d;

$Go.package("hello_world");
$Go.import("server/hello_world");

interface Greeter {
   greet @0 (person: Person) -> (response: Greeting);
}

struct Person { 
   name @0 :Text;
}

struct Greeting {
   text @0 :Text;
}
