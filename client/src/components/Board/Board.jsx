import "./Board.css";
import Ranks from "./bits/Ranks";
import Files from "./bits/Files";
import Pieces from "./Pieces/Pieces";
import { useState } from "react";

const Board = () => {
  const [highlightBoard, setHighlightBoard] = useState(
    Array(8)
      .fill()
      .map(() => Array(8).fill(0)),
  );

  const updateHighlight = (row, col) => {
    const newHighlightBoard = [...highlightBoard];
    newHighlightBoard[row][col] = 1;
    setHighlightBoard(newHighlightBoard);
  };

  const nothigh = () => {
    const a = [...highlightBoard];
    for (let i = 0; i < a.length; i++) {
      for (let j = 0; j < a[i].length; j++) {
        a[i][j] = 0;
      }
    }
    setHighlightBoard(a);
  };

  const ranks = Array(8)
    .fill()
    .map((x, i) => 8 - i);

  const files = Array(8)
    .fill()
    .map((x, i) => i + 1);

  const getClassName = (i, j) => {
    let c = "tile";
    c += (i + j) % 2 === 0 ? " tile--dark " : " tile--light ";
    return c;
  };

  return (
    <div className="board">
      <Ranks ranks={ranks} />

      <div className="tiles">
        {ranks.map((rank, i) =>
          files.map((file, j) => (
            <div
              key={file + "" + rank}
              i={i}
              j={j}
              className={
                highlightBoard[i][j] === 1 ? "dot" : getClassName(9 - i, j)
              }
            ></div>
          )),
        )}
      </div>
      <Pieces updateHighlight={updateHighlight} nothigh={nothigh} />
      <Files files={files} />
    </div>
  );
};
export default Board;
