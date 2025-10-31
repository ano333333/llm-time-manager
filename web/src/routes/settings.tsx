import { Link } from "react-router-dom";

export default function Settings() {
  return (
    <div style={{ padding: "2rem" }}>
      <h1>⚙️ ローカル設定</h1>
      <p>アプリケーションの設定画面です。</p>
      <div style={{ marginTop: "2rem" }}>
        <Link to="/">← ホームに戻る</Link>
      </div>
    </div>
  );
}
