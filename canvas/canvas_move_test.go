package canvas

import (
	"testing"
	"context"
	"github.com/strongo/slices"
	"time"
)

func TestMakeMove(t *testing.T) {
	testCases := []struct {
		round int
		userID          string
		move            string
		expectedUserIDs []string
		expectedMoves   []string
	}{
		{1, "u1", "rock", []string{"u1"}, []string{"rock"}},
		{1, "u1", "paper", []string{"u1"}, []string{"paper"}},
		{1, "u2", "scissors", []string{"u1", "u2"}, []string{"paper", "scissors"}},
		{expectedUserIDs: []string{"u1", "u2"}, expectedMoves: []string{}},
		{2, "u2", "rock", []string{"u1", "u2"}, []string{"", "rock"}},
		{2, "u1", "paper", []string{"u1", "u2"}, []string{"paper", "rock"}},
	}

	var board Board

	c := context.Background()

	var err error

	database := newMockDB(c)

	for i, testCase := range testCases {
		if testCase.round == 0 {
			NextRound(board)
		} else {
			board, err = MakeMove(c, time.Now(), database, testCase.round, "abc", testCase.userID, testCase.move)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		}
		if board.BoardEntity == nil {
			t.Fatalf("case #%v: board.BoardEntity == nil", i+1)
		}
		if !slices.EqualStrings(board.UserIDs, testCase.expectedUserIDs) {
			t.Fatalf("case #%v: Unexpected UserIDs=%v, expected: %v", i+1, board.UserIDs, testCase.expectedUserIDs)
		}
		database.Update(c, &board)
	}
}
