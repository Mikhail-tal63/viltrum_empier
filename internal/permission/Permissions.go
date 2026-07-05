package permission

const (
	PermCreateBoard = "create_board"
	PermDeleteBoard = "delete_board"
	PermEditBoard   = "edit_board"

	PermCreateColumn = "create_column"
	PermDeleteColumn = "delete_column"
	PermEditColumn   = "edit_column"

	PermCreateTask = "create_task"
	PermEditTask   = "edit_task"
	PermDeleteTask = "delete_task"

	PermCreateMessage = "create_message"
	PermDeleteMessage = "delete_message"

	PermManageMembers = "manage_members"
	PermManageRoles   = "manage_roles"
)

var DefaultRolePermissions = map[string][]string{
	"member": {
		PermCreateTask,
		PermEditTask,
	},

	"leader": {
		PermCreateTask,
		PermEditTask,
		PermDeleteTask,

		PermCreateColumn,
		PermEditColumn,

		PermEditBoard,
	},

	"admin": {
		PermCreateBoard,
		PermEditBoard,
		PermDeleteBoard,

		PermCreateColumn,
		PermEditColumn,
		PermDeleteColumn,

		PermCreateTask,
		PermEditTask,
		PermDeleteTask,

		PermManageMembers,
	},
}
