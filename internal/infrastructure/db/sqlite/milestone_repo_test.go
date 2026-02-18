package sqlite

import (
	"testing"

	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

func TestMilestoneRepoImplementsInterface(t *testing.T) {
	var _ repositories.MilestoneRepository = (*MilestoneRepo)(nil)
}
