# aix — AI eXchange CLI

AI CLI 도구에 컨텍스트 파일을 주입하여 작업을 위임하는 CLI.

```
aix "이 코드에서 버그 찾아줘" -c expert-prompt.md -f target.go
```

## 왜 만들었나

MCP 기반 AI 위임은 오버헤드가 크다 — 서버 설정, 프로토콜 협상, 도구 정의로 인한 컨텍스트 윈도우 소모. `aix`는 이를 직접 서브프로세스 호출로 대체한다. Expert 프롬프트는 프로토콜 레이어가 아닌 컨텍스트 파일(`-c` 플래그)로 주입.

## 설치

```bash
go install github.com/t2sc0m/aix@latest
```

필수: [Codex CLI](https://github.com/openai/codex) 설치 및 인증 완료.

## 사용법

```bash
# 기본 프롬프트
aix "이 에러 설명해줘: connection refused"

# Expert 컨텍스트 파일 (MCP developer-instructions 대체)
aix "이 플랜 리뷰해줘" -c prompts/plan-reviewer.md -f plan.md

# 파일 여러 개 첨부
aix "보안 이슈 찾아줘" -c prompts/security-analyst.md -f auth.go -f handler.go

# stdin 파이프
echo "왜 이렇게 느려?" | aix -f slow-query.sql

# 모델 지정
aix "분석해줘" -m o3

# 샌드박스 모드
aix "이 버그 고쳐줘" -s workspace-write -f broken.go

# 상태 확인
aix status
```

## 플래그

| 플래그 | 단축 | 기본값 | 설명 |
|--------|------|--------|------|
| `--context` | `-c` | | 컨텍스트 파일 (Expert 프롬프트 주입) |
| `--file` | `-f` | | 파일 첨부 (반복 가능) |
| `--model` | `-m` | | 모델 지정 |
| `--sandbox` | `-s` | `read-only` | `read-only`, `workspace-write`, `danger-full-access` |
| `--cwd` | | | codex 작업 디렉토리 |
| `--raw` | | `false` | stdout/stderr 그대로 출력 |
| `--timeout` | `-t` | `300` | 타임아웃 (초) |

## 설정

선택사항 `~/.config/aix/config.yaml`:

```yaml
timeout: 600
sandbox: read-only
adapters:
  codex:
    enabled: true
    model: o3
```

우선순위: CLI 플래그 > config.yaml > 기본값.

## Exit Code

| 코드 | 의미 |
|------|------|
| 0 | 성공 |
| 1 | 일반 에러 |
| 2 | Codex 미설치 |
| 3 | 인증 실패 |
| 4 | Codex 실행 실패 |
| 5 | 타임아웃 |

## 아키텍처

```
cmd/        Cobra CLI 커맨드 (ask, status)
adapter/    AI CLI 백엔드 인터페이스 + Codex 구현
prompt/     프롬프트 조립 + 사이즈 검증
runner/     서브프로세스 추상화 (mock 테스트 지원)
config/     YAML 설정 로더 + 머지 우선순위
```

## 라이선스

MIT
