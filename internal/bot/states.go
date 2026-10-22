package bot

// FSM состояния для бота
type UserState string

const (
	StateIdle         UserState = "idle"          // Ожидание команды
	StateWaitLMP      UserState = "wait_lmp"      // Ожидание даты последних месячных
	StateWaitSalary   UserState = "wait_salary"   // Ожидание ввода зарплаты
	StateWaitExcluded UserState = "wait_excluded" // Ожидание количества исключаемых дней
)

// UserSession хранит данные пользователя в сессии
type UserSession struct {
	State     UserState
	LeaveType string // "single", "compl", "multiple"
	LeaveDays int    // Количество дней декрета
}
