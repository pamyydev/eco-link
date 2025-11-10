package ports

import "github.com/pamelamiranda/eco-link/internal/domain"

// ReportRepository define as operações para persistência de relatórios
type ReportRepository interface {
	Save(report domain.Report) error
	FindByURL(url string) ([]domain.Report, error)
	FindRecent() ([]domain.Report, error)
}
