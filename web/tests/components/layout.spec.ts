import { expect, test } from "@playwright/test";

/**
 * コンポーネントテスト: Layout
 * レイアウトコンポーネントの動作を確認
 */
test.describe("Layoutコンポーネント", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("レイアウトの基本構造が正しい", async ({ page }) => {
    // ヘッダーが存在
    const header = page.getByRole("banner");
    await expect(header).toBeVisible();

    // ナビゲーションが存在
    const nav = page.getByRole("navigation");
    await expect(nav).toBeVisible();

    // メインコンテンツエリアが存在
    const main = page.getByRole("main");
    await expect(main).toBeVisible();
  });

  test("レスポンシブデザインが機能する", async ({ page }) => {
    // デスクトップサイズでの確認
    await page.setViewportSize({ width: 1920, height: 1080 });
    const header = page.getByRole("banner");
    await expect(header).toBeVisible();

    // ナビゲーションが通常表示されることを確認
    const nav = page.getByRole("navigation");
    await expect(nav).toBeVisible();

    // モバイルサイズでの確認
    await page.setViewportSize({ width: 375, height: 667 });
    await expect(header).toBeVisible();
    // 注: 現在の実装ではモバイル特有の動作はないが、将来的にハンバーガーメニュー等を追加する場合はここに検証を追加
  });
});
