// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package nic

import (
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
	"github.com/talos-systems/talos/internal/app/networkd/pkg/address"
	"github.com/talos-systems/talos/pkg/config/machine"
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
	Parent        string
	Id            uint16
	Link          *net.Interface
	VlanSettings  *netlink.AttributeEncoder
	AddressMethod []address.Addressing
}

// WithParent sets the parent index device
func (v *Vlan) WithParent(index uint32) {
	//v.Uint32(uint16(unix.IFLA_LINK), index)
}

// WithVlan defines the VLAN id to use
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
		n.Vlans = append(n.Vlans, vlan)
		return nil
	}
}

func WithVlanDhcp(id uint16) Option {
	return func(n *NetworkInterface) (err error) {
		for _, vlan := range n.Vlans {
			if vlan.Id == id {
				vlan.AddressMethod = append(vlan.AddressMethod, &address.DHCP{})
				return nil
			}
		}
		return fmt.Errorf("VLAN id not found for DHCP. Vlan ID  %v given", id)
	}

}

// WithVlanCIDR defines if the interface have static CIDRs added
func WithVlanCIDR(id uint16, cidr string, routeList []machine.Route) Option {
	return func(n *NetworkInterface) (err error) {

		for _, vlan := range n.Vlans {
			if vlan.Id == id {
				vlan.AddressMethod = append(vlan.AddressMethod, &address.Static{CIDR: cidr, RouteList: routeList})
				return nil
			}
		}
		return fmt.Errorf("VLAN id not found for CIDR setting  %v given", id)
	}

}
