migrate:
	go run github.com/prisma/prisma-client-go migrate dev

generate:
	go run github.com/prisma/prisma-client-go generate

push:
	go run github.com/prisma/prisma-client-go db push

format:
	go run github.com/prisma/prisma-client-go format
