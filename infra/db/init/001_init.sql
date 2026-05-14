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

INSERT INTO templates (type, name, title_format, body_format, is_default)
VALUES (
    'ALGORITHM',
    '기본 알고리즘 풀이 템플릿',
    '[{{platform}}] {{problemTitle}} - {{language}}',
    $body$# [{{platform}}] {{problemTitle}} - {{language}}

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
$body$,
    TRUE
)
ON CONFLICT (type, name) DO NOTHING;

INSERT INTO settings (key, value)
VALUES
    ('tistory.blog_url', 'https://hyeonyway.tistory.com'),
    ('tistory.category_id', ''),
    ('tistory.default_visibility', '0'),
    ('tistory.default_tags', '["알고리즘"]')
ON CONFLICT (key) DO NOTHING;

