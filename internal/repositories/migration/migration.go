package migration

type Migration interface {
	Up() error
	Down() error
}
