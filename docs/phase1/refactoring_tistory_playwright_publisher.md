# Phase 1 Tistory 발행 리팩토링 설계

## 1. 배경

현재 Phase 1의 Tistory 발행은 Selenium Remote WebDriver로 Tistory 에디터 화면을 직접 조작한다.

현재 방식의 문제:

- Kakao 추가 인증을 자동화 흐름 안에서 안정적으로 처리하기 어렵다.
- Tistory 에디터 DOM과 CSS selector가 자주 바뀔 수 있다.
- Markdown 모드 전환, 제목 입력, 본문 입력, 완료 버튼 클릭이 모두 화면 구조에 강하게 의존한다.
- 로그인 성공 기준이 명확하지 않아 사용자가 인증을 마쳐도 백엔드가 성공 상태를 안정적으로 판단하기 어렵다.

Viruagent는 Playwright로 로그인 세션을 확보하고, 이후 Tistory 내부 HTTP 엔드포인트를 호출하는 방식을 사용한다. 이 프로젝트도 같은 방향으로 리팩토링한다.

참고:

- https://github.com/greekr4/Viruagent
- https://raw.githubusercontent.com/greekr4/Viruagent/main/src/lib/login.js
- https://raw.githubusercontent.com/greekr4/Viruagent/main/src/lib/tistory.js

## 2. 목표

Selenium 기반 에디터 DOM 조작을 제거하고, 다음 구조로 전환한다.

```text
사용자 수동 로그인
        ↓
Playwright가 Tistory/Kakao 세션 저장
        ↓
백엔드가 저장된 세션 상태 확인
        ↓
카테고리 조회와 글 발행은 Tistory 내부 HTTP API 호출
```

리팩토링 후 목표:

- 사용자가 Kakao 추가 인증을 직접 처리할 수 있어야 한다.
- 앱은 Kakao/Tistory ID/PW를 입력받거나 저장하지 않는다.
- 카테고리 조회와 글 발행은 에디터 DOM 조작 없이 처리한다.
- 한 번 불러온 Tistory 카테고리 목록은 Settings 화면에 남겨두고, 사용자가 기본 카테고리를 계속 변경할 수 있어야 한다.
- 실패 원인은 로그인 필요, 세션 만료, 카테고리 조회 실패, 발행 실패로 구분한다.

## 3. 비목표

- Tistory 내부 API를 공식 API처럼 추상화하지 않는다.
- 다중 사용자 세션 관리는 Phase 1 범위에 넣지 않는다.
- 공개 배포용 인증/보안 모델은 만들지 않는다.
- AI 글 생성 흐름은 이 리팩토링 범위에 넣지 않는다.

## 4. 아키텍처 변경

### 4.1 현재 구조

```text
Go backend
  └─ Selenium Remote WebDriver
       └─ Tistory editor DOM 조작
```

### 4.2 변경 구조

```text
Go backend
  ├─ posts/settings/template API
  └─ tistory publisher adapter
       ├─ Playwright session helper 호출
       ├─ storage state 또는 cookie 파일 읽기
       └─ Tistory 내부 HTTP API 호출
```

Playwright는 로그인 세션 확보에만 사용한다. 글 발행 자체는 HTTP 요청으로 처리한다.

## 5. 세션 관리 설계

### 5.1 세션 저장 위치

로컬 전용 세션 파일을 사용한다.

```text
runtime/tistory/storage-state.json
```

이 파일은 Git에 커밋하지 않는다.

`.gitignore`에 다음 항목을 추가한다.

```text
runtime/
*.storage-state.json
```

### 5.2 로그인 흐름

1. 사용자가 Settings 화면에서 `Tistory 로그인 연결` 버튼을 누른다.
2. 프론트는 `POST /api/tistory/session/start`를 호출한다.
3. 백엔드는 Playwright helper를 실행해 headful Chromium을 연다.
4. 사용자는 열린 브라우저에서 Tistory/Kakao 로그인을 직접 완료한다.
5. 사용자가 Settings 화면에서 `인증 완료 확인`을 누른다.
6. 프론트는 `POST /api/tistory/session/confirm`을 호출한다.
7. 백엔드는 현재 브라우저 세션을 `storage-state.json`으로 저장한다.
8. 백엔드는 `myBlogs` 또는 `manage/newpost` 접근으로 세션 유효성을 확인한다.

### 5.3 로그인 성공 기준

로그인 성공은 다음 중 하나가 성공하면 인정한다.

- `https://www.tistory.com/legacy/member/blog/api/myBlogs` 호출 성공
- `{blogUrl}/manage/newpost` 접근 후 `window.Config.blog` 확인 성공

Kakao 로그인 화면에 머물러 있거나 Tistory 관리 페이지에 접근하지 못하면 실패로 처리한다.

## 6. Tistory 내부 API 설계

### 6.1 블로그 조회

저장된 세션 쿠키로 다음 엔드포인트를 호출한다.

```http
GET https://www.tistory.com/legacy/member/blog/api/myBlogs
```

사용 목적:

- 현재 로그인된 사용자의 블로그 목록 확인
- 기본 blogUrl 자동 감지 후보 제공
- 세션 유효성 확인

### 6.2 카테고리 조회

다음 페이지 HTML을 가져온다.

```http
GET {blogUrl}/manage/newpost
```

HTML 안의 `window.Config.blog.categories`를 파싱해 카테고리 목록을 만든다.

기존 `POST /api/tistory/categories/fetch`는 유지하되, 요청 body에서 login 정보를 제거한다.

변경 전:

```json
{
    "blogUrl": "https://hyeonyway.tistory.com",
    "login": {
        "id": "kakao-account@example.com",
        "password": "input-at-fetch-time"
    }
}
```

변경 후:

```json
{
    "blogUrl": "https://hyeonyway.tistory.com"
}
```

조회에 성공하면 카테고리 목록을 settings에 캐시한다. 이후 Settings 화면은 저장된 목록을 계속 보여주며, 사용자는 다시 조회하지 않아도 기본 카테고리를 변경할 수 있다.

저장할 카테고리 목록 예:

```json
[
    { "categoryId": "0", "label": "카테고리 없음" },
    { "categoryId": "1550884", "label": "알고리즘 문제 풀이" },
    { "categoryId": "1550885", "label": "취준" }
]
```

저장 정책:

- 카테고리 목록은 `settings`에 JSON 문자열로 저장한다.
- key는 `tistory.categories`를 사용한다.
- 선택된 기본 카테고리는 기존처럼 `tistory.category_id`에 저장한다.
- 카테고리 목록을 다시 불러오면 `tistory.categories`를 최신 목록으로 덮어쓴다.
- 기존 기본 카테고리가 새 목록에 없으면 Settings 화면에서 선택 필요 상태로 표시한다.

### 6.3 글 발행

Tistory 에디터 DOM을 조작하지 않고 다음 엔드포인트를 호출한다.

```http
POST {blogUrl}/manage/post.json
```

요청 payload는 Tistory 내부 API가 요구하는 필드에 맞춘다.

초기 후보 필드:

```json
{
    "title": "글 제목",
    "content": "Markdown 본문",
    "category": "1550884",
    "visibility": 0,
    "tag": "알고리즘,Java",
    "published": 0,
    "slogan": "",
    "acceptComment": 1
}
```

visibility 기준은 Viruagent 구현을 따른다.

```text
0 = 비공개
15 = 보호
20 = 공개
```

Phase 1 기본값은 계속 비공개다.

## 7. API 변경

### 7.1 추가 API

```http
POST /api/tistory/session/start
POST /api/tistory/session/confirm
GET /api/tistory/session/status
DELETE /api/tistory/session
```

역할:

- `start`: Playwright 로그인 브라우저 실행
- `confirm`: 현재 로그인 세션 저장 및 유효성 확인
- `status`: 저장된 세션 유효성 확인
- `delete`: 로컬 세션 파일 삭제

### 7.2 변경 API

```http
POST /api/tistory/categories/fetch
POST /api/posts/{id}/publish/tistory
```

변경점:

- 요청 body에서 Kakao/Tistory login 제거
- 저장된 세션이 없으면 `TISTORY_SESSION_REQUIRED`
- 세션이 만료되었으면 `TISTORY_SESSION_EXPIRED`
- 카테고리 조회 성공 시 `tistory.categories`에 목록 저장
- 발행 실패 시 내부 API 응답 메시지를 마스킹해 반환

## 8. 프론트 변경

Settings 화면에 Tistory 세션 영역을 추가한다.

상태:

- 연결 안 됨
- 브라우저 로그인 대기 중
- 연결됨
- 세션 만료

버튼:

- `로그인 브라우저 열기`
- `인증 완료 확인`
- `세션 연결 해제`
- `카테고리 불러오기`

카테고리 영역:

- 마지막으로 불러온 카테고리 목록을 select 또는 radio list로 표시한다.
- 사용자는 저장된 목록 안에서 기본 카테고리를 언제든 변경할 수 있다.
- `카테고리 불러오기` 버튼은 목록을 새로 동기화한다.
- 마지막 동기화 시간을 함께 표시한다.
- 저장된 목록이 없으면 `카테고리 불러오기`를 먼저 안내한다.

발행 모달에서는 Kakao ID/PW 입력란을 제거한다. 사용자는 발행 직전에 카테고리, 태그, 공개 범위만 확인한다.

## 9. 백엔드 변경

### 9.1 제거 대상

다음 책임을 제거한다.

- Selenium Remote WebDriver 세션 생성
- Tistory 에디터 DOM selector 관리
- 제목/본문/태그/버튼 클릭 자동화
- 요청 body로 받은 login ID/PW 처리

### 9.2 추가 대상

추가할 패키지 후보:

```text
backend/internal/tistory/session
backend/internal/tistory/client
backend/internal/publisher/tistory
```

역할:

- `session`: Playwright helper 실행, storage state 파일 관리
- `client`: 저장된 쿠키로 Tistory 내부 HTTP API 호출
- `publisher`: post 데이터를 Tistory publish payload로 변환

settings 저장 key:

```text
tistory.blog_url
tistory.category_id
tistory.default_visibility
tistory.default_tags
tistory.categories
tistory.categories_synced_at
```

## 10. Playwright helper

Playwright는 Node 스크립트로 둔다.

후보 경로:

```text
tools/tistory-playwright/package.json
tools/tistory-playwright/src/session.ts
```

명령:

```bash
npm run tistory:login:start
npm run tistory:login:confirm
```

Go 백엔드는 `exec.CommandContext`로 helper를 호출한다. helper는 JSON stdout으로 결과를 반환한다.

예:

```json
{
    "ok": true,
    "storageStatePath": "runtime/tistory/storage-state.json"
}
```

## 11. 보안 기준

- Kakao/Tistory ID/PW를 앱 화면에서 입력받지 않는다.
- 세션 파일은 Git에 커밋하지 않는다.
- 세션 파일 경로는 README에 명시한다.
- 로그에 쿠키, authorization 값, session key를 출력하지 않는다.
- 세션 삭제 API를 제공한다.

## 12. 리스크

- Tistory 내부 API는 공식 API가 아니므로 언제든 변경될 수 있다.
- `manage/post.json` payload 필드는 실제 요청을 보며 보정해야 한다.
- Playwright helper를 추가하면 Node 런타임 의존성이 생긴다.
- 로컬 브라우저를 띄우는 방식은 서버 배포 환경에는 맞지 않는다.

## 13. 구현 순서

1. Selenium publisher 코드를 제거하지 말고 새 Playwright 기반 publisher를 병렬 추가
2. Playwright helper 스캐폴딩
3. 세션 파일 저장/확인 API 추가
4. Tistory 내부 API client 추가
5. 카테고리 조회를 내부 API 방식으로 전환
6. 조회한 카테고리 목록을 settings에 캐시
7. 발행을 `manage/post.json` 방식으로 전환
8. 프론트 Settings 세션 UI 수정
9. Settings 화면에서 저장된 카테고리 목록을 계속 선택 가능하게 표시
10. 발행 모달에서 ID/PW 입력 제거
11. README와 `known_issues.md` 업데이트
12. Selenium 의존성 제거 여부 결정

## 14. 완료 기준

- 사용자가 Playwright 브라우저에서 Kakao/Tistory 로그인을 직접 완료할 수 있다.
- 저장된 세션으로 블로그 목록 또는 관리 페이지 접근을 확인할 수 있다.
- 카테고리 조회가 Selenium DOM 조작 없이 동작한다.
- 한 번 조회한 카테고리 목록이 Settings 화면에 남아 있고, 사용자가 기본 카테고리를 계속 변경할 수 있다.
- Tistory 비공개 발행이 에디터 DOM 조작 없이 동작한다.
- 발행 모달에 Kakao/Tistory ID/PW 입력란이 없다.
- 세션 없음, 세션 만료, 발행 실패가 구분되어 표시된다.
