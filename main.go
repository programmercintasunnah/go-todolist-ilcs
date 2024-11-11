package main

func main() {
	// todo-api/
	// ├── cmd/
	// │   └── main.go
	// ├── config/
	// │   └── config.go
	// ├── controllers/
	// │   └── task_controller.go
	// |   └── auth_controller.go
	// ├── models/
	// │   └── task.go
	// ├── repositories/
	// │   └── task_repository.go
	// |   └── redis.go
	// ├── services/
	// │   └── task_service.go
	// |   └── mock_service.go
	// ├── middlewares/
	// │   └── auth.go
	// │   └── logger.go
	// ├── utils/
	// │   └── jwt.go
	// │   └── validator.go
	// ├── tests/
	// │   └── task_test.go
	// ├── go.mod
	// ├── go.sum

	// mkdir go-todolist-ilcs
	// cd go-todolist-ilcs
	// go mod init github.com/programmercintasunnah/go-todolist-ilcs

	// go get -u github.com/gin-gonic/gin
	// go get -u gorm.io/gorm
	// go get -u gorm.io/driver/oracle
	// go get -u github.com/go-redis/redis/v8
	// go get -u github.com/dgrijalva/jwt-go
	// go get -u github.com/go-playground/validator/v10
	// go get -u github.com/sirupsen/logrus
	// go get -u github.com/stretchr/testify
	// dan banyak lagi

	// table nya
	// CREATE TABLE tasks (
	// 	id SERIAL PRIMARY KEY,
	// 	title VARCHAR(255) NOT NULL,
	// 	description TEXT,
	// 	status VARCHAR(20) CHECK (status IN ('pending', 'completed')) NOT NULL,
	// 	due_date TIMESTAMP NOT NULL,
	// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	// 	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	// 	deleted_at TIMESTAMP
	// );
}
