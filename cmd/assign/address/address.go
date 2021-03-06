package address

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/admin"
	"github.com/fionera/TeamDriveManager/api/drive"
	. "github.com/fionera/TeamDriveManager/config"
)

var (
	supportedTypes = []string{
		"user",
		"group",
	}

	supportedRoles = []string{
		"organizer",
		"fileOrganizer",
		"writer",
		"commenter",
		"reader",
	}
)

func NewAssignAddressCmd() cli.Command {
	return cli.Command{
		Name:      "address",
		Usage:     "Assign an address to a specified teamdrive",
		Action:    CmdAssignAddress,
		Flags:     []cli.Flag{},
		UsageText: "<TEAMDRIVE-NAME> <ADDRESS> <TYPE> <ROLE>",
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func CmdAssignAddress(c *cli.Context) {
	teamDriveName := c.Args().Get(0)
	address := c.Args().Get(1)
	addressType := c.Args().Get(2)
	role := c.Args().Get(3)

	if teamDriveName == "" {
		logrus.Error("Please supply a teamdrive name")
		return
	}

	if address == "" {
		logrus.Error("Please supply an address")
		return
	}

	if addressType == "" || !contains(supportedTypes, addressType) {
		logrus.Errorf("Unsupported or empty address type (allowed: %s)", strings.Join(supportedTypes, ", "))
		return
	}

	if role == "" {
		logrus.Info("No role supplied. Setting 'reader' permission...")
		role = "reader"
	}

	if !contains(supportedRoles, role) {
		logrus.Error("Unsupported role: '"+role+"' (allowed: %s)", strings.Join(supportedRoles, ", "))
		return
	}

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{drive.DriveScope, admin.AdminDirectoryGroupScope})
	if err != nil {
		logrus.Error(err)
		return
	}

	driveApi, err := drive.NewApi(client)
	if err != nil {
		logrus.Error(err)
		return
	}

	teamDrives, err := driveApi.ListTeamDrives()
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, teamDrive := range teamDrives {
		if teamDrive.Name == teamDriveName {
			_, err := driveApi.CreatePermission(teamDrive.Id, role, address, addressType)
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Info("Added Permission")

			break
		}
	}
}
