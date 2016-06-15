// Download and print floodgate status
package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	//"strings"
)

func main() {
	resp, err := http.Get("http://floodstatus.lcra.org")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	headers := []string{}
	data := []string{}
	// this for loop searches for table#GridView2
top:
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			break
		case tt == html.StartTagToken:
			t := z.Token()

			isTable := t.Data == "table"
			if !isTable {
				continue
			}
			id := ""
			for _, attr := range t.Attr {
				if attr.Key == "id" {
					id = attr.Val
					break
				}
			}
			if id != "GridView2" {
				continue
			}
			// This for loop processes until </table> is reached
			for {
				tt := z.Next()
				t := z.Token()
				switch {
				case tt == html.EndTagToken && t.Data == "table":
					break top
				case tt == html.StartTagToken && t.Data == "th":
					z.Next()
					t := z.Token()
					headers = append(headers, t.Data)
				case tt == html.StartTagToken && t.Data == "td":
					z.Next()
					t := z.Token()
					data = append(data, t.Data)
				}
			}
		}
	}
	for x := range data {
		fmt.Printf("%s: %s\n", headers[x%len(headers)], data[x])
	}
}
