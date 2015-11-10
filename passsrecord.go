package main

import (
	"fmt"

	"github.com/atotto/clipboard"
)

type PasssRecord struct {
	Website  string
	Username string
	Password string
}

func (pr *PasssRecord) PrintHead() {
	fmt.Printf("%16v @ %v\n", pr.Username, pr.Website)
}

func (pr *PasssRecord) LoadPassword() {
	clipboard.WriteAll(pr.Password)
}

type RecordsSlice []PasssRecord

func (rs RecordsSlice) Len() int {
	return len(rs)
}

func (rs RecordsSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs RecordsSlice) Less(i, j int) bool {
	if rs[i].Website != rs[j].Website {
		return rs[i].Website < rs[j].Website
	}
	return rs[i].Username < rs[j].Username
}
