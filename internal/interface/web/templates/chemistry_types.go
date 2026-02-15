package templates

import (
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

type ChemistryListData struct {
	Result     *repositories.PagedResult[entities.ChemistryLog]
	SortBy     string
	SortDir    string
	OutOfRange bool
	DateFrom   string
	DateTo     string
}
