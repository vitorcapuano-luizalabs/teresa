package cmd

import (
	"fmt"
	"os"

	context "golang.org/x/net/context"

	"github.com/luizalabs/teresa-api/cmd/client/connection"
	"github.com/luizalabs/teresa-api/pkg/client"
	_ "github.com/prometheus/common/log"
	"github.com/spf13/cobra"

	userpb "github.com/luizalabs/teresa-api/pkg/protobuf"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Create a user",
	Long: `Create a user.

Note that the user's password must be at least 8 characters long. eg.:

	$ teresa create user --email user@mydomain.com --name john --password foobarfoo
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if userNameFlag == "" || userEmailFlag == "" || userPasswordFlag == "" {
			Usage(cmd)
			return
		}
		tc := NewTeresa()
		user, err := tc.CreateUser(userNameFlag, userEmailFlag, userPasswordFlag, isAdminFlag)
		if err != nil {
			log.Fatalf("Failed to create user: %s", err)
		}
		log.Infof("User created. Name: %s Email: %s\n", *user.Name, *user.Email)
	},
}

// delete user
var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Delete an user",
	Long: `Delete an user.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if userIDFlag == 0 {
			Usage(cmd)
			return
		}
		if err := NewTeresa().DeleteUser(userIDFlag); err != nil {
			Fatalf(cmd, "Failed to delete user, err: %s\n", err)
		}
		log.Infof("User deleted.")
	},
}

// set password for an user
var setUserPasswordCmd = &cobra.Command{
	Use:   "set-password",
	Short: "Set password for current user",
	Long:  `Set password for current user.`,
	Run:   setPassword,
}

func setPassword(cmd *cobra.Command, args []string) {
	p, err := client.GetMaskedPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error trying to get the user password: ", err)
	}

	conn, err := connection.New(cfgFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to server: ", err)
	}
	defer conn.Close()

	cli := userpb.NewUserClient(conn)
	if _, err := cli.SetPassword(context.Background(), &userpb.SetPasswordRequest{Password: p}); err != nil {
		fmt.Fprintln(os.Stderr, "Error trying to set new password: ", err)
		return
	}
	fmt.Println("Password updated")
}

func init() {
	createCmd.AddCommand(userCmd)
	userCmd.Flags().StringVar(&userNameFlag, "name", "", "user name [required]")
	userCmd.Flags().StringVar(&userEmailFlag, "email", "", "user email [required]")
	userCmd.Flags().StringVar(&userPasswordFlag, "password", "", "user password [required]")
	userCmd.Flags().BoolVar(&isAdminFlag, "admin", false, "admin")

	deleteCmd.AddCommand(deleteUserCmd)
	deleteUserCmd.Flags().Int64Var(&userIDFlag, "id", 0, "user ID [required]")

	RootCmd.AddCommand(setUserPasswordCmd)
}