# WebView-ネイティブブリッジインターフェース

## 概要

WebView（フロントエンド）とネイティブ層（OS固有機能）を繋ぐ JavaScript ブリッジの仕様を定義する。

## ブリッジ I/F 定義

```ts
interface Bridge {
  /**
   * スクリーンショットを撮影する
   * @param params キャプチャパラメータ
   * @returns 保存先パスとサムネイルパス
   */
  captureScreenshot(params?: CaptureParams): Promise<CaptureResult>;

  /**
   * 権限をリクエストする
   * @param kind 権限の種類
   * @returns 権限状態
   */
  requestPermission(kind: PermissionKind): Promise<PermissionStatus>;

  /**
   * アプリ情報を取得する
   * @returns プラットフォーム情報
   */
  getAppInfo(): Promise<AppInfo>;

  /**
   * ファイルパスを開く（OS のファイルマネージャで表示）
   * @param path ファイルパス
   */
  openPath(path: string): Promise<void>;

  /**
   * 通知を表示する
   * @param notification 通知内容
   */
  showNotification(notification: Notification): Promise<void>;
}

interface CaptureParams {
  /** キャプチャモード */
  mode?: 'full' | 'window' | 'region';
  /** 領域指定（region モードの場合） */
  region?: {
    x: number;
    y: number;
    width: number;
    height: number;
  };
  /** 除外ウィンドウID（任意） */
  excludeWindows?: string[];
}

interface CaptureResult {
  /** 保存先パス */
  path: string;
  /** サムネイルパス（任意） */
  thumbPath?: string;
  /** キャプチャ日時 */
  capturedAt: string; // ISO 8601
  /** メタデータ */
  meta?: {
    width: number;
    height: number;
    format: string;
    displayId?: string;
  };
}

type PermissionKind = 'capture' | 'notification';

type PermissionStatus = 'granted' | 'denied' | 'not_determined';

interface AppInfo {
  /** プラットフォーム */
  platform: 'ios' | 'android' | 'mac' | 'win' | 'linux';
  /** OS バージョン */
  osVersion: string;
  /** アプリバージョン */
  version: string;
  /** アーキテクチャ */
  arch: 'x64' | 'arm64';
}

interface Notification {
  /** タイトル */
  title: string;
  /** メッセージ */
  message: string;
  /** アイコン（任意） */
  icon?: string;
}
```

## プラットフォーム別実装

### macOS (Swift)

```swift
import WebKit

class BridgeHandler: NSObject, WKScriptMessageHandler {
    func userContentController(
        _ userContentController: WKUserContentController,
        didReceive message: WKScriptMessage
    ) {
        guard let body = message.body as? [String: Any],
              let method = body["method"] as? String else {
            return
        }
        
        switch method {
        case "captureScreenshot":
            handleCaptureScreenshot(params: body["params"], callback: body["callback"])
        case "requestPermission":
            handleRequestPermission(kind: body["kind"], callback: body["callback"])
        case "getAppInfo":
            handleGetAppInfo(callback: body["callback"])
        default:
            break
        }
    }
}

// WebView 初期化時に登録
webView.configuration.userContentController.add(bridgeHandler, name: "bridge")
```

### Windows (C++)

```cpp
#include <wil/com.h>
#include <WebView2.h>

// WebView2 の AddWebMessageReceived を使用
webview->add_WebMessageReceived(
    Callback<ICoreWebView2WebMessageReceivedEventHandler>(
        [](ICoreWebView2* sender, ICoreWebView2WebMessageReceivedEventArgs* args) {
            wil::unique_cotaskmem_string message;
            args->get_WebMessageAsJson(&message);
            
            // JSON パース
            auto json = nlohmann::json::parse(message.get());
            std::string method = json["method"];
            
            if (method == "captureScreenshot") {
                handleCaptureScreenshot(json["params"], json["callback"]);
            }
            // ...
            
            return S_OK;
        }
    ).Get(),
    &token
);
```

### Linux (Rust)

```rust
use webkit2gtk::WebView;
use webkit2gtk::UserContentManager;

let context = webview.context().unwrap();
let manager = context.user_content_manager().unwrap();

manager.register_script_message_handler("bridge");
manager.connect_script_message_received(None, |_manager, result| {
    let value = result.value().unwrap();
    let json: serde_json::Value = serde_json::from_str(&value.to_string()).unwrap();
    
    let method = json["method"].as_str().unwrap();
    match method {
        "captureScreenshot" => handle_capture_screenshot(json["params"].clone()),
        "requestPermission" => handle_request_permission(json["kind"].clone()),
        _ => {}
    }
});
```

## JavaScript 側の使用例

```ts
// WebView 内での呼び出し
async function takeScreenshot() {
  try {
    const result = await window.bridge.captureScreenshot({
      mode: 'full'
    });
    
    console.log('Captured:', result.path);
    
    // サーバーにパスを送信して DB 保存
    await fetch('/screenshots', {
      method: 'POST',
      body: JSON.stringify(result)
    });
  } catch (error) {
    console.error('Capture failed:', error);
  }
}

// 権限確認
async function checkPermission() {
  const status = await window.bridge.requestPermission('capture');
  
  if (status === 'denied') {
    alert('画面収録権限が拒否されています。設定から許可してください。');
  }
}

// アプリ情報取得
async function getInfo() {
  const info = await window.bridge.getAppInfo();
  console.log('Platform:', info.platform);
  console.log('Version:', info.version);
}
```

## ルーティング/URL 設計（WebView 内 SPA）

### ルート定義

- `/` - ホーム
- `/chat` - チャット
- `/goals` - 目標一覧
- `/tasks` - タスク一覧
  - クエリパラメータ: `?filter=today|week|open`
- `/tasks/:id` - タスク詳細
- `/capture` - キャプチャ設定
- `/settings/local` - ローカル設定

### モーダル

URL クエリパラメータで制御：

- `?modal=new-task` - 新規タスク作成モーダル
- `?modal=new-goal` - 新規目標作成モーダル
- `?modal=permissions` - 権限説明モーダル
- `?modal=shortcuts` - ショートカット一覧

### ディープリンク/スキーム

ネイティブ層からの起動時に使用：

- `mytime://capture` - キャプチャ画面を開く
- `mytime://task/<id>` - 特定タスクの詳細を開く
- `mytime://chat` - チャット画面を開く

### React Router 実装例

```tsx
import { BrowserRouter, Routes, Route } from 'react-router-dom';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/chat" element={<ChatPage />} />
        <Route path="/goals" element={<GoalsPage />} />
        <Route path="/tasks" element={<TasksPage />} />
        <Route path="/tasks/:id" element={<TaskDetailPage />} />
        <Route path="/capture" element={<CapturePage />} />
        <Route path="/settings/local" element={<SettingsPage />} />
      </Routes>
      <ModalRouter /> {/* モーダルはクエリパラメータで制御 */}
    </BrowserRouter>
  );
}
```

## エラーハンドリング

### ブリッジエラー

```ts
class BridgeError extends Error {
  constructor(
    public code: string,
    message: string
  ) {
    super(message);
    this.name = 'BridgeError';
  }
}

// 使用例
try {
  await window.bridge.captureScreenshot();
} catch (error) {
  if (error instanceof BridgeError) {
    switch (error.code) {
      case 'PERMISSION_DENIED':
        // 権限エラー処理
        break;
      case 'BRIDGE_UNAVAILABLE':
        // ブリッジ未初期化エラー
        break;
      default:
        // その他のエラー
    }
  }
}
```

### タイムアウト

```ts
function withTimeout<T>(
  promise: Promise<T>,
  ms: number
): Promise<T> {
  return Promise.race([
    promise,
    new Promise<T>((_, reject) =>
      setTimeout(() => reject(new Error('Timeout')), ms)
    )
  ]);
}

// 使用
await withTimeout(
  window.bridge.captureScreenshot(),
  5000 // 5秒でタイムアウト
);
```

## セキュリティ考慮事項

- ブリッジは `localhost` からのリクエストのみ受け付ける
- CSP（Content Security Policy）を設定
- ブリッジメソッドの入力バリデーション
- ファイルパスのサニタイゼーション（パストラバーサル対策）

```html
<!-- CSP 例 -->
<meta http-equiv="Content-Security-Policy" 
      content="default-src 'self'; connect-src 'self' http://localhost:*">
```

