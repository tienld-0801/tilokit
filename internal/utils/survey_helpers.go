package utils

import "github.com/AlecAivazis/survey/v2"

// MultiSelect shows a multi-select prompt. Title should be a question.
func MultiSelect(title string, dest *[]string, options []string) error {
    prompt := &survey.MultiSelect{
        Message: title,
        Options: options,
    }
    return survey.AskOne(prompt, dest)
}
