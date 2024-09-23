package db

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // Import the sqlite3 driver
)

func Connection() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/example.db")
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}
	return db
}

func Init(conn *sql.DB) {
	_, err := conn.Exec(createTableUsers)
	if err != nil {
		panic("Error creating table users: " + err.Error())
	}

	_, err = conn.Exec(createTableProducts)
	if err != nil {
		panic("Error creating table products: " + err.Error())
	}

	_, err = conn.Exec(createTableTransactions)
	if err != nil {
		panic("Error creating table transactions: " + err.Error())
	}

	// Seeding if table is empty
	isTableEmpty := func(tableName string) bool {
		row := conn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName))
		var count int
		_ = row.Scan(&count)
		return count == 0
	}

	if isTableEmpty("users") {
		q := "INSERT INTO users (id, name) VALUES (?, ?)"
		_, err = conn.Exec(
			q,
			uuid.New().String(),
			"John Doe",
		)
		if err != nil {
			panic("Error inserting user: " + err.Error())
		}
		fmt.Println("User seeded")
	}

	if isTableEmpty("products") {
		_, err = conn.Exec(
			"INSERT INTO products (name, price, quantity) VALUES ('Product 1', 10000, 10)",
		)
		if err != nil {
			panic("Error inserting product: " + err.Error())
		}
		fmt.Println("Product seeded")
	}

	fmt.Println("Database initialized")
}

var createTableUsers = `
CREATE TABLE IF NOT EXISTS users (
	id VARCHAR(100) PRIMARY KEY,
	name TEXT NOT NULL
);
`

var createTableProducts = `
CREATE TABLE IF NOT EXISTS products (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	price INT NOT NULL,
	quantity INT NOT NULL
);
`

var createTableTransactions = `
CREATE TABLE IF NOT EXISTS transactions (
	id VARCHAR(100) PRIMARY KEY,
	user_id VARCHAR(100) NOT NULL,
	product_id INT NOT NULL,
	quantity INT NOT NULL,
	amount INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (product_id) REFERENCES products(id)
);
`
