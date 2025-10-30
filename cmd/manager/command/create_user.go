package command

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/app/managerapp"
	postgres2 "github.com/syntaxfa/quick-connect/app/managerapp/repository/postgres"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/types"
)

type CreateUser struct {
	cfg    managerapp.Config
	logger *slog.Logger
	// Flag values
	username    string
	password    string
	fullname    string
	email       string
	phoneNumber string
	roles       []string
}

func (c CreateUser) Command(cfg managerapp.Config, logger *slog.Logger, trap chan os.Signal) *cobra.Command {
	c.cfg = cfg
	c.logger = logger

	cmd := &cobra.Command{
		Use:   "create-user",
		Short: "Create a new user directly in the database",
		Long: `This command creates a new user, instead for bootstrapping
the first superuser or for administrative purpose.
Example:
go run ./cmd/manager create-user -u admin -p 'password' -r superuser -f Admin user`,
		Run: func(_ *cobra.Command, _ []string) {
			c.run()
		},
	}

	// bind flags to the struct fields
	// Sample is :
	// go run cmd/manager/main.go create-user --username=alireza --password=Password --fullname="alireza feizi" --email=alireza@gmail.com --phone_number=00441215485 --role=support --role=superuser
	cmd.Flags().StringVarP(&c.username, "username", "u", "", "Username for the new user (required)")
	cmd.Flags().StringVarP(&c.password, "password", "p", "", "Password for the new user (required)")
	cmd.Flags().StringVarP(&c.fullname, "fullname", "f", "", "Fullname for the new user (required)")
	cmd.Flags().StringVarP(&c.email, "email", "e", "", "email for the new user (required)")
	cmd.Flags().StringVarP(&c.phoneNumber, "phone_number", "n", "", "phone number for the new user (required)")
	cmd.Flags().StringSliceVarP(&c.roles, "role", "r", []string{}, "Roles to assign (e.g., --role=superuser --role=support (required)")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("fullname")
	cmd.MarkFlagRequired("email")
	cmd.MarkFlagRequired("phone_number")
	cmd.MarkFlagRequired("role")

	return cmd
}

func (c CreateUser) run() {
	const op = "command.CreateUser.run"
	ctx := context.Background()
	logger := c.logger.With(slog.String("op", op))

	db := postgres.New(c.cfg.Postgres, c.logger)
	defer db.Close()

	repo := postgres2.New(db)

	tokenSvc := tokenservice.New(c.cfg.Token, c.logger)

	t, tErr := translation.New(translation.DefaultLanguages...)
	if tErr != nil {
		errlog.WithoutErr(tErr, logger)

		return
	}

	userVld := userservice.NewValidate(t)
	userSvc := userservice.New(tokenSvc, userVld, repo, c.logger)

	var roles = make([]types.Role, 0)
	for _, role := range c.roles {
		roles = append(roles, types.Role(role))
	}

	user, cErr := userSvc.CreateUser(ctx, userservice.UserCreateRequest{
		Username:    c.username,
		Password:    c.password,
		Fullname:    c.fullname,
		Email:       c.email,
		PhoneNumber: c.phoneNumber,
		Roles:       roles,
	})
	if cErr != nil {
		errlog.WithoutErr(cErr, logger)

		return
	}

	logger.Info("✅ User created successfully!",
		slog.String("id", string(user.ID)),
		slog.String("username", user.Username),
		slog.String("fullname", user.Fullname),
		slog.String("email", user.Email),
		slog.String("phone number", user.PhoneNumber),
		slog.Any("roles", user.Roles),
	)
}
