gen.aes.key:
	openssl enc -aes-128-cbc -k secret -P -pbkdf2

encrypt:
	cat tmp/creds.json | go run cmd/cli/main.go -encrypt config/ijnrge.env.enc -key $(AES_KEY)

decrypt:
	go run cmd/cli/main.go -decrypt config/ijnrge.env.enc -key $(AES_KEY)


run.server:
	DOMAIN_NAME=localhost go run -race cmd/server/main.go

docker.run:
	docker build --platform linux/amd64  -t thebackendcompany:latest .
	docker run -e MASTER_KEY=$(MASTER_KEY) -e EMAIL_LEADS_SHEET_ID=$(SHEET_ID) -p 8080:8080 -it thebackendcompany:latest
