module MessagesService

go 1.25.4

require (
	github.com/bouk/monkey v1.0.1
	github.com/gocql/gocql v1.7.0
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo/v4 v4.1.17
	github.com/stretchr/testify v1.11.1
)

require (
	bou.ke/monkey v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.17.0
