echo "Building WASM..."
GOOS=js GOARCH=wasm go build -o main.wasm .
echo "Done"
