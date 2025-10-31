import { Link } from "react-router-dom";

export default function Chat() {
  return (
    <div style={{ padding: "2rem" }}>
      <h1>💬 LLMチャット</h1>
      <p>LLMとの対話画面です。</p>
      <div style={{ marginTop: "2rem" }}>
        <Link to="/">← ホームに戻る</Link>
      </div>
    </div>
  );
}
