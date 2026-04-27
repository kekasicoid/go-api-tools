// internal/domain/formatter.go
package domain

type Formatter interface {
	Format(input string) (string, error)
}
