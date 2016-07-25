package main

import (
	"flag"
	"runtime"
	"strings"

	"github.com/RedisLabs/RediSearchBenchmark/index"
	"github.com/RedisLabs/RediSearchBenchmark/index/redisearch"
	"github.com/RedisLabs/RediSearchBenchmark/ingest"
)

const IndexName = "wik"

var indexMetadata = index.NewMetadata().
	AddField(index.NewTextField("title", 10)).
	AddField(index.NewTextField("body", 1))

func selectIndex(engine string, hosts []string, partitions int) index.Index {

	switch engine {
	case "redis":
		return redisearch.NewDistributedIndex(IndexName, hosts, partitions, indexMetadata)

	}
	panic("could not find index type " + engine)
}

func main() {

	hosts := flag.String("hosts", "localhost:6379", "comma separated list of host:port to redis nodes")
	partitions := flag.Int("shards", 1, "the number of partitions we want (AT LEAST the number of cluster shards)")
	fileName := flag.String("file", "", "Input file to ingest data from (wikipedia abstracts)")
	engine := flag.String("engine", "redis", "The search backend to run")
	//benchmark := flag.Bool("benchmark", false, "if set, we run a benchmark")
	//conc := flag.Int("c", 4, "benchmark concurrency")
	//qs := flag.String("queries", "hello world", "comma separated list of queries to benchmark")

	flag.Parse()
	servers := strings.Split(*hosts, ",")
	if len(servers) == 0 {
		panic("No servers given")
	}
	//	queries := strings.Split(*qs, ",")

	runtime.GOMAXPROCS(runtime.NumCPU() / 2)

	// select index to run
	idx := selectIndex(*engine, servers, *partitions)

	if *fileName != "" {
		idx.Drop()
		idx.Create()
		if err := ingest.IngestDocuments(*fileName, ingest.ReadWikipediaExtracts, idx, nil, 10000); err != nil {
			panic(err)
		}
	}
	//	//LoadWikipedia(*fileName, modIdx, store)
	//	if *benchmark {

	//		f, err := os.Create("prof.out")
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		pprof.StartCPUProfile(f)
	//		defer pprof.StopCPUProfile()

	//		Benchmark(queries, *conc, modIdx, store)
	//	}

}