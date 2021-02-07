package config

// ReadConfig parses commandline arguments, reads parameters from config and from environment
func (p *Provider) ReadConfig(args []string) error {
	return p.Provider.ReadConfig(args)
}
