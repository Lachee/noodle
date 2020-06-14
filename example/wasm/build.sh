echo "Building WASM..."
GOOS=js GOARCH=wasm go build -o resources/noodle.wasm .
echo "Done"
