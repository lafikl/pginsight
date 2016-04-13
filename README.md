# PGInsight
CLI tool to easily dig deep inside your PostgreSQL database.

# Features
- Index usage
- Find unused indexes
- Find duplicate indexes
- Disk usage for every database
- Disk usage for all tables in a given database
- Cache hits ratio for a given database
- Cache hits per table
- Find slow queries
- Self-contained binary, no dependencies

# Installation
- **Step 1**: Download binary from https://github.com/lafikl/pginsight/releases
- **Step 2**: There's no step 2!

# Usage
```
Usage:
PGINSIGHT_DBURL="postgres://username:password@localhost/dbname?sslmode=disable" ./pginsight <cmdname> [subcommand]

Commands:
    index usage
        Shows which indexes are being scanned and how many tuples are fetched
    index unused
        Shows the indexes which haven't been scanned
    index duplicate
        Finds indexes which index on the same key(s)

    disk db
        Shows disk usage for all databases accessible from the configured user
    disk tables
        Show disk usage for all tables in the configured database

    cache total
        Shows the total number of cache hits in the database
    cache tables
        Breakdown of cache hits per table

    queries
        Shows the slowest 10 queries.
```


# Example
PGINSIGHT_DBURL="postgres://klafi:eee@localhost/test?sslmode=disable" ./pginsight queries


# License
The MIT License (MIT)

Copyright (c) 2016 Khalid Lafi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
