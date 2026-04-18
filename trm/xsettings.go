package trm

// XSettings extends trm.Settings with experimental hook controls.
type XSettings interface {
	Settings
	// EnableSavepointHooks enables running savepoint-scoped hooks at nested boundaries. Default: false.
	EnableSavepointHooks() bool
	// InitialHooks are registered after manager-level defaults when the transaction starts. Default: empty.
	InitialHooks() []Hooks
}

// DefaultXSettings returns XSettings
// that wraps basewith experimental defaults
// (EnableSavepointHooks=false, InitialHooks=nil).
func DefaultXSettings(base Settings) XSettings {
	return &xSettingsWrap{Settings: base, enableSavepointHooks: nil, initialHooks: nil}
}

// NewXSettings returns XSettings that wraps base with the given options.
func NewXSettings(base Settings, opts ...XSettingsOpt) XSettings {
	xSettings := &xSettingsWrap{
		Settings:             base,
		enableSavepointHooks: nil,
		initialHooks:         nil,
	}
	for _, o := range opts {
		o(xSettings)
	}

	return xSettings
}

type xSettingsWrap struct {
	Settings

	enableSavepointHooks *bool
	initialHooks         *[]Hooks
}

func (x *xSettingsWrap) EnableSavepointHooks() bool {
	if x.enableSavepointHooks != nil {
		return *x.enableSavepointHooks
	}

	return false
}

func (x *xSettingsWrap) InitialHooks() []Hooks {
	if x.initialHooks != nil {
		return *x.initialHooks
	}

	return nil
}

// XSettingsOpt configures XSettings.
type XSettingsOpt func(*xSettingsWrap)

// WithEnableSavepointHooks sets EnableSavepointHooks.
func WithEnableSavepointHooks(v bool) XSettingsOpt {
	return func(x *xSettingsWrap) {
		x.enableSavepointHooks = &v
	}
}

// WithInitialHooks sets InitialHooks.
func WithInitialHooks(h []Hooks) XSettingsOpt {
	return func(x *xSettingsWrap) {
		x.initialHooks = &h
	}
}
