// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package nic

import (
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
)

const (
	IFLA_VLAN_UNSPEC = iota
	IFLA_VLAN_ID
	IFLA_VLAN_FLAGS
	IFLA_VLAN_EGRESS_QOS
	IFLA_VLAN_INGRESS_QOS
	IFLA_VLAN_PROTOCOL
	IFLA_VLAN_MAX = IFLA_VLAN_PROTOCOL
)

// VlanProtocol possible values
const (
	VLAN_PROTOCOL_UNKNOWN = 0
	VLAN_PROTOCOL_8021Q   = 0x8100
	VLAN_PROTOCOL_8021AD  = 0x88A8
)

type Vlan struct {
	Parent       string
	Id           uint16
	CIDRs        []string
	Link         *net.Interface
	VlanSettings *netlink.AttributeEncoder
}

// WithParent sets the parent index device
func (v *Vlan) WithParent(index uint32) {
}

// WithCIDR defines if the interface have static CIDRs added
func WithVlan(id uint16) Option {
	return func(n *NetworkInterface) (err error) {

		for _, vlan := range n.Vlans {
			if vlan.Id == id {
				return fmt.Errorf("Duplicate VLAN id  %v given", vlan)
			}
		}
		vlan := &Vlan{
			Id:           id,
			VlanSettings: netlink.NewAttributeEncoder(),
		}
		vlan.VlanSettings.Uint16(uint16(IFLA_VLAN_ID), uint16(vlan.Id))
		vlan.VlanSettings.Uint16(uint16(IFLA_VLAN_PROTOCOL), VLAN_PROTOCOL_8021Q)
		n.Vlans = append(n.Vlans, vlan)
		return nil
	}
}

// WithVlanCIDR defines if the interface have static CIDRs added
func WithVlanCIDR(id uint16, cidr string) Option {
	return func(n *NetworkInterface) (err error) {

		for _, vlan := range n.Vlans {
			if vlan.Id == id {
				vlan.CIDRs = append(vlan.CIDRs, cidr)
				return nil
			}
		}
		return fmt.Errorf("VLAN id not found for CIDR setting  %v given", id)
	}

}
