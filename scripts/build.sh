echo "Building src/client..."
go build -o ./build/client.out src/client/main.go
echo "Done!"
echo "Building src/server..."
go build -o ./build/server.out src/server/main.go
echo "Done!"
echo "Building src/bot..."
go build -o ./build/bot.out src/bot/main.go
echo "Done!"
