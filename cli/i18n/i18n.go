package i18n

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"strings"

	"github.com/essentialkaos/ek/v12/i18n"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type I18NBundle struct {
	INFO  *I18NInfo
	USAGE *I18NUsage

	ERRORS *I18NErrors
}

type I18NErrors struct {
	OPTION_PARSING      i18n.String
	UNSUPPORTED_COMMAND i18n.String
	UNKNOWN_ARCH        i18n.String

	EMPTY_STRUCT_NAME i18n.String
	NO_STRUCT         i18n.String
	NO_ANY_STRUCTS    i18n.String
}

type I18NInfo struct {
	ALL_OPTIMAL     i18n.String
	OPTIMIZE_ADVICE i18n.String
	WITH_OPTIMAL    i18n.String
	ALREADY_OPTIMAL i18n.String
}

type I18NUsage struct {
	DESC      i18n.String
	ARGUMENTS i18n.String

	COMMANDS *I18NCommands
	OPTIONS  *I18NOptions
	EXAMPLES *I18NExamples
}

type I18NCommands struct {
	CHECK i18n.String
	VIEW  i18n.String
}

type I18NOptions struct {
	ARCH       i18n.String
	ARCH_VAL   i18n.String
	STRUCT     i18n.String
	STRUCT_VAL i18n.String
	TAGS       i18n.String
	TAGS_VAL   i18n.String
	PAGER      i18n.String
	NO_COLOR   i18n.String
	HELP       i18n.String
	VER        i18n.String
}

type I18NExamples struct {
	EXAMPLE_1 i18n.String
	EXAMPLE_2 i18n.String
	EXAMPLE_3 i18n.String
	EXAMPLE_4 i18n.String
	EXAMPLE_5 i18n.String
}

// ////////////////////////////////////////////////////////////////////////////////// //

// UI is a bundle with data for used language
var UI = getEN()

// ////////////////////////////////////////////////////////////////////////////////// //

// SetLanguage sets app language
func SetLanguage() {
	lang, _, _ := strings.Cut(os.Getenv("LANG"), "_")

	switch strings.ToLower(lang) {
	case "ru":
		l, _ := i18n.Fallback(getEN(), getRU())
		UI = l.(*I18NBundle)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getEN returns bundle for English language
func getEN() *I18NBundle {
	return &I18NBundle{
		INFO: &I18NInfo{
			ALL_OPTIMAL:     "{g}All structs are well aligned{!}",
			OPTIMIZE_ADVICE: "Struct {*}%s{!} {s-}(%s:%d){!} fields order can be optimized (%d → %d)",
			WITH_OPTIMAL:    "{s-}// %s:%d | Size: %d (Optimal: %d){!}",
			ALREADY_OPTIMAL: "{s-}// %s:%d | Size: %d{!}",
		},
		ERRORS: &I18NErrors{
			OPTION_PARSING:      "Options parsing errors",
			UNSUPPORTED_COMMAND: "Command %s is unsupported",
			UNKNOWN_ARCH:        "Unknown arch %s",

			NO_ANY_STRUCTS:    "Given package doesn't have any structs",
			NO_STRUCT:         "Can't find struct with name %q",
			EMPTY_STRUCT_NAME: "You should define struct name",
		},
		USAGE: &I18NUsage{
			DESC:      "Utility for viewing and checking Go struct alignment",
			ARGUMENTS: "path…",
			COMMANDS: &I18NCommands{
				CHECK: "Check package for alignment problems",
				VIEW:  "Print alignment info for all structs",
			},
			OPTIONS: &I18NOptions{
				ARCH:       "Architecture name",
				ARCH_VAL:   "name",
				STRUCT:     "Print info only about struct with given name",
				STRUCT_VAL: "name",
				TAGS:       "Build tags {s-}(mergeble){!}",
				TAGS_VAL:   "tag…",
				PAGER:      "Use pager for long output",
				NO_COLOR:   "Disable colors in output",
				HELP:       "Show this help message",
				VER:        "Show version",
			},
			EXAMPLES: &I18NExamples{
				EXAMPLE_1: "Show info about all structs in current package",
				EXAMPLE_2: "Check current package",
				EXAMPLE_3: "Check current package and all sub-packages",
				EXAMPLE_4: "Check current package and all sub-packages with custom build tags",
				EXAMPLE_5: "Show info about PostMessageParameters struct",
			},
		},
	}
}

// getRU returns bundle for Russian language
func getRU() *I18NBundle {
	return &I18NBundle{
		INFO: &I18NInfo{
			ALL_OPTIMAL:     "{g}Проблем с выравниванием структур не обнаружено{!}",
			OPTIMIZE_ADVICE: "Поля структуры {*}%s{!} {s-}(%s:%d){!} могут быть оптимизированны (%d → %d)",
			WITH_OPTIMAL:    "{s-}// %s:%d | Размер: %d (Оптимальный: %d){!}",
			ALREADY_OPTIMAL: "{s-}// %s:%d | Размер: %d{!}",
		},
		ERRORS: &I18NErrors{
			OPTION_PARSING:      "Ошибки обработки опций",
			UNSUPPORTED_COMMAND: "Команда %s не поддерживается",
			UNKNOWN_ARCH:        "Неизвестная архитектура %s",

			NO_ANY_STRUCTS:    "Указанный пакет не содержит структур",
			NO_STRUCT:         "Структура с именем %q не найдена",
			EMPTY_STRUCT_NAME: "Вы должны указать имя структуры",
		},
		USAGE: &I18NUsage{
			DESC:      "Утилита для просмотра и проверки выравнивания полей в структрах Go",
			ARGUMENTS: "путь…",
			COMMANDS: &I18NCommands{
				CHECK: "Проверка на наличие проблем с выравниванием",
				VIEW:  "Отображние информации о выравнивании",
			},
			OPTIONS: &I18NOptions{
				ARCH:       "Название архитектуры",
				ARCH_VAL:   "имя",
				STRUCT:     "Отображение информации только для указанной структуры",
				STRUCT_VAL: "имя",
				TAGS:       "Тэги сборки {s-}(повторяемая опция){!}",
				TAGS_VAL:   "тэг…",
				PAGER:      "Использовать паджинацию для длинного вывода",
				NO_COLOR:   "Отключение цветного вывода",
				HELP:       "Показать это справочное сообщение",
				VER:        "Показать версию",
			},
			EXAMPLES: &I18NExamples{
				EXAMPLE_1: "Просмотр информации о всех структурах пакета",
				EXAMPLE_2: "Проверка текущей директории",
				EXAMPLE_3: "Проверка текущей директории и всех дочерних",
				EXAMPLE_4: "Проверка текущей директории и всех дочерних с использованием тэгов",
				EXAMPLE_5: "Отображение информации о структуре PostMessageParameters",
			},
		},
	}
}
