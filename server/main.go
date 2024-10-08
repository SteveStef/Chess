package main

import (
	"fmt"
	"net/http"
	. "server/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

type MoveRes struct {
	File uint8 `json:"File"`
	Rank uint8 `json:"Rank"`
}

var PieceMap = map[string]uint8{
	"wp": WHITE_PAWN,
	"wr": WHITE_ROOK,
	"wn": WHITE_KNIGHT,
	"wb": WHITE_BISHOP,
	"wq": WHITE_QUEEN,
	"wk": WHITE_KING,
	"bp": BLACK_PAWN,
	"br": BLACK_ROOK,
	"bn": BLACK_KNIGHT,
	"bb": BLACK_BISHOP,
	"bq": BLACK_QUEEN,
	"bk": BLACK_KING,
}

var bitboard Bitboard

func positionFromRowCol(row uint8, col uint8) uint64 {
	p := uint64(1) << 63
	p = p >> (8 * row)
	p = p >> (8 - col - 1)
	return p
}

func rowColFromPosition(pos uint64) (row, col uint8) {
	var r, c uint8
	for r = 0; r <= 8; r++ {
		for c = 0; c <= 8; c++ {
			if ((pos >> (64 - (8*r + 8 - c))) & 1) == 1 {
				return r, c
			}
		}
	}
	return 0, 0
}

func Moves(context *gin.Context) {
	var move struct {
		Piece string `json:"Piece"`
		File  uint8  `json:"File"`
		Rank  uint8  `json:"Rank"`
	}

	if err := context.BindJSON(&move); err != nil {
		fmt.Println("Invalid request body")
		fmt.Println(err)
		return
	}

	intPos := positionFromRowCol(move.Rank, move.File)
	validMoves := GetValidMoves(PieceMap[move.Piece], intPos, &bitboard)

	moveList := make([]MoveRes, len(validMoves))
	for i, pos := range validMoves {
		row, col := rowColFromPosition(pos)
		moveList[i] = MoveRes{Rank: row, File: col}
	}

	context.IndentedJSON(http.StatusOK, moveList)
}

func GenerateBoard(context *gin.Context) {
	var fen string
	if err := context.BindJSON(&fen); err != nil {
		fmt.Println("Invalid request body")
		fmt.Println(err)
		return
	}
	context.IndentedJSON(http.StatusOK, "Board generated")
}

func MovePiece(context *gin.Context) {
	var move struct {
		Piece   string `json:"Piece"`
		File    uint8  `json:"File"`
		Rank    uint8  `json:"Rank"`
		NewFile uint8  `json:"NewFile"`
		NewRank uint8  `json:"NewRank"`
	}

	if err := context.BindJSON(&move); err != nil {
		fmt.Println("Invalid request body")
		fmt.Println(err)
		context.IndentedJSON(http.StatusBadRequest, "Invalid request body")
		return
	}

	intPos := positionFromRowCol(move.Rank, move.File)
	newPos := positionFromRowCol(move.NewRank, move.NewFile)

	pieceType := PieceMap[move.Piece]

	MakeMove(pieceType, intPos, newPos, &bitboard)
	board := GetBoardState(&bitboard)

	context.IndentedJSON(http.StatusOK, board)
}

func main() {

	InitBoard(&bitboard)
	PrintGame(&bitboard)

	//tmp := uint64(1) << 3 << 8

	//fmt.Println(tmp)
	//DisplayPieceLocation(tmp)

	ShowMailbox(&bitboard)

	router := gin.Default()

	router.POST("/moves", Moves)
	router.POST("/place", MovePiece)
	//router.POST("initboard", GenerateBoard)

	handler := cors.Default().Handler(router)

	http.ListenAndServe("localhost:8080", handler)
}
