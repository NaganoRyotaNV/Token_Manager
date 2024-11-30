import { useState } from "react";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;

const AddToken = ({ onTokenAdded }) => {
  const [formData, setFormData] = useState({
    projectName: "",
    token: "",
    permission: "",
    userId: "",
    userName: "",
    expiryDate: "",
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (
      !formData.projectName ||
      !formData.token ||
      !formData.permission ||
      !formData.userId ||
      !formData.userName ||
      !formData.expiryDate
    ) {
      alert("すべてのフィールドを入力してください。");
      return;
    }

    try {
      const response = await axios.post(`${API_URL}/api/tokens`, formData);
      console.log("トークンが追加されました:", response.data);

      setFormData({
        projectName: "",
        token: "",
        permission: "",
        userId: "",
        userName: "",
        expiryDate: "",
      });

      onTokenAdded();
    } catch (error) {
      console.error("トークンの追加中にエラーが発生しました:", error);
      alert("トークンの追加中にエラーが発生しました。もう一度お試しください。");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>トークンを追加</h2>
      {[
        "projectName",
        "token",
        "permission",
        "userId",
        "userName",
        "expiryDate",
      ].map((field) => (
        <input
          key={field}
          type="text"
          name={field}
          placeholder={
            field === "expiryDate"
              ? "有効期限"
              : field === "userId"
              ? "ユーザーID"
              : field === "userName"
              ? "ユーザー名"
              : field === "permission"
              ? "権限"
              : field === "token"
              ? "トークン"
              : "プロジェクト名"
          }
          value={formData[field]}
          onChange={handleChange}
        />
      ))}
      <button type="submit">トークンを追加</button>
    </form>
  );
};

export default AddToken;
