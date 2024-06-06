package utils 

import (
  "fmt"
  "math"
)

type BitboardConstants struct {
  A_File uint64
  B_File uint64
  G_File uint64
  H_File uint64
  AB_File uint64
  GH_File uint64

  RANK_1 uint64
  RANK_2 uint64
  RANK_3 uint64
  RANK_4 uint64
  RANK_5 uint64
  RANK_6 uint64
  RANK_7 uint64
  RANK_8 uint64

  OnBoard uint64
}

type Bitboard struct {
  whitePawns, whiteKnights, whiteBishops, whiteRooks, whiteQueens, whiteKing uint64
  blackPawns, blackKnights, blackBishops, blackRooks, blackQueens, blackKing uint64
  Constants BitboardConstants
  castlingRights uint8
  whiteOnBottom bool
  enPassant uint64
  whiteTurn bool
  mailbox []string
}


func PrintBoard(bitboard *Bitboard) {
  pieceFields := []*uint64{
    &bitboard.whitePawns,
    &bitboard.whiteKnights,
    &bitboard.whiteBishops,
    &bitboard.whiteRooks,
    &bitboard.whiteQueens,
    &bitboard.whiteKing,
    &bitboard.blackPawns,
    &bitboard.blackKnights,
    &bitboard.blackBishops,
    &bitboard.blackRooks,
    &bitboard.blackQueens,
    &bitboard.blackKing,
  }

  pieceChars := []rune{'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k'}
  for i := 0; i < 8; i++ {
    for j := 0; j < 8; j++ {
      bit := uint64(0)
      for _, piece := range pieceFields {
        bit |= (*piece >> (i*8 + j)) & 1
      }
      switch bit {
      case 0:
        fmt.Print(" . ")
      default:
        for idx, piece := range pieceFields {
          if (*piece >> (i*8 + j)) & 1 != 0 {
            fmt.Print(" ")
            fmt.Printf("%c", pieceChars[idx])
            fmt.Print(" ")
            break
          }
        }
      }
    }
    fmt.Println()
  }
  fmt.Println("White turn:", bitboard.whiteTurn)
}

func GetBoardState(bitboard *Bitboard) [8][8]string {
  var creatingBoard [8][8]string
  pieceFields := []*uint64{
    &bitboard.whitePawns,
    &bitboard.whiteKnights,
    &bitboard.whiteBishops,
    &bitboard.whiteRooks,
    &bitboard.whiteQueens,
    &bitboard.whiteKing,
    &bitboard.blackPawns,
    &bitboard.blackKnights,
    &bitboard.blackBishops,
    &bitboard.blackRooks,
    &bitboard.blackQueens,
    &bitboard.blackKing,
  }

  pieceChars := []string{"wp", "wn", "wb", "wr", "wq", "wk", "bp", "bn", "bb", "br", "bq", "bk"}
  for i := 0; i < 8; i++ {
    for j := 0; j < 8; j++ {
      bit := uint64(0)
      for _, piece := range pieceFields {
        bit |= (*piece >> (i*8 + j)) & 1
      }
      switch bit {
      case 0:
        creatingBoard[i][j] = ""
      default:
        for idx, piece := range pieceFields {
          if (*piece >> (i*8 + j)) & 1 != 0 {
            creatingBoard[i][j] = pieceChars[idx]
            break
          }
        }
      }
    }
  }

  // invert the board so that white is on the bottom
  for i := 0; i < 4; i++ {
    for j := 0; j < 8; j++ {
      creatingBoard[i][j], creatingBoard[7-i][j] = creatingBoard[7-i][j], creatingBoard[i][j]
    }
  }


  return creatingBoard
}

func InitBoard(bitboard *Bitboard) {
  Constants := BitboardConstants {
    A_File: 0x0101010101010101,
    B_File: 0x0202020202020202,
    G_File: 0x4040404040404040,
    H_File: 0x8080808080808080,
    AB_File: 0x0303030303030303,
    GH_File: 0xC0C0C0C0C0C0C0C0,

    RANK_1: 0xFF00000000000000,
    RANK_2: 0x00FF000000000000,
    RANK_3: 0x0000FF0000000000,
    RANK_4: 0x000000FF00000000,
    RANK_5: 0x00000000FF000000,
    RANK_6: 0x0000000000FF0000,
    RANK_7: 0x000000000000FF00,
    RANK_8: 0x00000000000000FF,

    OnBoard: 0xFFFFFFFFFFFFFFFF,
  }
  bitboard.whitePawns = uint64(0xFF) << 48
  bitboard.blackPawns = uint64(0xFF) << 8
  bitboard.whiteKnights = uint64(0x42) << 56
  bitboard.blackKnights = uint64(0x42)
  bitboard.whiteBishops = uint64(0x24) << 56
  bitboard.blackBishops = uint64(0x24)
  bitboard.whiteRooks = uint64(0x81) << 56
  bitboard.blackRooks = uint64(0x81)
  bitboard.whiteQueens = uint64(0x08) << 56
  bitboard.blackQueens = uint64(0x08)
  bitboard.whiteKing = uint64(0x10) << 56
  bitboard.blackKing = uint64(0x10)
  bitboard.enPassant = 0
  bitboard.castlingRights = 0xC3 // 11000011
  bitboard.whiteTurn = true;
  bitboard.whiteOnBottom = true;
  bitboard.Constants = Constants
  // let this be an array of 64 strings where each string is the piece at that location on the board and let it be the starting position of the board
  bitboard.mailbox = make([]string, 64)
  for i := 0; i < 64; i++ {
    bitboard.mailbox[i] = ".."
  }
  for i := 0; i < 8; i++ {
    bitboard.mailbox[i + 8] = "bp"
    bitboard.mailbox[i + 48] = "wp"
  }
  bitboard.mailbox[0] = "br"
  bitboard.mailbox[1] = "bn"
  bitboard.mailbox[2] = "bb"
  bitboard.mailbox[3] = "bq"
  bitboard.mailbox[4] = "bk"
  bitboard.mailbox[5] = "bb"
  bitboard.mailbox[6] = "bn"
  bitboard.mailbox[7] = "br"
  bitboard.mailbox[56] = "wr"
  bitboard.mailbox[57] = "wn"
  bitboard.mailbox[58] = "wb"
  bitboard.mailbox[59] = "wq"
  bitboard.mailbox[60] = "wk"
  bitboard.mailbox[61] = "wb"
  bitboard.mailbox[62] = "wn"
  bitboard.mailbox[63] = "wr"
}

func PrintPieceLocations(bitboard *Bitboard) {
  for i := 0; i < 8; i++ {
    for j := 0; j < 8; j++ {
      fmt.Printf("%s ", bitboard.mailbox[i*8 + j])
    }
    fmt.Println()
  }
}

func GenerateBoardFromFen(fen string, bitboard *Bitboard) {

}

var PieceMoveFuncs = map[string]func(*Bitboard, uint64, uint64) {
	"wp": func(b *Bitboard, from, to uint64) { b.whitePawns ^= from; b.whitePawns |= to },
	"bp": func(b *Bitboard, from, to uint64) { b.blackPawns ^= from; b.blackPawns |= to },

	"wn": func(b *Bitboard, from, to uint64) { b.whiteKnights ^= from; b.whiteKnights |= to },
	"bn": func(b *Bitboard, from, to uint64) { b.blackKnights ^= from; b.blackKnights |= to },

  "wb": func(b *Bitboard, from, to uint64) { b.whiteBishops ^= from; b.whiteBishops |= to },
  "bb": func(b *Bitboard, from, to uint64) { b.blackBishops ^= from; b.blackBishops |= to },

	"wr": func(b *Bitboard, from, to uint64) { b.whiteRooks ^= from; b.whiteRooks |= to },
	"br": func(b *Bitboard, from, to uint64) { b.blackRooks ^= from; b.blackRooks |= to },

	"wq": func(b *Bitboard, from, to uint64) { b.whiteQueens ^= from; b.whiteQueens |= to },
	"bq": func(b *Bitboard, from, to uint64) { b.blackQueens ^= from; b.blackQueens |= to },

	"wk": func(b *Bitboard, from, to uint64) { b.whiteKing ^= from; b.whiteKing |= to },
	"bk": func(b *Bitboard, from, to uint64) { b.blackKing ^= from; b.blackKing |= to },
}

var PieceCaptureFuncs = map[string]func(*Bitboard, uint64) {
  "wp": func(b *Bitboard, to uint64) { b.blackPawns ^= to },
  "bp": func(b *Bitboard, to uint64) { b.whitePawns ^= to },

  "wn": func(b *Bitboard, to uint64) { b.blackKnights ^= to },
  "bn": func(b *Bitboard, to uint64) { b.whiteKnights ^= to },

  "wb": func(b *Bitboard, to uint64) { b.blackBishops ^= to },
  "bb": func(b *Bitboard, to uint64) { b.whiteBishops ^= to },

  "wr": func(b *Bitboard, to uint64) { b.blackRooks ^= to },
  "br": func(b *Bitboard, to uint64) { b.whiteRooks ^= to },

  "wq": func(b *Bitboard, to uint64) { b.blackQueens ^= to },
  "bq": func(b *Bitboard, to uint64) { b.whiteQueens ^= to },

  "wk": func(b *Bitboard, to uint64) { b.blackKing ^= to },
  "bk": func(b *Bitboard, to uint64) { b.whiteKing ^= to },
}

func MakeMove(piece string, from uint64, to uint64, bitboard *Bitboard) {
	if moveFunc, ok := PieceMoveFuncs[piece]; ok {

    // =================================== removing enemy pawn from enpassant ===================================
    if (piece == "wp" || piece == "bp") && (to & bitboard.enPassant) != 0 {
      if (to & bitboard.Constants.RANK_6) != 0 {
        if piece == "wp" {
          bitboard.blackPawns &= ^(to << 8)
          bitboard.mailbox[int(math.Log2(float64(to << 8)))] = ".."
        } else {
          bitboard.whitePawns &= ^(to << 8)
          bitboard.mailbox[int(math.Log2(float64(to << 8)))] = ".."
        }
      } else if (to & bitboard.Constants.RANK_3) != 0 {
        if piece == "wp" {
          bitboard.blackPawns &= ^(to >> 8)
          bitboard.mailbox[int(math.Log2(float64(to >> 8)))] = ".."
        } else {
          bitboard.whitePawns &= ^(to >> 8)
          bitboard.mailbox[int(math.Log2(float64(to >> 8)))] = ".."
        }
      }
    }

    // =================================== updates enpassant rights ===================================
    bitboard.enPassant = 0
    if piece == "wp" && (from >> 16) == to {
      bitboard.enPassant = from >> 8
    } else if piece == "bp" && (from << 16) == to {
      bitboard.enPassant = from << 8
    } 

    //fmt.Println("En passant square:")
    //DisplayPieceLocation(bitboard.enPassant)

    // ===================================== updating castling rights =====================================
    bottomLeft := uint64(1) << 56
    bottomRight := uint64(1) << 63
    topLeft := uint64(1)
    topRight := uint64(1) << 7

    if piece == "wk" {
      bitboard.castlingRights &= 0xC0 // 11000000

    } else if piece == "bk" {
      bitboard.castlingRights &= 0x3 // 00000011

    } else if piece == "wr" && bitboard.whiteOnBottom && from == bottomRight {
      bitboard.castlingRights &= 0xC2 // 11000010

    } else if piece == "wr" && bitboard.whiteOnBottom && from == bottomLeft {
      bitboard.castlingRights &= 0xC1 // 11000001

    } else if piece == "wr" && !bitboard.whiteOnBottom && from == topLeft { 
      bitboard.castlingRights &= 0xC2 // 11000010

    } else if piece == "wr" && !bitboard.whiteOnBottom && from == topRight { 
      bitboard.castlingRights &= 0xC1 // 11000001

    } else if piece == "br" && bitboard.whiteOnBottom && from == topLeft {
      bitboard.castlingRights &= 0x83 // 10000011

    } else if piece == "br" && bitboard.whiteOnBottom && from == topRight {
      bitboard.castlingRights &= 0x43 // 01000011

    } else if piece == "br" && !bitboard.whiteOnBottom && from == bottomLeft {
      bitboard.castlingRights &= 0x83 // 10000011

    } else if piece == "br" && !bitboard.whiteOnBottom && from == bottomRight {
      bitboard.castlingRights &= 0x43 // 01000011
    }

    // ===================================== Moving the rook after castling =====================================
    if piece == "wk" && (from << 2) == to {
      if moveRook, ok := PieceMoveFuncs["wr"]; ok {
        moveRook(bitboard, bottomRight, bottomRight >> 2)
        bitboard.mailbox[int(math.Log2(float64(bottomRight)))] = ".."
        bitboard.mailbox[int(math.Log2(float64(bottomRight >> 2)))] = "wr"
      }
    } else if piece == "wk" && (from >> 2) == to {
      if moveRook, ok := PieceMoveFuncs["wr"]; ok {
        moveRook(bitboard, bottomLeft, bottomLeft << 3)
        bitboard.mailbox[int(math.Log2(float64(bottomLeft)))] = ".."
        bitboard.mailbox[int(math.Log2(float64(bottomLeft << 3)))] = "wr"
      }
    } else if piece == "bk" && (from << 2) == to {
      if moveRook, ok := PieceMoveFuncs["br"]; ok {
        moveRook(bitboard, topRight, topRight >> 2)
        bitboard.mailbox[int(math.Log2(float64(topRight)))] = ".."
        bitboard.mailbox[int(math.Log2(float64(topRight >> 2)))] = "br"
      }
    } else if piece == "bk" && (from >> 2) == to {
      if moveRook, ok := PieceMoveFuncs["br"]; ok {
        moveRook(bitboard, topLeft, topLeft << 3)
        bitboard.mailbox[int(math.Log2(float64(topLeft)))] = ".."
        bitboard.mailbox[int(math.Log2(float64(topLeft << 3)))] = "br"
      }
    }

    bitboard.mailbox[int(math.Log2(float64(from)))] = ".."
    bitboard.mailbox[int(math.Log2(float64(to)))] = piece

    // get rid of this soon
    if bitboard.whiteTurn {
      bitboard.blackBishops &= ^to
      bitboard.blackKnights &= ^to
      bitboard.blackPawns &= ^to
      bitboard.blackQueens &= ^to
      bitboard.blackRooks &= ^to
      bitboard.blackKing &= ^to
    } else {
      bitboard.whiteBishops &= ^to 
      bitboard.whiteKnights &= ^to
      bitboard.whitePawns &= ^to
      bitboard.whiteQueens &= ^to
      bitboard.whiteRooks &= ^to
      bitboard.whiteKing &= ^to
    }


    PrintPieceLocations(bitboard)

		moveFunc(bitboard, from, to)
    bitboard.whiteTurn = !bitboard.whiteTurn
	} else {
		fmt.Printf("Unknown piece: %s\n", piece)
	}
}


func DisplayPieceLocation(piece uint64) {
  for i := 0; i < 8; i++ {
    for j := 0; j < 8; j++ {
      bit := (piece >> (i*8 + j)) & 1
      switch bit {
      case 0:
        fmt.Print(" . ")
      case 1:
        fmt.Print(" X ")
      }
    }
    fmt.Println()
  }
  fmt.Println()
}

func GetValidMoves(typeOfPiece string, piece uint64, bitboard *Bitboard) []uint64 {
  if typeOfPiece == "wn" || typeOfPiece == "bn" {
    return GetKnightMoves(piece, bitboard, typeOfPiece == "wn")

  } else if typeOfPiece == "wp" || typeOfPiece == "bp" {
    return GetPawnMoves(piece, bitboard, (typeOfPiece == "wp") == bitboard.whiteOnBottom, typeOfPiece == "wp")

  } else if typeOfPiece == "wr" || typeOfPiece == "br" {
    return GetRookMoves(piece, bitboard, typeOfPiece == "wr")

  } else if typeOfPiece == "wb" || typeOfPiece == "bb" {
    return GetBishopMoves(piece, bitboard, typeOfPiece == "wb")

  } else if typeOfPiece == "wq" || typeOfPiece == "bq" {
    return GetQueenMoves(piece, bitboard, typeOfPiece == "wq")

  } else if typeOfPiece == "wk" || typeOfPiece == "bk" {
    return GetKingMoves(piece, bitboard, typeOfPiece == "wk", (typeOfPiece == "wk") == bitboard.whiteOnBottom)

  }

  return []uint64{}
}


