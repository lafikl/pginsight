package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lafikl/gotabulate"
)

func (in *Insight) Index(args []string) {
	if len(args) == 0 || len(args) > 2 {
		fmt.Println(indexUsage)
		os.Exit(0)
	}
	command := strings.ToLower(args[0])
	switch command {
	case "usage":
		in.idxUsage(args[1:])
		break
	case "hits":
		in.idxHits(args[1:])
		break
	case "unused":
		in.idxUnused(args[1:])
		break
	case "bloat":
		in.idxBloat(args[1:])
		break
	default:
		fmt.Println("Unknown command ", command)
		fmt.Printf("\n%s", indexUsage)
		os.Exit(1)
	}
}

func (in *Insight) idxUsage(args []string) {
	query := `
    SELECT
        t.tablename,
        indexname,
        c.reltuples AS num_rows,
        pg_size_pretty(pg_relation_size(quote_ident(t.tablename)::text)) AS table_size,
        pg_size_pretty(pg_relation_size(quote_ident(indexrelname)::text)) AS index_size,
        CASE WHEN indisunique THEN 'Y'
           ELSE 'N'
        END AS UNIQUE,
        idx_scan AS number_of_scans,
        idx_tup_read AS tuples_read,
        idx_tup_fetch AS tuples_fetched
    FROM pg_tables t
    LEFT OUTER JOIN pg_class c ON t.tablename=c.relname
    LEFT OUTER JOIN
        ( SELECT c.relname AS ctablename, ipg.relname AS indexname, x.indnatts AS number_of_columns, idx_scan, idx_tup_read, idx_tup_fetch, indexrelname, indisunique FROM pg_index x
               JOIN pg_class c ON c.oid = x.indrelid
               JOIN pg_class ipg ON ipg.oid = x.indexrelid
               JOIN pg_stat_all_indexes psai ON x.indexrelid = psai.indexrelid )
        AS foo
        ON t.tablename = foo.ctablename
    WHERE t.schemaname NOT IN ('pg_catalog', 'information_schema')
    ORDER BY 1,2;
    `
	rows, err := in.db.Query(query)
	if err != nil {
		fmt.Println("Couldn't get results\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	results := [][]string{}
	for rows.Next() {
		r := []string{"", "", "", "", "", "", "", "", ""}
		err = rows.Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5],
			&r[6], &r[7], &r[8])
		if err != nil {
			fmt.Println("error occured while parsing results", err)
			os.Exit(1)
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	t := gotabulate.Create(results)
	t.SetHeaders([]string{"Table", "Index", "Number of rows", "Table Size",
		"Index Size", "Unique?", "Number of Scans", "Tuples Read", "Tuples Fetched"})
	t.SetAlign("center")
	fmt.Println()
	fmt.Println(t.Render("grid"))
	fmt.Println()

}

func (in *Insight) idxHits(args []string) {
	query := `
	SELECT
	  relname,
	  n_live_tup rows_in_table,
	  seq_scan,
	  idx_scan,
	  100 * idx_scan/(idx_scan + seq_scan) percent_index_use,
	  coalesce(last_vacuum::TEXT, 'Never'),
	  coalesce(last_autovacuum::TEXT, 'Never'),
	  coalesce(last_analyze::TEXT, 'Never'),
	  coalesce(last_autoanalyze::TEXT, 'Never'),
	  autovacuum_count,
	  autoanalyze_count
	FROM
	  pg_stat_user_tables
	WHERE
	    seq_scan + idx_scan > 0
	ORDER BY
	  n_live_tup DESC;
    `
	rows, err := in.db.Query(query)
	if err != nil {
		fmt.Println("Couldn't get results\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	results := [][]string{}
	for rows.Next() {
		r := []string{"", "", "", "", "", "", "", "", "", "", ""}
		err = rows.Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5],
			&r[6], &r[7], &r[8], &r[9], &r[10])
		if err != nil {
			fmt.Println("error occured while parsing results", err)
			os.Exit(1)
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	t := gotabulate.Create(results)
	t.SetHeaders([]string{"Table", "Number of Rows", "Sequential Scans", "Index Scans",
		"Index Usage %", "Last Vacuum", "Last Auto-Vacuum", "Last Analyze",
		"Last Auto-Analyze", "Auto-Vacuums", "Auto-Analyzes"})
	t.SetAlign("center")
	fmt.Println()
	fmt.Println(t.Render("grid"))
	fmt.Println()
}

// columns order
// tablename, indexname, num_rows, table_size, index_size,
// unique, number_of_scans, tuples_read, tuples_fetched string
//
func (in *Insight) idxUnused(args []string) {
	query := `
    SELECT
        t.tablename,
        indexname,
        c.reltuples AS num_rows,
        pg_size_pretty(pg_relation_size(quote_ident(t.tablename)::text)) AS table_size,
        pg_size_pretty(pg_relation_size(quote_ident(indexrelname)::text)) AS index_size,
        CASE WHEN indisunique THEN 'Y'
           ELSE 'N'
        END AS UNIQUE,
        idx_scan AS number_of_scans,
        idx_tup_read AS tuples_read,
        idx_tup_fetch AS tuples_fetched
    FROM pg_tables t
    LEFT OUTER JOIN pg_class c ON t.tablename=c.relname
    LEFT OUTER JOIN
        ( SELECT c.relname AS ctablename, ipg.relname AS indexname, x.indnatts AS number_of_columns, idx_scan, idx_tup_read, idx_tup_fetch, indexrelname, indisunique FROM pg_index x
               JOIN pg_class c ON c.oid = x.indrelid
               JOIN pg_class ipg ON ipg.oid = x.indexrelid
               JOIN pg_stat_all_indexes psai ON x.indexrelid = psai.indexrelid )
        AS foo
        ON t.tablename = foo.ctablename
    WHERE t.schemaname NOT IN ('pg_catalog', 'information_schema') AND idx_scan = 0;
    `
	rows, err := in.db.Query(query)
	if err != nil {
		fmt.Println("Couldn't get results\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	results := [][]string{}
	for rows.Next() {
		r := []string{"", "", "", "", "", "", "", "", ""}
		err = rows.Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5],
			&r[6], &r[7], &r[8])
		if err != nil {
			fmt.Println("error occured while parsing results", err)
			os.Exit(1)
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	t := gotabulate.Create(results)
	t.SetHeaders([]string{"Table", "Index", "Number of rows", "Table Size",
		"Index Size", "Unique?", "Number of Scans", "Tuples Read", "Tuples Fetched"})
	t.SetAlign("center")
	fmt.Println()
	fmt.Println(t.Render("grid"))
	fmt.Println()
}

func (in *Insight) idxBloat(args []string) {
	query := `
    SELECT current_database(), schemaname, tblname, bs*tblpages AS real_size,
      (tblpages-est_tblpages)*bs AS extra_size,
      CASE WHEN tblpages - est_tblpages > 0
        THEN 100 * (tblpages - est_tblpages)/tblpages::float
        ELSE 0
      END AS extra_ratio, fillfactor, (tblpages-est_tblpages_ff)*bs AS bloat_size,
      CASE WHEN tblpages - est_tblpages_ff > 0
        THEN 100 * (tblpages - est_tblpages_ff)/tblpages::float
        ELSE 0
      END AS bloat_ratio, is_na
      -- , (pst).free_percent + (pst).dead_tuple_percent AS real_frag
    FROM (
      SELECT ceil( reltuples / ( (bs-page_hdr)/tpl_size ) ) + ceil( toasttuples / 4 ) AS est_tblpages,
        ceil( reltuples / ( (bs-page_hdr)*fillfactor/(tpl_size*100) ) ) + ceil( toasttuples / 4 ) AS est_tblpages_ff,
        tblpages, fillfactor, bs, tblid, schemaname, tblname, heappages, toastpages, is_na
        -- , stattuple.pgstattuple(tblid) AS pst
      FROM (
        SELECT
          ( 4 + tpl_hdr_size + tpl_data_size + (2*ma)
            - CASE WHEN tpl_hdr_size%ma = 0 THEN ma ELSE tpl_hdr_size%ma END
            - CASE WHEN ceil(tpl_data_size)::int%ma = 0 THEN ma ELSE ceil(tpl_data_size)::int%ma END
          ) AS tpl_size, bs - page_hdr AS size_per_block, (heappages + toastpages) AS tblpages, heappages,
          toastpages, reltuples, toasttuples, bs, page_hdr, tblid, schemaname, tblname, fillfactor, is_na
        FROM (
          SELECT
            tbl.oid AS tblid, ns.nspname AS schemaname, tbl.relname AS tblname, tbl.reltuples,
            tbl.relpages AS heappages, coalesce(toast.relpages, 0) AS toastpages,
            coalesce(toast.reltuples, 0) AS toasttuples,
            coalesce(substring(
              array_to_string(tbl.reloptions, ' ')
              FROM '%fillfactor=#"__#"%' FOR '#')::smallint, 100) AS fillfactor,
            current_setting('block_size')::numeric AS bs,
            CASE WHEN version()~'mingw32' OR version()~'64-bit|x86_64|ppc64|ia64|amd64' THEN 8 ELSE 4 END AS ma,
            24 AS page_hdr,
            23 + CASE WHEN MAX(coalesce(null_frac,0)) > 0 THEN ( 7 + count(*) ) / 8 ELSE 0::int END
              + CASE WHEN tbl.relhasoids THEN 4 ELSE 0 END AS tpl_hdr_size,
            sum( (1-coalesce(s.null_frac, 0)) * coalesce(s.avg_width, 1024) ) AS tpl_data_size,
            bool_or(att.atttypid = 'pg_catalog.name'::regtype) AS is_na
          FROM pg_attribute AS att
            JOIN pg_class AS tbl ON att.attrelid = tbl.oid
            JOIN pg_namespace AS ns ON ns.oid = tbl.relnamespace
            JOIN pg_stats AS s ON s.schemaname=ns.nspname
              AND s.tablename = tbl.relname AND s.inherited=false AND s.attname=att.attname
            LEFT JOIN pg_class AS toast ON tbl.reltoastrelid = toast.oid
          WHERE att.attnum > 0 AND NOT att.attisdropped
            AND tbl.relkind = 'r' And schemaname NOT IN ('pg_catalog', 'information_schema')
          GROUP BY 1,2,3,4,5,6,7,8,9,10, tbl.relhasoids
          ORDER BY 2,3
        ) AS s
      ) AS s2
    ) AS s3;
    `
	rows, err := in.db.Query(query)
	if err != nil {
		fmt.Println("Couldn't get results\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	results := [][]string{}
	for rows.Next() {
		r := []string{"", "", "", "", "", "", "", "", "", ""}
		err = rows.Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5],
			&r[6], &r[7], &r[8], &r[9])
		if err != nil {
			fmt.Println("error occured while parsing results", err)
			os.Exit(1)
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	t := gotabulate.Create(results)
	t.SetHeaders([]string{"Database", "Schema Name", "Table", "Real Table Size",
		"Extra Size", "Extra Ratio %", "Fillfactor", "Bloat Size", "Bloat Ratio %", "Not Applicable?"})
	t.SetAlign("center")
	fmt.Println()
	fmt.Println(idxBloatFields)
	fmt.Println(t.Render("grid"))
	fmt.Println()
}
