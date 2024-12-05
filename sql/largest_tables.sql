SELECT
    st.schemaname || '.' || st.relname AS table_name,
    pg_size_pretty(pg_total_relation_size(st.relid)) AS total_size,
    pg_size_pretty(pg_relation_size(st.relid)) AS table_size,
    pg_size_pretty(pg_total_relation_size(st.relid) - pg_relation_size(st.relid)) AS index_size,
    pg_size_pretty(pg_total_relation_size(st.relid) - pg_table_size(st.relid)) AS toast_size,
    pg_stat_get_live_tuples(st.relid) AS live_tuples,
    pg_stat_get_dead_tuples(st.relid) AS dead_tuples,
    stat.n_tup_ins AS tuples_inserted,
    stat.n_tup_upd AS tuples_updated,
    stat.n_tup_del AS tuples_deleted,
    stat.seq_scan AS sequential_scans,
    stat.idx_scan AS index_scans
FROM pg_catalog.pg_statio_user_tables st
JOIN pg_catalog.pg_stat_user_tables stat
    ON st.relid = stat.relid
ORDER BY pg_total_relation_size(st.relid) DESC
LIMIT $query_limit;
