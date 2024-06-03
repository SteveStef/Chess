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
var Constants BitboardConstants 

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
  validMoves := GetValidMoves(move.Piece, intPos, &bitboard, &Constants)

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

  fmt.Println(move)
  context.IndentedJSON(http.StatusOK, "Piece moved")
}

func main() {
  Constants = BitboardConstants {
    A_File: 0x0101010101010101,
    B_File: 0x0202020202020202,
    G_File: 0x4040404040404040,
    H_File: 0x8080808080808080,
    AB_File: 0x0303030303030303,
    GH_File: 0xC0C0C0C0C0C0C0C0,

    RANK_1: 0xFF00000000000000,
    RANK_8: 0x00000000000000FF,

    OnBoard: 0xFFFFFFFFFFFFFFFF,
  }

  InitBoard(&bitboard)
  PrintBoard(&bitboard)

  router := gin.Default()
  router.POST("/moves", Moves)
  router.POST("/place", MovePiece)
  //router.POST("initboard", GenerateBoard)
  handler := cors.Default().Handler(router)

  http.ListenAndServe("localhost:8080", handler)
}

