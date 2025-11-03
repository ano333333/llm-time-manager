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
    const header = page.locator("header");
    await expect(header).toBeVisible();

    // ナビゲーションが存在
    const nav = page.locator("nav");
    await expect(nav).toBeVisible();

    // メインコンテンツエリアが存在
    const main = page.locator("main");
    await expect(main).toBeVisible();
  });

  test("レスポンシブデザインが機能する", async ({ page }) => {
    // デスクトップサイズでの確認
    await page.setViewportSize({ width: 1920, height: 1080 });
    const header = page.locator("header");
    await expect(header).toBeVisible();

    // モバイルサイズでの確認
    await page.setViewportSize({ width: 375, height: 667 });
    await expect(header).toBeVisible();
  });
});
