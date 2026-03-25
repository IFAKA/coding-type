package snippets

// Snippet represents a single code typing exercise.
type Snippet struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Language   string   `json:"language"`   // python, javascript, typescript, go, cpp
	Difficulty string   `json:"difficulty"` // easy, medium, hard
	Tags       []string `json:"tags"`
	Code       string   `json:"code"`
}

// Config holds the user's selected session parameters.
type Config struct {
	Language   string
	Difficulty string
	Mode       string // "practice" or "timed"
}

// Languages is the ordered list of supported languages.
var Languages = []string{"python", "javascript", "typescript", "go", "cpp"}

// Difficulties is the ordered list of difficulty levels.
var Difficulties = []string{"easy", "medium", "hard"}

// Modes is the ordered list of practice modes.
var Modes = []string{"practice", "timed"}

// LangDisplay maps language IDs to display names.
var LangDisplay = map[string]string{
	"python":     "python",
	"javascript": "javascript",
	"typescript": "typescript",
	"go":         "go",
	"cpp":        "c++",
}

// ChromaLang maps snippet language IDs to chroma lexer names.
var ChromaLang = map[string]string{
	"python":     "python",
	"javascript": "javascript",
	"typescript": "typescript",
	"go":         "go",
	"cpp":        "c++",
}
