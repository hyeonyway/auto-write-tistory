# DevLog Studio Phase 1 문제점 및 대응 방향

## 1. 초안 수정 기능 부재

Status: Resolved in Phase 1 follow-up.

현재 Phase 1 구현에는 초안을 수정하는 화면과 API가 없다.  
대시보드는 글 목록 조회와 삭제만 제공하고, 새 글 작성만 가능하다. 또한 백엔드에도 `PUT /api/posts/{id}` 또는 `PATCH /api/posts/{id}`가 없다.

대응 방향:

- 초안 상세/편집 화면 추가
- `GET /api/posts/{id}`로 기존 글과 원본 입력값을 불러오기
- `PUT /api/posts/{id}` 또는 `PATCH /api/posts/{id}`로 수정 저장 지원
- `post_inputs`에 저장된 원본 JSON을 기반으로 폼 복원

구현 결과:

- `PUT /api/posts/{id}` 추가
- `GET /api/posts/{id}` 응답에 최신 `post_inputs.input_json` 포함
- 대시보드에서 수정 화면으로 이동 가능
- 기존 작성 폼을 편집 모드에서도 재사용
- 수정 저장 시 Markdown을 다시 렌더링하고 상태를 `DRAFT`로 되돌림

## 2. 카카오 로그인 2단계 인증 문제

Status: Resolved for local Phase 1 workflow.

현재 Selenium은 요청마다 새 세션을 만들고, 브라우저 프로필을 재사용하지 않는다.  
이 때문에 카카오 로그인 후 2단계 인증이 필요한 환경에서는 인증 완료 전에 세션이 끊기거나 다시 로그인 상태가 되어 실패할 수 있다.

대응 방향:

- Selenium 브라우저 프로필을 영구화
- 로그인/인증을 발행 요청과 분리
- Settings 화면에 세션 연결 또는 인증 준비 단계 추가
- 로그인 실패, 2FA 필요, 세션 만료를 구분하는 에러 코드 추가

구현 결과:

- Selenium 세션을 30분간 재사용
- 발행/카테고리 조회 후 세션을 즉시 삭제하지 않도록 변경
- 로그인 성공 기준을 카카오 로그인 요청 성공이 아니라 Tistory 글쓰기 화면의 `#category-btn` 로드로 변경
- 카카오 추가 인증이 남아 있으면 `TISTORY_AUTH_REQUIRED`로 구분
- 수동 인증을 위해 Selenium noVNC 포트 `7900` 노출

## 3. 카테고리 조회와 발행의 인증 의존성

Status: Resolved for local Phase 1 workflow.

카테고리 조회도 `https://hyeonyway.tistory.com/manage/newpost`에 진입해야 하므로 로그인 세션에 의존한다.  
현재 구조에서는 카테고리 조회와 발행이 각각 독립적으로 로그인 세션을 만들기 때문에 인증 부담이 커진다.

대응 방향:

- 카테고리 조회와 발행이 같은 Selenium 인증 세션을 재사용하도록 변경
- 필요한 경우 카테고리만 미리 불러와 저장하고 발행 시에는 저장된 `categoryId`만 사용

구현 결과:

- 카테고리 조회와 발행이 같은 Publisher 세션 캐시를 사용
- 같은 blogUrl/loginId 조합이면 기존 Selenium 세션을 재사용
- 기존 세션이 만료되었거나 글쓰기 화면 진입에 실패하면 새 세션으로 로그인

## 4. 후속 보완 원칙

위 문제들은 Phase 1의 기본 흐름을 막는 수준은 아니지만 실제 사용성에는 직접 영향을 준다.  
따라서 Phase 1 완료 후 보완 작업으로 분리하고, 우선은 저장, 미리보기, 발행의 기본 흐름을 안정화한다.
