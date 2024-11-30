import { useState } from "react";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;

const DeleteToken = ({ onTokenDeleted }) => {
  const [lineNum, setLineNum] = useState("");

  const handleChange = (e) => {
    setLineNum(e.target.value);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const indexToDelete = parseInt(lineNum, 10) - 1;

    if (isNaN(indexToDelete) || indexToDelete < 0) {
      alert("有効な行番号を入力してください。");
      return;
    }

    const confirmDelete = window.confirm("本当にこのトークンを削除しますか?");
    if (!confirmDelete) return;

    try {
      const response = await axios.delete(`${API_URL}/api/tokens`, {
        params: { line: indexToDelete },
      });
      console.log("トークンが削除されました:", response.data);

      onTokenDeleted();
    } catch (error) {
      console.error("トークンの削除中にエラーが発生しました:", error);
      alert("トークンの削除中にエラーが発生しました。もう一度お試しください。");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>トークンを削除</h2>
      <div>
        <input
          type="number"
          name="lineNum"
          placeholder="行番号"
          value={lineNum}
          onChange={handleChange}
          min="1"
        />
        <button type="submit">トークンを削除</button>
      </div>
    </form>
  );
};

export default DeleteToken;
