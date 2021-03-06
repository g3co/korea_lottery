package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

type Game struct {
	Number int   `json:"number"`
	Digits []int `json:"digits"`
}

type GamesResults struct {
	Games []Game `json:"games"`
}

func (g GamesResults) Len() int           { return len(g.Games) }
func (g GamesResults) Swap(i, j int)      { g.Games[i], g.Games[j] = g.Games[j], g.Games[i] }
func (g GamesResults) Less(i, j int) bool { return g.Games[i].Number < g.Games[j].Number }

func main() {
	lastGame := 909

	rgx, err := regexp.Compile(`<div class="lot_num">[^.]+(?:<\/div>)`)
	if err != nil {
		panic(err.Error())
	}
	rgxNum, err := regexp.Compile(`>([\d]+)<`)
	if err != nil {
		panic(err.Error())
	}

	res := GamesResults{}

	var wg sync.WaitGroup

	res.Games = make([]Game, lastGame+1)

	for i := 1; i <= lastGame; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			link := "https://m.search.naver.com/search.naver?where=m&sm=mtb_etc&query=" + strconv.Itoa(num) + "%ED%9A%8C%EB%A1%9C%EB%98%90/"

			var resp *http.Response
			var err error
			var counter int

			for {
				counter++
				resp, err = http.Get(link)
				if err != nil {
					fmt.Println(err.Error())
					if counter > 10 {
						panic(err)
					}
					continue
				}
				break
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			s := rgx.Find(body)

			x := rgxNum.FindAll(s, len(s))

			g := Game{
				Digits: make([]int, 7),
				Number: num,
			}

			for n, y := range x {
				g.Digits[n], _ = strconv.Atoi(string(y[1 : len(y)-1]))
			}

			res.Games[num] = g
		}(i)
	}

	wg.Wait()

	sort.Sort(res)

	jsonRes, _ := json.Marshal(res)

	fmt.Println(string(jsonRes))
}
