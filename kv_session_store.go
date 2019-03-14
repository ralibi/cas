package cas

type KvSessionStore interface {
	Read(key string) (string, bool)
	Write(key, value string)
	Delete(key string)
}
