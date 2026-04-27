// internal/usecase/formatter_usecase.go
package usecase

import "github.com/kekasicoid/go-api-tools/internal/domain"

type FormatterUsecase struct {
	formatter domain.Formatter
}

func NewFormatterUsecase(f domain.Formatter) *FormatterUsecase {
	return &FormatterUsecase{formatter: f}
}

func (u *FormatterUsecase) FormatJSON(input string) (string, error) {
	return u.formatter.Format(input)
}
