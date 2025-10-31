import { Link } from "react-router-dom";

export default function Goals() {
  return (
    <div style={{ padding: "2rem" }}>
      <h1>🎯 目標一覧</h1>
      <p>目標の管理画面です。</p>
      <div style={{ marginTop: "2rem" }}>
        <Link to="/">← ホームに戻る</Link>
      </div>
    </div>
  );
}
