package storage

const (
	msgTimePing          = "Пинг к БД %s выполнен за %v"
	msgMigrationsNotNeed = "Нет изменений схемы БД. Миграции не требуются"
	msgMigrationsDone    = "Миграции применены"
	msgTimeInsert        = "[%s] информация внесена в БД за время: %v"
	msgTimeSelect        = "[%s] информация найдена в БД за время: %v"
	msgResAffected       = "[%s] не выполнены изменения в БД: %v"
	msgStart             = "Запущена функция %s"
	msgEnd               = "Завершена функция %s"
)

const (
	errCreateDB           = "ошибка при создании БД: %v"
	errCantOpen           = "ошибка при открытии БД %s: %w"
	errPing               = "ошибка пинга БД %s: %w"
	errDriver             = "не удалось загрузить драйвер: %w"
	isntExistMigrations   = "нет каталога с миграциями %s: %w"
	errInitMigrations     = "ошибка инициализации миграции: %w"
	errExecMigrations     = "ошибка применения миграций: %w"
	errLaunchMigrations   = "ошибка обновления схемы БД %s: %w"
	errParseNotActiveConn = "ошибка при парсинге времени ожидания неактивного соединения: %s"
	errParseLifeConn      = "ошибка при парсинге времени жизни соединения: %s"
)

const (
	errCreateTx  = "не смогли создать транзакцию: %w"
	errCommitTx  = "не смогли подтвердить транзакцию: %w"
	errStmt      = "не удалось подготовить sql-запрос: %w"
	errExec      = "не удалось выполнить sql-запрос: %w"
	errRead      = "ошибка чтения данных %v"
	errRes       = "не смогли получить результаты sql-запроса: %w"
	errNoContent = "отсутствуют данные в БД: %v"
)
