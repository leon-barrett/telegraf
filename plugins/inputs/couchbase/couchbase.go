package couchbase

import (
	"regexp"
	"strconv"
	"sync"

	couchbase "github.com/couchbase/go-couchbase"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Couchbase struct {
	Servers         []string
	BucketNodeStats []string
}

var sampleConfig = `
  ## specify servers via a url matching:
  ##  [protocol://][:password]@address[:port]
  ##  e.g.
  ##    http://couchbase-0.example.com/
  ##    http://admin:secret@couchbase-0.example.com:8091/
  ##
  ## If no servers are specified, then localhost is used as the host.
  ## If no protocol is specifed, HTTP is used.
  ## If no port is specified, 8091 is used.
  servers = ["http://localhost:8091"]
  ## Detailed per-bucket-per-node stats may be collected also, although they
  ## are not included by default since these detailed stats entail an extra
  ## request for each bucket.
  ## The default is to collect none.
  # BucketNodeStats = []
  ## Specific per-bucket-per-node stats may be collected.
  # BucketNodeStats = ["curr_items", "ep_diskqueue_items", "ep_bg_fetched"]
  ## All per-bucket-per-node stats may be collected.
  # BucketNodeStats = ["*"]
`

var regexpURI = regexp.MustCompile(`(\S+://)?(\S+\:\S+@)`)

func (r *Couchbase) SampleConfig() string {
	return sampleConfig
}

func (r *Couchbase) Description() string {
	return "Read metrics from one or many couchbase clusters"
}

// Reads stats from all configured clusters. Accumulates stats.
// Returns one of the errors encountered while gathering stats (if any).
func (r *Couchbase) Gather(acc telegraf.Accumulator) error {
	if len(r.Servers) == 0 {
		r.gatherServer("http://localhost:8091/", acc, nil)
		return nil
	}

	var wg sync.WaitGroup

	for _, serv := range r.Servers {
		wg.Add(1)
		go func(serv string) {
			defer wg.Done()
			acc.AddError(r.gatherServer(serv, acc, nil))
		}(serv)
	}

	wg.Wait()

	return nil
}

func (r *Couchbase) gatherServer(addr string, acc telegraf.Accumulator, pool *couchbase.Pool) error {
	if pool == nil {
		client, err := couchbase.Connect(addr)
		if err != nil {
			return err
		}

		// `default` is the only possible pool name. It's a
		// placeholder for a possible future Couchbase feature. See
		// http://stackoverflow.com/a/16990911/17498.
		p, err := client.GetPool("default")
		if err != nil {
			return err
		}
		pool = &p
	}

	for i := 0; i < len(pool.Nodes); i++ {
		node := pool.Nodes[i]
		tags := map[string]string{"cluster": regexpURI.ReplaceAllString(addr, "${1}"), "hostname": node.Hostname}
		fields := make(map[string]interface{})
		fields["memory_free"] = node.MemoryFree
		fields["memory_total"] = node.MemoryTotal
		acc.AddFields("couchbase_node", fields, tags)
	}

	for bucketName := range pool.BucketMap {
		tags := map[string]string{"cluster": addr, "bucket": bucketName}
		bs := pool.BucketMap[bucketName].BasicStats
		fields := make(map[string]interface{})
		fields["quota_percent_used"] = bs["quotaPercentUsed"]
		fields["ops_per_sec"] = bs["opsPerSec"]
		fields["disk_fetches"] = bs["diskFetches"]
		fields["item_count"] = bs["itemCount"]
		fields["disk_used"] = bs["diskUsed"]
		fields["data_used"] = bs["dataUsed"]
		fields["mem_used"] = bs["memUsed"]
		acc.AddFields("couchbase_bucket", fields, tags)
	}

	if len(r.BucketNodeStats) > 0 {
		includedStats := map[string]bool{}
		for _, k := range r.BucketNodeStats {
			includedStats[k] = true
		}
		for bucketName, bucket := range pool.BucketMap {
			// We must refresh the bucket to get the list of nodes. Plus, if we
			// don't refresh the bucket, GetStats often segfaults.
			bucket.Refresh()
			allStats := bucket.GetStats("")
			for nodeName, nodeStats := range allStats {
				tags := map[string]string{"cluster": addr, "bucket": bucketName, "hostname": nodeName}
				fields := make(map[string]interface{})
				for k, v := range nodeStats {
					if includedStats["*"] || includedStats[k] {
						f, err := strconv.ParseFloat(v, 64)
						if err == nil {
							fields[k] = f
						}
					}
				}
				acc.AddFields("couchbase_bucket_node", fields, tags)
			}
		}
	}
	return nil
}

func init() {
	inputs.Add("couchbase", func() telegraf.Input {
		return &Couchbase{}
	})
}
