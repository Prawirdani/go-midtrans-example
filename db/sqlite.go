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
	ensureTables(conn)

	// Seeding if table is empty
	isTableEmpty := func(tableName string) bool {
		row := conn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName))
		var count int
		_ = row.Scan(&count)
		return count == 0
	}

	if isTableEmpty("users") {
		query := "INSERT INTO users (id, first_name, last_name, email, phone) VALUES (?, ?, ?, ?, ?)"
		_, err := conn.Exec(
			query,
			uuid.New().String(),
			"John",
			"Doe",
			"doe@mail.com",
			"081234567890",
		)
		if err != nil {
			panic("Error inserting user: " + err.Error())
		}
		fmt.Println("User seeded")
	}

	if isTableEmpty("products") {
		_, err := conn.Exec(
			`INSERT INTO products (name, price) VALUES 
			('Product 1', 5000),
			('Product 2', 10000);
			`,
		)
		if err != nil {
			panic("Error inserting product: " + err.Error())
		}
		fmt.Println("Product seeded")
	}

	if isTableEmpty("transaction_status") {
		query := `
			INSERT INTO transaction_status (status) VALUES 
			('pending'),
			('completed'),
			('failed'),
			('cancelled');
		`
		_, err := conn.Exec(query)
		if err != nil {
			panic(err)
		}
		fmt.Println("Transaction status seeded")
	}

	fmt.Println("Database initialized")
}

func ensureTables(conn *sql.DB) {
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
}

const createTableUsers = `
CREATE TABLE IF NOT EXISTS users (
	id VARCHAR(100) PRIMARY KEY,
	first_name VARCHAR(100) NOT NULL,
	last_name VARCHAR(100) NOT NULL,
	email VARCHAR(100) NOT NULL UNIQUE,
	phone VARCHAR(20) NOT NULL
);
`

const createTableProducts = `
CREATE TABLE IF NOT EXISTS products (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	price INT NOT NULL
);
`

const createTableTransactions = `
CREATE TABLE IF NOT EXISTS transaction_status(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	status VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS transactions (
	id VARCHAR(100) PRIMARY KEY,
	user_id VARCHAR(100) NOT NULL,
	status_id INT NOT NULL DEFAULT 1,
	total INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (status_id) REFERENCES transaction_status(id)
);

CREATE TABLE IF NOT EXISTS transaction_details (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	transaction_id VARCHAR(100) NOT NULL,
	product_id INT NOT NULL,
	product_price INT NOT NULL,
	quantity INT NOT NULL,
	subtotal INT NOT NULL,
	FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);
`
