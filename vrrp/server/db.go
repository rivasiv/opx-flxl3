//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package server

import (
	"l3/vrrp/config"
	"l3/vrrp/debug"
	"models/objects"
	"strings"
)

func (svr *VrrpServer) readVrrpGblCfg() {
	debug.Logger.Info("Reading Vrrp Global Config from DB")
	var dbObj objects.VrrpGlobal
	objList, err := svr.dmnBase.DbHdl.GetAllObjFromDb(dbObj)
	if err != nil {
		debug.Logger.Warning("Vrrp Global DB read returned:", err)
		return
	}
	debug.Logger.Info("Global Objects reterived from DB are", objList)
	for _, obj := range objList {
		cfg := obj.(objects.VrrpGlobal)
		gblCfg := &config.GlobalConfig{
			Vrf:       cfg.Vrf,
			Enable:    cfg.Enable,
			Operation: config.UPDATE,
		}
		svr.HandleGlobalConfig(gblCfg)
	}
	debug.Logger.Info("Global Config Read is done")
}

func (svr *VrrpServer) readVrrpV4IntfCfg() {
	debug.Logger.Info("Reading Vrrp V4 Intf Config from DB")
	var dbObj objects.VrrpV4Intf
	objList, err := svr.dmnBase.DbHdl.GetAllObjFromDb(dbObj)
	if err != nil {
		debug.Logger.Warning("Vrrp v4 Interface read returned:", err)
		return
	}
	debug.Logger.Info("Vrrp V4 Intf Objects reterived from DB are", objList)
	for _, obj := range objList {
		cfg := obj.(objects.VrrpV4Intf)
		if !strings.Contains(cfg.Address, config.NETMASK_DELIMITER) {
			cfg.Address += config.NETMASK_DELIMITER + config.SLASH_32
		}
		v4Cfg := &config.IntfCfg{
			IntfRef:               cfg.IntfRef,
			VRID:                  cfg.VRID,
			Priority:              cfg.Priority,
			VirtualIPAddr:         cfg.Address,
			AdvertisementInterval: cfg.AdvertisementInterval,
			PreemptMode:           cfg.PreemptMode,
			AcceptMode:            cfg.AcceptMode,
			AdminState:            cfg.AdminState,
			Version:               config.VERSION2,
			Operation:             config.CREATE,
		}
		svr.HandleVrrpIntfConfig(v4Cfg)
	}
	debug.Logger.Info("Vrrp v4 interface Config Read is done")
}

func (svr *VrrpServer) readVrrpV6IntfCfg() {
	debug.Logger.Info("Reading Vrrp V6 Intf Config from DB")
	var dbObj objects.VrrpV6Intf
	objList, err := svr.dmnBase.DbHdl.GetAllObjFromDb(dbObj)
	if err != nil {
		debug.Logger.Warning("Vrrp v4 Interface read returned:", err)
		return
	}
	debug.Logger.Info("Vrrp V4 Intf Objects reterived from DB are", objList)
	for _, obj := range objList {
		cfg := obj.(objects.VrrpV6Intf)
		v6Cfg := &config.IntfCfg{
			IntfRef:               cfg.IntfRef,
			VRID:                  cfg.VRID,
			Priority:              cfg.Priority,
			VirtualIPAddr:         cfg.Address,
			AdvertisementInterval: cfg.AdvertisementInterval,
			PreemptMode:           cfg.PreemptMode,
			AcceptMode:            cfg.AcceptMode,
			AdminState:            cfg.AdminState,
			Version:               config.VERSION3,
			Operation:             config.CREATE,
		}
		svr.HandleVrrpIntfConfig(v6Cfg)
	}
	debug.Logger.Info("Vrrp v6 interface Config Read is done")
}

func (svr *VrrpServer) ReadDB() {
	if svr.dmnBase == nil {
		return
	}
	if svr.dmnBase.DbHdl == nil {
		debug.Logger.Err("DB Handler is nil and hence cannot read anything from DATABASE")
		return
	}
	debug.Logger.Info("Reading Config from DB")
	svr.readVrrpGblCfg()
	svr.readVrrpV4IntfCfg()
	svr.readVrrpV6IntfCfg()
}
