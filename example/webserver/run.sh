cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
go build
./webserver.exe -dir ../../ -dir ../wasm/ -dir ../wasm/resources/shader/ -cmd "update-content.bat" -filter **/*.go -resources ../wasm/resources/