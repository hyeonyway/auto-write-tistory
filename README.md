# DevLog Studio

DevLog Studio는 알고리즘 풀이 글을 Markdown으로 작성하고, Playwright 로그인 세션과 Tistory 내부 요청을 통해 비공개 초안으로 발행하는 로컬 우선 도구입니다.

## 기술 스택

- 프론트엔드: React, Vite, TypeScript, Tailwind CSS
- 백엔드: Go, PostgreSQL, Playwright helper
- 로컬 인프라: Docker Compose

## 로컬 실행

전체 로컬 스택을 실행합니다.

```bash
docker compose up --build
```

접속 주소:

```text
Frontend: http://localhost:5173
Backend health: http://localhost:8080/health
Adminer: http://localhost:8081
PostgreSQL host port: localhost:5433
```

Adminer는 `tools` 프로필로 실행합니다.

```bash
docker compose --profile tools up --build
```

기본 데이터베이스 값:

```text
POSTGRES_DB=devlog
POSTGRES_USER=devlog
POSTGRES_PASSWORD=devlog
```

## 환경 변수

로컬 기본값은 `docker-compose.yml`에 정의되어 있습니다.

`.env.example`에는 비밀값이 아닌 로컬 실행 변수만 문서화합니다.

```text
APP_ENV=local
SERVER_PORT=8080
DATABASE_URL=postgres://devlog:devlog@localhost:5433/devlog?sslmode=disable
CORS_ALLOWED_ORIGINS=http://localhost:5173
VITE_API_BASE_URL=http://localhost:8080
```

Kakao/Tistory 로그인 정보는 `.env`에 넣지 않습니다.

## API 응답 형식

API 응답은 다음 envelope 형식을 사용합니다.

```json
{
    "data": {},
    "error": null
}
```

Tistory 발행 API는 저장된 로그인 세션을 사용하므로 요청 본문에 Kakao/Tistory 로그인 정보를 받지 않습니다.

```json
{
    "visibility": 0,
    "categoryId": "1550884",
    "tags": ["알고리즘", "Java"]
}
```

Tistory 카테고리는 `POST /api/tistory/categories/fetch`로 조회합니다. 백엔드는 저장된 Playwright 세션으로 `https://hyeonyway.tistory.com/manage/newpost`를 읽고, 카테고리 목록과 사용자가 선택한 `categoryId`를 설정에 저장합니다.

## 프로젝트 규칙

개발 전 다음 문서를 먼저 확인합니다.

- `GROUND_RULES.md`: 브랜치, 커밋, 병합, 포맷팅, pre-commit, 검증 규칙
- `AGENTS.md`: Codex, Claude 등 코딩 에이전트용 작업 규칙

## 로컬 도구 설정

프론트엔드 의존성을 설치합니다.

```bash
cd frontend
npm install
```

Tistory 로그인 helper 의존성을 설치합니다.

```bash
cd tools/tistory-playwright
npm install
npx playwright install chromium
```

pre-commit을 한 번 설치합니다.

```bash
pre-commit install
```

모든 hook을 수동으로 실행합니다.

```bash
pre-commit run --all-files
```

## Tistory 로그인 주의사항

Tistory Open API가 종료되었기 때문에 발행 기능은 Playwright로 확보한 로그인 세션과 Tistory 내부 요청을 사용합니다.

Kakao 로그인 정보는 앱 화면에 입력하지 않습니다. Settings 화면에서 로그인 브라우저를 열고, 사용자가 Tistory/Kakao 로그인을 직접 완료합니다.

로그인 세션은 `runtime/tistory/storage-state.json`에 저장됩니다. 이 파일은 로컬 런타임 파일이며 Git에 커밋하지 않습니다.

새 PC, 새 브라우저 환경에서는 2단계 인증, CAPTCHA, 본인확인이 발생할 수 있습니다. 이 경우 로그인 브라우저에서 인증을 완료한 뒤 `인증 완료 확인`을 누릅니다.

## 검증

프론트엔드:

```bash
cd frontend
npm run format:check
npm run build
```

백엔드:

```bash
cd backend
go test ./...
```

로컬에 Go가 없다면 Docker로 테스트합니다.

```bash
docker run --rm -v "$PWD/backend:/app" -w /app golang:1.23-alpine go test ./...
```
