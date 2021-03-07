package slack

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func SetDnd(durationMinutes int) {
	slackToken := os.Getenv("SLACK_TOKEN")
	fmt.Println("SetDnd")

	urlS := fmt.Sprintf("https://slack.com/api/dnd.setSnooze?num_minutes=%d", durationMinutes)

	form := url.Values{}
	form.Add("token", slackToken)

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlS, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	if err != nil {
		fmt.Println("Error when calling slack dnd.")
	}

	fmt.Println("starting dnd on slack")
}

func EndDnd() {
	slackToken := os.Getenv("SLACK_TOKEN")
	fmt.Println("SetDnd")

	urlS := fmt.Sprintf("https://slack.com/api/dnd.endSnooze")

	form := url.Values{}
	form.Add("token", slackToken)

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlS, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	if err != nil {
		fmt.Println("Error when calling slack dnd.")
	}

	fmt.Println("starting end dnd on slack")
}
