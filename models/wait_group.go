package models

import "sync"

var ChartWG = &sync.WaitGroup{}
var HistoryWG = &sync.WaitGroup{}
var RealtimeWG = &sync.WaitGroup{}
