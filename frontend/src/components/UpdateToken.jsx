import { useState } from "react";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;

const UpdateToken = ({ onTokenUpdated }) => {
  const [projectName, setProjectName] = useState("");
  const [tokens, setTokens] = useState([]);

  const handleProjectNameChange = (e) => {
    setProjectName(e.target.value);
  };

  const handleTokenChange = (index, e) => {
    const updatedTokens = [...tokens];
    updatedTokens[index][e.target.name] = e.target.value;
    setTokens(updatedTokens);
  };

  const handleFetchTokens = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/tokens`, {
        params: { projectName },
      });
      setTokens(response.data);
    } catch (error) {
      console.error("トークンの取得中にエラーが発生しました:", error);
      alert("トークンの取得中にエラーが発生しました。もう一度お試しください。");
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.put(`${API_URL}/api/tokens`, {
        projectName,
        tokens,
      });
      console.log("トークンが更新されました:", response.data);
      onTokenUpdated();
    } catch (error) {
      console.error("トークンの更新中にエラーが発生しました:", error);
      alert("トークンの更新中にエラーが発生しました。もう一度お試しください。");
    }
  };

  return (
    <div>
      <h2>トークンを更新</h2>
      <div>
        <input
          type="text"
          placeholder="プロジェクト名"
          value={projectName}
          onChange={handleProjectNameChange}
        />
        <button onClick={handleFetchTokens}>トークンを取得</button>
      </div>
      {tokens.length > 0 && (
        <form onSubmit={handleSubmit}>
          {tokens.map((token, index) => (
            <div key={index} className="token-row">
              <input
                type="text"
                name="token"
                placeholder="トークン"
                value={token.token}
                onChange={(e) => handleTokenChange(index, e)}
              />
              <input
                type="text"
                name="permission"
                placeholder="権限"
                value={token.permission}
                onChange={(e) => handleTokenChange(index, e)}
              />
              <input
                type="text"
                name="userId"
                placeholder="ユーザーID"
                value={token.userId}
                onChange={(e) => handleTokenChange(index, e)}
              />
              <input
                type="text"
                name="userName"
                placeholder="ユーザー名"
                value={token.userName}
                onChange={(e) => handleTokenChange(index, e)}
              />
              <input
                type="text"
                name="expiryDate"
                placeholder="有効期限"
                value={token.expiryDate}
                onChange={(e) => handleTokenChange(index, e)}
              />
            </div>
          ))}
          <button type="submit">トークンを更新</button>
        </form>
      )}
    </div>
  );
};

export default UpdateToken;
