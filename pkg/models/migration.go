package models

type Migration interface {
	Migrate()
}
