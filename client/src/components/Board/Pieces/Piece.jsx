const Piece = ({ updateHighlight, rank, file, piece }) => {

  const getValidMoves = async () => {
    try {
      const requestOptions = {
        method: "POST",
        body: JSON.stringify({ Piece: piece, Rank: rank, File: file }),
        headers: { "Content-Type": "application/json" },
      };
      console.log(piece, rank, file);
      const response = await fetch(
        "http://localhost:8080/moves",
        requestOptions,
      );
      const data = await response.json();
      console.log(data);
      for (const move of data) {
        updateHighlight(7 - move.Rank, move.File);
      }
    } catch (err) {
      console.log(err);
    }
  };

  const start = (e) => {
    getValidMoves();
    e.dataTransfer.setData("text/plain", `${piece},${rank},${file}`);
    console.log("picked up", piece, rank, file);
    setTimeout(() => {
      e.target.style.display = "none";
    }, 0);
  };

  const end = (e) => (e.target.style.display = "block");

  return (
    <div
      draggable={true}
      onDragStart={start}
      onDragEnd={end}
      className={`piece ${piece} p-${file}${rank}`}
    />
  );
};

export default Piece;
