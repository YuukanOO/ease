package testdepdata

import "os"

type (
	Logger interface{}

	Service struct {
		logger Logger
	}

	Database interface{}

	OtherService struct {
		logger Logger
	}

	DBOptions struct {
		ConnectionString string
	}

	logger struct{}
)

func LoadOptions() DBOptions {
	return DBOptions{
		ConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}
}

func NewService(logger Logger, db Database, other *OtherService) *Service {
	return &Service{
		logger: logger,
	}
}

func NewLogger(Database) Logger {
	return &logger{}
}

func NewOtherService(logger Logger) *OtherService {
	return &OtherService{
		logger: logger,
	}
}

func OpenDatabase(DBOptions) Database {
	return nil
}
