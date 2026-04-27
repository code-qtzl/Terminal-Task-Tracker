package store

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
	IsComplete  bool
}

type Store struct {
	f     *os.File
	tasks []Task
	dirty bool
}

var ErrNotFound = errors.New("task not found")

var csvHeader = []string{"ID", "Description", "CreatedAt", "IsComplete"}

func Open(path string) (*Store, error) {
	f, err := loadFile(path)
	if err != nil {
		return nil, err
	}
	s := &Store{f: f}
	if err := s.load(); err != nil {
		_ = closeFile(f)
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	if _, err := s.f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	r := csv.NewReader(s.f)
	rows, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv: %w", err)
	}
	if len(rows) == 0 {
		return nil
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 4 {
			return fmt.Errorf("malformed row %d", i)
		}
		id, err := strconv.Atoi(row[0])
		if err != nil {
			return fmt.Errorf("parse id row %d: %w", i, err)
		}
		created, err := time.Parse(time.RFC3339, row[2])
		if err != nil {
			return fmt.Errorf("parse time row %d: %w", i, err)
		}
		done, err := strconv.ParseBool(row[3])
		if err != nil {
			return fmt.Errorf("parse done row %d: %w", i, err)
		}
		s.tasks = append(s.tasks, Task{
			ID:          id,
			Description: row[1],
			CreatedAt:   created,
			IsComplete:  done,
		})
	}
	return nil
}

func (s *Store) Close() error {
	if s.dirty {
		if err := s.save(); err != nil {
			_ = closeFile(s.f)
			return err
		}
	}
	return closeFile(s.f)
}

func (s *Store) save() error {
	if _, err := s.f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := s.f.Truncate(0); err != nil {
		return err
	}
	w := csv.NewWriter(s.f)
	if err := w.Write(csvHeader); err != nil {
		return err
	}
	for _, t := range s.tasks {
		row := []string{
			strconv.Itoa(t.ID),
			t.Description,
			t.CreatedAt.Format(time.RFC3339),
			strconv.FormatBool(t.IsComplete),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func (s *Store) Add(description string) Task {
	id := 1
	for _, t := range s.tasks {
		if t.ID >= id {
			id = t.ID + 1
		}
	}
	t := Task{
		ID:          id,
		Description: description,
		CreatedAt:   time.Now(),
		IsComplete:  false,
	}
	s.tasks = append(s.tasks, t)
	s.dirty = true
	return t
}

func (s *Store) Complete(id int) error {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].IsComplete = true
			s.dirty = true
			return nil
		}
	}
	return ErrNotFound
}

func (s *Store) Delete(id int) error {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			s.dirty = true
			return nil
		}
	}
	return ErrNotFound
}

func (s *Store) List(showAll bool) []Task {
	if showAll {
		out := make([]Task, len(s.tasks))
		copy(out, s.tasks)
		return out
	}
	out := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		if !t.IsComplete {
			out = append(out, t)
		}
	}
	return out
}
