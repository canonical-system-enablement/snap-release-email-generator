/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"bytes"
	"strings"
)

type MessageBuilder struct {
	Publisher string
	Changes   map[string][]string
	Bileto    map[string]string
}

func NewMessageBuilder() *MessageBuilder {
	p := new(MessageBuilder)
	p.Changes = make(map[string][]string)
	p.Bileto = make(map[string]string)
	return p
}

func (b *MessageBuilder) SetPublisher(publisher string) {
	b.Publisher = publisher
}

func (b *MessageBuilder) AddChange(snap string, version string, change string) {
	key := snap + " " + version
	b.Changes[key] = append(b.Changes[key], change)
}

func (b *MessageBuilder) AddBiletoURL(snap string, bileto string) {
	b.Bileto[snap] = bileto
}

func (b *MessageBuilder) ConstructMessage() string {
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

	for snap, changes := range b.Changes {
		changelog.WriteString("\n")
		changelog.WriteString(snap + ":\n")
		changelog.WriteString("\n")
		for _, change := range changes {
			changelog.WriteString("* " + change + "\n")
		}
	}

	for snap, url := range b.Bileto {
		bileto.WriteString("- " + snap + ": " + url + "\n")
	}

	emailBody = strings.Replace(emailBody, "{snap changelog}", changelog.String(), 1)
	emailBody = strings.Replace(emailBody, "{bileto url}", bileto.String(), 1)
	emailBody = strings.Replace(emailBody, "{publisher}", b.Publisher, 1)

	return emailBody
}
