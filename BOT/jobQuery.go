package BOT

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	md "github.com/JohannesKaufmann/html-to-markdown"

	"neko-bot/API"
)

const baseURL = "https://hacker-news.firebaseio.com/v0/item/"
const parentURL = baseURL + "31582796.json?print=pretty"

func HackerNewsJobs(keySentence string, howMany int) []string {
	post := queryParentAPI()

	chAPI := make(chan int)
	chTrs := make(chan API.HackerNewsContentChild, 25)
	chStr := make(chan API.HackerNewsContentChild, 50)
	doneStr := make(chan bool)
	doneTrs := make(chan bool)
	var wg sync.WaitGroup

	var childSelectedRes []string

	for i := 0; i < 20; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			fmt.Println("Waiting for ID")
			for id := range chAPI {
				chTrs <- queryChildAPI(id)
			}
		}()
	}

	go func() {
		fmt.Println("Adding to storage")
		for resp := range chStr {
			childSelectedRes = append(childSelectedRes, resp.Text)
		}
		close(doneStr)
	}()

	go func() {
		fmt.Println("Filtering and translating to Markdown")
		for child := range chTrs {
			if strings.Contains(child.Text, keySentence) {
				child.Text = translateHTMLToMarkdown(child.Text)
				chStr <- child
			}
		}
		close(doneTrs)
	}()

	for _, id := range post.Kids {
		chAPI <- id
	}

	close(chAPI)
	wg.Wait()
	close(chTrs)
	<-doneTrs
	close(chStr)
	<-doneStr

	return childSelectedRes[:howMany]

}

func translateHTMLToMarkdown(html string) string {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}

	return markdown
}

func queryParentAPI() API.HackerNewsContentParent {
	resp, err := http.Get(parentURL)

	if err != nil {
		log.Fatalln(err)
	}

	var response API.HackerNewsContentParent
	json.NewDecoder(resp.Body).Decode(&response)
	return response
}

func queryChildAPI(postId int) API.HackerNewsContentChild {
	url := fmt.Sprintf("%s%d.json?print=pretty", baseURL, postId)

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	var response API.HackerNewsContentChild
	json.NewDecoder(resp.Body).Decode(&response)
	return response
}
