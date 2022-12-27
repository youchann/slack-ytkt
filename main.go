package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	var s = flag.String("token", "", "The Slack User OAuth Token.")
	var u = flag.String("user", "", "The Slack User Name")
	flag.Parse()

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://slack.com/api/search.messages", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+*s)
	q := url.Values{}
	q.Add("query", "from:@"+*u+" ytkt*")
	q.Add("count", "5")
	q.Add("sort", "timestamp")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	text := ""
	for _, m := range result.Messages.Matches {
		text += m.Text + "\n"
	}
	text = strings.Replace(text, "ytkt", "", -1)
	text = regexp.MustCompile("(•|◦|▪︎)").ReplaceAllString(text, "-")
	if err := copyToClipboard(text); err != nil {
		panic(err)
	}

	fmt.Println("Successfully copied to clipboard!!!!!")
}

// NOTE: Only for Mac OS.
func copyToClipboard(content string) error {
	cmd := exec.Command("pbcopy")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(content)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return cmd.Wait()
}

type Response struct {
	Ok       bool   `json:"ok"`
	Query    string `json:"query"`
	Messages struct {
		Total      int `json:"total"`
		Pagination struct {
			TotalCount int `json:"total_count"`
			Page       int `json:"page"`
			PerPage    int `json:"per_page"`
			PageCount  int `json:"page_count"`
			First      int `json:"first"`
			Last       int `json:"last"`
		} `json:"pagination"`
		Paging struct {
			Count int `json:"count"`
			Total int `json:"total"`
			Page  int `json:"page"`
			Pages int `json:"pages"`
		} `json:"paging"`
		Matches []struct {
			Iid     string  `json:"iid"`
			Team    string  `json:"team"`
			Score   float64 `json:"score"`
			Channel struct {
				ID                 string        `json:"id"`
				IsChannel          bool          `json:"is_channel"`
				IsGroup            bool          `json:"is_group"`
				IsIm               bool          `json:"is_im"`
				Name               string        `json:"name"`
				IsShared           bool          `json:"is_shared"`
				IsOrgShared        bool          `json:"is_org_shared"`
				IsExtShared        bool          `json:"is_ext_shared"`
				IsPrivate          bool          `json:"is_private"`
				IsMpim             bool          `json:"is_mpim"`
				PendingShared      []interface{} `json:"pending_shared"`
				IsPendingExtShared bool          `json:"is_pending_ext_shared"`
			} `json:"channel"`
			Type     string `json:"type"`
			User     string `json:"user"`
			Username string `json:"username"`
			Ts       string `json:"ts"`
			Blocks   []struct {
				Type     string `json:"type"`
				BlockID  string `json:"block_id"`
				Elements []struct {
					Type     string `json:"type"`
					Elements []struct {
						Type  string `json:"type"`
						Text  string `json:"text"`
						Style struct {
							ClientHighlight bool `json:"client_highlight"`
						} `json:"style"`
					} `json:"elements"`
				} `json:"elements"`
			} `json:"blocks"`
			Text        string `json:"text"`
			Permalink   string `json:"permalink"`
			NoReactions bool   `json:"no_reactions"`
		} `json:"matches"`
	} `json:"messages"`
}
