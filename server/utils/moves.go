package utils
import "fmt"

// ===================================BISHOP=======================================================================
func GetKnightMoves(knightPosition uint64, bitboard *Bitboard, isWhite bool) []uint64 {
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
      appendIfValid(func(pos uint64) uint64 { return pos << shiftAmount }, shiftAmount, A_File)
    case 10:
      appendIfValid(func(pos uint64) uint64 { return pos << shiftAmount }, shiftAmount, AB_File)
    case 15:
      appendIfValid(func(pos uint64) uint64 { return pos << shiftAmount }, shiftAmount, H_File)
    case 6:
      appendIfValid(func(pos uint64) uint64 { return pos << shiftAmount }, shiftAmount, GH_File)
  }
  }

  // upshifting
  for _, shiftAmount := range shiftAmounts {
    switch shiftAmount {
    case 17:
      appendIfValid(func(pos uint64) uint64 { return pos >> shiftAmount }, shiftAmount, H_File)
    case 10:
      appendIfValid(func(pos uint64) uint64 { return pos >> shiftAmount }, shiftAmount, GH_File)
    case 15:
      appendIfValid(func(pos uint64) uint64 { return pos >> shiftAmount }, shiftAmount, A_File)
    case 6:
      appendIfValid(func(pos uint64) uint64 { return pos >> shiftAmount }, shiftAmount, AB_File)
  }
  }

  return moves
}

// =================================PAWN=============================================
func GetPawnMoves(pawnPosition uint64, bitboard *Bitboard, isWhite bool) []uint64 {
  var moves []uint64
  var sameColorPieces uint64
  var oppositeColorPieces uint64
  var allPieces uint64

  shiftAmount := uint64(8)

  if isWhite {
    sameColorPieces = bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing
    oppositeColorPieces = bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing
  } else {
    sameColorPieces = bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing
    oppositeColorPieces = bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing
  }

  allPieces = sameColorPieces | oppositeColorPieces

  if isWhite && ((pawnPosition >> shiftAmount) & allPieces) == 0 && pawnPosition & RANK_8 == 0 {
    moves = append(moves, pawnPosition >> shiftAmount)

  } else if !isWhite && ((pawnPosition << shiftAmount) & allPieces) == 0 && pawnPosition & RANK_1 == 0 {
    moves = append(moves, pawnPosition << shiftAmount)
  }

  // double move if pawn hasn't moved
  upshift2 := pawnPosition >> 16
  downshift2 := pawnPosition << 16


  if isWhite {
    hasMoved := pawnPosition & (uint64(0xFF) << 48) == 0
    if !hasMoved && len(moves) > 0 && upshift2 & allPieces == 0 {
      moves = append(moves, upshift2)
    }
  } else {
    hasMoved := pawnPosition & (uint64(0xFF) << 8) == 0
    if !hasMoved && len(moves) > 0 && downshift2 & allPieces == 0 {
      moves = append(moves, downshift2)
    }
  }

  // capture moves
  tmpMove := pawnPosition >> 7
  if isWhite && (tmpMove & oppositeColorPieces != 0 || tmpMove & bitboard.enPassant != 0) {
    moves = append(moves, tmpMove)
  }
  tmpMove = pawnPosition >> 9
  if isWhite && (tmpMove & oppositeColorPieces != 0 || tmpMove & bitboard.enPassant != 0) {
    moves = append(moves, tmpMove)
  }
  tmpMove = pawnPosition << 7
  if !isWhite && (tmpMove & oppositeColorPieces != 0 || tmpMove & bitboard.enPassant != 0) {
    moves = append(moves, tmpMove)
  }
  tmpMove = pawnPosition << 9
  if !isWhite && (tmpMove & oppositeColorPieces != 0 || tmpMove & bitboard.enPassant != 0) {
    moves = append(moves, tmpMove)
  }

  return moves
}

// =================================ROOK===========================================================================
// add go routine to check for moves in all directions
func GetRookMoves(rookPosition uint64, bitboard *Bitboard, isWhite bool) []uint64 {
  var moves []uint64
  var sameColorPieces uint64
  var oppositeColorPieces uint64

  if isWhite {
    sameColorPieces = bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing
    oppositeColorPieces = bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing
  } else {
    sameColorPieces = bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing
    oppositeColorPieces = bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing
  }

  shifted := rookPosition
  for (shifted & RANK_8) == 0 { // while we're not at the top
    shifted = shifted >> 8
    if (shifted & sameColorPieces) != 0 { break }
    moves = append(moves, shifted)
    if (shifted & oppositeColorPieces) != 0 { break }
  }

  shifted = rookPosition
  for (shifted & RANK_1 == 0) { // while we're not at the bottom
    shifted = shifted << 8
    if (shifted & sameColorPieces) != 0 { break }
    moves = append(moves, shifted)
    if (shifted & oppositeColorPieces) != 0 { break }
  }

  shifted = rookPosition
  for (shifted & H_File == 0) {
    shifted = shifted << 1
    if (shifted & sameColorPieces) != 0 { break }
    moves = append(moves, shifted)
    if (shifted & oppositeColorPieces) != 0 || (shifted & H_File) != 0 { break }
  }

  shifted = rookPosition
  for (shifted & A_File == 0) {
    shifted = shifted >> 1
    if (shifted & sameColorPieces) != 0 { break }
    moves = append(moves, shifted)
    if (shifted & oppositeColorPieces) != 0 || (shifted & A_File) != 0 { break }
  }

  return moves
}

func GetBishopMoves(bishopPosition uint64, bitboard *Bitboard, isWhite bool) []uint64 {
  var moves []uint64
  var sameColorPieces uint64
  var oppositeColorPieces uint64

  if isWhite {
    sameColorPieces = bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing
    oppositeColorPieces = bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing
  } else {
    sameColorPieces = bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing
    oppositeColorPieces = bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing
  }

  shifted := bishopPosition
  for (shifted & RANK_8) == 0 && (shifted & H_File) == 0 { // while we're not at the top or right
    shifted = shifted >> 7
    if (shifted & sameColorPieces) != 0 { break }
    moves = append(moves, shifted)
    if (shifted & oppositeColorPieces) != 0 { break }
  }

  shifted = bishopPosition
  for (shifted & RANK_8) == 0 && (shifted & A_File) == 0 { // while we're not at the top or left
    shifted = shifted >> 9
    if (shifted & sameColorPieces) != 0 { break }
    moves = append(moves, shifted)
    if (shifted & oppositeColorPieces) != 0 { break }
  }

  shifted = bishopPosition
  for (shifted & RANK_1 == 0) && (shifted & H_File == 0) { // while we're not at the bottom or right
    shifted = shifted << 9
    if (shifted & sameColorPieces) != 0 { break }
    moves = append(moves, shifted)
    if (shifted & oppositeColorPieces) != 0 { break }
  }
  
  shifted = bishopPosition
  for (shifted & RANK_1 == 0) && (shifted & A_File == 0) { // while we're not at the bottom or left
    shifted = shifted << 7
    if (shifted & sameColorPieces) != 0 { break }
    moves = append(moves, shifted)
    if (shifted & oppositeColorPieces) != 0 { break }
  }

  return moves
}

func GetQueenMoves(queenPosition uint64, bitboard *Bitboard, isWhite bool) []uint64 {
  crossMoves := GetRookMoves(queenPosition, bitboard, isWhite)
  diagonalMoves := GetBishopMoves(queenPosition, bitboard, isWhite)

  allMoves := append(crossMoves, diagonalMoves...)
  return allMoves
}

func GetKingMoves(kingPosition uint64, bitboard *Bitboard, isWhite bool) []uint64 { // onbottom = white on bottom
  var moves []uint64
  var sameColorPieces uint64
  var allPieces uint64

  if isWhite {
    sameColorPieces = bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing
  } else {
    sameColorPieces = bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing
  }

  allPieces = bitboard.whitePawns | bitboard.whiteKnights | bitboard.whiteBishops | bitboard.whiteRooks | bitboard.whiteQueens | bitboard.whiteKing | bitboard.blackPawns | bitboard.blackKnights | bitboard.blackBishops | bitboard.blackRooks | bitboard.blackQueens | bitboard.blackKing

  possibleMoves := []uint64{
    kingPosition << 8, kingPosition >> 8,
    kingPosition << 1, kingPosition >> 1,
    kingPosition << 7, kingPosition << 9,
    kingPosition >> 7, kingPosition >> 9,
  }

  for _, move := range possibleMoves {
    if (move & OnBoard) != 0 && (move & sameColorPieces) == 0 {
      moves = append(moves, move)
    }
  }

  if isWhite && (bitboard.castlingRights & 1 != 0) { // white king side
    attacked := squaresAreGettingAttacked(bitboard, isWhite, true)

    fmt.Print("Getting Attacked: ")
    fmt.Println(attacked)

    if !attacked && (kingPosition << 1) & allPieces == 0 && (kingPosition << 2) & allPieces == 0 { moves = append(moves, kingPosition << 2) }
  }

  if isWhite && (bitboard.castlingRights & 2 != 0) { // white queen side 
    if (kingPosition >> 1) & allPieces == 0 && (kingPosition >> 2) & allPieces == 0 {
      moves = append(moves, kingPosition >> 2) 
    }
  }

  if !isWhite && (bitboard.castlingRights & 0x80 != 0) { // black king side
    if (kingPosition << 1) & allPieces == 0 && (kingPosition << 2) & allPieces == 0 { moves = append(moves, kingPosition << 2) }
  }

  if !isWhite && (bitboard.castlingRights & 0x40 != 0) { // black queen side
    if (kingPosition >> 1) & allPieces == 0 && (kingPosition >> 2) & allPieces == 0 { moves = append(moves, kingPosition >> 2) }
  }

  return moves
}

func squaresAreGettingAttacked(bitboard *Bitboard, isWhite bool, isKingside bool) bool {
  if isWhite {
    // ==================== Checking the knights ============================
    kingSideSquareAttackers := uint64(62197173760032768)
    queenSideSquareAttackers := uint64(35816591274803200)
    if isKingside && kingSideSquareAttackers & bitboard.blackKnights > 0 {
      return  true
    } else if !isKingside && queenSideSquareAttackers & bitboard.blackKnights > 0 {
      return true
    }

    fmt.Println("Knights are not attacking")

    // ========================== Checking the rooks/queens ====================
    checkRank := func(start int) bool {
      for start > 0 {
        if bitboard.mailbox[start] == BLACK_ROOK || bitboard.mailbox[start] == BLACK_QUEEN {
          return true
        } else if bitboard.mailbox[start] != 0 {
          return false
        }
        start -= 8
      }
      return false
    }

    if isKingside {
      if checkRank(62-8) || checkRank(61-8) || checkRank(60-8) { return true }
    } else {
      if checkRank(57-8) || checkRank(58-8) || checkRank(59-8) || checkRank(60-8) { return true }
    }

    fmt.Println("no rooks or queens attacking")

    // ========================== Checking the bishops/queens ====================
    checkDiagonal := func(start int, increment int) bool {
      for start > 0 {
        if bitboard.mailbox[start] == BLACK_BISHOP || bitboard.mailbox[start] == BLACK_QUEEN {
          return true
        } else if bitboard.mailbox[start] != 0 {
          return false
        }
        start += increment
      }
      return false
    }

    if isKingside {
      if checkDiagonal(61, -7) || checkDiagonal(61, -9) || checkDiagonal(62, -7) || checkDiagonal(62, -9) { return true }
    } else {
      if checkDiagonal(57, -7) || checkDiagonal(57, -9) || checkDiagonal(58, -7) || checkDiagonal(58, -9) || checkDiagonal(59, -7) || checkDiagonal(59, -9) { return true }
    }

    fmt.Println("no bishops or queens attacking")

    // ========================= pawn ========================================
    if isKingside {
      pawnAttackers := uint64(69805794224242688)
      if pawnAttackers & bitboard.blackPawns > 0 { return true }
    } else {
      pawnAttackers := uint64(17732923532771328)
      if pawnAttackers & bitboard.blackPawns > 0 { return true }
    }
    fmt.Println("no pawns attacking")

    // ========================= king ========================================
    if isKingside {
      kingAttackers := uint64(54043195528445952)
      if kingAttackers & bitboard.blackKing > 0 { return true }
    } else {
      kingAttackers := uint64(1970324836974592)
      if kingAttackers & bitboard.blackKing > 0 { return true }
    }
    fmt.Println("no king attacking")

  } else { // this is for black

  }

  return false
}
