package utils

func RandomAI_move(bitboard *Bitboard, whiteTurn bool) (uint64, uint64) {
  return 0, 0
}

func AI_move(bitboard *Bitboard, whiteTurn bool, depth int) (uint64, uint64) {
  from := uint64(0)
  to := uint64(0)

  if whiteTurn { // we want a score that is high position for white
    //var bestValue float32 = -500
    //allWhiteMoves := GenerateAllMoves(bitboard, true)
    //for _, move := range allWhiteMoves {
      //bitboard.MakeMove(move)
      //value := Minimax(bitboard, false, depth - 1)
      //if value > bestValue {
       // bestValue = value
       // from = move.from
       // to = move.to
      //}
      //bitboard.UndoMove(move)
    //}
  }

  return from, to
}

func Minimax(bitboard *Bitboard, whiteTurn bool, depth int) int32 {
  if depth == 0 {
    return Evaluate(bitboard)
  }

  if whiteTurn {

    var bestValue int32 = -500
    //allMoveForWhite := GenerateAllMoves(bitboard, true)
    //for _, move := range allMoveForWhite {
      //MakeMove(move)
      //bestValue = math.Max(bestValue, Minimax(bitboard, false, depth - 1))
      //bitboard.UndoMove(move)
    //}
    return bestValue

  } else {

    var bestValue int32 = 500

    return bestValue

  }
}

func Evaluate(bitboard *Bitboard) int32 { // positive is good for white, negative is good for black
  var score int32 = 0

  for i := 0; i < 64; i++ {
    if bitboard.mailbox[i] == 0 {
      continue
    }
    score += PIECE_TO_VALUE[bitboard.mailbox[i]] + PIECE_TO_TABLE[bitboard.mailbox[i]][i]
  }

  // calculate the pawn structure score (isolated, doubled, passed)
  // calculate the mobility score
  // calculate the rook on open file score

  return score
}
