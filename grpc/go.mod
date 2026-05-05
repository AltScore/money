module github.com/AltScore/money/grpc/v2

go 1.18

require (
	github.com/AltScore/money/v2 v2.0.0
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4
)

replace github.com/AltScore/money/v2 => ../

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	go.mongodb.org/mongo-driver/v2 v2.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
