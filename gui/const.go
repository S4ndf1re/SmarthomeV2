package gui

// Gui constants
const (
	sessionLogin       = "session_login"
	sessionUsername    = "username"
	sessionPassword    = "password"
	sessionKeepAlive   = 8000
	sessionEnvKey      = "COOKIE_KEY"
	sessionDelete      = -1
	emptyPath          = "/"
	pathSeparator      = emptyPath
	loginPath          = "/login"
	loginApiPath       = "/util/login"
	logoutPath         = "/logout"
	loginPagePath      = "html/login.html"
	fileServeDirectory = "html/"
	apiMountPath       = "/api/"
	guiPath            = "/gui"
)

// Gui Type constants
const (
	ButtonType    = "gui.Button"
	CheckboxType  = "gui.Checkbox"
	TextFieldType = "gui.TextField"
	AlertType     = "gui.Alert"
	DataType      = "gui.Data"
)

// Button constants
const (
	buttonOnClickExtension = "/button/click"
)

// Container constants
const (
	containerInitPath   = "/container/init"
	containerUnloadPath = "/container/unload"
)

// Data constants
const (
	dataRequestPath = "/data/request"
	dataSocketPath  = "/data/socket"
)

// TextField constants
const (
	textFieldTextInputRequest = "/textfield/input"
)

// Checkbox constants
const (
	checkboxOnOnState  = "/onstate/click"
	checkboxOnOffState = "/offstate/click"
	checkboxOnGetState = "/state/get"
)
