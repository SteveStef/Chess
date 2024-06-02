import "./App.css";
import { useEffect } from 'react';
import Board from './components/Board/Board';

function App() {

  const init = () => {
    try {
      const url = `http://localhost:8000/initboard`;
      console.log(url);
    } catch(err) {
      console.log(err);
    }
  }

  useEffect(() => {
    init();
  }, []);

  return <div className="App">
    <Board />
  </div>;
}

export default App;
