package utils 

import (
	"fmt"
)

type BitboardConstants struct {
   NotAFile uint64
   NotBFile uint64
   NotGFile uint64
   NotHFile uint64
   NotABFile uint64
   NotGHFile uint64
}

type Bitboard struct {
	whitePawns, whiteKnights, whiteBishops, whiteRooks, whiteQueens, whiteKing uint64
	blackPawns, blackKnights, blackBishops, blackRooks, blackQueens, blackKing uint64
  whiteTurn bool
  castlingRights uint64
  elPassantSquare uint64
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
  bitboard.whiteTurn = true;
}


func D(piece uint64) {
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

func GetValidMoves(knightPosition uint64, bitboard *Bitboard, Constants *BitboardConstants) []uint64 {
  var moves []uint64
  allPieces := bitboard.whitePawns | bitboard.blackPawns | bitboard.whiteKnights | bitboard.blackKnights | bitboard.whiteBishops | bitboard.blackBishops | bitboard.whiteRooks | bitboard.blackRooks | bitboard.whiteQueens | bitboard.blackQueens | bitboard.whiteKing | bitboard.blackKing
  shiftAmounts := []uint64{ 15, 17, 10, 6 }

  // Helper function for shifting and appending moves
  appendIfValid := func(shiftOperation func(uint64) uint64, shiftAmount uint64, fileMask uint64) {
    shiftedPosition := shiftOperation(knightPosition)
    if (shiftedPosition & allPieces) == 0 && (shiftedPosition & fileMask) == 0 && shiftedPosition != 0 {
      moves = append(moves, shiftedPosition)
    }
  }

  // down shifting
  for _, shiftAmount := range shiftAmounts {
    switch shiftAmount {
    case 17:
      appendIfValid(func(pos uint64) uint64 { return pos << shiftAmount }, shiftAmount, Constants.NotAFile)
    case 10:
      appendIfValid(func(pos uint64) uint64 { return pos << shiftAmount }, shiftAmount, Constants.NotABFile)
    case 15:
      appendIfValid(func(pos uint64) uint64 { return pos << shiftAmount }, shiftAmount, Constants.NotHFile)
    case 6:
      appendIfValid(func(pos uint64) uint64 { return pos << shiftAmount }, shiftAmount, Constants.NotGHFile)
    }
  }

  // upshifting
  for _, shiftAmount := range shiftAmounts {
    switch shiftAmount {
    case 17:
      appendIfValid(func(pos uint64) uint64 { return pos >> shiftAmount }, shiftAmount, Constants.NotHFile)
    case 10:
      appendIfValid(func(pos uint64) uint64 { return pos >> shiftAmount }, shiftAmount, Constants.NotGHFile)
    case 15:
      appendIfValid(func(pos uint64) uint64 { return pos >> shiftAmount }, shiftAmount, Constants.NotAFile)
    case 6:
      appendIfValid(func(pos uint64) uint64 { return pos >> shiftAmount }, shiftAmount, Constants.NotABFile)
    }
  }

  return moves
}
