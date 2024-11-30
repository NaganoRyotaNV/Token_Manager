import { useState, useEffect } from "react";
import axios from "axios";
import { saveAs } from "file-saver";
import AddToken from "./components/AddToken";
import TokenList from "./components/TokenList";
import UpdateToken from "./components/UpdateToken";
import DeleteToken from "./components/DeleteToken";
import UploadCSV from "./components/UploadCSV";
import "./App.css";

const API_URL = import.meta.env.VITE_API_URL;

function App() {
  const [tokens, setTokens] = useState([]);
  const [csvUploaded, setCsvUploaded] = useState(false);

  const fetchTokens = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/tokens`);
      setTokens(response.data);
    } catch (error) {
      console.error("トークンの取得中にエラーが発生しました:", error);
      alert(`トークンの取得中にエラーが発生しました: ${error.message}`);
    }
  };

  const saveTokensAsCSV = () => {
    if (tokens.length === 0) {
      alert("保存するトークンがありません。");
      return;
    }

    const csvContent = tokens
      .map(
        (token) =>
          `${token.projectName},${token.token},${token.permission},${token.userId},${token.userName},${token.expiryDate}`
      )
      .join("\n");

    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
    const fileName = window.prompt(
      "保存するファイル名を入力してください:",
      "tokens.csv"
    );
    if (fileName) {
      saveAs(blob, fileName);
    }
  };

  const handleUploadComplete = () => {
    setCsvUploaded(true);
    fetchTokens();
  };

  useEffect(() => {
    fetchTokens();
  }, []);

  return (
    <div className="app-container">
      <h1>トークン管理</h1>
      <div className="actions-container">
        <UploadCSV onUploadComplete={handleUploadComplete} />
        <AddToken onTokenAdded={fetchTokens} />
        <UpdateToken onTokenUpdated={fetchTokens} />
        <DeleteToken onTokenDeleted={fetchTokens} />
      </div>
      {csvUploaded ? (
        <TokenList tokens={tokens} />
      ) : (
        <p className="no-data-message">CSVをアップロードしてください。</p>
      )}
      <button className="save-button" onClick={saveTokensAsCSV}>
        CSVとして保存
      </button>
    </div>
  );
}

export default App;
