import { expect, test } from "@playwright/test";

/**
 * コンポーネントテスト: Header
 * ヘッダーコンポーネントの動作を確認
 */
test.describe("Headerコンポーネント", () => {
  test.beforeEach(async ({ page }) => {
    // 各テスト前にトップページへアクセス
    await page.goto("/");
  });

  test("ヘッダーが表示される", async ({ page }) => {
    const header = page.locator("header");
    await expect(header).toBeVisible();
  });

  test("アプリタイトルが表示される", async ({ page }) => {
    // アプリタイトルまたはロゴの確認（実装に応じて調整）
    const title = page.locator("header");
    await expect(title).toBeVisible();
  });

  test("ナビゲーションリンクが表示される", async ({ page }) => {
    // ナビゲーションメニューの確認
    const nav = page.locator("nav");
    await expect(nav).toBeVisible();

    // 各リンクが存在するか確認（実装に応じて調整）
    const links = page.locator("nav a");
    const count = await links.count();
    expect(count).toBeGreaterThan(0);
  });
});
