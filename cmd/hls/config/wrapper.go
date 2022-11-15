package config

var (
	_ IWrapper = &sWrapper{}
)

type sWrapper struct {
	fConfig IConfig
	fEditor IEditor
}

func NewWrapper(cfg IConfig) IWrapper {
	return &sWrapper{
		fConfig: cfg,
		fEditor: newEditor(cfg),
	}
}

func (w *sWrapper) Config() IConfig {
	return w.fConfig
}

func (w *sWrapper) Editor() IEditor {
	return w.fEditor
}
