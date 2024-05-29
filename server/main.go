package main

import (
  "fmt"
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/rs/cors"
  . "server/utils"
)


type MoveReq struct {
  Piece string `json:"Piece"`
  File uint8 `json:"File"` 
  Rank uint8`json:"Rank"`
}

type MoveRes struct {
  File uint8 `json:"File"` 
  Rank uint8`json:"Rank"`
}

var bitboard Bitboard 
var Constants BitboardConstants 

func Moves(context *gin.Context) {
  var move MoveReq

  if err := context.BindJSON(&move); err != nil {
    fmt.Println("Invalid request body")
    fmt.Println(err)
    return
  }

  intPos := positionFromRowCol(move.Rank, move.File)

  fmt.Println("Piece:", move.Piece)
  //DisplayPieceLocation(intPos)

  validMoves := GetValidMoves(move.Piece, intPos, &bitboard, &Constants)

  moveList := make([]MoveRes, len(validMoves))
  for i, pos := range validMoves {
    row, col := rowColFromPosition(pos)
    moveList[i] = MoveRes{Rank: row, File: col} 
  }

  context.IndentedJSON(http.StatusOK, moveList) 
}

func positionFromRowCol(row, col uint8) uint64 {
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


func main() {
  Constants = BitboardConstants {
    NotAFile: 0x0101010101010101,
    NotBFile: 0x0202020202020202,
    NotGFile: 0x4040404040404040,
    NotHFile: 0x8080808080808080,
    NotABFile: 0x0303030303030303,
    NotGHFile: 0xC0C0C0C0C0C0C0C0,
  }

  InitBoard(&bitboard)
  PrintBoard(&bitboard)

  router := gin.Default()
  router.POST("/moves", Moves)
  handler := cors.Default().Handler(router)

  http.ListenAndServe("localhost:8080", handler)
}

