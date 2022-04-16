package credential

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func AuthForURLFromHelper(url string, cfg config.Credential) (http.BasicAuth, error) {
	for _, helper := range cfg.Helper {
		if string(helper[0]) == "/" {
			cmd := exec.Command(helper, "get")
			cmd.Stdin = strings.NewReader(fmt.Sprintf("url=%s\n\n", url))

			output, err := cmd.Output()
			if err != nil {
				fmt.Printf("Error while running helper %v: %w", helper, err)
				continue
			}

			auth, err := parseCredential(string(output))
			if err != nil {
				fmt.Printf("Error while parsing helper %v output: %w", helper, err)
				continue
			}

			return auth, nil
		}
	}

	return nil, fmt.Errorf("Credential not found for: %q", url)
}

func parseCredential(cred string) (http.BasicAuth, error) {
	var auth http.BasicAuth

	for _, l := range strings.Split(strings.TrimSpace(cred), "\n") {
		kv := strings.SplitN(l, "=", 2)
		if len(kv) < 2 {
			continue
		}

		name, value := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		switch name {
		case "username":
			auth.Username = value
		case "password":
			auth.Password = value
		}
	}

	if auth.Username != "" && auth.Password != "" {
		return auth, nil
	}

	return auth, fmt.Errorf("Username and password not found")
}
