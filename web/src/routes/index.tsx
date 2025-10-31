import { Link } from "react-router-dom";

export default function Home() {
  return (
    <div style={{ padding: "2rem" }}>
      <h1>LLM時間管理ツール</h1>
      <p>ホーム画面へようこそ</p>
      <nav style={{ marginTop: "2rem" }}>
        <ul style={{ listStyle: "none", display: "flex", flexDirection: "column", gap: "1rem" }}>
          <li>
            <Link to="/chat">💬 LLMチャット</Link>
          </li>
          <li>
            <Link to="/goals">🎯 目標一覧</Link>
          </li>
          <li>
            <Link to="/tasks">✅ タスク一覧</Link>
          </li>
          <li>
            <Link to="/capture">📸 キャプチャ設定</Link>
          </li>
          <li>
            <Link to="/settings/local">⚙️ ローカル設定</Link>
          </li>
        </ul>
      </nav>
    </div>
  );
}
