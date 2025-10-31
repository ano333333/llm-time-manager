import { Link } from "react-router-dom";

export default function Tasks() {
  return (
    <div style={{ padding: "2rem" }}>
      <h1>✅ タスク一覧</h1>
      <p>タスクの管理画面です。</p>
      <div style={{ marginTop: "2rem" }}>
        <Link to="/">← ホームに戻る</Link>
      </div>
    </div>
  );
}
