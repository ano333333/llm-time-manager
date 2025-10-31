import { useState } from "react";
import "./header.css";

/**
 * 共通ヘッダーコンポーネント
 * 現在の日付、検索、グローバルメニューを表示
 */
export default function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isSearchOpen, setIsSearchOpen] = useState(false);

  // 現在の日付をフォーマット
  const formatDate = () => {
    const now = new Date();
    const year = now.getFullYear();
    const month = now.getMonth() + 1;
    const day = now.getDate();
    const weekdays = ["日", "月", "火", "水", "木", "金", "土"];
    const weekday = weekdays[now.getDay()];

    return `${year}年${month}月${day}日（${weekday}）`;
  };

  return (
    <header className="header">
      <div className="header-content">
        <div className="header-left">
          <h1 className="header-title">LLM時間管理</h1>
          <span className="header-date">{formatDate()}</span>
        </div>

        <div className="header-right">
          {/* 検索ボタン */}
          <button
            className="header-button"
            onClick={() => setIsSearchOpen(!isSearchOpen)}
            aria-label="検索"
            title="検索 (Cmd/Ctrl+K)"
          >
            <svg
              width="20"
              height="20"
              viewBox="0 0 20 20"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M9 17A8 8 0 1 0 9 1a8 8 0 0 0 0 16zM19 19l-4.35-4.35"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>

          {/* メニューボタン */}
          <button
            className="header-button"
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            aria-label="メニュー"
            title="メニュー"
          >
            <svg
              width="20"
              height="20"
              viewBox="0 0 20 20"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M3 10h14M3 5h14M3 15h14"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
          </button>
        </div>
      </div>

      {/* 検索パネル */}
      {isSearchOpen && (
        <div className="search-panel">
          <input
            type="text"
            className="search-input"
            placeholder="タスク、目標、設定を検索..."
            autoFocus
          />
        </div>
      )}

      {/* グローバルメニュー */}
      {isMenuOpen && (
        <div className="global-menu">
          <nav>
            <ul className="menu-list">
              <li>
                <a href="#" className="menu-item">
                  <span>ショートカット</span>
                  <span className="menu-shortcut">Cmd+K</span>
                </a>
              </li>
              <li>
                <a href="#" className="menu-item">
                  <span>ヘルプ</span>
                </a>
              </li>
              <li>
                <a href="#" className="menu-item">
                  <span>バージョン情報</span>
                </a>
              </li>
            </ul>
          </nav>
        </div>
      )}
    </header>
  );
}

