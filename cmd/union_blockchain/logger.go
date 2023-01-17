package main

import (
	"log"
)

type SLogger struct {
	fReset   string
	fMessage string
}

func Log() *SLogger {
	return &SLogger{
		fReset:   "\033[0m",
		fMessage: "[%c] %-10sheight=%012d hash=%016X...%016X mempool=%06d txs=%d conn=%d",
	}
}

func (lg *SLogger) Warning(name string, height uint64, hash []byte, mempool uint64, txs int, conns int) {
	colorYellow := "\033[33m"
	log.Printf(colorYellow+lg.fMessage+lg.fReset,
		'W', name, height, hash[:8], hash[24:], mempool, txs, conns)
}

func (lg *SLogger) Error(name string, height uint64, mempool uint64, txs int, conns int) {
	colorRed := "\033[31m"
	log.Printf(colorRed+lg.fMessage+lg.fReset,
		'E', name, height, []byte{0}, []byte{0}, mempool, txs, conns)
}

func (lg *SLogger) Info(name string, height uint64, hash []byte, mempool uint64, txs int, conns int) {
	log.Printf(lg.fMessage+lg.fReset,
		'I', name, height, hash[:8], hash[24:], mempool, txs, conns)
}
