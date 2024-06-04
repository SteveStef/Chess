package utils 

import "fmt"

type BitboardConstants struct {
  A_File uint64
  B_File uint64
  G_File uint64
  H_File uint64
  AB_File uint64
  GH_File uint64

  RANK_1 uint64
  RANK_8 uint64

  OnBoard uint64
}

type Bitboard struct {
  whitePawns, whiteKnights, whiteBishops, whiteRooks, whiteQueens, whiteKing uint64
  blackPawns, blackKnights, blackBishops, blackRooks, blackQueens, blackKing uint64
  enPassant uint64
  castlingRights uint8 // 0x1 white king side, 0x2 white queen side, 0x4 black king side, 0x8 black queen side
  whiteTurn bool
  whiteOnBottom bool
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

func MakeMove(piece string, from uint64, to uint64, bitboard *Bitboard) {
	if moveFunc, ok := PieceMoveFuncs[piece]; ok {
    // =================================== updates enpassant rights ===================================
    bitboard.enPassant = 0
    // if the piece is a pawn and it is moving two squares forward then set the enPassant square to the square behind the pawn
    if piece == "wp" && (from >> 16) == to {
      bitboard.enPassant = from >> 8
    } else if piece == "bp" && (from << 16) == to {
      bitboard.enPassant = from << 8
    } 
    //fmt.Println("En passant square:")
    //DisplayPieceLocation(bitboard.enPassant)

    // if the move is en passant then remove the captured pawn (test this)
    if piece == "wp" && (from << 1) == to {
      bitboard.blackPawns &= ^(to >> 8)
    } else if piece == "bp" && (from >> 1) == to {
      bitboard.whitePawns &= ^(to << 8)
    }

    // ===================================== updating castling rights =====================================
    if piece == "wk" {
      bitboard.castlingRights &= 0xC0 // 11000000

    } else if piece == "bk" {
      bitboard.castlingRights &= 0x3 // 00000011

    } else if piece == "wr" && from == uint64(1) << 63 { // rook is on the king side
      bitboard.castlingRights &= 0xC2 // 11000010

    } else if piece == "wr" && from == uint64(1) << 56 {
      bitboard.castlingRights &= 0xC1 // 11000001

    } else if piece == "br" && from == uint64(1) {
      bitboard.castlingRights &= 0x83 // 10000011

    } else if piece == "br" && from == uint64(1) << 7 {
      bitboard.castlingRights &= 0x43 // 01000011
    }

    fmt.Println("Castling rights:")
    fmt.Printf("%b\n", bitboard.castlingRights)

    // if the move is about castling then move the rook as well
    // if the king is being shifted to the left or right by 2 squares then the rook is being moved as well
    if piece == "wk" && (from << 2) == to {

    } else if piece == "wk" && (from >> 2) == to {

    } else if piece == "bk" && (from << 2) == to {

    } else if piece == "bk" && (from >> 2) == to {

    }

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

func GetValidMoves(typeOfPiece string, piece uint64, bitboard *Bitboard, Constants *BitboardConstants) []uint64 {
  if typeOfPiece == "wn" || typeOfPiece == "bn" {
    return GetKnightMoves(piece, bitboard, Constants, typeOfPiece == "wn")

  } else if typeOfPiece == "wp" || typeOfPiece == "bp" {
    return GetPawnMoves(piece, bitboard, Constants, (typeOfPiece == "wp") == bitboard.whiteOnBottom, typeOfPiece == "wp")

  } else if typeOfPiece == "wr" || typeOfPiece == "br" {
    return GetRookMoves(piece, bitboard, Constants, typeOfPiece == "wr")

  } else if typeOfPiece == "wb" || typeOfPiece == "bb" {
    return GetBishopMoves(piece, bitboard, Constants, typeOfPiece == "wb")

  } else if typeOfPiece == "wq" || typeOfPiece == "bq" {
    return GetQueenMoves(piece, bitboard, Constants, typeOfPiece == "wq")

  } else if typeOfPiece == "wk" || typeOfPiece == "bk" {
    return GetKingMoves(piece, bitboard, Constants, typeOfPiece == "wk", (typeOfPiece == "wk") == bitboard.whiteOnBottom)

  }

  return []uint64{}
}


