// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package kubeconfig_test

import (
	"bytes"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/talos-systems/talos/internal/pkg/kubeconfig"
	"github.com/talos-systems/talos/pkg/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/crypto/x509"
)

type AdminSuite struct {
	suite.Suite
}

func (suite *AdminSuite) TestGenerate() {
	ca, err := x509.NewSelfSignedCertificateAuthority(x509.RSA(true))
	suite.Require().NoError(err)

	u, err := url.Parse("http://localhost:3333/api")
	suite.Require().NoError(err)

	cfg := &v1alpha1.ClusterConfig{
		ClusterName: "talos1",
		ClusterCA: &x509.PEMEncodedCertificateAndKey{
			Crt: ca.CrtPEM,
			Key: ca.KeyPEM,
		},
		ControlPlane: &v1alpha1.ControlPlaneConfig{
			Endpoint: &v1alpha1.Endpoint{
				URL: u,
			},
		},
		AdminKubeconfigConfig: v1alpha1.AdminKubeconfigConfig{
			AdminKubeconfigCertLifetime: time.Hour,
		},
	}

	var buf bytes.Buffer

	suite.Require().NoError(kubeconfig.GenerateAdmin(cfg, &buf))

	// verify config via k8s client
	config, err := clientcmd.Load(buf.Bytes())
	suite.Require().NoError(err)

	suite.Assert().NoError(clientcmd.ConfirmUsable(*config, fmt.Sprintf("admin@%s", cfg.ClusterName)))
}

func TestAdminSuite(t *testing.T) {
	suite.Run(t, new(AdminSuite))
}
