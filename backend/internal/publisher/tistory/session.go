package tistory

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultRuntimeDir = "runtime/tistory"
	storageStateFile  = "storage-state.json"
	userDataDirName   = "browser-profile"
)

var ErrSessionRequired = errors.New("tistory session required")

type SessionManager struct {
	mu         sync.Mutex
	runtimeDir string
	process    *exec.Cmd
}

type SessionStatus struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
}

func NewSessionManager(runtimeDir string) *SessionManager {
	if runtimeDir == "" {
		runtimeDir = defaultRuntimeDir
	}
	return &SessionManager{runtimeDir: runtimeDir}
}

func (m *SessionManager) Start(ctx context.Context) (*SessionStatus, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := os.MkdirAll(m.runtimeDir, 0o700); err != nil {
		return nil, err
	}
	if m.process != nil && m.process.Process != nil {
		return &SessionStatus{Connected: false, Message: "로그인 브라우저가 이미 열려 있습니다."}, nil
	}

	cmd := exec.CommandContext(
		ctx,
		"npm",
		"--prefix",
		helperDir(),
		"run",
		"session:start",
		"--",
		"--state",
		m.StorageStatePath(),
		"--user-data-dir",
		filepath.Join(m.runtimeDir, userDataDirName),
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	m.process = cmd

	ready := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var event struct {
				Event string `json:"event"`
			}
			if err := json.Unmarshal(scanner.Bytes(), &event); err == nil && event.Event == "ready" {
				ready <- nil
				return
			}
		}
		if err := scanner.Err(); err != nil {
			ready <- err
			return
		}
		ready <- errors.New("playwright helper exited before ready")
	}()

	go func() {
		_ = cmd.Wait()
		m.mu.Lock()
		if m.process == cmd {
			m.process = nil
		}
		m.mu.Unlock()
	}()

	select {
	case err := <-ready:
		if err != nil {
			return nil, err
		}
	case <-time.After(15 * time.Second):
		return nil, errors.New("playwright helper did not become ready")
	}

	return &SessionStatus{Connected: false, Message: "브라우저에서 Tistory/Kakao 로그인을 완료해주세요."}, nil
}

func (m *SessionManager) Confirm(ctx context.Context, client *Client, blogURL string) (*SessionStatus, error) {
	if _, err := os.Stat(m.StorageStatePath()); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &SessionStatus{Connected: false, Message: "아직 저장된 Tistory 세션이 없습니다."}, nil
		}
		return nil, err
	}
	if err := client.Validate(ctx, blogURL); err != nil {
		return &SessionStatus{Connected: false, Message: "Tistory 세션을 확인하지 못했습니다."}, err
	}
	return &SessionStatus{Connected: true, Message: "Tistory 세션이 연결되었습니다."}, nil
}

func (m *SessionManager) Status(ctx context.Context, client *Client, blogURL string) (*SessionStatus, error) {
	if _, err := os.Stat(m.StorageStatePath()); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &SessionStatus{Connected: false, Message: "저장된 Tistory 세션이 없습니다."}, nil
		}
		return nil, err
	}
	if err := client.Validate(ctx, blogURL); err != nil {
		return &SessionStatus{Connected: false, Message: "Tistory 세션이 만료되었거나 유효하지 않습니다."}, nil
	}
	return &SessionStatus{Connected: true, Message: "Tistory 세션이 연결되어 있습니다."}, nil
}

func (m *SessionManager) Delete() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.process != nil && m.process.Process != nil {
		_ = m.process.Process.Kill()
		m.process = nil
	}
	if err := os.Remove(m.StorageStatePath()); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.RemoveAll(filepath.Join(m.runtimeDir, userDataDirName))
}

func (m *SessionManager) StorageStatePath() string {
	return filepath.Join(m.runtimeDir, storageStateFile)
}

func (m *SessionManager) RequireState() (string, error) {
	path := m.StorageStatePath()
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ErrSessionRequired
		}
		return "", fmt.Errorf("read tistory session: %w", err)
	}
	return path, nil
}

func helperDir() string {
	candidates := []string{
		"tools/tistory-playwright",
		"../tools/tistory-playwright",
		"/tools/tistory-playwright",
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(filepath.Join(candidate, "package.json")); err == nil {
			return candidate
		}
	}
	return "tools/tistory-playwright"
}
