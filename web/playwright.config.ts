import { defineConfig, devices } from "@playwright/test";

/**
 * Playwright設定ファイル
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: "./tests",

  // テスト全体のタイムアウト
  timeout: 30 * 1000,

  // 各アクションのタイムアウト
  expect: {
    timeout: 5000,
  },

  // テスト失敗時の動作
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 2 : undefined,

  // レポーター
  reporter: [["html"], ["list"]],

  // すべてのテストで共通の設定
  use: {
    // ベースURL（開発サーバーが起動している場合）
    baseURL: "http://localhost:5173",

    // トレース記録（失敗時のみ）
    trace: "on-first-retry",

    // スクリーンショット（失敗時のみ）
    screenshot: "only-on-failure",
  },

  // テスト対象のブラウザ設定
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },

    {
      name: "webkit",
      use: { ...devices["Desktop Safari"] },
    },

    // モバイルブラウザテスト（オプション）
    // {
    //   name: 'Mobile Chrome',
    //   use: { ...devices['Pixel 5'] },
    // },
    // {
    //   name: 'Mobile Safari',
    //   use: { ...devices['iPhone 12'] },
    // },
  ],

  // 開発サーバーの自動起動設定（オプション）
  webServer: {
    command: "pnpm run dev",
    url: "http://localhost:5173",
    reuseExistingServer: !process.env.CI,
    stdout: "pipe",
    stderr: "pipe",
  },
});
