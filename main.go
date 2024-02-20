package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	ProductName       string
	ProductID         int
	OrderID           int
	Quantity          int
	ShelfID           string
	AdditionalShelfID []string
}

func main() {
	db, err := sql.Open("mysql", "root:f7kmXohh!@tcp(127.0.0.1:3306)/store")
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return
	}

	defer db.Close()

	if len(os.Args) < 2 {
		fmt.Println("Необходимо указать номера заказов через запятую")
		return
	}

	orderIDs := strings.Split(os.Args[1], ",")

	query := "SELECT p.product_name, p.product_id, o.order_id, o.quantity, o.shelf_id, ash.shelf_id FROM Products p JOIN Orders o ON p.product_id = o.product_id LEFT JOIN additional_shelves ash ON p.product_id = ash.product_id WHERE o.order_id IN ("
	for i, orderID := range orderIDs {
		if i > 0 {
			query += ","
		}
		query += orderID
	}
	query += ")"

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer rows.Close()

	products := make(map[string]map[int]*Product)

	for rows.Next() {
		var productName string
		var productID int
		var orderID int
		var quantity int
		var shelfID string
		var additionalShelfID sql.NullString
		err := rows.Scan(&productName, &productID, &orderID, &quantity, &shelfID, &additionalShelfID)
		if err != nil {
			fmt.Println("Ошибка при сканировании строки:", err)
			return
		}

		if products[shelfID] == nil {
			products[shelfID] = make(map[int]*Product)
		}

		element := products[shelfID][productID]
		if element == nil {

			temporaryProduct := &Product{
				ProductName:       productName,
				ProductID:         productID,
				OrderID:           orderID,
				Quantity:          quantity,
				ShelfID:           shelfID,
				AdditionalShelfID: []string{},
			}

			if additionalShelfID.Valid {
				temporaryProduct.AdditionalShelfID = append(temporaryProduct.AdditionalShelfID, additionalShelfID.String)
			}
			products[shelfID][productID] = temporaryProduct
		} else {
			if additionalShelfID.Valid {
				element.AdditionalShelfID = append(element.AdditionalShelfID, additionalShelfID.String)
			}
		}
	}

	fmt.Println("=+=+=+=")
	fmt.Println("Страница сборки заказов", strings.Join(orderIDs, ","))
	fmt.Println("======================")
	for key, value := range products {
		fmt.Printf("Стеллаж: %s\n", key)
		fmt.Println("---------------------------")
		for _, elementValue := range value {
			fmt.Printf("Название товара: %s\n", elementValue.ProductName)
			fmt.Printf("ID товара: %d\n", elementValue.ProductID)
			fmt.Printf("Номер заказа: %d\n", elementValue.OrderID)
			fmt.Printf("Количество: %d\n", elementValue.Quantity)
			if len(elementValue.AdditionalShelfID) > 0 {
				fmt.Printf("Дополнительные стеллажи: %v\n", strings.Join(elementValue.AdditionalShelfID, ", "))
			}
			fmt.Println("---------------------------")
		}
		fmt.Println("======================")
	}
}
