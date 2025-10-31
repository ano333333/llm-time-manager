import { Link } from "react-router-dom";

export default function Capture() {
  return (
    <div style={{ padding: "2rem" }}>
      <h1>📸 キャプチャ設定</h1>
      <p>画面キャプチャの設定画面です。</p>
      <div style={{ marginTop: "2rem" }}>
        <Link to="/">← ホームに戻る</Link>
      </div>
    </div>
  );
}
