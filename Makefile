local-up:
	docker compose up -d;

backend:
	go build -o ~/.local/bin/photos_api cmd/photos_api/main.go
	~/.local/bin/photos_api

backend-local:
	go build -o ~/.local/bin/photos_api cmd/photos_api/main.go
	~/.local/bin/photos_api -s local
