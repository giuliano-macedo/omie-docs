package crawler

import (
	"fmt"
	"sync"

	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/internal/docs_fetcher"
)

type getPageWorkerTask struct {
	index int
	url   string
}
type getPageWorker struct {
	tasks    []getPageWorkerTask
	initial  int
	step     int
	pages    *[]core.Page
	entities []core.Entity
	fetcher  docs_fetcher.DocsFetcher
	group    *sync.WaitGroup
	errch    chan error
}

func (worker getPageWorker) run() {
	defer worker.group.Done()
	for i := worker.initial; i < len(worker.tasks); i += worker.step {
		task := worker.tasks[i]
		fmt.Println("crawling", task.url)
		page, err := GetPage(worker.fetcher, task.url)
		if err != nil {
			worker.errch <- fmt.Errorf("failed getting %v: %v", task.url, err)
			return
		}
		page.EntityName = worker.entities[task.index].Name
		(*worker.pages)[task.index] = page
	}
	worker.errch <- nil
}

func GetPages(fetcher docs_fetcher.DocsFetcher, home core.HomePage, numberOfWorkers int) ([]core.Page, error) {
	pagesSize := 0
	for _, feature := range home.Features {
		pagesSize += len(feature.MainEntities) + len(feature.AuxiliaryEntites)
	}

	pages := make([]core.Page, pagesSize)
	tasks := make([]getPageWorkerTask, 0, pagesSize)
	entities := make([]core.Entity, 0, pagesSize)
	index := 0
	for _, feature := range home.Features {
		for _, entity := range feature.AllEntities() {
			tasks = append(tasks, getPageWorkerTask{url: entity.Url, index: index})
			entities = append(entities, entity)
			index++
		}
	}
	fmt.Println("using", numberOfWorkers, "workers")

	var group sync.WaitGroup
	group.Add(numberOfWorkers)
	errch := make(chan error, numberOfWorkers)

	for i := 0; i < numberOfWorkers; i++ {
		worker := getPageWorker{
			tasks:    tasks,
			initial:  i,
			step:     numberOfWorkers,
			group:    &group,
			pages:    &pages,
			fetcher:  fetcher,
			errch:    errch,
			entities: entities,
		}
		go worker.run()
	}
	group.Wait()

	for i := 0; i < numberOfWorkers; i++ {
		if err := <-errch; err != nil {
			return pages, err
		}
	}

	return pages, nil
}
