# Контракты тестов

Этот документ фиксирует полный набор поведенческих тестов до реализации их тел.
Каждый пункт становится отдельным `Test...` либо именованным `t.Run`; наборы
однотипных входов оформляются таблицами. Тесты не зависят от порядка запуска,
не используют сеть и создают файлы только под `t.TempDir()`.

## `cmd/dslparser`

- `TestParseArguments` проверяет обязательный единственный путь, `--replace`,
  `--depth N`, обе формы записи depth, ноль, отрицательное, нецелое значение,
  отсутствующее значение, неизвестный флаг и отсутствие публичного `--help`.
- `TestRunRejectsArgumentsBeforeFileSystem` доказывает, что exit code `2`
  возвращается до доступа к пути и создания файлов.
- `TestCLIResultGolden` сравнивает exit code, stdout, stderr, имена и точный JSON
  для корректного документа, нефатальной ошибки, пустого файла, неверного UTF-8,
  неподдерживаемой версии, фатальной ошибки структуры, актуального результата и
  принудительной замены устаревшего результата.
- `TestCLIDoesNotLeaveBinaryArtifacts` собирает программу только в `t.TempDir()`
  и проверяет отсутствие исполняемых и тестовых бинарников в fixture-каталогах.

## `internal/discovery`

- `TestDiscoverPathKinds` покрывает обычный TXT, TXT со смешанным регистром,
  каталог, отсутствующий путь, обычный не-TXT, symlink и специальный файл.
- `TestDiscoverDepthAndOrder` покрывает глубину `0`, `1` и отсутствие лимита,
  лексикографический порядок абсолютных очищенных путей и игнорирование symlink.
- `TestDiscoverScanErrors` проверяет продолжение после ошибки подкаталога,
  нормативный порядок ошибок и отдельную ошибку недоступного корня.
- `TestDiscoverJSONScope` проверяет сбор обычных `.json` без учёта регистра
  только из каталогов с найденными TXT, дедупликацию каталогов, сортировку,
  игнорирование посторонних каталогов и JSON-symlink; для явного TXT областью
  является только его родительский каталог.

## `internal/storage`

- `TestReadSource` проверяет абсолютный очищенный путь, имя, отдельную копию
  точных байтов, byte length, SHA-256 с BOM/EOL и ошибки не-обычного файла.
- `TestTargetPath` проверяет последнее расширение, смешанный регистр `.txt`,
  имя с несколькими точками, отсутствие расширения и недопустимый путь.
- `TestWriteJSONCreateOnly` проверяет атомарное создание, полный объём,
  sync/close и отказ без изменения существующего target.
- `TestWriteJSONReplaceExisting` проверяет backup-rename, установку нового файла,
  удаление backup после успеха и восстановление старых байтов при сбое каждого
  этапа. Временные файлы и backup после завершения не остаются.
- `TestWriteError` проверяет `Operation`, `Path`, текст и `errors.Is/As`.

## `internal/cache`

- `TestReadHeader` проверяет чтение только format version, DSL version,
  document ID и SHA-256, игнорирование неизвестных полей, `null`, повреждённый
  JSON, ошибку чтения и независимость возвращённых строк.
- `TestBuildIndex` проверяет сортировку, отсутствие мутации входа, ошибки reader,
  регистронезависимый индекс ID и независимые копии `Entries`.
- `TestMatchSource` покрывает missing, current, target conflict и duplicate ID;
  сравнивает ID без учёта регистра, но версию и SHA-256 точно; повреждённый или
  неподдерживаемый target не считается current. Конфликты упорядочены по пути.
- `TestMatchCanSkip` возвращает `true` только для `MatchCurrent`.

## `internal/diagnostics`

- `TestRegistry` проверяет полный порядок P001–P014, IO001, допустимые scope,
  политику fatal и нормативные сообщения из `Errors.md` со всеми Details.
- `TestBuilderValidation` проверяет пустой source, неизвестный код, scope,
  location, related locations, обязательные Details и сочетание severity/fatal.
- `TestBuilderOrderAndDeduplication` проверяет порядок обнаружения, ID,
  дедупликацию только одинаковых code/location и сохранение разных кодов.
- `TestBuilderOwnership` проверяет копирование Details.Supported, location,
  related locations и результатов `Items` без скрытой мутации.

## `internal/grammar` и `internal/grammar/v1_2`

- `TestRegistry` проверяет nil grammar, пустую и повторную версию, порядок
  регистрации, lookup, независимость `Versions` и типизированные причины ошибок.
- `TestTagForms` проверяет все 34 имени и каждую строковую/блочную форму таблицы
  P005–P007, регистр, обязательный U+0020, TAB, отсутствующие и лишние части.
- `TestHeadings` проверяет уровни `#`–`###`, название, пробел, лишний `#`,
  экранирование и всегда корневой parent decision.
- `TestContexts` проверяет root, opaque, example, wordlist, table, text, editor,
  instruction, fragment, choice, matching, multifill и variants; отдельно
  проверяются разрешённые separator, дочерние теги, переходы и родители.
- `TestEscaping` проверяет только контекстно значимые escape-последовательности,
  отсутствие повторного разбора, media quotes/backslash и буквальный backslash
  `resource-dir`.
- `TestResourceDirs` проверяет пустой список, пробелы, quoted/unquoted пути,
  запятую внутри кавычек, остаток и исходные диапазоны.
- `TestMultifillPlaceholders` проверяет только защиту диапазона `_____{...}`,
  пустую content-строку, запрет структурного `@answer` и отсутствие элементов
  `placeholder`.
- `TestGrammarProblemsAndRecovery` проверяет P003–P008, P010, P012, их scope,
  element index, Details, fatal и решение продолжения без изменения состояния.

## `internal/parser`

- `TestDecodeSource` покрывает UTF-8, BOM только в начале, LF, CRLF, одиночный
  CR, пустой файл, завершающий EOL и EOF без EOL; P001 не сопровождается P013,
  P002 сохраняет все доступные строки в безопасном режиме.
- `TestUnicodeColumns` проверяет byte-to-rune перевод на ASCII, кириллице,
  составных символах и emoji, включая пустые диапазоны и конец строки.
- `TestProbeSource` проверяет корректную первую строку, неподдерживаемую версию,
  пустой файл, единственный/отсутствующий/повторный document ID и отсутствие
  диагностик и processing ID на лёгком этапе.
- `TestParseEveryTagForm` использует таблицу всех форм DSL v1.2 и проверяет
  line type, raw/eol, элементы, значения, диапазоны и error IDs.
- `TestParserState` проверяет стек блоков, вложение, отдельное закрытие, P009,
  P011, автоматическое завершение task/step/variants/variant и parent/nesting.
- `TestParserRecovery` проверяет каждую P003–P010/P012 в восстанавливаемом и,
  где допустимо, неоднозначном варианте; после fatal весь хвост представлен
  корневыми blank/content без снятия экранирования.
- `TestOpaqueContent` проверяет text/table/multifill, `---` внутри table,
  пустую multifill-строку, editor и отсутствие placeholder/HTML-разбора.
- `TestMetadata` проверяет однозначность ID/title/subtitle/section/order,
  исходный регистр/ведущие нули, порядок resource dirs, `[]` и `null`.
- `TestParserDiagnosticIntegration` проверяет единый diagnostics registry,
  отсутствие дублей и соответствие element error IDs.

## `internal/model`

- `TestJSONContract` фиксирует все имена полей, перечисления, `null` и `[]`.
- `TestValidateDocument` проверяет format version, доступные/недоступные поля,
  SHA-256, byte/line count, metadata и вычисленный `hasErrors`.
- `TestValidateProcessing` проверяет уникальность ID, порядок, ссылки source,
  tool/version, RFC3339 UTC и неотрицательную duration.
- `TestValidateLines` проверяет номера, parent-before-child, nesting, heading,
  EOL, Unicode-диапазоны, точный raw, сортировку и непересечение элементов.
- `TestValidateDiagnostics` проверяет severity/scope/fatal, location/related,
  source, error IDs, обратные element-ссылки, дубликаты и line/document flags.

## `internal/report`

- `TestBuildOwnershipAndArrays` проверяет глубокое копирование всех агрегатов и
  замену обязательных nil-срезов на независимые `[]`.
- `TestBuildHasErrors` проверяет document и line flags для element/line/block/
  document errors, начала многострочного диапазона, дочерней строки и warning.
- `TestMarshal` проверяет вызов model validation, два пробела, один финальный LF,
  отсутствие HTML-экранирования и сохранение нормативного порядка.

## `internal/console`

- `TestFormatFile` проверяет `УСПЕХ`/`ОШИБКА`, action, errors, quoting path,
  необязательный message и отсутствие LF.
- `TestFormatSummary` проверяет точный порядок всех девяти счётчиков.
- `TestWriterRouting` проверяет file/summary в stdout, operational error в stderr,
  ровно один LF и возврат ошибки writer без изменения статуса.

## `internal/app`

- `TestOptionsValidate` проверяет путь, depth, версию инструмента и отсутствие
  обращения к файловой системе.
- `TestSummary` проверяет каждый status/action, diagnostics, scan errors и exit
  codes `0`/`1`; input/argument error `2` остаётся обязанностью CLI.
- `TestRunWithDependencies` проверяет discovery, единый cache до обработки,
  последовательный порядок файлов, время начала, ровно одну итоговую строку,
  scan errors и отсутствие загрузки байтов всех источников одновременно.
- `TestProcessorCacheDecisions` проверяет create, skip current, target conflict,
  replace target и запрет replace при duplicate ID. Skip не создаёт processing.
- `TestProcessorPipelineOrder` фиксирует read, hash, decode, probe, match, parse,
  report, validate, marshal, write; target не открывается до готового JSON.
- `TestProcessorIO001` проверяет частичный документ, единственную fatal IO001,
  попытку записи и переход к операционной ошибке при невозможности записи.
- `TestProcessingAttempt` проверяет уникальный ID, `dsl-parser`, непустую version,
  UTC RFC3339, duration включая формирование отчёта и отсутствие отрицательного
  значения при немонотонных тестовых часах.
