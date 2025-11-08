echo "Building src/client..."
go build -o ./build/client.out src/client/main.go
echo "Done!"
echo "Building src/server..."
go build -o ./build/server.out src/server/main.go
echo "Done!"
