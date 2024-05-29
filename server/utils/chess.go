package utils 

import "fmt"

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
    //                                              we on the bottom
    return GetPawnMoves(piece, bitboard, Constants, typeOfPiece == "wp", typeOfPiece == "wp")
  } else if typeOfPiece == "wr" || typeOfPiece == "br" {
    return GetRookMoves(piece, bitboard, Constants, typeOfPiece == "wr")
  }
  return []uint64{}
}

func GetKnightMoves(knightPosition uint64, bitboard *Bitboard, Constants *BitboardConstants, isWhite bool) []uint64 {
  var moves []uint64
  var allPieces uint64

  if isWhite {
    allPieces = bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing
  } else {
    allPieces = bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing
  }

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














func GetPawnMoves(pawnPosition uint64, bitboard *Bitboard, Constants *BitboardConstants, onBottom bool, isWhite bool) []uint64 {
  var moves []uint64
  var allPieces uint64
  shiftAmount := uint64(8)

  whitePieces := bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing
  blackPieces := bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing
  allPieces = whitePieces | blackPieces
  DisplayPieceLocation(whitePieces)

  if onBottom && ((pawnPosition >> shiftAmount) & allPieces) == 0 {
    moves = append(moves, pawnPosition >> shiftAmount)

  } else if !onBottom && ((pawnPosition << shiftAmount) & allPieces) == 0 {
    moves = append(moves, pawnPosition << shiftAmount)
  }

  // double move if pawn hasn't moved
  if onBottom && !pawnHasMoved(pawnPosition, onBottom) && len(moves) > 0 {
    moves = append(moves, pawnPosition >> 16)
  } else if !onBottom && !pawnHasMoved(pawnPosition, onBottom) && len(moves) > 0 {
    moves = append(moves, pawnPosition << 16)
  }

  return moves
}

func pawnHasMoved(pawnPosition uint64, onBottom bool) bool {
  var initialPawnPosition uint64
  if onBottom {
    initialPawnPosition = uint64(0xFF) << 48
  } else {
    initialPawnPosition = uint64(0xFF) << 8
  }
  result := initialPawnPosition & pawnPosition
  return result == 0
}









func GetRookMoves(rookPosition uint64, bitboard *Bitboard, Constants *BitboardConstants, isWhite bool) []uint64 {
  var moves []uint64
  //var allPieces uint64
  blackPieces := bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing
  whitePieces := bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing
  //allPieces = blackPieces | whitePieces

  // go up
  for i := 1; i < 8; i++ {
    shiftedPosition := rookPosition >> (8 * i)
    // find a way to check if the position is on the board
    if isWhite && ((shiftedPosition & blackPieces) != 0 || (shiftedPosition == 0)) {
      moves = append(moves, shiftedPosition)
      if shiftedPosition & blackPieces != 0 {
        break
      }
    } else if !isWhite && ((shiftedPosition & whitePieces) != 0 || (shiftedPosition == 0)) {
      moves = append(moves, shiftedPosition)
      if shiftedPosition & whitePieces != 0 {
        break
      }
    }
  }

  return moves
}


