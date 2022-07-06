package main

import (
	"log"
)

type Logger struct {
	reset   string
	message string
}

func Log() *Logger {
	return &Logger{
		reset:   "\033[0m",
		message: "[%c] %-10sheight=%012d hash=%016X...%016X mempool=%06d txs=%d conn=%d",
	}
}

func (lg *Logger) Warning(name string, height uint64, hash []byte, mempool uint64, txs int, conns int) {
	colorYellow := "\033[33m"
	log.Printf(colorYellow+lg.message+lg.reset,
		'W', name, height, hash[:8], hash[24:], mempool, txs, conns)
}

func (lg *Logger) Error(name string, height uint64, mempool uint64, txs int, conns int) {
	colorRed := "\033[31m"
	log.Printf(colorRed+lg.message+lg.reset,
		'E', name, height, []byte{0}, []byte{0}, mempool, txs, conns)
}

func (lg *Logger) Info(name string, height uint64, hash []byte, mempool uint64, txs int, conns int) {
	log.Printf(lg.message+lg.reset,
		'I', name, height, hash[:8], hash[24:], mempool, txs, conns)
}
