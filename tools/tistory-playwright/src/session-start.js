import { chromium } from 'playwright';
import fs from 'node:fs/promises';
import path from 'node:path';

const args = process.argv.slice(2);

function argValue(name, fallback) {
    const index = args.indexOf(name);
    if (index >= 0 && args[index + 1]) {
        return args[index + 1];
    }
    return fallback;
}

const statePath = argValue('--state', 'runtime/tistory/storage-state.json');
const userDataDir = argValue('--user-data-dir', 'runtime/tistory/browser-profile');

await fs.mkdir(path.dirname(statePath), { recursive: true });
await fs.mkdir(userDataDir, { recursive: true });

const context = await chromium.launchPersistentContext(userDataDir, {
    headless: false,
    viewport: { width: 1440, height: 1200 },
});

let page = context.pages()[0];
if (!page) {
    page = await context.newPage();
}

await page.goto('https://www.tistory.com/auth/login', { waitUntil: 'domcontentloaded' });
console.log(JSON.stringify({ event: 'ready' }));

const deadline = Date.now() + 10 * 60 * 1000;
let saved = false;

while (Date.now() < deadline) {
    try {
        const res = await context.request.get(
            'https://www.tistory.com/legacy/member/blog/api/myBlogs',
        );
        if (res.ok()) {
            await context.storageState({ path: statePath });
            console.log(JSON.stringify({ event: 'saved', statePath }));
            saved = true;
            break;
        }
    } catch {
        // Keep polling while the user completes Kakao/Tistory verification.
    }
    await new Promise((resolve) => setTimeout(resolve, 2000));
}

if (!saved) {
    console.log(JSON.stringify({ event: 'timeout' }));
}

await context.close();
