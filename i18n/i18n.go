package i18n

import (
	"os"
	"strings"
)

type Translations map[string]string

var messages = map[string]Translations{
	"en": {
		"start":                  "Start",
		"pause":                  "Pause",
		"stop":                   "Stop",
		"focus":                  "Focus",
		"break":                  "Break",
		"pomodoro":               "Pomodoro",
		"time_to_focus":          "It's time to focus!",
		"time_to_break":          "It's time for a break!",
		"bens_pomodoro":          "Ben's Pomodoro",
		"simple_pomodoro":        "A simple pomodoro for Ben",
		"settings":               "Settings",
		"start_on_launch":        "Start on launch",
		"auto_start_cycles":      "Auto start cycles",
		"inactive_period_1":      "Inactive Period 1",
		"inactive_period_2":      "Inactive Period 2",
		"start_time":             "Start:",
		"end_time":               "End:",
		"durations_in_minutes":   "Durations (minutes)",
		"focus_duration":         "Focus:",
		"short_break_duration":   "Short Break:",
		"long_break_duration":    "Long Break:",
	},
	"es": {
		"start":                  "Iniciar",
		"pause":                  "Pausar",
		"stop":                   "Parar",
		"focus":                  "Foco",
		"break":                  "Pausa",
		"pomodoro":               "Pomodoro",
		"time_to_focus":          "¡Es hora de concentrarse!",
		"time_to_break":          "¡Es hora de un descanso!",
		"bens_pomodoro":          "Pomodoro de Ben",
		"simple_pomodoro":        "Un simple pomodoro para Ben",
		"settings":               "Configuraciones",
		"start_on_launch":        "Iniciar al lanzar",
		"auto_start_cycles":      "Iniciar ciclos automáticamente",
		"inactive_period_1":      "Período inactivo 1",
		"inactive_period_2":      "Período inactivo 2",
		"start_time":             "Inicio:",
		"end_time":               "Fin:",
		"durations_in_minutes":   "Duraciones (minutos)",
		"focus_duration":         "Foco:",
		"short_break_duration":   "Pausa corta:",
		"long_break_duration":    "Pausa larga:",
	},
	"zh": {
		"start":                  "开始",
		"pause":                  "暂停",
		"stop":                   "停止",
		"focus":                  "专注",
		"break":                  "休息",
		"pomodoro":               "番茄钟",
		"time_to_focus":          "是时候专注了！",
		"time_to_break":          "是时候休息一下了！",
		"bens_pomodoro":          "Ben的番茄钟",
		"simple_pomodoro":        "一个简单的番茄钟",
		"settings":               "设置",
		"start_on_launch":        "启动时开始",
		"auto_start_cycles":      "自动开始循环",
		"inactive_period_1":      "非活动期间 1",
		"inactive_period_2":      "非活动期间 2",
		"start_time":             "开始：",
		"end_time":               "结束：",
		"durations_in_minutes":   "持续时间（分钟）",
		"focus_duration":         "专注：",
		"short_break_duration":   "短暂休息：",
		"long_break_duration":    "长期休息：",
	},
	"pt": {
		"start":                  "Iniciar",
		"pause":                  "Pausar",
		"stop":                   "Parar",
		"focus":                  "Foco",
		"break":                  "Pausa",
		"pomodoro":               "Pomodoro",
		"time_to_focus":          "É hora de focar!",
		"time_to_break":          "É hora de uma pausa!",
		"bens_pomodoro":          "Pomodoro do Ben",
		"simple_pomodoro":        "Um simples pomodoro para o Ben",
		"settings":               "Configurações",
		"start_on_launch":        "Iniciar no lançamento",
		"auto_start_cycles":      "Iniciar ciclos automaticamente",
		"inactive_period_1":      "Período inativo 1",
		"inactive_period_2":      "Período inativo 2",
		"start_time":             "Início:",
		"end_time":               "Fim:",
		"durations_in_minutes":   "Durações (minutos)",
		"focus_duration":         "Foco:",
		"short_break_duration":   "Pausa curta:",
		"long_break_duration":    "Pausa longa:",
	},
}

var lang string

func init() {
	lang = os.Getenv("LANG")
	if lang == "" {
		lang = "en" // default to English
	}
	lang = strings.Split(lang, ".")[0]
	lang = strings.Split(lang, "_")[0]
}

func T(key string) string {
	if translations, ok := messages[lang]; ok {
		if translation, ok := translations[key]; ok {
			return translation
		}
	}
	// Fallback to English if the language or key is not found
	if translations, ok := messages["en"]; ok {
		if translation, ok := translations[key]; ok {
			return translation
		}
	}
	return key // return key if no translation is found
}
