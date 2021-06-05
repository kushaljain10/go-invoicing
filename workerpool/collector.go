package workerpool

var WorkQueue = make(chan WorkRequest, 100)

func Collector(work WorkRequest) {
	WorkQueue <- work
}
