package utils

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
func GetPawnMoves(pawnPosition uint64, bitboard *Bitboard, onBottom bool, isWhite bool) []uint64 {
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
  if onBottom && ((pawnPosition >> shiftAmount) & allPieces) == 0 && pawnPosition & RANK_8 == 0 {
    moves = append(moves, pawnPosition >> shiftAmount)

  } else if !onBottom && ((pawnPosition << shiftAmount) & allPieces) == 0 && pawnPosition & RANK_1 == 0 {
    moves = append(moves, pawnPosition << shiftAmount)
  }

  // double move if pawn hasn't moved
  if onBottom && !pawnHasMoved(pawnPosition, onBottom) && len(moves) > 0 {
    moves = append(moves, pawnPosition >> 16)
  } else if !onBottom && !pawnHasMoved(pawnPosition, onBottom) && len(moves) > 0 {
    moves = append(moves, pawnPosition << 16)
  }

  // capture moves
  tmpMove := pawnPosition >> 7
  if onBottom && (tmpMove & oppositeColorPieces != 0 || tmpMove & bitboard.enPassant != 0) {
    moves = append(moves, tmpMove)
  }
  tmpMove = pawnPosition >> 9
  if onBottom && (tmpMove & oppositeColorPieces != 0 || tmpMove & bitboard.enPassant != 0) {
    moves = append(moves, tmpMove)
  }
  tmpMove = pawnPosition << 7
  if !onBottom && (tmpMove & oppositeColorPieces != 0 || tmpMove & bitboard.enPassant != 0) {
    moves = append(moves, tmpMove)
  }
  tmpMove = pawnPosition << 9
  if !onBottom && (tmpMove & oppositeColorPieces != 0 || tmpMove & bitboard.enPassant != 0) {
    moves = append(moves, tmpMove)
  }

  return moves
}

// helper function to check if pawn has moved
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

// IF your are not going to switch black and white from bottom to top then you can remove the onBottom parameter and extra logic
func GetKingMoves(kingPosition uint64, bitboard *Bitboard, isWhite bool, onBottom bool) []uint64 { // onbottom = white on bottom
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
    if onBottom {
      if (kingPosition << 1) & allPieces == 0 && (kingPosition << 2) & allPieces == 0 { moves = append(moves, kingPosition << 2) }
    } else {
      if (kingPosition >> 1) & allPieces == 0 && (kingPosition >> 2) & allPieces == 0 { moves = append(moves, kingPosition >> 2) }
    }
  }

  if isWhite && (bitboard.castlingRights & 2 != 0) { // white queen side 
    if onBottom {
      if (kingPosition >> 1) & allPieces == 0 && (kingPosition >> 2) & allPieces == 0 {
        moves = append(moves, kingPosition >> 2) 
      }
    } else {
      if (kingPosition << 1) & allPieces == 0 && (kingPosition << 2) & allPieces == 0 { moves = append(moves, kingPosition << 2) }
    }
  }

  if !isWhite && (bitboard.castlingRights & 0x80 != 0) { // black king side
    if !onBottom { // onBottom = white on bottom
      if (kingPosition << 1) & allPieces == 0 && (kingPosition << 2) & allPieces == 0 { moves = append(moves, kingPosition << 2) }
    } else {
      if (kingPosition >> 1) & allPieces == 0 && (kingPosition >> 2) & allPieces == 0 { moves = append(moves, kingPosition >> 2) }
    }
  }

  if !isWhite && (bitboard.castlingRights & 0x40 != 0) { // black queen side
    if !onBottom {
      if (kingPosition >> 1) & allPieces == 0 && (kingPosition >> 2) & allPieces == 0 { moves = append(moves, kingPosition >> 2) }
    } else {
      if (kingPosition << 1) & allPieces == 0 && (kingPosition << 2) & allPieces == 0 { moves = append(moves, kingPosition << 2) }
    }
  }

  return moves
}

