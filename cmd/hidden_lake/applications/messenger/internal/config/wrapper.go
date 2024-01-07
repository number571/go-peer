package config

var (
	_ IWrapper = &sWrapper{}
)

type sWrapper struct {
	fConfig IConfig
	fEditor IEditor
}

func NewWrapper(pCfg IConfig) IWrapper {
	return &sWrapper{
		fConfig: pCfg,
		fEditor: newEditor(pCfg),
	}
}

func (p *sWrapper) GetConfig() IConfig {
	return p.fConfig
}

func (p *sWrapper) GetEditor() IEditor {
	return p.fEditor
}
