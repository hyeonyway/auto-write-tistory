# DevLog Studio 전체 설계 문서

## 1. 프로젝트 개요

### 1.1 서비스명

**DevLog Studio**

### 1.2 한 줄 설명

개발자가 코드, 문제 링크, 프로젝트 자료, GitHub 저장소, Notion 문서 등을 입력하면 기술 블로그 형식의 글 초안을 생성하고, Tistory 등 외부 플랫폼에 발행할 수 있도록 돕는 템플릿 기반 개발자 콘텐츠 발행 서비스.

### 1.3 서비스 목표

초기 목표는 사용자가 알고리즘 풀이 내용을 직접 입력하면, 정해진 템플릿에 맞춰 Tistory 블로그 초안을 생성하는 것이다. 이후 GitHub Repository 분석, AI 기반 설명 생성, Notion/MCP 연동, 프로젝트 회고 글 자동 생성 등으로 확장한다.

### 1.4 기존 Chrome Extension 방식과의 차이

초기 아이디어는 프로그래머스 페이지에서 문제 정보와 코드를 직접 추출하는 Chrome Extension이었다. 그러나 이 방식은 특정 사이트 DOM 구조에 종속되고, 백준/소프티어/LeetCode/GitHub/Notion 등으로 확장하기 어렵다.

DevLog Studio는 웹 서비스 중심으로 설계한다.

```text
사용자 입력 / 외부 소스 연동
        ↓
글 유형 선택
        ↓
템플릿 기반 초안 생성
        ↓
미리보기 및 수정
        ↓
Tistory / Notion / GitHub / Markdown Export 발행
```

초기에는 수동 입력 기반으로 시작하고, 이후 각 입력 소스를 독립 서비스로 추가한다.

---

## 2. 핵심 사용자 시나리오

### 2.1 Phase 1 시나리오

사용자가 알고리즘 문제 풀이를 블로그 글로 정리하고 싶다.

1. 사용자가 웹 서비스에 접속한다.
2. `새 글 만들기`를 선택한다.
3. 글 유형으로 `알고리즘 풀이`를 선택한다.
4. 문제 플랫폼, 문제 제목, 문제 링크, 언어, 코드, 접근 방식, 복잡도, 회고를 입력한다.
5. 서비스가 Tistory용 Markdown 초안을 생성한다.
6. 사용자가 미리보기를 확인하고 수정한다.
7. 사용자가 `Tistory 비공개 초안 생성`을 누른다.
8. 서비스가 Selenium으로 Tistory 글쓰기 화면을 열고 마크다운 모드에서 제목과 본문을 입력한 뒤 비공개 글을 생성한다.
9. 작성 이력에 생성된 글과 Tistory URL이 저장된다.

### 2.2 확장 시나리오

#### GitHub Repository 기반 프로젝트 글 생성

1. 사용자가 GitHub Repository URL을 입력한다.
2. GitHub Collector Service가 README, 주요 설정 파일, 디렉터리 구조, 사용 언어를 수집한다.
3. Content Service가 프로젝트 회고 템플릿을 선택한다.
4. AI Serving Service가 프로젝트 개요, 기술 스택, 핵심 기능, 트러블슈팅 초안을 생성한다.
5. 사용자가 내용을 수정하고 Tistory에 발행한다.

#### Notion/MCP 기반 문서 글 생성

1. 사용자가 Notion 페이지 또는 MCP 연결 문서를 선택한다.
2. Notion Connector Service가 문서 내용을 가져온다.
3. AI Serving Service가 기술 블로그 형식으로 재구성한다.
4. 사용자가 미리보기 후 Tistory 또는 Notion에 재발행한다.

---

## 3. 설계 방향

### 3.1 서비스 분리 원칙

로컬 Docker Compose 기반으로 실행하지만, 서비스는 처음부터 느슨하게 분리한다. 실제 운영용 MSA까지는 아니더라도, 기능 단위로 나누어 향후 확장과 교체가 쉽도록 한다.

분리 기준은 다음과 같다.

| 영역 | 분리 이유 |
|---|---|
| Frontend | UI와 API 서버 분리 |
| API Gateway / BFF | 프론트 요청 진입점 단순화 |
| Content Service | 글/템플릿/작성 이력 핵심 도메인 |
| Publisher Service | Tistory, Notion, GitHub 등 외부 발행 분리 |
| AI Serving Service | Python/FastAPI 기반 AI 기능 독립 운영 |
| Connector Service | GitHub, Notion, MCP 등 외부 소스별 수집 기능 분리 |
| Database | 서비스 간 데이터 저장소 분리 가능성 확보 |

### 3.2 저사양 로컬 환경 우선

사용자의 컴퓨터 스펙이 높지 않으므로, Phase 1에서는 최소 컨테이너 수와 낮은 리소스 사용량을 우선한다.

Phase 1 기본 실행 컨테이너는 다음 정도로 제한한다.

```text
frontend
api-server
postgres
selenium-chrome
```

선택적으로 추가한다.

```text
publisher-service
adminer 또는 pgadmin
```

AI, GitHub, Notion, Redis, Queue 등은 Phase 1에서는 제외한다. Selenium Chrome은 Tistory Open API 종료에 따른 발행 자동화용 컨테이너로만 사용한다.

---

## 4. 전체 아키텍처

## 4.1 최종 지향 아키텍처

```text
[Browser]
   ↓
[Frontend - Next.js]
   ↓
[API Gateway / BFF - Go]
   ↓
 ┌─────────────────────────────────────────────┐
 │                                             │
 │ [Content Service - Spring Boot]             │
 │  - posts                                    │
 │  - templates                                │
 │  - histories                                │
 │                                             │
 │ [Publisher Service - Go]                    │
 │  - Tistory browser automation               │
 │  - Notion publish                           │
 │  - Markdown export                          │
 │                                             │
 │ [Connector Service - Go]                    │
 │  - GitHub collector                         │
 │  - Notion collector                         │
 │  - MCP bridge                               │
 │                                             │
 │ [AI Serving Service - FastAPI]              │
 │  - summary                                  │
 │  - title suggestion                         │
 │  - code explanation                         │
 │  - project article draft                    │
 │                                             │
 └─────────────────────────────────────────────┘
   ↓
[PostgreSQL]
[Redis - optional]
[Object Storage - optional]
```

## 4.2 Phase 1 경량 아키텍처

Phase 1은 로컬에서 가볍게 돌리는 것이 중요하므로, 다음처럼 줄인다.

```text
[Browser]
   ↓
[Frontend - Next.js]
   ↓
[Backend API - Spring Boot 또는 Go]
   ├─ 글 입력 저장
   ├─ 템플릿 렌더링
   ├─ Selenium 기반 Tistory 발행
   └─ 작성 이력 관리
   ↓
[PostgreSQL]
```

Phase 1에서는 Publisher Service를 별도 애플리케이션으로 분리하지 않아도 된다. 단, 코드 구조는 `publisher` 모듈로 분리해두고, Phase 2에서 독립 서비스로 떼어낼 수 있게 한다. Tistory Open API는 종료되었으므로 Phase 1의 Tistory 발행은 Selenium 기반 브라우저 자동화로 처리한다.

---

## 5. 기술 스택 선정

## 5.1 Frontend 추천

### 추천: Next.js + TypeScript

이유:

- 입력 폼과 미리보기 UI 구현이 편하다.
- Markdown Preview, Monaco Editor, CodeMirror 등과 궁합이 좋다.
- 추후 서버 사이드 렌더링이 필요하지 않아도 React 기반 SPA처럼 사용할 수 있다.
- 포트폴리오에서 프론트 구조를 설명하기 좋다.

대안:

| 스택 | 장점 | 단점 |
|---|---|---|
| React + Vite | 가장 가볍고 빠름 | 라우팅/구조 직접 설계 필요 |
| Next.js | 구조화 좋고 확장성 좋음 | Vite보다 약간 무거움 |
| SvelteKit | 가볍고 생산성 좋음 | 팀/채용 시장 인지도 낮을 수 있음 |

저사양 로컬에서는 **React + Vite**가 가장 가볍다. 포트폴리오와 확장성을 고려하면 **Next.js**가 좋다.

최종 추천:

```text
Phase 1: React + Vite + TypeScript
확장형: Next.js + TypeScript
```

이 문서에서는 저사양 로컬 환경을 고려해 **React + Vite**를 기본값으로 둔다.

## 5.2 Backend 추천

### Phase 1 추천: Spring Boot 또는 Go 중 하나 선택

사용자가 Go와 Spring Boot 중심을 원하므로 두 가지 안을 제시한다.

#### 안 A. Spring Boot 중심

장점:

- 이력서/포트폴리오에서 백엔드 경험으로 설명하기 좋다.
- JPA, Validation, REST API, 외부 API 연동 경험을 보여주기 좋다.
- 향후 인증, 템플릿 관리, 작성 이력 관리에 적합하다.

단점:

- JVM 기반이라 Go보다 메모리 사용량이 크다.
- 저사양 로컬에서는 컨테이너 메모리 제한을 신경 써야 한다.

권장 메모리:

```text
backend: 512MB ~ 768MB
```

#### 안 B. Go 중심

장점:

- 실행 파일이 가볍고 메모리 사용량이 낮다.
- 로컬 Docker Compose에 적합하다.
- API Gateway, Publisher Service, Connector Service에 잘 어울린다.

단점:

- Spring Boot 대비 프로젝트 구조를 직접 잡아야 한다.
- ORM, Validation, 설정 관리 등에서 개발자가 선택해야 할 것이 많다.

권장 메모리:

```text
backend: 128MB ~ 256MB
```

### 최종 추천

저사양 로컬 기준으로는 다음 구성이 가장 적합하다.

```text
Phase 1 Backend: Go
Phase 2 Content Service: Spring Boot로 분리 가능
Phase 2 Publisher/Connector Service: Go 유지
Phase 3 AI Serving Service: FastAPI
```

단, 취업 포트폴리오에서 Spring Boot 경험을 강조하고 싶다면 Phase 1부터 Spring Boot로 가도 된다.

이 문서에서는 다음 혼합 전략을 권장한다.

```text
Frontend: React + Vite + TypeScript
Backend API: Go
Database: PostgreSQL
Future Content Service: Spring Boot
Future AI Service: FastAPI
Future Publisher/Connector Service: Go
```

## 5.3 Database

추천: PostgreSQL

이유:

- JSONB를 활용해 글 입력 원본을 유연하게 저장할 수 있다.
- 작성 이력, 템플릿, 설정 관리에 적합하다.
- MySQL보다 JSON 기반 확장 데이터 저장이 편하다.

대안:

- MySQL: 익숙하고 포트폴리오 친화적이다.
- SQLite: 가장 가볍지만 Docker Compose 기반 다중 서비스 확장에는 약하다.

저사양 환경에서는 PostgreSQL 컨테이너 메모리를 작게 제한한다.

```text
postgres: 256MB ~ 384MB
```

---

## 6. 서비스 구성

## 6.1 Frontend Service

### 역할

- 글 유형 선택
- 입력 폼 제공
- Markdown 미리보기
- Tistory 발행 요청
- 작성 이력 조회
- 설정 화면 제공

### 주요 화면

| 화면 | 설명 |
|---|---|
| Dashboard | 작성 이력 목록 |
| New Post | 글 유형 선택 및 입력 |
| Algorithm Form | 알고리즘 풀이 입력 폼 |
| Preview | Markdown 미리보기 |
| Settings | Tistory 기본 발행 옵션 관리 |

### 권장 라이브러리

```text
React
TypeScript
Vite
React Router
React Hook Form
Zod
TanStack Query
Markdown Preview 라이브러리
CodeMirror 또는 Monaco Editor
Tailwind CSS
```

저사양 환경에서는 Monaco Editor보다 CodeMirror가 더 가볍다.

추천:

```text
CodeMirror 6
```

---

## 6.2 Backend API Service

### Phase 1 역할

- 글 생성 요청 처리
- 템플릿 렌더링
- 글 저장
- Selenium 기반 Tistory 발행
- 작성 이력 조회
- 설정 저장

### Phase 1 기술 스택 추천

```text
Go 1.22+
Gin 또는 Chi
pgx 또는 sqlc
PostgreSQL
Docker
```

Gin은 빠르게 만들기 좋고, Chi는 표준 net/http에 가까워 가볍다.

추천:

```text
Go + Chi + pgx
```

### 주요 모듈

```text
cmd/server
internal/post
internal/template
internal/publisher/tistory
internal/settings
internal/db
internal/http
```

---

## 6.3 Content Service

Phase 2 이후 분리 대상.

### 역할

- posts 도메인 관리
- templates 도메인 관리
- 작성 이력 관리
- 글 타입별 입력 스키마 관리

### 추천 스택

```text
Spring Boot 3
Java 21
Spring Web
Spring Data JPA
PostgreSQL
Flyway
```

분리 시점:

```text
- 글 유형이 3개 이상으로 늘어날 때
- 템플릿 편집 기능이 복잡해질 때
- 사용자/권한 기능이 들어갈 때
```

---

## 6.4 Publisher Service

Phase 2 이후 분리 대상.

### 역할

- Tistory 발행
- Notion 발행
- GitHub Markdown commit
- Markdown 파일 export

### 추천 스택

```text
Go
```

이유:

- 외부 발행 어댑터와 브라우저 자동화 제어를 가볍게 다루기 좋아 Go가 적합하다.
- 서비스별 발행 어댑터를 인터페이스로 분리하기 좋다.

### Adapter 구조

```text
Publisher interface
 ├─ TistoryPublisher
 ├─ NotionPublisher
 ├─ GitHubPublisher
 └─ MarkdownExporter
```

---

## 6.5 Connector Service

Phase 3 이후 분리 대상.

### 역할

- GitHub Repository 분석
- Notion Page 가져오기
- MCP 기반 외부 문서 수집
- Markdown 파일 업로드 분석

### 추천 스택

```text
Go
```

GitHub API, Notion API, HTTP 수집 중심이므로 Go가 적합하다.

---

## 6.6 AI Serving Service

Phase 3 이후 추가.

### 역할

- 코드 설명 생성
- 접근 방식 초안 생성
- 시간복잡도 추론
- 프로젝트 회고 초안 생성
- README 요약
- 제목/태그 추천

### 추천 스택

```text
Python
FastAPI
Pydantic
LangChain 또는 직접 OpenAI SDK 사용
```

AI 기능은 Python 생태계가 편하므로 별도 서비스로 분리한다.

### API 예시

```text
POST /ai/algorithm/explain
POST /ai/project/summarize
POST /ai/title/suggest
POST /ai/tags/suggest
```

---

## 7. 데이터 모델

## 7.1 posts

```sql
CREATE TABLE posts (
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
    external_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### status 값

```text
DRAFT
PUBLISHED
FAILED
```

## 7.2 post_inputs

사용자가 입력한 원본 데이터를 JSON으로 저장한다.

```sql
CREATE TABLE post_inputs (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    input_json JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## 7.3 templates

```sql
CREATE TABLE templates (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    title_format TEXT NOT NULL,
    body_format TEXT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## 7.4 settings

단일 사용자 MVP 기준으로 key-value 설정을 사용한다.

```sql
CREATE TABLE settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    encrypted BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

저장 예시:

```text
tistory.blog_url
tistory.category_id
tistory.default_visibility
tistory.default_tags
```

Tistory 로그인 ID/PW는 DB에 저장하지 않는다. 발행 시점에 웹 화면에서 1회 입력받고, 백엔드는 요청 처리 중 메모리에서만 사용한다.

---

## 8. API 설계

## 8.1 Posts API

### 글 미리보기 생성

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
    "approach": "오른쪽과 아래쪽으로만 이동 가능하므로...",
    "code": "class Solution { ... }",
    "timeComplexity": "O(NM)",
    "spaceComplexity": "O(NM)",
    "review": "처음에는 BFS로 접근했지만 DP가 적합했다."
  }
}
```

Response:

```json
{
  "title": "[프로그래머스] 등굣길 - Java",
  "contentMarkdown": "# [프로그래머스] 등굣길 - Java\n..."
}
```

### 글 저장

```http
POST /api/posts
```

### 글 목록 조회

```http
GET /api/posts
```

### 글 상세 조회

```http
GET /api/posts/{id}
```

### 글 삭제

```http
DELETE /api/posts/{id}
```

---

## 8.2 Publish API

### Tistory 비공개 발행

```http
POST /api/posts/{id}/publish/tistory
```

Request:

```json
{
  "visibility": 0,
  "categoryId": "1550884",
  "tags": ["알고리즘", "프로그래머스", "Java", "DP"]
}
```

`categoryId`는 Settings 화면에서 Selenium으로 Tistory 글쓰기 페이지의 카테고리 목록을 불러온 뒤 선택한 값이다.

Response:

```json
{
  "success": true,
  "externalUrl": "https://example.tistory.com/123"
}
```

---

## 8.3 Category Fetch API

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

## 8.4 Settings API

```http
GET /api/settings
PUT /api/settings
```

민감한 값은 화면에 그대로 노출하지 않는다.

예:

```json
{
  "tistory": {
    "blogUrl": "https://hyeonyway.tistory.com",
    "categoryId": "1550884",
    "defaultVisibility": 0
  }
}
```

---

## 9. 템플릿 설계

## 9.1 템플릿 변수

알고리즘 풀이 템플릿에서 사용할 변수:

```text
{{platform}}
{{problemTitle}}
{{problemUrl}}
{{language}}
{{categories}}
{{approach}}
{{code}}
{{timeComplexity}}
{{spaceComplexity}}
{{review}}
{{tags}}
```

## 9.2 기본 알고리즘 풀이 템플릿

```markdown
# [{{platform}}] {{problemTitle}} - {{language}}

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
```

---

## 10. 로컬 Docker Compose 설계

## 10.1 Phase 1 컨테이너

```text
frontend
backend
postgres
selenium-chrome
```

선택:

```text
adminer
```

## 10.2 권장 리소스 제한

저사양 컴퓨터 기준.

| 서비스 | CPU 제한 | 메모리 제한 |
|---|---:|---:|
| frontend | 0.25 CPU | 256MB |
| backend-go | 0.50 CPU | 256MB |
| backend-spring | 0.75 CPU | 768MB |
| postgres | 0.50 CPU | 384MB |
| selenium-chrome | 0.50 CPU | 512MB |
| adminer | 0.10 CPU | 128MB |
| ai-fastapi | 0.50 CPU | 512MB |

Phase 1에서 Spring Boot를 사용한다면 전체 메모리 사용량이 커질 수 있다. 저사양 환경에서는 Go Backend를 권장한다.

## 10.3 JVM 메모리 제한 예시

Spring Boot를 사용할 경우:

```text
JAVA_TOOL_OPTIONS=-XX:MaxRAMPercentage=70 -XX:InitialRAMPercentage=30
```

또는:

```text
JAVA_OPTS=-Xms128m -Xmx512m
```

## 10.4 PostgreSQL 경량 설정

개발용이므로 다음 환경을 권장한다.

```text
POSTGRES_DB=devlog
POSTGRES_USER=devlog
POSTGRES_PASSWORD=devlog
```

추가 튜닝은 Phase 1에서는 과하지 않게 한다.

---

## 11. 개발 단계 로드맵

## Phase 1. 수동 입력 기반 알고리즘 풀이 글 생성

목표:

```text
수동 입력 → 템플릿 렌더링 → 미리보기 → Selenium 기반 Tistory 비공개 발행 → 작성 이력 저장
```

포함 기능:

- 알고리즘 풀이 입력 폼
- Markdown 미리보기
- 글 저장
- Tistory 기본 발행 옵션 저장
- Selenium 기반 Tistory 비공개 발행
- 작성 이력 조회

제외 기능:

- 로그인
- AI 초안 생성
- GitHub 자동 분석
- Notion/MCP 연동
- 크롬 익스텐션
- Queue/Redis

## Phase 2. 글 유형 확장 및 Publisher 분리

- 프로젝트 회고 템플릿 추가
- 트러블슈팅 템플릿 추가
- Publisher Service 분리
- Markdown Export 추가

## Phase 3. AI Serving Service 추가

- FastAPI 기반 AI Serving Service 추가
- 코드 설명 생성
- 시간복잡도 추론
- 제목/태그 추천
- 프로젝트 설명 초안 생성

## Phase 4. GitHub Connector 추가

- GitHub Repository URL 입력
- README 및 설정 파일 수집
- 기술 스택 분석
- 프로젝트 글 초안 생성

## Phase 5. Notion/MCP Connector 추가

- Notion 문서 가져오기
- MCP 기반 문서 선택
- 회의록/기획서 → 블로그 글 변환

---

## 12. 에러 처리 전략

| 상황 | 처리 |
|---|---|
| Tistory 로그인 정보 미입력 | 발행 확인 모달에서 입력 안내 |
| Tistory Selenium 발행 실패 | 실패 상태 저장, 재시도 버튼 제공 |
| 템플릿 렌더링 실패 | 필수 입력값 누락 메시지 표시 |
| DB 연결 실패 | 서버 시작 실패 로그 명확히 출력 |
| 외부 URL 유효하지 않음 | URL 검증 후 사용자에게 안내 |
| 발행 성공했지만 URL 누락 | 외부 postId 또는 응답 원문 일부 저장 |

---

## 13. 보안 고려사항

### 13.1 Tistory 로그인 정보 저장

Phase 1은 단일 사용자 로컬 서비스이지만 Tistory ID/PW는 DB에 저장하지 않는다. 웹 화면에서 발행 시점에만 1회 입력받고, API 요청 처리 중 메모리에서만 사용한다. 화면 저장소, API 응답, 로그에도 노출하지 않는다.

Phase 1 권장:

```text
사용자가 Tistory 발행 버튼 클릭
→ 발행 확인 모달에서 카카오 로그인 ID/PW 입력
→ 백엔드가 Selenium 로그인에만 사용
→ 요청 완료 후 폐기
```

Tistory 로그인은 카카오 계정 로그인을 기준으로 한다. Selenium 발행은 카카오 추가 인증이 완료된 본인 PC/로컬 환경에서만 안정적으로 동작하며, 새 PC나 새 브라우저 환경에서는 2단계 인증, CAPTCHA, 본인확인이 발생해 자동 발행이 실패할 수 있다. 카카오 로그인 ID/PW를 저장하지 않고 발행 시점에만 입력받는다는 점과 인증 PC 제약은 README에 명확히 안내한다.

Phase 2 이상:

```text
브라우저 프로필/세션 재사용 검토
사용자별 계정 분리 시 서버 측 암호화 검토
추가 인증, CAPTCHA, 로그인 실패 처리 강화
```

### 13.2 기본 발행 정책

초기에는 반드시 비공개로 발행한다.

```text
Tistory visibility = 0
```

### 13.3 문제 본문 복사 방지

알고리즘 문제 본문 전문을 자동으로 복사하지 않는다. 문제 링크만 기록한다.

---

## 14. 포트폴리오 표현

### Phase 1 경험 문장

```text
개발자 기술 글 작성 자동화를 위한 템플릿 기반 콘텐츠 발행 서비스를 개발했습니다.
알고리즘 풀이 정보를 입력하면 Markdown 초안을 생성하고, Selenium 기반 브라우저 자동화로 Tistory 비공개 글 발행을 지원하는 기능을 구현했습니다.
Go 기반 API 서버와 PostgreSQL을 Docker Compose 환경에서 구성하고, 글 작성 이력과 사용자 설정을 관리했습니다.
```

### 확장 후 경험 문장

```text
GitHub Repository, Notion 문서 등 다양한 개발 자료를 기술 블로그 글로 변환하는 콘텐츠 자동화 플랫폼으로 확장했습니다.
외부 소스 수집, 템플릿 렌더링, AI 초안 생성, Tistory 발행 기능을 서비스별로 분리하여 유지보수성과 확장성을 높였습니다.
```

---

## 15. 최종 권장 구현 방향

저사양 로컬 개발 환경을 고려한 최종 추천은 다음과 같다.

```text
Phase 1:
- Frontend: React + Vite + TypeScript
- Backend: Go + Chi + pgx
- Database: PostgreSQL
- Infra: Docker Compose
- External Publish: Selenium 기반 Tistory 글쓰기 자동화
```

향후 확장:

```text
- Content Service: Spring Boot
- Publisher Service: Go
- Connector Service: Go
- AI Serving Service: FastAPI
- Queue/Cache: Redis, 필요 시 추가
```

처음부터 완전한 MSA를 만들지는 않는다. 단, 코드와 Docker Compose 구조를 서비스 단위로 분리해두어, 새로운 기능을 별도 서비스로 추가하기 쉽게 설계한다.
