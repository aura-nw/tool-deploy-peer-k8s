package util

import "github.com/manifoldco/promptui"

func GetInputFromKeyboard(label string, defaultValue string) string {
	promt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	result, err := promt.Run()
	if err != nil {
		panic(err)
	}
	return result
}

func GetInputFromSelect(label string, list []string) string {
	promt := promptui.Select{
		Label: label,
		Items: list,
	}
	_, result, err := promt.Run()
	if err != nil {
		panic(err)
	}
	return result
}
