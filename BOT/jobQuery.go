package BOT

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	md "github.com/JohannesKaufmann/html-to-markdown"

	"neko-bot/API"
	"neko-bot/MSG"
)

const baseUserURL = "https://hacker-news.firebaseio.com/v0/user/"
const profileId = "whoishiring"
const baseItemURL = "https://hacker-news.firebaseio.com/v0/item/"
const complementURL = ".json?print=pretty"
const commentURL = "https://news.ycombinator.com/item?id="
const profileURL = baseUserURL + profileId + complementURL

type id int

type comment struct {
	body API.HackerNewsContentComment
	url  string
}

func HackerNewsJobs(keySentence string, howMany int) []MSG.Message {
	id, err := queryProfileAPI()

	if err != nil {
		log.Print(err)
		return []MSG.Message{{Body: err.Error(), Kind: "error"}}
	}

	post, err := querySubmissionAPI(id)

	if err != nil {
		log.Print(err)
		return []MSG.Message{{Body: err.Error(), Kind: "error"}}
	}

	chAPI := make(chan int)
	chTrs := make(chan comment, 25)
	chStr := make(chan comment, 50)
	doneStr := make(chan bool)
	doneTrs := make(chan bool)
	doneCnt := make(chan bool)
	var wg sync.WaitGroup

	var childSelectedRes []MSG.Message
	var requests uint64

	for i := 0; i < 20; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			fmt.Println("Waiting for ID")
			for id := range chAPI {
				content, err := queryCommentAPI(id)
				if err != nil {
					log.Print(err)
					return
				}
				if atomic.LoadUint64(&requests) >= uint64(howMany) {
					log.Print("I reached my limit!")
					return
				}
				chTrs <- comment{body: content, url: commentURL + fmt.Sprint(id)}
			}
		}()
	}

	go func() {
		fmt.Println("Adding to storage")
		for resp := range chStr {
			childSelectedRes = append(childSelectedRes, MSG.Message{Body: resp.body.Text, TitleLink: resp.url, Kind: "jobs"})
		}
		close(doneStr)
	}()

	go func() {
		fmt.Println("Filtering and translating to Markdown")
		for child := range chTrs {
			if strings.Contains(strings.ToLower(child.body.Text), strings.ToLower(keySentence)) {
				child.body.Text = translateHTMLToMarkdown(child.body.Text)
				log.Print("Adding one to the counter!")
				atomic.AddUint64(&requests, 1)
				if atomic.LoadUint64(&requests) >= uint64(howMany) {
					close(doneCnt)
				}
				chStr <- child
			}
		}
		close(doneTrs)
	}()

	for _, id := range post.Kids {
		select {
		case <-doneCnt:
			break
		case chAPI <- id:
			fmt.Println("ID Sent!")
		}
	}

	wg.Wait()
	close(chAPI)
	close(chTrs)
	<-doneTrs
	close(chStr)
	<-doneStr

	fmt.Println(len(childSelectedRes))
	if len(childSelectedRes) == 0 {
		return []MSG.Message{{Kind: "lackOfJobs"}}
	} else {
		return childSelectedRes
	}
}

func translateHTMLToMarkdown(html string) string {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}

	return markdown
}

func queryProfileAPI() (id, error) {
	resp, err := http.Get(profileURL)

	if err != nil {
		log.Print(err)
		return -1, errors.New("I couldn't reach profile data!")
	}

	var response API.HackerNewsProfile
	json.NewDecoder(resp.Body).Decode(&response)

	firstSubmission := -1
	for _, submission := range response.Submitted {
		potentialJobPost, err := querySubmissionAPI(id(submission))
		if err != nil {
			log.Print(err)
			break
		}
		if strings.Contains(potentialJobPost.Title, "Who is hiring?") {
			firstSubmission = submission
			break
		}
	}

	if firstSubmission == -1 {
		return -1, errors.New("I couldn't find any posts for hiring!")
	} else {
		return id(firstSubmission), nil
	}
}

func querySubmissionAPI(submissionId id) (API.HackerNewsContentSubmission, error) {
	url := baseItemURL + fmt.Sprint(submissionId) + complementURL
	log.Print(url)
	resp, err := http.Get(url)

	var response API.HackerNewsContentSubmission

	if err != nil {
		log.Print(err)
		return response, errors.New("I couldn't reach submission data!")
	}

	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}

func queryCommentAPI(postId int) (API.HackerNewsContentComment, error) {
	url := fmt.Sprintf("%s%d.json?print=pretty", baseItemURL, postId)

	resp, err := http.Get(url)

	var response API.HackerNewsContentComment

	if err != nil {
		log.Print(err)
		return response, errors.New("I couldn't reach comments data!")
	}

	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}
