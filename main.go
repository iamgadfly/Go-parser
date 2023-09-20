package main

import (
	"bufio"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goParser/src/wb"
	"os"
	"strings"
)

var link string

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type WbProduct struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	SalePrice int64  `db:"sale_price"`
	Price     int64  `db:"price"`
	Color     string `db:"color"`
	Type      string `db:"type_product"`
	ShopId    int    `db:"shop_id"`
	CatId     int    `db:"category_id"`
}

func main() {
	fmt.Print("Enter link: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return
	}
	link = strings.TrimSuffix(input, "\n")
	//ozon.Parse(&link)

	//wb.ParseOneProduct(&link) // парссинг одного товара
	wb.ParseProducts(&link) // паррсинг категории

	//https://catalog.wb.ru/catalog/electronic15/catalog?appType=64&curr=rub&dest=-3889739&lang=ru&locale=ru&page=1&subject=2290
	//id := "3889739"
	//wb.ParseCategory(&id)

}
