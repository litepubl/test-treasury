package importer

type State int

const (
	Empty State = iota + 1
	Updating
	Ok
)

func (s State) String() string {
	return [...]string{"empty", "updating", "ok"}[s-1]
}

func (s State) EnumIndex() int {
	return int(s)
}
