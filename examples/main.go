/**
 * Copyright (c) 2017 eBay Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 **/

package main

import (
	"fmt"

	"github.com/ebay/go-ovn/goovn"
	"os"
)

const (
	OVS_RUNDIR = "/var/run/openvswitch"
	OVNNB_SOCKET = "ovnnb_db.sock"
	MATCH = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.61"

	MATCH_SECOND = "outport == \"96d44061-1823-428b-a7ce-f473d10eb3d0\" && ip && ip.dst == 10.97.183.62"

)

var ovndbapi goovn.OVNDBApi

func init() {
	var ovs_rundir = os.Getenv("OVS_RUNDIR")
	if ovs_rundir == "" {
		ovs_rundir = OVS_RUNDIR
	}
	var socket = ovs_rundir + "/" + OVNNB_SOCKET
	ovndbapi = goovn.GetInstance(socket, goovn.UNIX, "", 0, nil)
}

func main() {
	ocmd := ovndbapi.LSWAdd("ls1")
	ovndbapi.Execute(ocmd)
	ocmd = ovndbapi.LSPAdd("ls1", "test")
	ovndbapi.Execute(ocmd)
	ocmd = ovndbapi.LSPSetAddress("test", "12:34:56:78:90 10.10.10.1")
	ovndbapi.Execute(ocmd)

	lports := ovndbapi.GetLogicPortsBySwitch("ls1")
	for _, lp := range(lports) {
		fmt.Printf("%v\n", *lp)
	}

	ocmd = ovndbapi.ACLAdd("ls1", "to-lport", MATCH, "drop", 1001, nil, true, "")
	ovndbapi.Execute(ocmd)

	ocmd = ovndbapi.ACLAdd("ls1", "to-lport", MATCH_SECOND, "drop", 1001, map[string]string{"A": "a", "B": "b"}, false, "")
	ovndbapi.Execute(ocmd)


	ocmd = ovndbapi.ACLAdd("ls1", "to-lport", MATCH_SECOND, "drop", 1001, map[string]string{"A": "b", "B": "b"}, false, "")
	ovndbapi.Execute(ocmd)

	acls := ovndbapi.GetACLsBySwitch("ls1")
	for _, acl := range(acls) {
		fmt.Printf("%v\n", *acl)
	}
	fmt.Println()

	ocmd = ovndbapi.ACLDel("ls1", "to-lport", MATCH, 1001,  map[string]string{})
	ovndbapi.Execute(ocmd)
	acls = ovndbapi.GetACLsBySwitch("ls1")
	for _, acl := range(acls) {
		fmt.Printf("%v\n", *acl)
	}

	fmt.Println()
	ocmd = ovndbapi.ACLDel("ls1", "to-lport", MATCH_SECOND, 1001, map[string]string{"A": "a"})
	ovndbapi.Execute(ocmd)
	acls = ovndbapi.GetACLsBySwitch("ls1")
	for _, acl := range(acls) {
		fmt.Printf("%v\n", *acl)
	}
	fmt.Println()
	ocmd = ovndbapi.ACLDel("ls1", "to-lport", MATCH_SECOND, 1001, map[string]string{"A": "b"})
	ovndbapi.Execute(ocmd)
	acls = ovndbapi.GetACLsBySwitch("ls1")
	for _, acl := range(acls) {
		fmt.Printf("%v\n", *acl)
	}
	ocmd = ovndbapi.LSPDel("test")
	ovndbapi.Execute(ocmd)
	ocmd = ovndbapi.LSWDel("ls1")
	ovndbapi.Execute(ocmd)

	fmt.Printf("\nAdding LB to OVN")
	ocmd = ovndbapi.LBAdd("lb1", "192.168.0.19:80", "tcp", []string{"10.0.0.11:80","10.0.0.12:80"})
	err := ovndbapi.Execute(ocmd)
	if err != nil {
		fmt.Printf("\nAdding LB OVN failed with err %v", err)
		return
	}
	fmt.Printf("\nAdding LB to OVN Done")


	fmt.Printf("\nUpdating LB to OVN")
	ocmd = ovndbapi.LBUpdate("lb1", "192.168.0.10:80", "tcp", []string{"10.10.10.127:8080","10.10.10.120:8080"})
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		fmt.Printf("\nUpdating LB OVN failed with err %v", err)
		return
	}
	fmt.Printf("\nUpdating LB to OVN done")

	fmt.Printf("\nGettting LB by name")
	lb := ovndbapi.GetLB("lb1")
	if len(lb) != 1 {
		fmt.Printf("err getting lbs, total:%v", len(lb))
	}
	fmt.Printf("\n Lb found:%+v", lb[0])

	fmt.Printf("\nDeleting LB")
	ocmd = ovndbapi.LBDel("lb1")
	err = ovndbapi.Execute(ocmd)
	if err != nil {
		fmt.Printf("err executing command:%v", err)
		return
	}
	fmt.Printf("\nDeleting LB Done")
}
