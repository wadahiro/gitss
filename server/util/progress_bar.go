package util

import (
	pb "gopkg.in/cheggaaa/pb.v1"
	"sync"
)

type ProgressBar struct {
	bar  *pb.ProgressBar
	lock sync.Mutex
}

func NewProgressBar(total int) *ProgressBar {
	bar := &ProgressBar{bar: pb.StartNew(total)}
	bar.bar.ShowPercent = true

	return bar
}

func (p *ProgressBar) SetTotal(total int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.bar.Total = int64(total)
}

func (p *ProgressBar) AddTotal(addTotal int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.bar.Total += int64(addTotal)
}

func (p *ProgressBar) Add(count int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.bar.Add(count)
}

func (p *ProgressBar) FinishPrint(text string) {
	p.bar.FinishPrint(text)
}
