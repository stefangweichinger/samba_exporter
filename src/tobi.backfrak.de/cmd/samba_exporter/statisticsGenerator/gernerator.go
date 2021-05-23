package statisticsGenerator

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"tobi.backfrak.de/cmd/samba_exporter/smbstatusreader"
)

// Type for numeric statistic values from the samba server
type SmbStatisticsNumeric struct {
	Name   string
	Value  float64
	Help   string
	Labels map[string]string
}

// GetSmbStatistics - Get the statistic data for prometheus out of the response data arrays
func GetSmbStatistics(lockData []smbstatusreader.LockData, processData []smbstatusreader.ProcessData, shareData []smbstatusreader.ShareData) []SmbStatisticsNumeric {
	var ret []SmbStatisticsNumeric

	var users []int
	var pids []int
	var shares []string
	var clients []string
	locksPerShare := make(map[string]int, 0)

	for _, lock := range lockData {
		if !intArrContains(users, lock.UserID) {
			users = append(users, lock.UserID)
		}

		if !intArrContains(pids, lock.PID) {
			pids = append(pids, lock.PID)
		}

		locksOfShare, found := locksPerShare[lock.SharePath]
		if found == false {
			locksPerShare[lock.SharePath] = 1
		} else {
			locksPerShare[lock.SharePath] = locksOfShare + 1
		}
	}

	for _, process := range processData {
		if !intArrContains(users, process.UserID) {
			users = append(users, process.UserID)
		}

		if !intArrContains(pids, process.PID) {
			pids = append(pids, process.PID)
		}
	}

	for _, share := range shareData {
		if !intArrContains(pids, share.PID) {
			pids = append(pids, share.PID)
		}

		if !strArrContains(shares, share.Service) {
			shares = append(shares, share.Service)
		}

		if !strArrContains(clients, share.Machine) {
			clients = append(clients, share.Machine)
		}
	}

	ret = append(ret, SmbStatisticsNumeric{"individual_user_count", float64(len(users)), "The number of users connected to this samba server", nil})
	ret = append(ret, SmbStatisticsNumeric{"locked_file_count", float64(len(lockData)), "Number of files locked by the samba server", nil})
	ret = append(ret, SmbStatisticsNumeric{"pid_count", float64(len(pids)), "Number of processes running by the samba server", nil})
	ret = append(ret, SmbStatisticsNumeric{"share_count", float64(len(shares)), "Number of shares used by clients of the samba server", nil})
	ret = append(ret, SmbStatisticsNumeric{"client_count", float64(len(clients)), "Number of clients using the samba server", nil})

	if len(locksPerShare) > 0 {
		for share, locks := range locksPerShare {
			ret = append(ret, SmbStatisticsNumeric{"locks_per_share", float64(locks), "Number of locks on share", map[string]string{"share": share}})
		}
	} else {
		// Add this value even if no locks found, so prometheus description will be created
		ret = append(ret, SmbStatisticsNumeric{"locks_per_share", float64(0), "Number of locks on share", map[string]string{"share": ""}})
	}

	return ret
}

func intArrContains(arr []int, value int) bool {
	for _, field := range arr {
		if field == value {
			return true
		}
	}

	return false
}

func strArrContains(arr []string, value string) bool {
	for _, field := range arr {
		if field == value {
			return true
		}
	}

	return false
}