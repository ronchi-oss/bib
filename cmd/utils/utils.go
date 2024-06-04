package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/ronchi-oss/bib/conf"
	"github.com/spf13/cobra"
)

func GetTargetDir(targetDirFlag, profileFlag string) (string, error) {
	var (
		err       error
		profile   *conf.Profile
		targetDir string
	)
	if len(profileFlag) > 0 {
		profile, err = conf.GetProfile(profileFlag)
		if err != nil {
			return targetDir, err
		}
		targetDir = profile.TargetDir
	}
	if len(targetDirFlag) > 0 {
		targetDir = targetDirFlag
	}
	if len(targetDir) == 0 {
		return targetDir, fmt.Errorf("no target directory specified")
	}
	targetDir, err = ExpandPath(targetDir)
	if err != nil {
		return "", fmt.Errorf("failed to expand path %s: %v", targetDir, err)
	}
	if _, err := os.Stat(targetDir); err != nil {
		return targetDir, fmt.Errorf("target directory %s does not exist", targetDir)
	}
	return targetDir, nil
}

func ExpandPath(targetDir string) (string, error) {
	result := targetDir
	if strings.HasPrefix(targetDir, "~/") {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not determine user home directory: %v", err)
		}
		result = userHomeDir + "/" + targetDir[2:]
	}
	return result, nil
}

func InitTargetDirScopedFlags(cmd *cobra.Command, targetDir, profileName *string) {
	cmd.Flags().StringVarP(profileName, "profile", "p", os.Getenv("BIB_PROFILE"), "Profile")
	cmd.Flags().StringVarP(targetDir, "target-dir", "d", "", "Target directory")
	cmd.RegisterFlagCompletionFunc("profile", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		c, err := conf.GetGlobalConf()
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveDefault
		}
		pNames := []string{}
		for _, p := range c.Profiles {
			pNames = append(pNames, p.Name)
		}
		return pNames, cobra.ShellCompDirectiveDefault
	})
}

func GetPreferredEditor() (string, error) {
	e := os.Getenv("BIB_EDITOR")
	if len(e) > 0 {
		return e, nil
	}
	e = os.Getenv("EDITOR")
	if len(e) > 0 {
		return e, nil
	}
	return "", fmt.Errorf("could not determine preferred editor command (both BIB_EDITOR and EDITOR variables are empty)")
}

func TemplateNameShellComp(targetDir, profileName string) ([]string, cobra.ShellCompDirective) {
	targetDir, err := GetTargetDir(targetDir, profileName)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveError
	}
	templates, err := conf.GetTemplates(targetDir)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveError
	}
	return templates, cobra.ShellCompDirectiveDefault
}
