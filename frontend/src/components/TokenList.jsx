import React from "react";

const TokenList = ({ tokens }) => {
  if (!tokens || tokens.length === 0) {
    return <p>表示するトークンがありません。</p>;
  }

  return (
    <div>
      <h2>トークンリスト</h2>
      <table>
        <thead>
          <tr>
            <th>#</th>
            <th>プロジェクト名</th>
            <th>トークン</th>
            <th>権限</th>
            <th>ユーザーID</th>
            <th>ユーザー名</th>
            <th>有効期限</th>
          </tr>
        </thead>
        <tbody>
          {tokens.map((token, index) => (
            <tr key={index}>
              <td>{index + 1}</td>
              <td>{token.projectName}</td>
              <td>{token.token}</td>
              <td>{token.permission}</td>
              <td>{token.userId}</td>
              <td>{token.userName}</td>
              <td>{token.expiryDate}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default TokenList;
