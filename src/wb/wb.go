package wb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type WbProduct struct {
	Id        int         `db:"id"`
	Name      string      `db:"name"`
	SalePrice int64       `db:"sale_price"`
	Price     int64       `db:"price"`
	Color     string      `db:"color"`
	Type      string      `db:"type_product"`
	ShopId    int         `db:"shop_id"`
	CatId     interface{} `db:"category_id"`
}

func ParseCategory(id *string) {
	fmt.Println(*id)
}

func ParseOneProduct(link *string) {
	id := GetId(link)
	url := "https://card.wb.ru/cards/detail?appType=1&curr=rub&dest=-1257786&regions=80,38,83,4,64,33,68,70,30,40,86,75,69,1,31,66,110,48,22,71,114&spp=35&nm=" + id
	resp := sendData(url)
	ShopId, _ := strconv.Atoi(id)

	body, _ := ioutil.ReadAll(resp.Body)

	Name := gjson.Get(string(body), "data.products.0.name").String()
	SalePrice := gjson.Get(string(body), "data.products.0.salePriceU").Int() / 100
	Price := gjson.Get(string(body), "data.products.0.priceU").Int() / 100
	Color := gjson.Get(string(body), "data.products.0.colors.0.name").String()

	db, err := sqlx.Connect("mysql", "root:@(127.0.0.1:3306)/go_parser")
	if err != nil {
		panic(err)
	}

	check := WbProduct{}
	//db.Select(&check, "SELECT * FROM products WHERE shop_id=")
	err = db.Get(&check, "SELECT * FROM products WHERE shop_id=?", id)
	data := getDataFromJson(Name, int(Price), int(SalePrice), Color, ShopId)

	if err != nil {
		_, err = db.NamedExec(`INSERT INTO products (name,price,sale_price,color,shop_id) VALUES (:name,:price,:sale_price,:color,:shop_id)`, data)
	} else {
		fmt.Println("этот товар уже есть!")
		_, err = db.NamedExec(`UPDATE products SET price=:price,sale_price=:sale_price WHERE shop_id=:shop_id`, data)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("%+v\n", data)
	//fmt.Printf("%+v\n", product)
}

func ParseProducts(link *string) {
	//id := GetId(link)
	//url := "https://card.wb.ru/cards/detail?appType=1&curr=rub&dest=-1257786&regions=80,38,83,4,64,33,68,70,30,40,86,75,69,1,31,66,110,48,22,71,114&spp=35&nm=" + id
	//"https://catalog.wb.ru/catalog/electronic15/catalog?appType=64&curr=rub&dest=-3889739&lang=ru&locale=ru&page=1&subject=2290"

	db, _ := sqlx.Connect("mysql", "root:@(127.0.0.1:3306)/go_parser")

	resp := sendData(*link) // + "&page=" + strconv.Itoa(i)
	body, _ := ioutil.ReadAll(resp.Body)
	data, products := gjson.Get(string(body), "data.products").String(), []map[string]interface{}{}
	for i := 0; i < 2; i++ {
		iString := strconv.Itoa(i)
		check := gjson.Get(data, iString).String()
		if check == "" {
			i = 0
			continue
		}
		Name := gjson.Get(data, iString+".name").String()
		SalePrice := gjson.Get(data, iString+".salePriceU").Int() / 100
		Price := gjson.Get(data, iString+".priceU").Int() / 100
		Color := gjson.Get(data, iString+".colors.0.name").String()
		ShopId := gjson.Get(data, iString+".id").Int()
		product := getDataFromJson(Name, int(Price), int(SalePrice), Color, int(ShopId))

		products = append(products, product)

	}
	fmt.Println(products)
	_, err := db.NamedExec(`INSERT INTO products (name,price,sale_price,color,shop_id)
	VALUES (:name, :price, :sale_price, :color, :shop_id) ON DUPLICATE KEY UPDATE sale_price=sale_price AND price=price`, products)
	if err != nil {
		panic(err)
	}
	//
	//fmt.Printf("%+v\n", products)

}

//func getDataFromUrl(respData chan string, link *string) {
//	for i := 1; i > 0; i++ {
//		resp := sendData(*link + "&page=" + strconv.Itoa(i))
//		body, _ := ioutil.ReadAll(resp.Body)
//		data := gjson.Get(string(body), "data.products").String()
//		respData <- data
//	}
//}
//
//func insertData(respData chan string) {
//	db, _ := sqlx.Connect("mysql", "root:@(127.0.0.1:3306)/go_parser")
//	data := <-respData
//	products := []map[string]interface{}
//	for i := 0; i < 10; i++ {
//		iString := strconv.Itoa(i)
//		check := gjson.Get(data, iString).String()
//		if check == "" {
//			i = 0
//			continue
//		}
//		Name := gjson.Get(data, iString+".name").String()
//		SalePrice := gjson.Get(data, iString+".salePriceU").Int() / 100
//		Price := gjson.Get(data, iString+".priceU").Int() / 100
//		Color := gjson.Get(data, iString+".colors.0.name").String()
//		ShopId := gjson.Get(data, iString+".id").Int()
//		product := getDataFromJson(Name, int(Price), int(SalePrice), Color, int(ShopId))
//		products = append(products, product)
//	}
//
//	_, err := db.NamedExec(`INSERT INTO products (name,price,sale_price,color,shop_id)
//	   VALUES (:name, :price, :sale_price, :color, :shop_id) ON DUPLICATE KEY UPDATE sale_price=sale_price AND price=price`, products)
//	if err != nil {
//		panic(err)
//	}
//}

func createProduct(Name string, Price int64, SalePrice int64, Color string, ShopId int) WbProduct {
	return WbProduct{
		Name:      Name,
		Price:     Price,
		SalePrice: SalePrice,
		Color:     Color,
		ShopId:    ShopId,
	}
}

func getDataFromJson(Name string, Price int, SalePrice int, Color string, ShopId int) map[string]interface{} {
	return map[string]interface{}{
		"name":       Name,
		"price":      Price,
		"sale_price": SalePrice,
		"color":      Color,
		"shop_id":    ShopId,
	}
}

func GetId(link *string) string {
	raw := strings.Split(*link, "/")
	return raw[len(raw)-2]
}

func sendData(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	return resp
}
