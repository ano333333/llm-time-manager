import { expect, test } from "@playwright/test";

/**
 * E2Eテスト: ナビゲーション
 * アプリ全体のページ遷移と基本的な動作を確認
 */
test.describe("ナビゲーション", () => {
  test("トップページが正しく表示される", async ({ page }) => {
    await page.goto("/");

    // ページタイトルの確認
    await expect(page).toHaveTitle(/LLM時間管理ツール/);

    // ヘッダーが存在する確認
    const header = page.locator("header");
    await expect(header).toBeVisible();
  });

  test("チャット、タスク、目標ページへ遷移できる", async ({ page }) => {
    await page.goto("/");

    // チャットページへの遷移
    const chatLink = page.getByRole("link", { name: /chat/i });
    await chatLink.click();
    await expect(page).toHaveURL(/.*chat/);

    // タスクページへの遷移
    const tasksLink = page.getByRole("link", { name: /task/i });
    await tasksLink.click();
    await expect(page).toHaveURL(/.*tasks/);

    // 目標ページへの遷移
    const goalsLink = page.getByRole("link", { name: /goal/i });
    await goalsLink.click();
    await expect(page).toHaveURL(/.*goals/);
  });

  test("404ページが正しく表示される", async ({ page }) => {
    await page.goto("/non-existent-page");

    // 404ページ固有のコンテンツを確認
    await expect(
      page.getByRole("heading", { name: /404.*ページが見つかりません/ }),
    ).toBeVisible();
    await expect(
      page.getByText(/お探しのページは存在しません/),
    ).toBeVisible();

    // ホームに戻るリンクの確認
    const homeLink = page.getByRole("link", { name: /ホームに戻る/ });
    await expect(homeLink).toBeVisible();
  });
});
