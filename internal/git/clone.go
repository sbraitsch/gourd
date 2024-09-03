package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
)

func Clone(url string, pat string) string {
	_, err := git.PlainClone("../", false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return fmt.Sprintf("%e", err)
	}
	return "Repository cloned successfully"
}
