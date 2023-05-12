package importer

// State  перечислимый тип состояний обновления данных
type State int

// все доступные значения для типа State
const (
	Empty State = iota + 1
	Updating
	Ok
)

// String возвращает состояние в виде текста
func (s State) String() string {
	return [...]string{"empty", "updating", "ok"}[s-1]
}

// EnumIndex возвращает состояние в виде числа
func (s State) EnumIndex() int {
	return int(s)
}
