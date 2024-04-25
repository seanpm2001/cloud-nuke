package resources

import (
	"context"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/networkfirewall"
	"github.com/aws/aws-sdk-go/service/networkfirewall/networkfirewalliface"
	"github.com/gruntwork-io/cloud-nuke/config"
	"github.com/gruntwork-io/go-commons/errors"
)

type NetworkFirewallPolicy struct {
	BaseAwsResource
	Client      networkfirewalliface.NetworkFirewallAPI
	Region      string
	Identifiers []string
}

func (nfw *NetworkFirewallPolicy) Init(session *session.Session) {
	nfw.BaseAwsResource.Init(session)
	nfw.Client = networkfirewall.New(session)
}

// ResourceName - the simple name of the aws resource
func (nfw *NetworkFirewallPolicy) ResourceName() string {
	return "network-firewall-policy"
}

func (nfw *NetworkFirewallPolicy) ResourceIdentifiers() []string {
	return nfw.Identifiers
}

func (nfw *NetworkFirewallPolicy) MaxBatchSize() int {
	// Tentative batch size to ensure AWS doesn't throttle. Note that nat gateway does not support bulk delete, so
	// we will be deleting this many in parallel using go routines. We conservatively pick 10 here, both to limit
	// overloading the runtime and to avoid AWS throttling with many API calls.
	return 10
}

func (nfw *NetworkFirewallPolicy) GetAndSetIdentifiers(c context.Context, configObj config.Config) ([]string, error) {
	identifiers, err := nfw.getAll(c, configObj)
	if err != nil {
		return nil, err
	}

	nfw.Identifiers = awsgo.StringValueSlice(identifiers)
	return nfw.Identifiers, nil
}

// Nuke - nuke 'em all!!!
func (nfw *NetworkFirewallPolicy) Nuke(identifiers []string) error {
	if err := nfw.nukeAll(awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
