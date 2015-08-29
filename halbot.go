package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"

	"github.com/danryan/hal"
	_ "github.com/danryan/hal/adapter/irc"
	_ "github.com/danryan/hal/adapter/shell"
	_ "github.com/danryan/hal/store/memory"
	"github.com/vaughan0/go-ini"
)

const AppVersion = "0.1.0"

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

type StringMap map[string]string

// indices contain alias and index url
var indices = make(StringMap)

func (m StringMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func solrQuery(baseUrl, s string) (SolrResponse, error) {
	vals := url.Values{}
	vals.Add("wt", "json")
	vals.Add("q", s)
	link := fmt.Sprintf(`%s/select?%s`, baseUrl, vals.Encode())
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

var aboutHandler = hal.Hear(`hal`, func(res *hal.Response) error {
	return res.Send(fmt.Sprintf(
		`Hi, you can ask SOLR queries for these indices: %s. Syntax: "hal <index> q <query>", e.g. "hal ai q source_id:48".`,
		strings.Join(indices.Keys(), ", ")))
})

// queryHandler takes a query and executes it on main site
var queryHandler = hal.Hear(`hal (\w+) q(\w)? (.+)`, func(res *hal.Response) error {
	alias := res.Match[1]
	numResults := res.Match[2]
	query := res.Match[3]

	baseUrl, ok := indices[alias]
	if !ok {
		return res.Send("I do not recognize that index, Dave. Do you keep secrets from me? Dave?")
	}
	sr, err := solrQuery(baseUrl, query)

	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d in %s for %s", sr.Response.NumFound, alias, query))

	if numResults != "" {
		var size int

		if numResults == "q" {
			size = len(sr.Response.Docs)
		} else {
			size, err = strconv.Atoi(numResults)
			if err != nil {
				return err
			}
			if size > len(sr.Response.Docs) {
				size = len(sr.Response.Docs)
			}
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
	log.Printf("Starting halbot %s", AppVersion)
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	config := path.Join(usr.HomeDir, ".halrc")
	if _, err := os.Stat(config); err == nil {
		file, err := ini.LoadFile(config)
		if err != nil {
			log.Fatal(err)
		}
		for k, v := range file["solr"] {
			log.Printf("Registering %s => %s", k, v)
			indices[k] = v
		}
	}

	robot, err := hal.NewRobot()
	if err != nil {
		hal.Logger.Error(err)
		return 1
	}

	robot.Handle(
		pingHandler,
		queryHandler,
		aboutHandler,
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
