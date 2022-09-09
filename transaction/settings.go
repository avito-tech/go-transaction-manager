package transaction

type SettingsOpt func(s *SettingsImpl)

// TODO rename.
type SettingsImpl struct {
	ctxKey      CtxKey
	isReadOnly  bool
	propagation Propagation
}

func NewSettings(oo ...SettingsOpt) SettingsImpl {
	s := &SettingsImpl{
		ctxKey:      ctxKey{},
		isReadOnly:  false,
		propagation: PropagationRequired,
	}

	for _, o := range oo {
		o(s)
	}

	return *s
}

func (s SettingsImpl) CtxKey() CtxKey {
	return s.ctxKey
}

func (s SettingsImpl) IsReadOnly() bool {
	return s.isReadOnly
}

func (s SettingsImpl) Propagation() Propagation {
	return s.propagation
}

// TODO fix long name.
func SettingsWithCtxKey(key CtxKey) SettingsOpt {
	return func(s *SettingsImpl) {
		s.ctxKey = key
	}
}

func SettingsWithReadOnly(is bool) SettingsOpt {
	return func(s *SettingsImpl) {
		s.isReadOnly = is
	}
}

func SettingsWithPropagation(p Propagation) SettingsOpt {
	return func(s *SettingsImpl) {
		s.propagation = p
	}
}
