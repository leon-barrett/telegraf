# Telegraf Plugin: Couchbase

## Configuration:

```
# Read per-node and per-bucket metrics from Couchbase
[[inputs.couchbase]]
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
```

## Measurements:

### couchbase_node

Tags:
- cluster: sanitized string from `servers` configuration field e.g.: `http://user:password@couchbase-0.example.com:8091/endpoint` -> `http://couchbase-0.example.com:8091/endpoint`
- hostname: Couchbase's name for the node and port, e.g., `172.16.10.187:8091`

Fields:
- memory_free (unit: bytes, example: 23181365248.0)
- memory_total (unit: bytes, example: 64424656896.0)

### couchbase_bucket

Tags:
- cluster: whatever you called it in `servers` in the configuration, e.g.: `http://couchbase-0.example.com/`)
- bucket: the name of the couchbase bucket, e.g., `blastro-df`

Fields:
- quota_percent_used (unit: percent, example: 68.85424936294555)
- ops_per_sec (unit: count, example: 5686.789686789687)
- disk_fetches (unit: count, example: 0.0)
- item_count (unit: count, example: 943239752.0)
- disk_used (unit: bytes, example: 409178772321.0)
- data_used (unit: bytes, example: 212179309111.0)
- mem_used (unit: bytes, example: 202156957464.0)

### couchbase_bucket_node

Tags:
- cluster: sanitized string from `servers` configuration field e.g.: `http://user:password@couchbase-0.example.com:8091/endpoint` -> `http://couchbase-0.example.com:8091/endpoint`
- bucket: the name of the couchbase bucket, e.g., `blastro-df`
- hostname: Couchbase's name for the node and port, e.g., `172.16.10.187:11210`

Fields:
- Many fields, not specified by couchbase documentation. Examples include:
- bytes
- ep_tap_ack_initial_sequence_number
- ep_item_begin_failed
- ep_num_non_resident
- vb_pending_ht_memory
- vb_active_ht_memory
- ep_diskqueue_pending
- ep_num_access_scanner_runs
- ep_alog_task_time
- vb_pending_itm_memory
- ...

## Example output

```
$ telegraf --config telegraf.conf --input-filter couchbase --test
* Plugin: couchbase, Collection 1
> couchbase_node,cluster=https://couchbase-0.example.com/,hostname=172.16.10.187:8091 memory_free=22927384576,memory_total=64424656896 1458381183695864929
> couchbase_node,cluster=https://couchbase-0.example.com/,hostname=172.16.10.65:8091 memory_free=23520161792,memory_total=64424656896 1458381183695972112
> couchbase_node,cluster=https://couchbase-0.example.com/,hostname=172.16.13.105:8091 memory_free=23531704320,memory_total=64424656896 1458381183695995259
> couchbase_node,cluster=https://couchbase-0.example.com/,hostname=172.16.13.173:8091 memory_free=23628767232,memory_total=64424656896 1458381183696010870
> couchbase_node,cluster=https://couchbase-0.example.com/,hostname=172.16.15.120:8091 memory_free=23616692224,memory_total=64424656896 1458381183696027406
> couchbase_node,cluster=https://couchbase-0.example.com/,hostname=172.16.8.127:8091 memory_free=23431770112,memory_total=64424656896 1458381183696041040
> couchbase_node,cluster=https://couchbase-0.example.com/,hostname=172.16.8.148:8091 memory_free=23811371008,memory_total=64424656896 1458381183696059060
> couchbase_bucket,bucket=default,cluster=https://couchbase-0.example.com/ data_used=25743360,disk_fetches=0,disk_used=31744886,item_count=0,mem_used=77729224,ops_per_sec=0,quota_percent_used=10.58976636614118 1458381183696210074
> couchbase_bucket,bucket=demoncat,cluster=https://couchbase-0.example.com/ data_used=38157584951,disk_fetches=0,disk_used=62730302441,item_count=14662532,mem_used=24015304256,ops_per_sec=1207.753207753208,quota_percent_used=79.87855353525707 1458381183696242695
> couchbase_bucket,bucket=blastro-df,cluster=https://couchbase-0.example.com/ data_used=212552491622,disk_fetches=0,disk_used=413323157621,item_count=944655680,mem_used=202421103760,ops_per_sec=1692.176692176692,quota_percent_used=68.9442170551845 1458381183696272206
> couchbase_bucket_node,cluster=https://couchbase-0.example.com/,bucket=default,hostname=172.16.8.127:11210,host=localhost bytes=10890985784,ep_tap_ack_initial_sequence_number=1,ep_item_begin_failed=0,ep_num_non_resident=15594899,vb_pending_ht_memory=0,vb_active_ht_memory=75385008,ep_diskqueue_pending=2118656,ep_num_access_scanner_runs=208,ep_alog_task_time=2,vb_pending_itm_memory=0,vb_pending_perc_mem_resident=100,ep_getl_max_timeout=30,ep_max_threads=0,ep_warmup_oom=0,curr_conns_on_port_11209=43,ep_dcp_conn_buffer_size=10485760,msgused_high_watermark=2,ep_bg_wait_avg=411,ep_dcp_conn_buffer_size_max=52428800,ep_dcp_enable_flow_control=1,curr_items_tot=24019509,ep_diskqueue_fill=471704508,vb_pending_ops_create=0,ep_storedval_num=24019509,ep_item_num_based_new_chk=1,ep_tap_throttle_threshold=99,ep_uncommitted_items=0,ep_item_num=10185,ep_total_new_items=67167069,vb_pending_expired=0,threads=6,ep_dcp_backfill_byte_limit=20971832,ep_num_ops_get_meta=0,ep_exp_pager_stime=3600,ep_diskqueue_drain=471704434,ep_data_traffic_enabled=0,ep_warmup_time=422,ep_keep_closed_chks=0,ep_num_ops_set_meta_res_fail=0,ep_warmup=1,ep_waitforwarmup=0,ep_dcp_max_unacked_bytes=524288,ep_commit_time=11,rbufs_loaned=550669024,ep_max_size=12884901888,ep_bg_max_load=29897603,ep_pending_ops_max_duration=0,ep_max_bg_remaining_jobs=0,uptime=4511994,ep_vb_total=684,vb_active_queue_fill=234903233,vb_pending_queue_pending=0,ep_diskqueue_memory=5328,ep_bg_meta_fetched=0,ep_warmup_dups=0,rbufs_allocated=16620729,delete_hits=0,ep_defragmenter_age_threshold=10,ep_defragmenter_num_visited=32734385439,ep_bg_min_load=44,vb_active_ops_delete=21584929,ep_max_num_workers=3,vb_replica_queue_fill=236801275,vb_pending_meta_data_memory=0,vb_pending_ops_reject=0,ep_vb0=0,ep_expired_pager=21535704,ep_item_flush_expired=0,cmd_flush=0,ep_total_cache_size=9426972182,ep_tap_backoff_period=5,vb_pending_queue_memory=0,vb_pending_eject=0,ep_item_commit_failed=0,ep_dcp_conn_buffer_size_perc=1,wbufs_allocated=0,vb_replica_meta_data_disk=440785478,vb_active_eject=41178213,ep_db_data_size=13406632548,ep_compaction_exp_mem_threshold=85,vb_active_num=342,vb_dead_num=0,rusage_system=1447389.779447,vb_replica_meta_data_memory=827272197,max_conns_on_port_11210=30000,vb_replica_queue_size=55,ep_defragmenter_num_moved=868374298,vb_pending_meta_data_disk=0,cmd_total_ops=9038258159,ep_chk_remover_stime=5,get_misses=1999702274,pid=9519,ep_access_scanner_num_items=2033952,ep_tap_backlog_limit=5000,incr_misses=0,ep_kv_size=10451190661,ep_compaction_write_queue_cap=10000,ep_pending_ops_max=0,ep_tap_keepalive=300,curr_conns_on_port_11210=87,ep_tap_backfill_resident=0.9,ep_bfilter_fp_prob=0.01,ep_tap_ack_window_size=10,ep_max_num_shards=4,cmd_total_sets=4440955621,vb_active_queue_memory=1368,vb_replica_curr_items=12000415,ep_expired_access=45906,vb_replica_ops_update=178487785,daemon_connections=4,ep_flush_duration_total=2154136,ep_chk_period=5,ep_vbucket_del_max_walltime=92377,ep_bg_wait=7038836817,vb_active_perc_mem_resident=53,ep_storage_age_highwat=391,ep_bg_min_wait=3,incr_hits=0,ep_dcp_scan_byte_limit=4194304,ep_dcp_conn_buffer_size_aggr_mem_threshold=10,vb_pending_num_non_resident=0,ep_enable_chk_merge=0,cas_misses=0,ep_meta_data_memory=1655831045,ep_dcp_scan_item_limit=4096,ep_blob_overhead=932123407,ep_num_pager_runs=130,ep_tap_requeue_sleep_time=0.1,ep_tap_bg_fetch_requeued=0,ep_meta_data_disk=877095288,vb_pending_ops_update=0,vb_active_curr_items=12019094,ep_startup_time=1500318296,ep_tmp_oom_errors=6,vb_active_num_non_resident=5629487,ep_dcp_enable_dynamic_conn_buffer_size=1,vb_active_meta_data_memory=828558848,ep_alog_sleep_time=1440,auth_cmds=191,curr_connections=130,ep_chk_persistence_timeout=10,ep_num_expiry_pager_runs=1252,ep_pending_ops=0,ep_max_num_auxio=0,ep_degraded_mode=0,vb_active_queue_age=0,ep_tap_bg_fetched=0,auth_errors=0,vb_active_itm_memory=6765521216,rusage_user=4518098.147224,ep_backfill_mem_threshold=96,ep_warmup_min_items_threshold=100,wbufs_loaned=191,ep_vbucket_del=170,get_hits=1949711112,ep_allow_data_loss_during_shutdown=0,vb_pending_ops_delete=0,bytes_written=2965116528919,ep_mutation_mem_threshold=93,vb_replica_ops_reject=0,ep_pager_active_vb_pcnt=40,ep_blob_num=8424614,vb_replica_queue_drain=236801220,conn_yields=71997165,ep_bfilter_residency_threshold=0.1,ep_dcp_enable_noop=1,ep_overhead=149912326,vb_replica_num_non_resident=9965412,ep_max_checkpoints=2,ep_defragmenter_enabled=1,ep_flushall_enabled=1,ep_storedval_size=1921551024,vb_replica_ht_memory=72041312,ep_vbucket_del_fail=0,vb_replica_queue_age=0,ep_rollback_count=0,ep_db_file_size=16439739172,vb_replica_queue_pending=1636162,ep_storedval_overhead=932123407,listen_disabled_num=0,vb_active_ops_create=33604035,ep_failpartialwarmup=0,ep_defragmenter_interval=10,ep_bg_load_avg=2347,iovused_high_watermark=6,ep_vb_snapshot_total=61800,ep_bg_max_wait=29410188,ep_max_num_readers=0,ep_num_not_my_vbuckets=717,ep_num_eject_failures=1003802636,ep_value_size=8795359616,ep_pending_ops_total=0,ep_total_persisted=471980369,ep_chk_persistence_remains=0,cas_badval=0,vb_active_ops_update=176905303,ep_alog_block_size=4096,ep_max_item_size=20971520,vb_active_meta_data_disk=436309810,ep_oom_errors=0,vb_pending_queue_size=0,vb_active_queue_drain=234903214,ep_chk_max_items=500,vb_replica_queue_memory=3960,ep_bfilter_enabled=1,decr_hits=0,bucket_active_conns=1,ep_warmup_min_memory_threshold=100,ep_commit_num=390100216,ep_max_num_nonio=0,ep_items_rm_from_checkpoints=383385997,ep_warmup_batch_size=1000,ep_mlog_compactor_runs=0,vb_pending_queue_fill=0,ep_commit_time_total=1987641804,ep_flusher_todo=0,mem_used=10890985784,ep_total_del_items=43147550,ep_bg_fetch_delay=0,ep_num_workers=10,vb_replica_expired=24,vb_replica_ops_create=33563034,vb_active_queue_pending=482494,ep_storage_age=0,ep_defragmenter_chunk_duration=20,ep_bg_load=40124554653,ep_ht_locks=47,bucket_conns=30,vb_replica_ops_delete=21562621,vb_replica_num=342,ep_tap_bg_max_pending=500,vb_replica_eject=53765202,ep_tap_ack_grace_period=300,cmd_set=203422402,vb_pending_queue_drain=0,ep_mem_low_wat=9663676416,vb_active_ops_reject=0,cmd_get=3949413386,ep_queue_size=73,ep_bfilter_key_count=10000,ep_num_ops_del_ret_meta=0,delete_misses=0,bytes_read=676984112504,vb_pending_curr_items=0,ep_max_failover_entries=25,ep_num_ops_del_meta_res_fail=0,ep_dcp_noop_interval=180,ep_ht_size=3079,decr_misses=0,cmd_total_gets=4597302538,ep_tap_throttle_queue_cap=-1,ep_total_enqueued=478817098,ep_item_flush_failed=0,ep_tap_throttle_cap_pcnt=10,ep_tap_ack_interval=1000,ep_num_ops_set_ret_meta=0,connection_structures=130,ep_num_ops_get_meta_on_set_meta=0,vb_replica_perc_mem_resident=16,accepting_conns=1,cas_hits=0,curr_temp_items=0,ep_access_scanner_enabled=1,ep_tap_noop_interval=20,ep_num_ops_set_meta=0,ep_max_vbuckets=1024,ep_diskqueue_items=73,vb_pending_num=0,vb_active_queue_size=19,ep_num_ops_del_meta=0,ep_mem_high_wat=10952166604,pointer_size=64,rbufs_existing=631390261,ep_vbucket_del_avg_walltime=6956,ep_num_value_ejects=95343207,ep_getl_default_timeout=15,vb_replica_itm_memory=2661450966,ep_max_num_writers=0,ep_pending_compactions=0,curr_items=12019094,total_connections=69704,rejected_conns=0,ep_bg_fetched=17092899,ep_bg_remaining_jobs=0,ep_bg_num_samples=17092899,vb_pending_queue_age=0,ep_access_scanner_last_runtime=10,time=1504828905,vb_active_expired=21581586,max_conns_on_port_11209=5000 1504828907000000000
```
