// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	rolebinding "github.com/sq3/provider-syseleven-metakube/internal/controller/namespaced/cluster/rolebinding"
	cronjob "github.com/sq3/provider-syseleven-metakube/internal/controller/namespaced/maintenance/cronjob"
	cluster "github.com/sq3/provider-syseleven-metakube/internal/controller/namespaced/metakube/cluster"
	sshkey "github.com/sq3/provider-syseleven-metakube/internal/controller/namespaced/metakube/sshkey"
	deployment "github.com/sq3/provider-syseleven-metakube/internal/controller/namespaced/node/deployment"
	providerconfig "github.com/sq3/provider-syseleven-metakube/internal/controller/namespaced/providerconfig"
	binding "github.com/sq3/provider-syseleven-metakube/internal/controller/namespaced/role/binding"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		rolebinding.Setup,
		cronjob.Setup,
		cluster.Setup,
		sshkey.Setup,
		deployment.Setup,
		providerconfig.Setup,
		binding.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		rolebinding.SetupGated,
		cronjob.SetupGated,
		cluster.SetupGated,
		sshkey.SetupGated,
		deployment.SetupGated,
		providerconfig.SetupGated,
		binding.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
