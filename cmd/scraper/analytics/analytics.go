package analytics

type ScrapeResult struct {
	Failed    bool
	NbDeleted int
	NbFound   int
	Retries   int
}
