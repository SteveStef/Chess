package main

import (
  "fmt"
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/rs/cors"
  . "server/utils"
)

type MoveRes struct {
  File uint8 `json:"File"` 
  Rank uint8`json:"Rank"`
}

var bitboard Bitboard 

func positionFromRowCol(row uint8, col uint8) uint64 {
  p := uint64(1) << 63 
  p = p >> (8 * row)
  p = p >> (8 - col-1)
  return p
}

func rowColFromPosition(pos uint64) (row, col uint8) {
  var r, c uint8
  for r = 0; r <= 8; r++ {
    for c = 0; c <= 8; c++ {
      if ((pos >> (64 - (8 * r + 8 - c))) & 1) == 1 {
        return r, c
      }
    }
  }
  return 0, 0
}

func Moves(context *gin.Context) {
  var move struct {
    Piece string `json:"Piece"`
    File uint8 `json:"File"` 
    Rank uint8`json:"Rank"`
  }

  if err := context.BindJSON(&move); err != nil {
    fmt.Println("Invalid request body")
    fmt.Println(err)
    return
  }

  intPos := positionFromRowCol(move.Rank, move.File)
  validMoves := GetValidMoves(move.Piece, intPos, &bitboard)

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
  GenerateBoardFromFen(fen, &bitboard)
  context.IndentedJSON(http.StatusOK, "Board generated")
}

func MovePiece(context *gin.Context) {
  var move  struct {
    Piece string `json:"Piece"`
    File uint8 `json:"File"`
    Rank uint8 `json:"Rank"`
    NewFile uint8 `json:"NewFile"`
    NewRank uint8 `json:"NewRank"`
  }

  if err := context.BindJSON(&move); err != nil {
    fmt.Println("Invalid request body")
    fmt.Println(err)
    context.IndentedJSON(http.StatusBadRequest, "Invalid request body")
    return
  }
  
  intPos := positionFromRowCol(move.Rank, move.File)
  newPos := positionFromRowCol(move.NewRank, move.NewFile)

  MakeMove(move.Piece, intPos, newPos, &bitboard)
  PrintBoard(&bitboard)
  board := GetBoardState(&bitboard)

  context.IndentedJSON(http.StatusOK, board)
}

func main() {

  InitBoard(&bitboard)
  PrintBoard(&bitboard)
  //DisplayPieceLocation(uint64(1) << 63)

  router := gin.Default()

  router.POST("/moves", Moves)
  router.POST("/place", MovePiece)
  //router.POST("initboard", GenerateBoard)

  handler := cors.Default().Handler(router)

  http.ListenAndServe("localhost:8080", handler)
}



