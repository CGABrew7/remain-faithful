package core

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type Session struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Store struct {
	mu    sync.Mutex
	dir   string
	sess  Session
	state LocalState
	ok    bool
}

func OpenStore(dir string) (*Store, error) {
	if dir == "" {
		base, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(base, "Remain Faithful")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	s := &Store{dir: dir}
	s.state.normalize()
	if raw, err := os.ReadFile(s.sessionPath()); err == nil {
		_ = json.Unmarshal(raw, &s.sess)
	}
	if raw, err := os.ReadFile(s.statePath()); err == nil {
		_ = json.Unmarshal(raw, &s.state)
	}
	s.state.normalize()
	s.ok = s.sess.Token != "" && s.sess.User.ID > 0
	return s, nil
}

func (s *Store) sessionPath() string { return filepath.Join(s.dir, "session.json") }
func (s *Store) statePath() string   { return filepath.Join(s.dir, "state.json") }

func (s *Store) SignedIn() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ok
}

func (s *Store) Token() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.ok {
		return ""
	}
	return s.sess.Token
}

func (s *Store) User() User {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sess.User
}

func (s *Store) SetSession(token string, user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sess = Session{Token: token, User: user}
	s.ok = token != "" && user.ID > 0
	return s.writeSession()
}

func (s *Store) UpdateToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sess.Token = token
	s.ok = token != "" && s.sess.User.ID > 0
	return s.writeSession()
}

func (s *Store) UpdateUser(user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user.ID == 0 {
		user.ID = s.sess.User.ID
	}
	s.sess.User = user
	return s.writeSession()
}

func (s *Store) ClearSession() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sess = Session{}
	s.ok = false
	ClearMonitoring(&s.state)
	if err := s.writeState(); err != nil {
		return err
	}
	if err := os.Remove(s.sessionPath()); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func (s *Store) State() LocalState {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.normalize()
	cp := s.state
	cp.Selected = append([]string(nil), s.state.Selected...)
	cp.Queue = append([]ProtectionNotice(nil), s.state.Queue...)
	return cp
}

func (s *Store) Update(fn func(*LocalState) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.normalize()
	prev := s.state
	prev.Selected = append([]string(nil), s.state.Selected...)
	prev.Queue = append([]ProtectionNotice(nil), s.state.Queue...)
	if err := fn(&s.state); err != nil {
		s.state = prev
		return err
	}
	s.state.normalize()
	if err := s.writeState(); err != nil {
		s.state = prev
		return err
	}
	return nil
}

func (s *Store) Enqueue(userID int64, notes []ProtectionNotice) error {
	if len(notes) == 0 {
		return nil
	}
	return s.Update(func(st *LocalState) error {
		if st.QueueUser != 0 && st.QueueUser != userID {
			st.Queue = nil
		}
		st.QueueUser = userID
		st.Queue = append(st.Queue, notes...)
		if len(st.Queue) > 40 {
			st.Queue = st.Queue[len(st.Queue)-40:]
		}
		return nil
	})
}

func (s *Store) TakeQueue(userID int64) []ProtectionNotice {
	s.mu.Lock()
	defer s.mu.Unlock()
	if st := s.state.QueueUser; st != 0 && st != userID {
		return nil
	}
	notes := append([]ProtectionNotice(nil), s.state.Queue...)
	return notes
}

func (s *Store) ReplaceQueue(userID int64, notes []ProtectionNotice) error {
	return s.Update(func(st *LocalState) error {
		st.QueueUser = userID
		st.Queue = append([]ProtectionNotice(nil), notes...)
		return nil
	})
}

func (s *Store) PanicTap(userID int64) (int, error) {
	var n int
	err := s.Update(func(st *LocalState) error {
		if st.PanicUser != userID {
			st.PanicUser = userID
			st.PanicTaps = 0
		}
		st.PanicTaps++
		n = st.PanicTaps
		return nil
	})
	return n, err
}

func (s *Store) writeSession() error {
	return writeFile(s.sessionPath(), s.sess, 0o600)
}

func (s *Store) writeState() error {
	return writeFile(s.statePath(), s.state, 0o600)
}

func writeFile(path string, v any, mode os.FileMode) error {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
