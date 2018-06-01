package canvas

import (
	"github.com/strongo/db"
	"context"
	"fmt"
	"time"
	"github.com/strongo-games/arena/arena-go"
	"errors"
)

var (
	ErrOldRound = errors.New("old round")
	ErrUnknownRound = errors.New("unknown round")
)

func MakeMove(c context.Context, now time.Time, database db.Database, round int, boardID, userID, move string) (board Board, err error) {
	board.ID = boardID
	if err = database.Get(c, &board); err != nil {
		if db.IsNotFound(err) {
			err = nil
			// New canvas
			if round != 1 {
				err = fmt.Errorf("round should be 1, got: %v", round)
				return
			}
			board.BoardEntity = &BoardEntity{
				Round: round,
				Created: now,
				UserIDs: []string{userID},
				UserMoves: arena.CommaSeparatedValuesList(move),
				UserTimes: []time.Time{now},
			}
		}
		return
	} else if err == nil {
		if round < board.Round {
			err = ErrOldRound
			return
		} else if round > board.Round {
			err = ErrUnknownRound
			return
		}
	}
	userIDsCount := len(board.UserIDs)
	userMovesCount := board.UserMoves.Count()
	if userMovesCount > userIDsCount  {
		err = fmt.Errorf("userMovesCount > userIDsCount: %v > %v", userMovesCount, userIDsCount)
		return
	}
	switch userIDsCount {
	case 1:
		if userID == board.UserIDs[0] {
			board.UserMoves = arena.CommaSeparatedValuesList(move)
		} else {
			board.UserIDs = append(board.UserIDs, userID)
			board.UserMoves = board.UserMoves.Add(move)
			board.UserTimes = append(board.UserTimes, now)
		}
	case 2:
		switch userMovesCount {
		case 0:
			switch userID {
			case board.UserIDs[0]:
				board.UserMoves = arena.CommaSeparatedValuesList(move)
			case board.UserIDs[1]:
				board.UserMoves = arena.CommaSeparatedValuesList(","+ move)
			default:

			}
		case 1:
			switch userID {
			case board.UserIDs[0]:
				board.UserMoves = board.UserMoves.Set(0, move)
			case board.UserIDs[1]:
				board.UserMoves = board.UserMoves.Add(move)
				board.UserTimes = append(board.UserTimes, now)
			default:
				err = fmt.Errorf("wrong user id=%v", userID)
			}
		case 2:
			switch userID {
			case board.UserIDs[0]:
				board.UserMoves = board.UserMoves.Set(0, move)
			case board.UserIDs[1]:
				board.UserMoves = board.UserMoves.Set(1, move)
			default:
				err = fmt.Errorf("wrong user id=%v", userID)
			}
		}
	}
	return
}

func NextRound(board Board) {
	board.UserMoves = ""
	board.UserTimes = []time.Time{}
	board.Round += 1
}