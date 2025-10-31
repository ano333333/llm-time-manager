import { ReactNode } from "react";
import Header from "./Header";
import Navigation from "./Navigation";
import "./layout.css";

interface LayoutProps {
  children: ReactNode;
}

/**
 * 共通レイアウトコンポーネント
 * ヘッダーとナビゲーションを含み、各ページのコンテンツを表示する
 */
export default function Layout({ children }: LayoutProps) {
  return (
    <div className="layout">
      <Header />
      <div className="layout-body">
        <Navigation />
        <main className="layout-content">{children}</main>
      </div>
    </div>
  );
}

