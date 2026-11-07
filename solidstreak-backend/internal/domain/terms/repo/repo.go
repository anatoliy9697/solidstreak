package repo

type Repo interface {
	GetByLangAndVersion(string, string) ([]byte, error)
}

func Init() (Repo, error) {
	return initEmbeddedRepo()
}
