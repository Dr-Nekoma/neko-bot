package BOT

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	md "github.com/JohannesKaufmann/html-to-markdown"

	"neko-bot/API"
)

const baseUserURL = "https://hacker-news.firebaseio.com/v0/user/"
const profileId = "whoishiring"
const baseItemURL = "https://hacker-news.firebaseio.com/v0/item/"
const complementURL = ".json?print=pretty"
const profileURL = baseUserURL + profileId + complementURL

type id int

func HackerNewsJobs(keySentence string, howMany int) []string {
	id, err := queryProfileAPI()

	if err != nil {
		log.Print(err)
		return []string{err.Error()}
	}

	post, err := querySubmissionAPI(id)

	if err != nil {
		log.Print(err)
		return []string{err.Error()}
	}

	chAPI := make(chan int)
	chTrs := make(chan API.HackerNewsContentComment, 25)
	chStr := make(chan API.HackerNewsContentComment, 50)
	doneStr := make(chan bool)
	doneTrs := make(chan bool)
	var wg sync.WaitGroup

	var childSelectedRes []string
	var messages []string

	for i := 0; i < 20; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			fmt.Println("Waiting for ID")
			for id := range chAPI {
				comment, err := queryCommentAPI(id)
				if err != nil {
					log.Print(err)
					messages = append(messages, err.Error())
					break
				}
				chTrs <- comment
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
			if strings.Contains(strings.ToLower(child.Text), strings.ToLower(keySentence)) {
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
