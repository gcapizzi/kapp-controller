package imgpkg

const (
	DockerHub int = iota
	PrivateRegistry
)

const (
	DockerHubBaseURL string = "index.docker.io"
	URLSeparator     string = "/"
)

type RegistryAuthDetails struct {
	RegistryURL string
	Username    string
	Password    string
}

func (createImgpkg CreateImgPkgStep) PopulateRegistryAuthDetails() (RegistryAuthDetails, error) {
	registry, err := createImgpkg.ui.AskForChoice("Where do you want to push the bundle", []string{"DockerHub", "Private Registry"})
	if err != nil {
		return RegistryAuthDetails{}, err
	}
	switch registry {
	case DockerHub:
		username, err := createImgpkg.ui.AskForText("Enter the username on the docker hub")
		if err != nil {
			return RegistryAuthDetails{}, err
		}
		return RegistryAuthDetails{RegistryURL: DockerHubBaseURL + URLSeparator + username}, nil
	case PrivateRegistry:
		registryURL, err := createImgpkg.ui.AskForText("Enter the registry URL")
		if err != nil {
			return RegistryAuthDetails{}, err
		}
		/*
			username, err := createImgpkg.Ui.AskForText("Registry UserName")
			if err != nil {
				return RegistryAuthDetails{}, err
			}
			password, err := createImgpkg.Ui.AskForPassword("Registry Password")
			if err != nil {
				return RegistryAuthDetails{}, err
			}
		*/

		return RegistryAuthDetails{RegistryURL: registryURL, Username: "", Password: ""}, nil
	}
	return RegistryAuthDetails{}, nil
}