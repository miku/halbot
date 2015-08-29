package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/danryan/hal"
	_ "github.com/danryan/hal/adapter/irc"
	// _ "github.com/danryan/hal/adapter/shell"
	_ "github.com/danryan/hal/store/memory"
)

type Document struct {
	Title    string `json:"title"`
	SourceID string `json:"source_id"`
}

type SolrResponse struct {
	Response struct {
		NumFound int        `json:"numFound"`
		Docs     []Document `json:"docs"`
	} `json:"response"`
}

func solrQuery(s string) (SolrResponse, error) {
	vals := url.Values{}
	vals.Add("wt", "json")
	vals.Add("q", s)
	link := fmt.Sprintf(`%s/select?%s`, os.Getenv("PARROT_SOLR_URL"), vals.Encode())
	r, err := http.Get(link)
	if err != nil {
		return SolrResponse{}, err
	}
	defer r.Body.Close()
	var sr SolrResponse
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&sr)
	if err != nil {
		return SolrResponse{}, err
	}
	return sr, nil
}

// queryHandler takes a query and executes it on main site
var queryHandler = hal.Hear(`hal ai q(\d)? (.+)`, func(res *hal.Response) error {
	sr, err := solrQuery(res.Match[2])

	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d in AI for %s", sr.Response.NumFound, res.Match[2]))

	if res.Match[1] != "" {
		size, err := strconv.Atoi(res.Match[1])
		if err != nil {
			return err
		}
		if size > len(sr.Response.Docs) {
			size = len(sr.Response.Docs)
		}
		if size > 0 {
			buf.WriteString(" -- ")
		}
		var items []string
		for i := 0; i < size; i++ {
			doc := sr.Response.Docs[i]
			items = append(items, fmt.Sprintf("(%d) %s [%s]", i+1, doc.Title, doc.SourceID))
		}
		buf.WriteString(strings.Join(items, ", "))
	}

	return res.Send(buf.String())
})

var pingHandler = hal.Hear(`ping`, func(res *hal.Response) error {
	return res.Send("PONG")
})

func run() int {
	robot, err := hal.NewRobot()
	if err != nil {
		hal.Logger.Error(err)
		return 1
	}

	robot.Handle(
		pingHandler,
		queryHandler,
	)

	if err := robot.Run(); err != nil {
		hal.Logger.Error(err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
