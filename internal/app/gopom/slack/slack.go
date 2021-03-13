package slack

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func callSlack(urlS string) error {
	form := url.Values{}
	token := viper.GetString("SLACK_TOKEN")
	form.Add("token", token)

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlS, strings.NewReader(form.Encode()))

	if err != nil {
		return fmt.Errorf("error when creating request to slack endpoint - %s", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	resp, err := client.Do(req)

	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Fatalf("error when closing response body - %s", err)
			}
		}()
	}

	if err != nil {
		return fmt.Errorf("error when calling slack endpoint - %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error when reading body - %s", err)
		}
		return fmt.Errorf("slack responded with error message - %s", bodyBytes)
	}

	// return empty error
	return nil
}

func SetDnd(durationMinutes int) {
	log.Println("Setting DND on slack.")

	urlS := fmt.Sprintf("https://slack.com/api/dnd.setSnooze?num_minutes=%d", durationMinutes)

	err := callSlack(urlS)
	if err != nil {
		fmt.Println(err)
	}
}

func EndDnd() {
	log.Println("Ending DND on slack.")

	urlS := "https://slack.com/api/dnd.endSnooze"

	err := callSlack(urlS)
	if err != nil {
		fmt.Println(err)
	}
}
