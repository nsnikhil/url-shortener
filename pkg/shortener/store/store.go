package store

type Store struct {
	ss ShortnerStore
}

func NewStore(ss ShortnerStore) *Store {
	return &Store{
		ss: ss,
	}
}

func (s *Store) GetShortnerStore() ShortnerStore {
	return s.ss
}
