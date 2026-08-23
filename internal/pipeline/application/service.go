package pipelineapp

import (
	"fmt"
	pipeline "github.com/acme/streamforge-cdc/internal/pipeline/domain"
	source "github.com/acme/streamforge-cdc/internal/source/domain"
)

type Service struct{ Registry *pipeline.Registry }

func (s *Service) Create(p pipeline.Pipeline) error {
	if p.Status == "" {
		p.Status = pipeline.PipelineDraft
	}
	return s.Registry.Put(p)
}
func (s *Service) Start(id string, src source.Source) error {
	p, ok := s.Registry.Get(id)
	if !ok {
		return fmt.Errorf("pipeline %s not found", id)
	}
	if p.SourceID != src.ID {
		return fmt.Errorf("pipeline source mismatch")
	}
	p.Status = pipeline.PipelineRunning
	return s.Registry.Put(p)
}
func (s *Service) Pause(id string) error {
	p, ok := s.Registry.Get(id)
	if !ok {
		return fmt.Errorf("pipeline %s not found", id)
	}
	p.Status = pipeline.PipelinePaused
	return s.Registry.Put(p)
}
