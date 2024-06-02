import "./Pieces.css";
import Piece from "./Piece";
import { createPosition, copyPosition } from "../../../helper";
import { useState, useRef } from "react";

const Pieces = ({ updateHighlight, nothigh }) => {
  const [state, setState] = useState(createPosition());
  const ref = useRef();

  const calculateCoords = (e) => {
    const { width, left, top } = ref.current.getBoundingClientRect();
    const size = width / 8;
    const y = Math.floor((e.clientX - left) / size);
    const x = 7 - Math.floor((e.clientY - top) / size);
    return { x, y };
  };

  const placePiece = async (piece, oldRank, oldFile, newRank, newFile) => {
    try {
      const requestOptions = {
        method: "POST",
        body: JSON.stringify({ Piece: piece, Rank: parseInt(oldRank), File: parseInt(oldFile), NewRank: parseInt(newRank), NewFile: parseInt(newFile) }),
        headers: { "Content-Type": "application/json" },
      };
      const response = await fetch(
        "http://localhost:8080/place",
        requestOptions,
      );
      const data = await response.json();
      console.log(data);
    } catch (err) {
      console.log(err);
    }
  };

  const drop = async (e) => {
    const newPosition = copyPosition(state);
    const { x, y } = calculateCoords(e);
    const [p, rank, file] = e.dataTransfer.getData("text").split(",");
    nothigh();
    console.log("dropping", p, x, y);
    newPosition[rank][file] = "";
    newPosition[x][y] = p;
    setState(newPosition);
    await placePiece(p, rank, file, x, y);
  };

  const onDragOver = (e) => {
    e.preventDefault();
  };

  return (
    <div ref={ref} className="pieces" onDrop={drop} onDragOver={onDragOver}>
      {state.map((r, rank) =>
        r.map((f, file) =>
          state[rank][file] ? (
            <Piece
              key={rank + "-" + file}
              updateHighlight={updateHighlight}
              rank={rank}
              file={file}
              piece={state[rank][file]}
            />
          ) : null,
        ),
      )}
    </div>
  );
};
export default Pieces;
