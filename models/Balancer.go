package models

import (
	"container/heap"
	"sync"
)

const ENDMESSAGE = "LastURLs"

var (
	WORKERS    = 10 //количество рабочих
	WORKERSCAP = 10 //размер очереди каждого рабочего
)

type Balancer struct {
	pool     Pool            // "куча" рабочих
	done     chan *Worker    // Канал уведомления для рабочих
	requests chan string     // Канал для получения новых заданий
	flowctrl chan bool       // Канал для PMFC
	queue    int             // Количество незавершенных заданий переданных рабочим
	wg       *sync.WaitGroup // Группа ожидания для рабочих
}

func (b *Balancer) Init(in chan string) {
	b.requests = make(chan string)
	b.flowctrl = make(chan bool)
	b.done = make(chan *Worker)
	b.wg = new(sync.WaitGroup)

	go func() {
		for {
			b.requests <- <-in // получаем новое задание и пересылаем его на внутренний канал
			<-b.flowctrl       // ждем получения подтверждения
		}
	}()

	heap.Init(&b.pool)
	for i := 0; i < WORKERS; i++ {
		w := &Worker{
			urls:    make(chan string, WORKERSCAP),
			pending: 0,
			index:   0,
			wg:      b.wg,
		}
		go w.work(b.done)
		heap.Push(&b.pool, w)
	}
}

func (b *Balancer) Balance(quit chan bool) {
	lastjobs := false
	for {
		select {
		case <-quit:
			b.wg.Wait()
			quit <- true
		case url := <-b.requests:
			if url != ENDMESSAGE {
				b.dispatch(url)
			} else {
				lastjobs = true
			}
		case w := <-b.done:
			b.completed(w)
			if lastjobs {
				if w.pending == 0 {
					heap.Remove(&b.pool, w.index)
				}
				if len(b.pool) == 0 {
					quit <- true
				}
			}
		}
	}
}

func (b *Balancer) dispatch(url string) {
	w := heap.Pop(&b.pool).(*Worker)
	w.urls <- url
	w.pending++
	heap.Push(&b.pool, w)
	if b.queue++; b.queue < WORKERS*WORKERSCAP {
		b.flowctrl <- true
	}
}

func (b *Balancer) completed(w *Worker) {
	w.pending--
	heap.Remove(&b.pool, w.index)
	heap.Push(&b.pool, w)
	if b.queue--; b.queue == WORKERS*WORKERSCAP-1 {
		b.flowctrl <- true
	}
}
