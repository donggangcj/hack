package pkg

import (
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
)

func PromptString(name string) (string, error) {
	prompt := promptui.Prompt{
		Label:    name,
		Validate: ValidateEmptyInput,
	}
	return prompt.Run()
}

func ValidateEmptyInput(input string) error {
	if len(strings.TrimSpace(input)) < 1 {
		return errors.New("this input must not be empty")
	}
	return nil
}
