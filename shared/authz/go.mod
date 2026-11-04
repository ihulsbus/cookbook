module github.com/ihulsbus/cookbook/shared/authz

go 1.26.0

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/gin-gonic/gin v1.12.0
	github.com/google/uuid v1.6.0
	github.com/ihulsbus/cookbook/shared/oidc v0.0.1
	github.com/ihulsbus/cookbook/shared/rabbitmq v0.0.1
	github.com/stretchr/testify v1.11.1
	github.com/wagslane/go-rabbitmq v0.15.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
)

replace github.com/ihulsbus/cookbook/shared/oidc => ../oidc

replace github.com/ihulsbus/cookbook/shared/rabbitmq => ../rabbitmq
