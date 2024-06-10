package utils 

import (
  "fmt"
  "math"
)

// w b k q r b n p
// 0 0 0 0 0 0 0 0

const (
  // THis is for the mailbox
  BLACK_PAWN      uint8 = 0x41 
  BLACK_KNIGHT    uint8 = 0x42
  BLACK_BISHOP    uint8 = 0x44
  BLACK_ROOK      uint8 = 0x48
  BLACK_QUEEN     uint8 = 0x50
  BLACK_KING      uint8 = 0x60

  WHITE_PAWN      uint8 = 0x81
  WHITE_KNIGHT    uint8 = 0x82
  WHITE_BISHOP    uint8 = 0x84
  WHITE_ROOK      uint8 = 0x88
  WHITE_QUEEN     uint8 = 0x90
  WHITE_KING      uint8 = 0xA0
  // this is for general use
  A_File uint64 = 0x0101010101010101
  B_File uint64 = 0x0202020202020202
  G_File uint64 = 0x4040404040404040
  H_File uint64 = 0x8080808080808080
  AB_File uint64 = 0x0303030303030303
  GH_File uint64 = 0xC0C0C0C0C0C0C0C0

  RANK_1 uint64 = 0xFF00000000000000
  RANK_2 uint64 = 0x00FF000000000000
  RANK_3 uint64 = 0x0000FF0000000000
  RANK_4 uint64 = 0x000000FF00000000
  RANK_5 uint64 = 0x00000000FF000000
  RANK_6 uint64 = 0x0000000000FF0000
  RANK_7 uint64 = 0x000000000000FF00
  RANK_8 uint64 = 0x00000000000000FF

  OnBoard uint64 = 0xFFFFFFFFFFFFFFFF
  MOVED uint8 = 0x80
)

type Bitboard struct {
  whitePawns, whiteKnights, whiteBishops, whiteRooks, whiteQueens, whiteKing uint64
  blackPawns, blackKnights, blackBishops, blackRooks, blackQueens, blackKing uint64
  castlingRights uint8
  whiteOnBottom bool
  enPassant uint64
  whiteTurn bool
  mailbox []uint8
}


var PieceMoveFuncs = map[uint8]func(*Bitboard, uint64, uint64) {
  WHITE_PAWN: func(b *Bitboard, from, to uint64) { b.whitePawns ^= from; b.whitePawns |= to },
  BLACK_PAWN: func(b *Bitboard, from, to uint64) { b.blackPawns ^= from; b.blackPawns |= to },

  WHITE_KNIGHT: func(b *Bitboard, from, to uint64) { b.whiteKnights ^= from; b.whiteKnights |= to },
  BLACK_KNIGHT: func(b *Bitboard, from, to uint64) { b.blackKnights ^= from; b.blackKnights |= to },

  WHITE_BISHOP: func(b *Bitboard, from, to uint64) { b.whiteBishops ^= from; b.whiteBishops |= to },
  BLACK_BISHOP: func(b *Bitboard, from, to uint64) { b.blackBishops ^= from; b.blackBishops |= to },

  WHITE_ROOK: func(b *Bitboard, from, to uint64) { b.whiteRooks ^= from; b.whiteRooks |= to },
  BLACK_ROOK: func(b *Bitboard, from, to uint64) { b.blackRooks ^= from; b.blackRooks |= to },

  WHITE_QUEEN: func(b *Bitboard, from, to uint64) { b.whiteQueens ^= from; b.whiteQueens |= to },
  BLACK_QUEEN: func(b *Bitboard, from, to uint64) { b.blackQueens ^= from; b.blackQueens |= to },

  WHITE_KING: func(b *Bitboard, from, to uint64) { b.whiteKing ^= from; b.whiteKing |= to },
  BLACK_KING: func(b *Bitboard, from, to uint64) { b.blackKing ^= from; b.blackKing |= to },
}

var PieceCaptureFuncs = map[uint8]func(*Bitboard, uint64) {
  0: func(b *Bitboard, to uint64) {},
  BLACK_PAWN: func(b *Bitboard, to uint64) { b.blackPawns ^= to },
  WHITE_PAWN: func(b *Bitboard, to uint64) { b.whitePawns ^= to },

  BLACK_KNIGHT: func(b *Bitboard, to uint64) { b.blackKnights ^= to },
  WHITE_KNIGHT: func(b *Bitboard, to uint64) { b.whiteKnights ^= to },

  BLACK_BISHOP: func(b *Bitboard, to uint64) { b.blackBishops ^= to },
  WHITE_BISHOP: func(b *Bitboard, to uint64) { b.whiteBishops ^= to },

  BLACK_ROOK: func(b *Bitboard, to uint64) { b.blackRooks ^= to },
  WHITE_ROOK: func(b *Bitboard, to uint64) { b.whiteRooks ^= to },

  BLACK_QUEEN: func(b *Bitboard, to uint64) { b.blackQueens ^= to },
  WHITE_QUEEN: func(b *Bitboard, to uint64) { b.whiteQueens ^= to },

  BLACK_KING: func(b *Bitboard, to uint64) { b.blackKing ^= to },
  WHITE_KING: func(b *Bitboard, to uint64) { b.whiteKing ^= to },
}

func GetMailBoxIndex(square uint64) int {
  return int(math.Log2(float64(square)))
}

func MakeMove(piece uint8, from uint64, to uint64, bitboard *Bitboard) {
  // =================================== removing enemy pawn from enpassant ===================================
  if (piece & 0x1 > 0) && (to & bitboard.enPassant) != 0 {
    if (to & RANK_6) != 0 {
      enemyPawn := to << 8
      if piece == WHITE_PAWN {
        bitboard.blackPawns &= ^enemyPawn
        bitboard.mailbox[GetMailBoxIndex(enemyPawn)] = 0
      } else {
        bitboard.whitePawns &= ^enemyPawn
        bitboard.mailbox[GetMailBoxIndex(enemyPawn)] = 0
      }
    } else if (to & RANK_3) != 0 {
      enemyPawn := to >> 8
      if piece == WHITE_PAWN {
        bitboard.blackPawns &= ^enemyPawn
        bitboard.mailbox[GetMailBoxIndex(enemyPawn)] = 0
      } else {
        bitboard.whitePawns &= ^enemyPawn
        bitboard.mailbox[GetMailBoxIndex(enemyPawn)] = 0
      }
    }
  }

  // =================================== updates enpassant rights ===================================
  bitboard.enPassant = 0
  if piece == WHITE_PAWN && (from >> 16) == to {
    bitboard.enPassant = from >> 8
  } else if piece == BLACK_PAWN && (from << 16) == to {
    bitboard.enPassant = from << 8
  } 

  // ===================================== updating castling rights =====================================
  bottomLeft := uint64(1) << 56
  bottomRight := uint64(1) << 63
  topLeft := uint64(1)
  topRight := uint64(1) << 7

  if piece == WHITE_KING {
    bitboard.castlingRights &= 0xC0 // 11000000

  } else if piece == BLACK_KING {
    bitboard.castlingRights &= 0x3 // 00000011

  } else if piece == WHITE_ROOK && bitboard.whiteOnBottom && from == bottomRight {
    bitboard.castlingRights &= 0xC2 // 11000010

  } else if piece == WHITE_ROOK && bitboard.whiteOnBottom && from == bottomLeft {
    bitboard.castlingRights &= 0xC1 // 11000001

  } else if piece == WHITE_ROOK && !bitboard.whiteOnBottom && from == topLeft { 
    bitboard.castlingRights &= 0xC2 // 11000010

  } else if piece == WHITE_ROOK && !bitboard.whiteOnBottom && from == topRight { 
    bitboard.castlingRights &= 0xC1 // 11000001

  } else if piece == BLACK_ROOK && bitboard.whiteOnBottom && from == topLeft {
    bitboard.castlingRights &= 0x83 // 10000011

  } else if piece == BLACK_ROOK && bitboard.whiteOnBottom && from == topRight {
    bitboard.castlingRights &= 0x43 // 01000011

  } else if piece == BLACK_ROOK && !bitboard.whiteOnBottom && from == bottomLeft {
    bitboard.castlingRights &= 0x83 // 10000011

  } else if piece == BLACK_ROOK && !bitboard.whiteOnBottom && from == bottomRight {
    bitboard.castlingRights &= 0x43 // 01000011
  }

  // ===================================== Moving the rook after castling =====================================
  if piece == WHITE_KING && (from << 2) == to {
    if moveRook, ok := PieceMoveFuncs[WHITE_ROOK]; ok {
      moveRook(bitboard, bottomRight, bottomRight >> 2)
      bitboard.mailbox[int(math.Log2(float64(bottomRight)))] = 0
      bitboard.mailbox[int(math.Log2(float64(bottomRight >> 2)))] = WHITE_ROOK
    }
  } else if piece == WHITE_KING && (from >> 2) == to {
    if moveRook, ok := PieceMoveFuncs[WHITE_ROOK]; ok {
      moveRook(bitboard, bottomLeft, bottomLeft << 3)
      bitboard.mailbox[int(math.Log2(float64(bottomLeft)))] = 0
      bitboard.mailbox[int(math.Log2(float64(bottomLeft << 3)))] = WHITE_ROOK
    }
  } else if piece == BLACK_KING && (from << 2) == to {
    if moveRook, ok := PieceMoveFuncs[BLACK_ROOK]; ok {
      moveRook(bitboard, topRight, topRight >> 2)
      bitboard.mailbox[int(math.Log2(float64(topRight)))] = 0
      bitboard.mailbox[int(math.Log2(float64(topRight >> 2)))] = BLACK_ROOK
    }
  } else if piece == BLACK_KING && (from >> 2) == to {
    if moveRook, ok := PieceMoveFuncs[BLACK_ROOK]; ok {
      moveRook(bitboard, topLeft, topLeft << 3)
      bitboard.mailbox[int(math.Log2(float64(topLeft)))] = 0
      bitboard.mailbox[int(math.Log2(float64(topLeft << 3)))] = BLACK_ROOK
    }
  }

  // Moving the pieces
  toLocation := GetMailBoxIndex(to)
  if capturePiece, ok := PieceCaptureFuncs[bitboard.mailbox[toLocation]]; ok {
    capturePiece(bitboard, to)
  }

  bitboard.mailbox[GetMailBoxIndex(from)] = 0
  bitboard.mailbox[toLocation] = piece

  if movePiece, ok := PieceMoveFuncs[piece]; ok {
    movePiece(bitboard, from, to)
    // =================================== Pawn Promotion ===================================
    if piece == WHITE_PAWN && (to & RANK_8) != 0 {
      bitboard.whitePawns ^= to
      bitboard.whiteQueens |= to
      bitboard.mailbox[toLocation] = WHITE_QUEEN
    } else if piece == BLACK_PAWN && (to & RANK_1) != 0 {
      bitboard.blackPawns ^= to
      bitboard.blackQueens |= to
      bitboard.mailbox[toLocation] = BLACK_QUEEN
    }
    // ====================================================================================== 
  }

  bitboard.whiteTurn = !bitboard.whiteTurn
  PrintGame(bitboard)
}

func GetValidMoves(typeOfPiece uint8, piece uint64, bitboard *Bitboard) []uint64 {
  if typeOfPiece & 0x2 > 0 {
    return GetKnightMoves(piece, bitboard, typeOfPiece == WHITE_KNIGHT)

  } else if typeOfPiece & 0x1 > 0 {
    return GetPawnMoves(piece, bitboard, typeOfPiece == WHITE_PAWN)

  } else if typeOfPiece & 0x8 > 0 {
    return GetRookMoves(piece, bitboard, typeOfPiece == WHITE_ROOK)

  } else if typeOfPiece & 0x4 > 0 {
    return GetBishopMoves(piece, bitboard, typeOfPiece == WHITE_BISHOP)

  } else if typeOfPiece & 0x10 > 0 {
    return GetQueenMoves(piece, bitboard, typeOfPiece == WHITE_QUEEN)

  } else if typeOfPiece & 0x20 > 0 {
    return GetKingMoves(piece, bitboard, typeOfPiece == WHITE_KING)

  }

  return []uint64{}
}
// ================================================== NO NEED TO TOUCH ==================================================
// =================================== GETTING BOARD STATE FOR FRONTEND ===================================
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

// =================================== INIT THE BOARD ===================================
func InitBoard(bitboard *Bitboard) {
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

  bitboard.mailbox = make([]uint8, 64)
  for i := 0; i < 8; i++ {
    bitboard.mailbox[i + 8] = BLACK_PAWN
    bitboard.mailbox[i + 48] = WHITE_PAWN
  }
  bitboard.mailbox[0] = BLACK_ROOK
  bitboard.mailbox[1] = BLACK_KNIGHT
  bitboard.mailbox[2] = BLACK_BISHOP
  bitboard.mailbox[3] = BLACK_QUEEN
  bitboard.mailbox[4] = BLACK_KING
  bitboard.mailbox[5] = BLACK_BISHOP
  bitboard.mailbox[6] = BLACK_KNIGHT
  bitboard.mailbox[7] = BLACK_ROOK

  bitboard.mailbox[56] = WHITE_ROOK
  bitboard.mailbox[57] = WHITE_KNIGHT
  bitboard.mailbox[58] = WHITE_BISHOP
  bitboard.mailbox[59] = WHITE_QUEEN
  bitboard.mailbox[60] = WHITE_KING
  bitboard.mailbox[61] = WHITE_BISHOP
  bitboard.mailbox[62] = WHITE_KNIGHT
  bitboard.mailbox[63] = WHITE_ROOK
}

// =================================== PRINTING STUFF ===================================
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

func PrintGame(bitboard *Bitboard) {
  var piecemap = map[uint8]string{
    WHITE_PAWN: "P",
    WHITE_KNIGHT: "N",
    WHITE_BISHOP: "B",
    WHITE_ROOK: "R",
    WHITE_QUEEN: "Q",
    WHITE_KING: "K",
    BLACK_PAWN: "p",
    BLACK_KNIGHT: "n",
    BLACK_BISHOP: "b",
    BLACK_ROOK: "r",
    BLACK_QUEEN: "q",
    BLACK_KING: "k",
    0: ".",
  }
  fmt.Println("==================MAILBOX==================")
  for i := 0; i < 8; i++ {
    for j := 0; j < 8; j++ {
      fmt.Printf("%s ", piecemap[bitboard.mailbox[i*8 + j]])
    }
    fmt.Println()
  }
  fmt.Println("=================BITBOARD==================")
  printBoard(bitboard)
  fmt.Println("===========================================")

  fmt.Println("Next Turn is White:", bitboard.whiteTurn)
  fmt.Printf("Castling Rights: %b", bitboard.castlingRights)
  fmt.Println("\n-------------------------------------------")

}

func printBoard(bitboard *Bitboard) {
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
        fmt.Print(". ")
      default:
        for idx, piece := range pieceFields {
          if (*piece >> (i*8 + j)) & 1 != 0 {
            fmt.Printf("%c ", pieceChars[idx])
            break
          }
        }
      }
    }
    fmt.Println()
  }
}

func ShowMailbox(bitboard *Bitboard) {
  for i := 0; i < 8; i++ {
    for j := 0; j < 8; j++ {
      fmt.Printf("%d ", i*8 + j)
    }
    fmt.Println()
  }
}


