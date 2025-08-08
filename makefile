path:
	export PATH="$PATH:$(go env GOPATH)/bin"
gen:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/service.proto
start:
	sudo docker compose up --build -d      
down:
	sudo docker compose down      