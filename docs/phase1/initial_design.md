# DevLog Studio Phase 1 설계 문서

## 1. Phase 1 목표

Phase 1의 목표는 **수동 입력 기반 알고리즘 풀이 글 생성 및 Tistory 비공개 발행**이다.

초기 기능은 작게 잡는다.

```text
알고리즘 풀이 정보 입력
        ↓
템플릿 기반 Markdown 생성
        ↓
미리보기
        ↓
Tistory 비공개 초안 발행
        ↓
작성 이력 저장
```

Phase 1에서는 프로그래머스 자동 추출, AI 설명 생성, GitHub 분석, Notion/MCP 연동은 구현하지 않는다.

---

## 2. Phase 1 범위

## 2.1 포함 기능

- 알고리즘 풀이 글 입력 폼
- Markdown 미리보기
- 기본 알고리즘 풀이 템플릿 제공
- 글 초안 저장
- Tistory 기본 발행 옵션 저장
- Tistory 비공개 글 발행
- 작성 이력 목록 조회
- 작성 이력 상세 조회
- Docker Compose 기반 로컬 실행

## 2.2 제외 기능

- 사용자 로그인
- 다중 사용자 권한 관리
- AI 기반 글 생성
- GitHub Repository 분석
- Notion/MCP 연동
- Chrome Extension
- Redis/Queue
- 파일 업로드
- 공개 배포용 인증/보안 강화

---

## 3. Phase 1 기술 스택

## 3.1 Frontend

추천:

```text
React + Vite + TypeScript
```

이유:

- Next.js보다 가볍다.
- 로컬 개발 서버가 빠르다.
- 입력 폼과 미리보기 중심 MVP에 충분하다.
- 저사양 컴퓨터에서 부담이 적다.

권장 라이브러리:

```text
react-router-dom
react-hook-form
zod
@tanstack/react-query
react-markdown
remark-gfm
codemirror 또는 textarea 기반 코드 입력
tailwindcss
```

저사양 환경에서는 Monaco Editor 대신 CodeMirror 또는 단순 textarea를 우선 사용한다.

Phase 1 추천:

```text
초기 구현: textarea
추후 개선: CodeMirror
```

## 3.2 Backend

추천:

```text
Go + Chi + pgx
```

이유:

- Spring Boot보다 메모리 사용량이 낮다.
- Docker Compose 로컬 환경에서 가볍다.
- Selenium 기반 Tistory 글쓰기 자동화, 템플릿 렌더링, DB CRUD 중심의 Phase 1에 적합하다.
- 이후 Publisher Service, Connector Service를 Go로 분리하기 쉽다.

대안:

```text
Spring Boot 3 + Java 21 + Spring Data JPA
```

Spring Boot를 쓰고 싶다면 Phase 1부터 가능하지만, 저사양 환경에서는 메모리 제한이 필요하다.

이 문서의 Phase 1 기본 설계는 **Go Backend** 기준이다.

## 3.3 Database

추천:

```text
PostgreSQL 16 Alpine
```

이유:

- JSONB로 입력 원본 저장이 쉽다.
- 템플릿, 글 이력, 설정 데이터를 안정적으로 관리할 수 있다.
- Docker Compose에서 쉽게 구성 가능하다.

---

## 4. Phase 1 로컬 아키텍처

```text
[Browser]
   ↓
[frontend: React + Vite]
   ↓ HTTP /api
[backend: Go API Server]
   ├─ template rendering
   ├─ post management
   ├─ settings management
   ├─ Tistory Selenium publisher
   ├─ [postgres]
   └─ [selenium-chrome] → [Tistory 글쓰기 화면]
```

컨테이너:

```text
frontend
backend
postgres
selenium-chrome
```

선택 컨테이너:

```text
adminer
```

---

## 5. Docker Compose 설계

## 5.1 서비스 목록

| 서비스 | 역할 | 포트 | 필수 여부 |
|---|---|---:|---|
| frontend | React 개발 서버 | 5173 | 필수 |
| backend | Go API 서버 | 8080 | 필수 |
| postgres | 데이터베이스 | 5432 | 필수 |
| selenium-chrome | Tistory 글쓰기 브라우저 자동화 | 4444 | 필수 |
| adminer | DB 확인 도구 | 8081 | 선택 |

## 5.2 저사양 리소스 기준

| 서비스 | CPU 제한 | 메모리 제한 | 비고 |
|---|---:|---:|---|
| frontend | 0.25 | 256MB | Vite 개발 서버 |
| backend | 0.50 | 256MB | Go API 서버 |
| postgres | 0.50 | 384MB | 개발용 DB |
| selenium-chrome | 0.50 | 512MB | Chrome WebDriver |
| adminer | 0.10 | 128MB | 선택 |

총 권장 사용량:

```text
CPU: 약 1.75 core 이하
RAM: 약 1.4GB ~ 1.6GB 내외
```

adminer를 끄면 더 가볍게 실행 가능하다. Selenium Chrome은 Tistory Open API 종료로 인해 Phase 1 발행 기능에 필요하지만, 발행 기능을 테스트하지 않는 개발 중에는 profile로 분리하거나 일시적으로 끌 수 있다.

## 5.3 docker-compose.yml 예시

```yaml
services:
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.dev
    ports:
      - "5173:5173"
    environment:
      - VITE_API_BASE_URL=http://localhost:8080
    volumes:
      - ./frontend:/app
      - frontend_node_modules:/app/node_modules
    depends_on:
      - backend
    deploy:
      resources:
        limits:
          cpus: "0.25"
          memory: 256M

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=local
      - SERVER_PORT=8080
      - DATABASE_URL=postgres://devlog:devlog@postgres:5432/devlog?sslmode=disable
      - SELENIUM_REMOTE_URL=http://selenium-chrome:4444/wd/hub
    volumes:
      - ./backend:/app
    depends_on:
      postgres:
        condition: service_healthy
      selenium-chrome:
        condition: service_started
    deploy:
      resources:
        limits:
          cpus: "0.50"
          memory: 256M

  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_DB=devlog
      - POSTGRES_USER=devlog
      - POSTGRES_PASSWORD=devlog
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./infra/db/init:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U devlog -d devlog"]
      interval: 5s
      timeout: 3s
      retries: 10
    deploy:
      resources:
        limits:
          cpus: "0.50"
          memory: 384M

  selenium-chrome:
    image: seleniarm/standalone-chromium:latest
    ports:
      - "4444:4444"
    shm_size: "2gb"
    deploy:
      resources:
        limits:
          cpus: "0.50"
          memory: 512M

  adminer:
    image: adminer:latest
    ports:
      - "8081:8080"
    depends_on:
      - postgres
    profiles:
      - tools
    deploy:
      resources:
        limits:
          cpus: "0.10"
          memory: 128M

volumes:
  postgres_data:
  frontend_node_modules:
```

주의:

- `deploy.resources`는 Docker Swarm이 아니면 일부 환경에서 강제되지 않을 수 있다.
- 로컬 Docker Desktop에서는 Compose 버전과 설정에 따라 제한 동작이 다를 수 있다.
- 그래도 문서화와 운영 기준으로 남겨둔다.

---

## 6. 레포지토리 구조

Monorepo를 추천한다.

```text
devlog-studio/
├─ docker-compose.yml
├─ .env.example
├─ README.md
├─ frontend/
│  ├─ Dockerfile.dev
│  ├─ package.json
│  ├─ vite.config.ts
│  └─ src/
│     ├─ app/
│     ├─ pages/
│     ├─ features/
│     │  ├─ posts/
│     │  ├─ preview/
│     │  └─ settings/
│     ├─ shared/
│     │  ├─ api/
│     │  ├─ components/
│     │  └─ utils/
│     └─ main.tsx
├─ backend/
│  ├─ Dockerfile.dev
│  ├─ go.mod
│  ├─ cmd/
│  │  └─ server/
│  │     └─ main.go
│  └─ internal/
│     ├─ http/
│     ├─ post/
│     ├─ template/
│     ├─ publisher/
│     │  └─ tistory/
│     ├─ settings/
│     ├─ db/
│     └─ config/
└─ infra/
   └─ db/
      └─ init/
         └─ 001_init.sql
```

`.env.example`에는 서버 포트, DB URL, Selenium URL, CORS origin처럼 비밀이 아닌 로컬 실행 설정만 둔다. Tistory/Kakao 로그인 ID/PW는 웹 발행 모달에서 입력받으므로 `.env.example`에 포함하지 않는다.

---

## 7. Backend 패키지 설계

## 7.1 cmd/server

애플리케이션 진입점.

역할:

- config 로드
- DB 연결
- router 생성
- HTTP server 실행

## 7.2 internal/http

역할:

- 라우터 구성
- 공통 미들웨어
- 에러 응답 포맷
- CORS 설정

## 7.3 internal/post

역할:

- 글 미리보기 생성
- 글 저장
- 글 목록 조회
- 글 상세 조회
- 글 삭제
- 발행 상태 업데이트

구성:

```text
handler.go
service.go
repository.go
model.go
dto.go
```

## 7.4 internal/template

역할:

- 글 유형별 템플릿 로드
- 입력값을 Markdown으로 렌더링
- 기본 알고리즘 풀이 템플릿 제공

Phase 1에서는 DB 템플릿 + 기본 fallback 템플릿을 지원한다.

## 7.5 internal/publisher/tistory

역할:

- Selenium Remote WebDriver 연결
- Tistory 로그인 화면 자동 입력
- 글쓰기 화면에서 마크다운 모드 선택
- 제목, Markdown 본문, 카테고리, 태그, 비공개 옵션 입력
- 저장/발행 결과 URL 파싱
- 실패 처리

Phase 2에서 별도 Publisher Service로 분리 가능하게 인터페이스를 둔다.

```go
type Publisher interface {
    Publish(ctx context.Context, req PublishRequest) (*PublishResult, error)
}
```

## 7.6 internal/settings

역할:

- blogUrl, categoryId(카테고리 목록에서 선택), visibility, defaultTags 저장/조회
- 발행 시 프론트에서 전달된 Tistory 로그인 정보 검증

Phase 1에서는 단일 사용자 key-value 방식.

Tistory 로그인 ID/PW는 DB에 저장하지 않는다. 발행 확인 모달에서 1회 입력받아 publish API 요청에 포함하고, 백엔드는 요청 처리 중 메모리에서만 사용한다. 로그와 API 응답에도 노출하지 않는다.

## 7.7 internal/db

역할:

- PostgreSQL 연결
- Query helper
- Migration 실행 여부는 선택

Phase 1에서는 `infra/db/init/001_init.sql`로 초기화한다.

---

## 8. Frontend 화면 설계

## 8.1 Dashboard

경로:

```text
/
/posts
```

기능:

- 작성 이력 목록 조회
- 글 제목, 유형, 상태, 생성일, 외부 URL 표시
- 새 글 만들기 버튼

목록 컬럼:

| 컬럼 | 설명 |
|---|---|
| 제목 | 글 제목 |
| 유형 | ALGORITHM |
| 플랫폼 | Programmers, Baekjoon 등 |
| 상태 | DRAFT, PUBLISHED, FAILED |
| 생성일 | 작성일 |
| 링크 | Tistory URL |

## 8.2 New Algorithm Post

경로:

```text
/posts/new/algorithm
```

입력 필드:

| 필드 | 타입 | 필수 |
|---|---|---|
| platform | select | Y |
| problemTitle | text | Y |
| problemUrl | url | Y |
| language | select/text | Y |
| categories | tag input | N |
| approach | textarea | Y |
| code | textarea | Y |
| timeComplexity | text | N |
| spaceComplexity | text | N |
| review | textarea | N |
| tags | tag input | N |

버튼:

```text
[미리보기 생성]
[초안 저장]
[티스토리 비공개 발행]
```

## 8.3 Preview Panel

화면 오른쪽 또는 하단에 Markdown 미리보기를 제공한다.

저사양 환경을 고려해 실시간 미리보기보다는 버튼 클릭 시 미리보기 갱신을 우선한다.

```text
입력 폼 변경마다 렌더링 X
[미리보기 생성] 클릭 시 렌더링 O
```

## 8.4 Settings

경로:

```text
/settings
```

필드:

| 필드 | 설명 |
|---|---|
| Blog URL | 기본값 `https://hyeonyway.tistory.com` |
| Category ID | 기본 카테고리 |
| Default Visibility | 기본값 0, 비공개 |
| Default Tags | 기본 태그 |

보안상 Tistory ID/PW는 설정 화면에서 저장하지 않는다. 발행 버튼을 눌렀을 때 확인 모달에서 카카오 로그인 ID/PW를 1회 입력받고, 요청 완료 후 폐기한다.

---

## 9. Database Schema

## 9.1 001_init.sql

```sql
CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    source_platform VARCHAR(50),
    source_url TEXT,
    language VARCHAR(50),
    content_markdown TEXT NOT NULL,
    content_html TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT',
    external_provider VARCHAR(50),
    external_post_id VARCHAR(100),
    external_url TEXT,
    error_message TEXT,
    published_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS post_inputs (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    input_json JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS templates (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    title_format TEXT NOT NULL,
    body_format TEXT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (type, name)
);

CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    encrypted BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_posts_type ON posts(type);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
```

settings 저장 예시:

```text
tistory.blog_url
tistory.category_id
tistory.default_visibility
tistory.default_tags
```

Tistory 로그인 ID/PW는 settings 테이블에 저장하지 않는다.

## 9.2 기본 템플릿 seed

```sql
INSERT INTO templates (type, name, title_format, body_format, is_default)
VALUES (
  'ALGORITHM',
  '기본 알고리즘 풀이 템플릿',
  '[{{platform}}] {{problemTitle}} - {{language}}',
  '# [{{platform}}] {{problemTitle}} - {{language}}

## 문제 링크

{{problemUrl}}

## 문제 유형

{{categories}}

## 접근

{{approach}}

## 풀이 코드

```{{language}}
{{code}}
```

## 복잡도

- 시간복잡도: {{timeComplexity}}
- 공간복잡도: {{spaceComplexity}}

## 회고

{{review}}
',
  TRUE
)
ON CONFLICT DO NOTHING;
```

주의:

PostgreSQL 문자열에 백틱과 줄바꿈이 포함되므로 실제 구현에서는 dollar-quoted string을 사용하는 편이 안전하다.

---

## 10. API 상세 설계

## 10.1 공통 응답 형식

성공:

```json
{
  "data": {},
  "error": null
}
```

실패:

```json
{
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "problemTitle is required"
  }
}
```

---

## 10.2 Preview API

```http
POST /api/posts/preview
```

Request:

```json
{
  "type": "ALGORITHM",
  "input": {
    "platform": "Programmers",
    "problemTitle": "등굣길",
    "problemUrl": "https://school.programmers.co.kr/...",
    "language": "Java",
    "categories": ["DP"],
    "approach": "오른쪽과 아래쪽으로만 이동할 수 있으므로...",
    "code": "class Solution { ... }",
    "timeComplexity": "O(NM)",
    "spaceComplexity": "O(NM)",
    "review": "처음에는 BFS로 접근했지만 DP가 더 적합했다.",
    "tags": ["알고리즘", "프로그래머스", "Java"]
  }
}
```

Response:

```json
{
  "data": {
    "title": "[Programmers] 등굣길 - Java",
    "contentMarkdown": "# [Programmers] 등굣길 - Java\n..."
  },
  "error": null
}
```

---

## 10.3 Create Post API

```http
POST /api/posts
```

역할:

- 입력값 저장
- 템플릿 렌더링 결과 저장
- status = DRAFT

Response:

```json
{
  "data": {
    "id": 1,
    "title": "[Programmers] 등굣길 - Java",
    "status": "DRAFT"
  },
  "error": null
}
```

---

## 10.4 List Posts API

```http
GET /api/posts?page=1&size=20
```

Response:

```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "title": "[Programmers] 등굣길 - Java",
        "type": "ALGORITHM",
        "sourcePlatform": "Programmers",
        "status": "PUBLISHED",
        "externalUrl": "https://example.tistory.com/1",
        "createdAt": "2026-05-13T10:00:00"
      }
    ],
    "page": 1,
    "size": 20,
    "total": 1
  },
  "error": null
}
```

---

## 10.5 Get Post API

```http
GET /api/posts/{id}
```

---

## 10.6 Publish to Tistory Endpoint

```http
POST /api/posts/{id}/publish/tistory
```

Request:

```json
{
  "visibility": 0,
  "categoryId": "1550884",
  "tags": ["알고리즘", "프로그래머스", "Java", "DP"],
  "login": {
    "id": "kakao-account@example.com",
    "password": "input-at-publish-time"
  }
}
```

처리 흐름:

1. post 조회
2. settings에서 blogUrl, 선택된 categoryId, visibility, defaultTags 조회
3. 요청 body에서 Tistory 로그인 ID/PW 확인
4. Selenium Remote WebDriver 세션 생성
5. Tistory 카카오 로그인 화면에서 ID/PW 입력
6. `https://hyeonyway.tistory.com/manage/newpost` 글쓰기 화면 진입
7. `#category-btn` 버튼을 눌러 `#category-list [role='option']` 목록 로드
8. settings에 저장된 categoryId에 해당하는 항목 선택
9. 마크다운 모드 선택
10. 제목, Markdown 본문, 태그, 비공개 옵션 입력
11. 저장 또는 발행 버튼 클릭
12. 성공 시 posts.external_provider = `TISTORY`
13. posts.external_url, posts.external_post_id, posts.published_at 업데이트
14. posts.status = `PUBLISHED`
15. 결과 반환

실패 시:

- posts.status = `FAILED`
- posts.error_message에 실패 사유 저장
- 브라우저 자동화 실패 원인은 로그에 남기되 ID/PW는 마스킹한다.
- 요청으로 받은 login.password는 저장하지 않고 응답에도 포함하지 않는다.

---

## 10.7 Tistory Categories Fetch API

```http
POST /api/tistory/categories/fetch
```

Request:

```json
{
  "blogUrl": "https://hyeonyway.tistory.com",
  "login": {
    "id": "kakao-account@example.com",
    "password": "input-at-fetch-time"
  }
}
```

Response:

```json
{
  "items": [
    { "categoryId": "0", "label": "카테고리 없음" },
    { "categoryId": "1550884", "label": "알고리즘 문제 풀이" },
    { "categoryId": "1550885", "label": "취준" }
  ]
}
```

Selenium은 `https://hyeonyway.tistory.com/manage/newpost`에 진입한 뒤 `#category-btn`을 누르고 `#category-list [role='option']` DOM을 읽어서 카테고리 목록을 반환한다.

## 10.8 Settings API

```http
GET /api/settings
PUT /api/settings
```

PUT Request:

```json
{
  "tistory": {
    "blogUrl": "https://hyeonyway.tistory.com",
    "categoryId": "1550884",
    "defaultVisibility": 0,
    "defaultTags": ["알고리즘", "Java"]
  }
}
```

GET Response:

```json
{
  "data": {
    "tistory": {
      "blogUrl": "https://hyeonyway.tistory.com",
      "categoryId": "1550884",
      "defaultVisibility": 0,
      "defaultTags": ["알고리즘", "Java"]
    }
  },
  "error": null
}
```

---

## 11. Tistory 발행 처리

Tistory Open API는 종료되었으므로 Phase 1에서는 공식 API 호출이 아니라 Selenium 기반 브라우저 자동화로 발행한다.

참고:

```text
https://notice.tistory.com/2664
https://github.com/tistory/document-tistory-apis
```

## 11.1 발행 기본 정책

Phase 1에서는 무조건 비공개 발행을 기본값으로 둔다.

```text
visibility = 0
```

사용자가 공개 발행을 선택할 수 있게 하더라도 기본값은 비공개다.

## 11.2 Tistory Selenium 입력 필드

Selenium 발행에 필요한 주요 값:

```text
blogUrl
login.id
login.password
title
contentMarkdown
visibility
categoryId
tags
```

`categoryId`는 Settings 화면에서 `POST /api/tistory/categories/fetch`로 불러온 목록 중에서 선택해 저장한다.

## 11.3 Content 형식

Phase 1의 콘텐츠 기준은 Markdown이다.

```text
DB 저장: Markdown
프론트 미리보기: Markdown
Tistory 입력: 글쓰기 화면을 마크다운 모드로 전환한 뒤 Markdown 본문 붙여넣기
```

HTML 변환은 Phase 1 범위에 넣지 않는다. Tistory 에디터의 마크다운 모드 선택이 실패하면 발행 실패로 처리하고 사용자에게 재시도 또는 수동 확인을 안내한다.

## 11.4 로그인 방식

Phase 1에서는 카카오 계정 로그인을 기본 전제로 한다. Tistory가 카카오 계정으로 연결되어 있고, 해당 PC/브라우저 환경에서 카카오 추가 인증이 완료되어 있어야 Selenium 발행이 안정적으로 동작한다.

카카오 로그인 정보는 웹 사이트에서 발행 시점에만 입력받는다.

```text
사용자가 [티스토리 비공개 발행] 클릭
        ↓
발행 확인 모달 표시
        ↓
카카오 로그인 ID/PW 입력
        ↓
POST /api/posts/{id}/publish/tistory 요청 body에 login 포함
        ↓
백엔드는 요청 처리 중 메모리에서만 사용
        ↓
요청 완료 후 폐기
```

Settings 화면에서는 별도로 카테고리 불러오기 모달을 열어 `#category-btn` / `#category-list [role='option']` DOM을 읽고, 원하는 항목의 `category-id`를 저장한다.

주의:

- ID/PW는 DB에 저장하지 않는다.
- ID/PW는 설정 API 응답에 포함하지 않는다.
- ID/PW는 로그에 출력하지 않는다.
- 프론트는 ID/PW를 localStorage, sessionStorage, IndexedDB에 저장하지 않는다.
- 비밀번호 input에는 autocomplete를 끄고, 발행 요청 완료 후 폼 상태를 초기화한다.
- 카카오 2단계 인증, CAPTCHA, 추가 본인확인이 발생하면 자동 발행은 실패 처리한다.
- 새 PC, 새 브라우저, 새 컨테이너 환경에서는 카카오 인증이 다시 필요할 수 있다.
- 이 프로젝트는 로컬 개인 자동화 도구이므로, README에 "인증이 완료된 본인 PC에서만 Selenium 발행을 테스트한다"는 안내를 반드시 적는다.
- Selenium 자동화는 Tistory 화면 구조 변경에 취약하므로 선택자 변경에 대비해 publisher 모듈 내부에 선택자를 모아둔다.

---

## 12. 구현 순서

## Step 1. 프로젝트 스캐폴딩

- monorepo 생성
- docker-compose.yml 작성
- frontend Vite 프로젝트 생성
- backend Go 프로젝트 생성
- postgres init sql 추가

## Step 2. Backend 기본 API

- health check
- DB 연결
- 공통 응답 포맷
- CORS 설정

API:

```http
GET /health
```

## Step 3. DB Schema 적용

- posts
- post_inputs
- templates
- settings
- 기본 템플릿 seed

## Step 4. Template Preview 구현

- POST /api/posts/preview
- 기본 알고리즘 템플릿 렌더링
- 필수 값 validation

## Step 5. Post 저장 구현

- POST /api/posts
- GET /api/posts
- GET /api/posts/{id}
- DELETE /api/posts/{id}

## Step 6. Settings 구현

- GET /api/settings
- PUT /api/settings
- categoryId(카테고리 목록에서 선택), defaultVisibility, defaultTags 저장

## Step 7. Tistory Selenium Publisher 구현

- internal/publisher/tistory 구현
- POST /api/posts/{id}/publish/tistory
- Selenium Remote WebDriver 연동
- 요청 body로 받은 카카오 로그인 정보 사용
- Tistory 로그인, 카테고리 목록 로드, categoryId 선택, 마크다운 모드 전환, 제목/본문/발행 옵션 입력
- 성공/실패 상태 업데이트

## Step 8. Frontend 구현

- Dashboard
- New Algorithm Post
- Preview Panel
- Settings
- API client

## Step 9. 통합 테스트

- Docker Compose로 전체 실행
- 설정 저장
- 글 미리보기
- 글 저장
- Tistory 비공개 발행
- 작성 이력 확인

---

## 13. 프론트 구현 상세

## 13.1 페이지 구조

```text
src/pages/DashboardPage.tsx
src/pages/NewAlgorithmPostPage.tsx
src/pages/PostDetailPage.tsx
src/pages/SettingsPage.tsx
```

## 13.2 feature 구조

```text
src/features/posts/
  api.ts
  types.ts
  components/PostList.tsx
  components/AlgorithmPostForm.tsx
  components/PostPreview.tsx

src/features/settings/
  api.ts
  types.ts
  components/TistorySettingsForm.tsx
```

## 13.3 UX 원칙

- 입력 폼과 미리보기를 한 화면에서 본다.
- 자동 발행하지 않는다.
- `미리보기 생성`과 `초안 저장`을 분리한다.
- Tistory 발행 전 확인 모달을 띄우고 카카오 로그인 ID/PW를 1회 입력받는다.
- Settings 화면에서 Tistory 글쓰기 화면의 카테고리 목록을 불러와 선택한다.
- 카카오 로그인 ID/PW는 브라우저 저장소에 저장하지 않고 발행 요청 후 즉시 폼에서 제거한다.
- 발행 성공 시 URL을 바로 보여준다.

---

## 14. Backend 구현 상세

## 14.1 환경 변수

```text
APP_ENV=local
SERVER_PORT=8080
DATABASE_URL=postgres://devlog:devlog@postgres:5432/devlog?sslmode=disable
SELENIUM_REMOTE_URL=http://selenium-chrome:4444/wd/hub
CORS_ALLOWED_ORIGINS=http://localhost:5173
```

## 14.2 Health Check

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

## 14.3 Validation

필수 입력값:

```text
type
platform
problemTitle
problemUrl
language
approach
code
```

선택 입력값:

```text
categories
timeComplexity
spaceComplexity
review
tags
```

## 14.4 Template Renderer

처음에는 단순 문자열 치환으로 충분하다.

예:

```text
{{problemTitle}} → 등굣길
```

주의:

- 코드 블록 안의 특수문자가 깨지지 않도록 한다.
- 비어 있는 선택값은 `-` 또는 빈 문자열로 처리한다.
- categories, tags 배열은 `- DP\n- BFS` 같은 Markdown list로 변환한다.

---

## 15. 테스트 전략

## 15.1 Backend Unit Test

대상:

- Template Renderer
- Post Service
- Settings Service
- Tistory Selenium Publisher 입력값 검증

## 15.2 Backend Integration Test

대상:

- DB 저장/조회
- Preview API
- Create Post API

## 15.3 Frontend Test

Phase 1에서는 필수는 아니다. 단, 다음 정도는 수동 테스트한다.

- 입력 폼 validation
- 미리보기 표시
- 설정 저장
- 발행 결과 표시

## 15.4 수동 E2E 시나리오

1. Docker Compose 실행
2. `/settings`에서 Tistory 기본 발행 옵션 저장
3. `/posts/new/algorithm`에서 글 입력
4. 미리보기 생성
5. 초안 저장
6. 발행 확인 모달에서 카카오 로그인 ID/PW 입력
7. Selenium으로 Tistory 마크다운 모드 비공개 발행
8. Dashboard에서 작성 이력 확인
9. Tistory에서 비공개 글 확인

---

## 16. 에러 처리

| 에러 코드 | 상황 | 사용자 메시지 |
|---|---|---|
| VALIDATION_ERROR | 필수 입력 누락 | 필수 입력값을 확인해주세요. |
| SETTINGS_NOT_FOUND | Tistory 기본 설정 없음 | Tistory 설정을 먼저 등록해주세요. |
| TISTORY_LOGIN_REQUIRED | 발행 요청에 로그인 정보 없음 | Tistory 로그인 정보를 입력해주세요. |
| TISTORY_LOGIN_FAILED | 로그인 실패 또는 추가 인증 필요 | Tistory 로그인 상태를 확인해주세요. |
| TISTORY_EDITOR_ERROR | 글쓰기 화면 또는 마크다운 모드 제어 실패 | Tistory 글쓰기 화면을 확인해주세요. |
| TISTORY_PUBLISH_ERROR | 저장/발행 실패 | Tistory 발행 중 오류가 발생했습니다. |
| POST_NOT_FOUND | 글 없음 | 글을 찾을 수 없습니다. |
| DATABASE_ERROR | DB 오류 | 데이터 처리 중 오류가 발생했습니다. |

---

## 17. 보안 및 주의사항

## 17.1 Tistory 로그인 정보

Phase 1은 로컬 단일 사용자 기준이지만 Tistory ID/PW는 DB에 저장하지 않는다. 웹 사이트에서 발행 시점에만 입력받고 다음을 지킨다.

- DB volume을 외부에 공유하지 않는다.
- 발행 확인 모달 외 화면에는 ID/PW를 보여주지 않는다.
- 프론트는 ID/PW를 localStorage, sessionStorage, IndexedDB에 저장하지 않는다.
- 백엔드는 요청 처리 중 메모리에서만 ID/PW를 사용하고 DB에 저장하지 않는다.
- 로그에 ID/PW를 출력하지 않는다.
- 실패 로그와 스크린샷에 로그인 정보가 포함되지 않도록 주의한다.

## 17.2 문제 본문 복사 금지

알고리즘 문제 본문 전문을 서비스가 자동으로 복사하지 않는다.

저장 대상:

```text
문제 제목
문제 링크
내 풀이 설명
내 코드
내 회고
```

## 17.3 기본 비공개 발행

실수로 공개 발행되는 것을 막기 위해 기본값은 비공개다.

---

## 18. Codex 구현 지침

Codex는 다음 순서로 구현한다.

```text
1. 전체 레포 구조 생성
2. docker-compose.yml 작성
3. PostgreSQL init SQL 작성
4. Go backend health check 구현
5. DB 연결 구현
6. Preview API 구현
7. Posts CRUD 구현
8. Settings API 구현
9. Tistory Selenium Publisher 구현
10. React frontend 구현
11. 통합 실행 README 작성
```

README에는 다음 내용을 반드시 포함한다.

```text
- Tistory Open API 종료로 인해 Selenium 기반 글쓰기 자동화를 사용한다.
- Tistory 로그인은 카카오 계정 로그인을 기준으로 한다.
- Selenium 발행은 카카오 추가 인증이 완료된 본인 PC/로컬 환경에서만 안정적으로 동작한다.
- 새 PC, 새 브라우저, 새 컨테이너 환경에서는 카카오 2단계 인증, CAPTCHA, 본인확인이 발생할 수 있으며 이 경우 자동 발행은 실패할 수 있다.
- 카카오 로그인 ID/PW는 발행 시점에 웹 화면에서 1회 입력받고 저장하지 않는다.
- 프론트는 ID/PW를 브라우저 저장소에 저장하지 않고, 백엔드는 DB/로그에 남기지 않는다.
- 자동 발행은 기본적으로 비공개 발행으로 동작한다.
```

병렬 처리가 가능한 경우 subagent를 사용한다.

예:

```text
- backend-api subagent: Go API, DB, Tistory Selenium publisher 구현
- frontend-ui subagent: React 화면, 폼, 미리보기 구현
- infra subagent: Docker Compose, DB init, README 작성
```

Codex 환경에 subagent 템플릿이 있다면 다음 경로를 참고한다.

```text
~/.codex/agents
```

각 subagent는 서로 충돌하지 않도록 작업 디렉터리와 책임 범위를 분리한다.

권장 분리:

```text
frontend-ui → frontend/**
backend-api → backend/**
infra-docs → docker-compose.yml, infra/**, README.md
```

---

## 20. 완료 기준

Phase 1 완료 기준은 다음과 같다.

- `docker compose up`으로 전체 실행 가능
- 프론트에서 알고리즘 풀이 입력 가능
- 미리보기 생성 가능
- 글 초안 저장 가능
- 작성 이력 조회 가능
- Tistory 기본 발행 옵션 저장 가능
- Selenium으로 Tistory 마크다운 모드 비공개 발행 가능
- 발행 성공 시 Tistory URL 표시
- README에 실행 방법, 카카오 로그인 입력 방식, 인증 PC 제약 설명 포함

---

## 21. Phase 1 이후 확장 포인트

Phase 1 완료 후 다음 순서로 확장한다.

1. 프로젝트 회고 템플릿 추가
2. 트러블슈팅 템플릿 추가
3. Publisher Service를 Go 별도 서비스로 분리
4. AI Serving Service를 FastAPI로 추가
5. GitHub Repository Connector 추가
6. Notion/MCP Connector 추가
7. 사용자 로그인 및 OAuth 기반 토큰 관리 추가
