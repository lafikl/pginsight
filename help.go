package main

const usage = `
Usage:
./pginsight <cmdname> [--flags]

Commands:
    index usage
        Shows which indexes are being scanned and how many tuples are fetched
    index unused
        Shows the indexes which haven't been scanned
    index duplicate
        Finds indexes which index on the same key(s)

    disk db
        Shows disk usage for all databases accessible from the configured user
    disk relations
        Show disk usage for all relations in the configured database

    cache total
        Shows the total number of cache hits in the database
    cache tables
        Breakdown of cache hits per table

    queries
        Shows the slowest 10 queries.

`

const indexUsage = `
Usage:
./pginsight index <subcommand>

Commands:
    index usage
        Shows which indexes are being scanned and how many tuples are fetched
    index unused
        Shows the indexes which haven't been scanned
    index bloat
        Measure index bloat for all tables
`

const idxBloatFields = `
Fields Info:
    Extra Size: Estitmation of the amount of bytes not used in the table. Which composed of fillfactor, bloat, and alignment padding spaces
    Extra Ratio: Estimated ratio of the real size used by Extra Size
    Fillfactor: http://www.postgresql.org/docs/9.4/static/sql-createtable.html#SQL-CREATETABLE-STORAGE-PARAMETERS
    Bloat Size: Estimated size of the bloat without the extra space kept for the fillfactor
    Bloat Ratio: Estimated ratio of the real size used by Bloat Size
    Not Applicable: If true, do not trust the stats

Credits:
    This report is based of https://github.com/ioguix/pgsql-bloat-estimation
`

const cacheUsage = `
Usage:
./pginsight cache <subcommand>

Commands:
    cache total
        Shows the total number of cache hits in the database
    cache tables
        Breakdown of cache hits per table
`

const cacheFields = `
Fields Info:
    Heap Read: Number of disk blocks read.
    Heap Hit: Number of buffer hits.
    Hit Ratio: Ratio of cache hits.
`

const queriesHelp = `
Fields Info:
    Calls: Number of times this query got executed
    Total Time: Cumulative sum of time spent executing this statement, in minutes
    Avg Time/Query: Average time spent executing this statement per query
    Rows: Total number of rows retrieved or affected by the statement
    Hit Ratio: Ratio of cache hits.
`

const diskUsage = `
Usage:
./pginsight disk <subcommand>

Commands:
    disk db
        Shows disk usage for all databases accessible from the configured user
    disk tables
        Show disk usage for all tables in the configured database
`
