brew install go

from the link(go lang clean)-> https://github.com/ilyakaznacheev/cleanenv

go get -u github.com/ilyakaznacheev/cleanenv

go run cmd/server/main.go -config config/local.yaml
-config is the flag which is must

https://github.com/go-playground/validator
go get github.com/go-playground/validator/v10
The above command is used to validate the request which is sent by the user.