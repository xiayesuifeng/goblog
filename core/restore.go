package core

func Restore(file string, useOldConfig bool) error {
	config, err := GetConfigForZip(file)
	if err != nil {
		return err
	}

	if !useOldConfig {
		config.DataDir = Conf.DataDir
		config.Db = Conf.Db
	}

	Conf = config
	if err := SaveConf(); err != nil {
		return err
	}

	if err := Unzip(file, Conf.DataDir); err != nil {
		return err
	}

	return nil
}
