// Download and print floodgate status
package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
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
	var spool *[]string
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
					spool = &headers
				case tt == html.StartTagToken && t.Data == "td":
					spool = &data
				case tt == html.TextToken && spool != nil:
					field := strings.TrimSpace(t.Data)
					if len(field) == 0 {
						continue
					}
					*spool = append(*spool, field)
				}
			}
		}
	}
	for x := range data {
		fmt.Printf("%s: %s\n", headers[x%len(headers)], data[x])
	}
}
