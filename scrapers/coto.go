package scrapers

import (
	"strings"

	"ratoneando/cores/html"
	"ratoneando/products"
	"ratoneando/utils/numbers"
)

func coalesce(strs ...string) string {
	for _, str := range strs {
		if str != "" {
			return str
		}
	}
	return ""
}

func extractor(element *html.ElementWrapper, doc *html.DocumentWrapper) products.ExtendedSchema {
	id, _ := element.Attr("id")
	id = strings.Replace(id, "li_prod", "", -1)

	name := element.Find("div.descrip_full").Text()

	// rawPrice := coalesce(element.Find(".price_discount").Text(), element.Find(".info_discount span.atg_store_newPrice").Text())
	rawPrice := coalesce(
		strings.ReplaceAll(element.Find(".price_discount_gde").Text(), ".", ","),
		element.Find("span.atg_store_newPrice").First().Text(),
	)

	price, _ := numbers.ParseMoney(rawPrice)

	link, _ := element.Find("a").Attr("href")
	link = "https://www.cotodigital3.com.ar" + link

	image, _ := element.Find("span.atg_store_productImage > img").Attr("src")
	unavailable := element.Find("div.product_not_available").Length() > 0

	return products.ExtendedSchema{
		ID:          id,
		Name:        name,
		Link:        link,
		Image:       image,
		Price:       price,
		Unavailable: unavailable,
	}
}

func Coto(query string) ([]products.Schema, error) {
	return html.Core(html.CoreProps{
		Query:   query,
		BaseUrl: "https://www.cotodigital3.com.ar",
		Source:  "coto",
		SearchPattern: func(q string) string {
			return "/sitios/cdigi/browse?Ntt=" + q + "&_DARGS=%2Fsitios%2Fcartridges%2FSearchBox%2FSearchBox.jsp"
		},
		ContainerSelector: "ul#products",
		ProductSelector:   "li[id^='li_prod']",
		Extractor:         extractor,
	})
}
