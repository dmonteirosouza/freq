package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	colorReset := "\033[0m"
	colorRed := "\033[31m"

	sc := bufio.NewScanner(os.Stdin)

	jobs := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for domain := range jobs {
				resp, err := http.Get(domain)

				if err != nil {
					continue
				}

				body, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					fmt.Println(err)
				}

				sb := string(body)
				check_result := strings.Contains(sb, "alert(1)")

				if check_result != false {
					fmt.Println(string(colorRed), "Vulnerable To XSS:", domain, string(colorReset))
				}
			}
		}()
	}

	for sc.Scan() {
		domain := sc.Text()
		jobs <- domain
	}

	close(jobs)
	wg.Wait()
}
