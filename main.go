package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/VojtechVitek/go-trello"
	"io/ioutil"
	"log"
	"strings"
)

type SnapToRelease struct {
	Name    string
	Version string
	Bileto  string
	Changes []string
}

type SnapsToRelease []SnapToRelease

type TrelloSecrets struct {
	AppKey string `json:"app_id"`
	Token  string `json:"token"`
}

var (
	trelloSecretsFile = flag.String("secrets", "trello_secrets.json", "Trello Secrets configuration")
	snapPublisher     = flag.String("publisher", "Simon", "Snap publisher's name")
)

func getSnapNameAndVersion(Name string) (string, string) {
	snapName := ""
	snapVersion := "null"

	tmp := strings.Split(Name, " - ")

	if len(tmp) > 0 {
		snapName = tmp[0]
	}

	if len(tmp) > 1 {
		snapVersion = tmp[1]
	}

	return snapName, snapVersion
}

func createEmail(str SnapsToRelease) string {

	var emailBody string = `
Hey everyone,

new releases of the following snaps were pushed into the candidate
channel.
{snap changelog}
---

An overview of which revisions / versions of the particular snaps are
available in which channel is available at
https://docs.google.com/document/d/1-rKAjr6FLUzt7oOtR_xcAEHJpntUPGpmixU6PV8K2KU/edit#

The snaps have passed our engineering QA and will now be tested by the
platform and commercial QA teams before the new versions are pushed to
the stable channel.

Bileto requests are:

{bileto url}

If you have any questions feel free to ping me.

regards,
{publisher}
`
	// 1. build changelog and bileto
	var changelog bytes.Buffer
	var bileto bytes.Buffer
	for _, snap := range str {
		bileto.WriteString("- " + snap.Name + ": " + snap.Bileto + "\n")
		changelog.WriteString("\n")
		changelog.WriteString(snap.Name + " " + snap.Version + ":\n")
		changelog.WriteString("\n")
		for _, change := range snap.Changes {
			changelog.WriteString("* " + change + "\n")
		}
	}
	emailBody = strings.Replace(emailBody, "{snap changelog}", changelog.String(), 1)
	emailBody = strings.Replace(emailBody, "{bileto url}", bileto.String(), 1)
	emailBody = strings.Replace(emailBody, "{publisher}", *snapPublisher, 1)

	return emailBody
}

func readConfig() (*TrelloSecrets, error) {
	trelloSecrets := new(TrelloSecrets)
	data, err := ioutil.ReadFile(*trelloSecretsFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &trelloSecrets)
	if err != nil {
		return nil, err
	}

	return trelloSecrets, nil
}

func main() {
	flag.Parse()

	tsec, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to Trello
	trello, err := trello.NewAuthClient(tsec.AppKey, &tsec.Token)
	if err != nil {
		log.Fatal(err)
	}

	// Get my board
	const sprintBoardId string = "5722667126680b7e86626557"
	board, err := trello.Board(sprintBoardId)
	if err != nil {
		log.Fatal(err)
	}

	// Get the lists
	lists, err := board.Lists()
	if err != nil {
		log.Fatal(err)
	}

	str := SnapsToRelease{}

	for _, l := range lists {

		// We will egzamine the contents of Snaps To Release list so
		// skip the rest
		if !strings.Contains(l.Name, "Snap to Release") {
			continue
		}

		cards, err := l.Cards()
		if err != nil {
			log.Fatal(err)
		}

		// For each card on the "Snap to Release" swimlane,
		// capture the snap name and it's version from card title.
		// Then traverse the checklist for each card and capture the
		// names of items on checklist called "Changes"

		for _, c := range cards {
			// Skip the README card
			if strings.Contains(c.Name, "How to use this column") {
				continue
			}

			tmp := SnapToRelease{}

			// Get the data from card

			snapName, snapVersion := getSnapNameAndVersion(c.Name)
			tmp.Name = snapName
			tmp.Version = snapVersion

			checklists, err := c.Checklists()
			if err != nil {
				log.Fatal(err)
			}

			for _, clist := range checklists {

				if strings.Contains(clist.Name, "Changes") {
					for _, citem := range clist.CheckItems {
						tmp.Changes = append(tmp.Changes, citem.Name)
					}
				}

				if strings.Contains(clist.Name, "Bileto") {
					for _, bitem := range clist.CheckItems {
						tmp.Bileto = bitem.Name
					}
				}

			}

			str = append(str, tmp)
		}
	}

	email := createEmail(str)
	fmt.Println(email)
}
