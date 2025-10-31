import { Link } from "react-router-dom";

export default function NotFound() {
  return (
    <div style={{ padding: "2rem" }}>
      <h1>404 - ページが見つかりません</h1>
      <p>お探しのページは存在しません。</p>
      <div style={{ marginTop: "2rem" }}>
        <Link to="/">← ホームに戻る</Link>
      </div>
    </div>
  );
}
