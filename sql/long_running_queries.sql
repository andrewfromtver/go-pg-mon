SELECT 
    application_name, 
    backend_type, 
    client_addr, 
    client_port, 
    datname AS pg_database, 
    usename AS pg_user, 
    pid, 
    query, 
    now() - query_start AS query_duration,
    query_start, 
    state, 
    wait_event 
FROM 
    pg_stat_activity
WHERE 
    now() - query_start > interval '$query_duration'
    AND datname = '$pg_database'
    AND (
        '$query_state' = '*' OR state = '$query_state'
    );
