package models

import (
	"go-vksave/utils"
	"sync"
)

type Worker struct {
	urls    chan string     // канал для заданий
	pending int             // кол-во оставшихся задач
	index   int             // позиция в куче
	wg      *sync.WaitGroup // указатель на группу ожидания
}

func (w *Worker) work(done chan *Worker) {
	for {
		url := <-w.urls
		w.wg.Add(1)
		utils.Download(url)
		w.wg.Done()
		done <- w
	}
}

// реализация "кучи" для простого и удобного
// получения наименее загруженного рабочего
type Pool []*Worker

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Len() int {
	return len(p)
}

func (p Pool) Swap(i, j int) {
	if i >= 0 && i < len(p) && j >= 0 && j < len(p) {
		p[i], p[j] = p[j], p[i]
		p[i].index, p[j].index = i, j
	}
}

func (p *Pool) Push(x interface{}) {
	n := len(*p)
	worker := x.(*Worker)
	worker.index = n
	*p = append(*p, worker)
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	item.index = -1
	*p = old[0 : n-1]
	return item
}
