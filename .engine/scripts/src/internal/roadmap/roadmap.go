package roadmap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spacecraft/internal/config"
	"spacecraft/internal/util"
)

type RoadmapStore interface {
	Create(r *Roadmap) error
	Load(id string) (*Roadmap, error)
	Save(r *Roadmap) error
	List() ([]*Roadmap, error)
	Delete(id string) error
}

type FSStore struct {
	cfg *config.Config
}

func NewFSStore(cfg *config.Config) *FSStore {
	return &FSStore{cfg: cfg}
}

func (s *FSStore) filePath(id string) string {
	return filepath.Join(s.cfg.RoadmapsDir(), id+".json")
}

func (s *FSStore) Create(r *Roadmap) error {
	if err := os.MkdirAll(s.cfg.RoadmapsDir(), 0755); err != nil {
		return err
	}
	return s.Save(r)
}

func (s *FSStore) Load(id string) (*Roadmap, error) {
	path := s.filePath(id)
	if !util.Exists(path) {
		return nil, fmt.Errorf("roadmap not found: %s", id)
	}
	var r Roadmap
	if err := util.ReadJson(path, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *FSStore) Save(r *Roadmap) error {
	return util.WriteJson(s.filePath(r.ID), r)
}

func (s *FSStore) List() ([]*Roadmap, error) {
	dir := s.cfg.RoadmapsDir()
	if !util.Exists(dir) {
		return nil, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []*Roadmap
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		r, err := s.Load(id)
		if err != nil {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}

func (s *FSStore) Delete(id string) error {
	path := s.filePath(id)
	if !util.Exists(path) {
		return fmt.Errorf("roadmap not found: %s", id)
	}
	return os.Remove(path)
}

func DeriveState(r *Roadmap, isShipped func(missionId string) bool) string {
	if len(r.Missions) == 0 {
		return "active"
	}
	for _, mid := range r.Missions {
		if !isShipped(mid.ID) {
			return "active"
		}
	}
	return "done"
}
