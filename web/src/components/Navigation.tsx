import { Link, useLocation } from "react-router-dom";
import "./navigation.css";

interface NavItem {
  path: string;
  label: string;
  icon: string;
}

const navItems: NavItem[] = [
  { path: "/", label: "ホーム", icon: "🏠" },
  { path: "/chat", label: "チャット", icon: "💬" },
  { path: "/tasks", label: "タスク", icon: "✅" },
  { path: "/goals", label: "目標", icon: "🎯" },
  { path: "/capture", label: "キャプチャ", icon: "📸" },
  { path: "/settings/local", label: "設定", icon: "⚙️" },
];

/**
 * ナビゲーションコンポーネント
 * サイドバーまたはボトムナビゲーションとして機能
 */
export default function Navigation() {
  const location = useLocation();

  const isActive = (path: string) => {
    if (path === "/") {
      return location.pathname === path;
    }
    return location.pathname.startsWith(path);
  };

  return (
    <nav className="navigation">
      <ul className="nav-list">
        {navItems.map((item) => (
          <li key={item.path} className="nav-item">
            <Link
              to={item.path}
              className={`nav-link ${isActive(item.path) ? "active" : ""}`}
              aria-current={isActive(item.path) ? "page" : undefined}
            >
              <span className="nav-icon" role="img" aria-label={item.label}>
                {item.icon}
              </span>
              <span className="nav-label">{item.label}</span>
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
}

