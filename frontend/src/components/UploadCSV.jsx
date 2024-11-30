import { useState } from "react";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;

const UploadCSV = ({ onUploadComplete }) => {
  const [file, setFile] = useState(null);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!file) {
      alert("アップロードするファイルを選択してください。");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await axios.post(`${API_URL}/api/upload`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      console.log("ファイルが正常にアップロードされました:", response.data);
      onUploadComplete();
    } catch (error) {
      console.error("ファイルのアップロード中にエラーが発生しました:", error);
      alert("アップロード中にエラーが発生しました。もう一度お試しください。");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>CSVをアップロード</h2>
      <input type="file" onChange={handleFileChange} />
      <button type="submit">アップロード</button>
    </form>
  );
};

export default UploadCSV;
