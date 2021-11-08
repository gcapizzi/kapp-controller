// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/logger"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	Wait         bool
	PollInterval time.Duration
	PollTimeout  time.Duration
}

func NewDeleteOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *DeleteOptions {
	return &DeleteOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewDeleteCmd(o *DeleteOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"l", "ls"},
		Short:   "Delete a package repository",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set package repository name")

	cmd.Flags().BoolVarP(&o.Wait, "wait", "", true, "Wait for the package repository reconciliation to complete, optional. To disable wait, specify --wait=false")
	cmd.Flags().DurationVarP(&o.PollInterval, "poll-interval", "", 1*time.Second, "Time interval between subsequent polls of package repository reconciliation status, optional")
	cmd.Flags().DurationVarP(&o.PollTimeout, "poll-timeout", "", 5*time.Minute, "Timeout value for polls of package repository reconciliation status, optional")

	return cmd
}

func (o *DeleteOptions) Run() error {
	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Deleting package repository '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)

	o.ui.AskForConfirmation()

	err = client.PackagingV1alpha1().PackageRepositories(
		o.NamespaceFlags.Name).Delete(context.Background(), o.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	//TODO: Use better approach for waiting
	if o.Wait {
		start := time.Now()
		o.ui.PrintLinef("Waiting for deletion to be completed")
		for {
			_, err := client.PackagingV1alpha1().PackageRepositories(
				o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					o.ui.PrintLinef("Resource deleted successfully")
					break
				}
				return err
			}

			if time.Now().After(start.Add(o.PollTimeout)) {
				return fmt.Errorf("Timed out waiting for resource deletion")
			}
			time.Sleep(o.PollInterval)
		}
	}

	return nil
}