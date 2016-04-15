package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lafikl/gotabulate"
)

func (in *Insight) Queries(args []string) {
	query := `
    SELECT query, calls, total_time/1000/60 total_minutes, (total_time/calls)/1000 avg_per_q, rows, coalesce((100.0 * shared_blks_hit /
                   nullif(shared_blks_hit + shared_blks_read, 0))::TEXT, '0%%') AS hit_percent
      FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;
	`
	rows, err := in.db.Query(query)
	if err != nil {
		fmt.Println("Couldn't get results\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	results := [][]string{}
	for rows.Next() {
		r := []string{"", "", "", "", "", ""}
		err = rows.Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5])
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
	t.SetHeaders([]string{"Query", "Calls", "Total Time", "Avg Time/Query", "Rows", "Cache Hits %"})
	// Set Max Cell Size
	t.SetMaxCellSize(20)
	// Turn On String Wrapping
	t.SetWrapStrings(true)
	t.SetAlign("center")
	fmt.Println()
	fmt.Println(queriesHelp)
	fmt.Println(t.Render("grid"))
	fmt.Println()
}
