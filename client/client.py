import capnp
import os
capnp.remove_import_hook()

path_to_go = os.path.expanduser('~') + "/go/src/capnproto.org/go/capnp/std/"
hello_world_capnp = capnp.load(file_name="../hello_world.capnp", imports=[path_to_go])


# hello_world_capnp = capnp.load('../hello_world.capnp')
client = capnp.TwoPartyClient('localhost:8080')
greeter = client.bootstrap().cast_as(hello_world_capnp.Greeter)

promise = greeter.greet({"name": "world"})

result = promise.wait()
print(result)
