package theme

import "github.com/charmbracelet/lipgloss"

// Built-in themes
var (
	Monokai = BaseTheme{
		ThemeName: "monokai",
		ThemeColors: ThemeColors{
			Debug:      lipgloss.Color("#75715E"),
			Info:       lipgloss.Color("#A6E22E"),
			Warn:       lipgloss.Color("#E6DB74"),
			Error:      lipgloss.Color("#F92672"),
			Fatal:      lipgloss.Color("#F92672"),
			Timestamp:  lipgloss.Color("#66D9EF"),
			Key:        lipgloss.Color("#FD971F"),
			Value:      lipgloss.Color("#E6DB74"),
			Background: lipgloss.Color("#272822"),
			Foreground: lipgloss.Color("#F8F8F2"),
			StatusBar:  lipgloss.Color("#3E3D32"),
			StatusText: lipgloss.Color("#F8F8F2"),
			Border:     lipgloss.Color("#75715E"),
			Highlight:  lipgloss.Color("#49483E"),
		},
	}

	Dracula = BaseTheme{
		ThemeName: "dracula",
		ThemeColors: ThemeColors{
			Debug:      lipgloss.Color("#6272A4"),
			Info:       lipgloss.Color("#50FA7B"),
			Warn:       lipgloss.Color("#F1FA8C"),
			Error:      lipgloss.Color("#FF5555"),
			Fatal:      lipgloss.Color("#FF5555"),
			Timestamp:  lipgloss.Color("#8BE9FD"),
			Key:        lipgloss.Color("#BD93F9"),
			Value:      lipgloss.Color("#F1FA8C"),
			Background: lipgloss.Color("#282A36"),
			Foreground: lipgloss.Color("#F8F8F2"),
			StatusBar:  lipgloss.Color("#44475A"),
			StatusText: lipgloss.Color("#F8F8F2"),
			Border:     lipgloss.Color("#6272A4"),
			Highlight:  lipgloss.Color("#44475A"),
		},
	}

	Gruvbox = BaseTheme{
		ThemeName: "gruvbox",
		ThemeColors: ThemeColors{
			Debug:      lipgloss.Color("#928374"),
			Info:       lipgloss.Color("#B8BB26"),
			Warn:       lipgloss.Color("#FABD2F"),
			Error:      lipgloss.Color("#FB4934"),
			Fatal:      lipgloss.Color("#FB4934"),
			Timestamp:  lipgloss.Color("#83A598"),
			Key:        lipgloss.Color("#FE8019"),
			Value:      lipgloss.Color("#FABD2F"),
			Background: lipgloss.Color("#282828"),
			Foreground: lipgloss.Color("#EBDBB2"),
			StatusBar:  lipgloss.Color("#3C3836"),
			StatusText: lipgloss.Color("#EBDBB2"),
			Border:     lipgloss.Color("#928374"),
			Highlight:  lipgloss.Color("#3C3836"),
		},
	}

	Nord = BaseTheme{
		ThemeName: "nord",
		ThemeColors: ThemeColors{
			Debug:      lipgloss.Color("#616E88"),
			Info:       lipgloss.Color("#A3BE8C"),
			Warn:       lipgloss.Color("#EBCB8B"),
			Error:      lipgloss.Color("#BF616A"),
			Fatal:      lipgloss.Color("#BF616A"),
			Timestamp:  lipgloss.Color("#88C0D0"),
			Key:        lipgloss.Color("#81A1C1"),
			Value:      lipgloss.Color("#EBCB8B"),
			Background: lipgloss.Color("#2E3440"),
			Foreground: lipgloss.Color("#ECEFF4"),
			StatusBar:  lipgloss.Color("#3B4252"),
			StatusText: lipgloss.Color("#ECEFF4"),
			Border:     lipgloss.Color("#4C566A"),
			Highlight:  lipgloss.Color("#3B4252"),
		},
	}

	Kanagawa = BaseTheme{
		ThemeName: "kanagawa",
		ThemeColors: ThemeColors{
			Debug:      lipgloss.Color("#727169"), // sumiInk4
			Info:       lipgloss.Color("#7FB4CA"), // springBlue
			Warn:       lipgloss.Color("#FF9E3B"), // roninYellow
			Error:      lipgloss.Color("#E82424"), // samuraiRed
			Fatal:      lipgloss.Color("#C34043"), // peachRed
			Timestamp:  lipgloss.Color("#7AA89F"), // waveAqua1
			Key:        lipgloss.Color("#957FB8"), // oniViolet
			Value:      lipgloss.Color("#DCA561"), // boatYellow2
			Background: lipgloss.Color("#1F1F28"), // sumiInk1
			Foreground: lipgloss.Color("#DCD7BA"), // fujiWhite
			StatusBar:  lipgloss.Color("#2A2A37"), // sumiInk3
			StatusText: lipgloss.Color("#DCD7BA"), // fujiWhite
			Border:     lipgloss.Color("#54546D"), // sumiInk5
			Highlight:  lipgloss.Color("#2D4F67"), // waveBlue2
		},
	}
)

// registry holds all available themes.
var registry = map[string]Theme{
	"monokai":  Monokai,
	"dracula":  Dracula,
	"gruvbox":  Gruvbox,
	"nord":     Nord,
	"kanagawa": Kanagawa,
}

// Get returns a theme by name. Falls back to Kanagawa if not found.
func Get(name string) Theme {
	if t, ok := registry[name]; ok {
		return t
	}
	return Kanagawa
}

// Names returns all available theme names.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
