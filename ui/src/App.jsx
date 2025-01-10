import { useState } from "react";
import "./App.css";

function App() {
  const [plainTextEdit, setPlainTextEdit] = useState("");
  const [shortenedUrl, setShortenedUrl] = useState("");

  const editPlainText = (e) => {
    setPlainTextEdit(e.target.value);
  };

  const shortenerHandler = () => {
    const options = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        long_url: plainTextEdit,
      }),
    };
    fetch("http://localhost:8090/s", options)
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        setShortenedUrl(data["short_url"]);
      })
      .catch((error) => console.log(error));
  };

  return (
    <>
      <h1>Url shortener service</h1>
      <div className="card">
        <div>
          <p>Paste Url</p>
          <div>
            <div>
              <textarea
                value={plainTextEdit}
                name="string_to_encode"
                rows="4"
                cols="50"
                onChange={(e) => editPlainText(e)}
              ></textarea>
              <br />
              <button onClick={shortenerHandler}>Encode</button>
            </div>
            <div>
              <p>Your url:</p>
              <a href={shortenedUrl}>{shortenedUrl}</a>
              {/* http://{shortenedUrl} */}
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

export default App;
