# Noodle WebGL Library
It does WebGL stuff. Experimental. Need more docs.

## How to run examples
The main example is in `example/wasm`, and you will want to run the `build.sh` in there to generate the WASM. 

If you wish to use the built in host (since WebAssemblies require a host), you need to run `example/webserver/run.sh` which will create a HTTP server at `localhost:8090`.

The webserver is able to listen to changes in the `example/wasm` folder and the root project folder. When it detects a change it will reload the page automatically.