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
	"fmt"
	"github.com/VojtechVitek/go-trello"
	"strings"
)

type TrelloProxy struct {
	trello trello.Client
	board  trello.Board
}

func (t *TrelloProxy) Connect(secrets TrelloSecrets) error {
	// Connect to Trello
	trello, err := trello.NewAuthClient(secrets.AppKey, &secrets.Token)
	if err != nil {
		return err
	}

	// Get my board
	const sprintBoardId string = "5722667126680b7e86626557"
	board, err := trello.Board(sprintBoardId)
	if err != nil {
		return err
	}

	t.trello = *trello
	t.board = *board

	return nil
}

func (t *TrelloProxy) CardsOfSnapsToRelease() (cards []trello.Card, err error) {

	// Get the lists
	lists, err := t.board.Lists()
	if err != nil {
		return nil, err
	}

	for _, l := range lists {

		// Find the correct swimlane
		if !strings.Contains(l.Name, "Snap to Release") {
			continue
		}

		cards, err := l.Cards()
		if err != nil {
			return nil, err
		}
		for i, c := range cards {
			// Remove the README card
			if strings.Contains(c.Name, "How to use this column") {
				cards = cards[:i+copy(cards[i:], cards[i+1:])]
			}
		}
		return cards, nil
	}

	return nil, fmt.Errorf("Snap to Release swimlanie is no more")
}
