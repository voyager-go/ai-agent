package db

import (
	"ai-agent/model"
	"ai-agent/shared"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func CreateTables() error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		age INT NOT NULL
	);`
	orderTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		amount DECIMAL(10, 2) NOT NULL,
		order_date DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err := shared.DB.Exec(userTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = shared.DB.Exec(orderTable)
	if err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}

	return nil
}

func QueryUsers() ([]model.User, error) {
	rows, err := shared.DB.Query("SELECT id, name, email, age FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func QueryUserByName(name string) ([]model.User, error) {
	query := "SELECT id, name, email, age FROM users WHERE name LIKE ?"
	rows, err := shared.DB.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// 检查是否发生错误
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func InsertUser(user model.User) (int64, error) {
	result, err := shared.DB.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", user.Name, user.Email, user.Age)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func DeleteUser(db *sql.DB, userID int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
	return err
}

func UpdateUser(db *sql.DB, user model.User) error {
	_, err := db.Exec("UPDATE users SET name = ?, email = ?, age = ? WHERE id = ?", user.Name, user.Email, user.Age, user.ID)
	return err
}

func QueryOrdersByUserName(namePattern string) (map[model.User][]model.Order, error) {
	query := `
	SELECT u.id, u.name, u.email, o.id, o.amount, o.order_date 
	FROM users u 
	LEFT JOIN orders o ON u.id = o.user_id 
	WHERE u.name LIKE ? 
	ORDER BY u.id, o.order_date`

	rows, err := shared.DB.Query(query, "%"+namePattern+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 结果存入 map，按用户分组
	usersOrders := make(map[model.User][]model.Order)

	for rows.Next() {
		var user model.User
		var order model.Order
		var orderDateStr sql.NullString // 订单时间可能为空

		// 读取数据
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &order.ID, &order.Amount, &orderDateStr)
		if err != nil {
			return nil, err
		}

		user.Age = 0

		// 解析订单时间
		if orderDateStr.Valid {
			order.UserID = user.ID
			order.OrderDate, _ = time.Parse("2006-01-02 15:04:05", orderDateStr.String)
			usersOrders[user] = append(usersOrders[user], order)
		} else {
			usersOrders[user] = []model.Order{}
		}
	}
	return usersOrders, nil
}

func QueryUserOrders(userID int) ([]model.Order, error) {
	rows, err := shared.DB.Query("SELECT id, user_id, amount, order_date FROM orders WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.Amount, &order.OrderDate)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// 生成模拟用户和订单数据
func GenerateMockData() {
	// 清空表
	shared.DB.Exec("TRUNCATE TABLE orders")
	shared.DB.Exec("TRUNCATE TABLE users")

	// 插入用户
	for i := 1; i <= 5; i++ {
		user := model.User{
			Name:  fmt.Sprintf("User%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   rand.Intn(40) + 20,
		}
		userID, _ := InsertUser(user)

		// 插入订单
		for j := 1; j <= rand.Intn(5)+1; j++ {
			order := model.Order{
				UserID:    int(userID),
				Amount:    rand.Float64() * 100,
				OrderDate: time.Now().AddDate(0, 0, -rand.Intn(30)),
			}
			_, err := shared.DB.Exec("INSERT INTO orders (user_id, amount, order_date) VALUES (?, ?, ?)", order.UserID, order.Amount, order.OrderDate)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
