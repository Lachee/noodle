cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
go build
./webserver.exe