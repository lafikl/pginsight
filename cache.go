package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lafikl/gotabulate"
)

func (in *Insight) Cache(args []string) {
	if len(args) == 0 || len(args) > 2 {
		fmt.Println(cacheUsage)
		os.Exit(0)
	}
	command := strings.ToLower(args[0])
	switch command {
	case "total":
		in.cacheTotal(args[1:])
		break
	case "tables":
		in.cacheTables(args[1:])
		break
	default:
		fmt.Println("Unknown command ", command)
		fmt.Printf("\n%s", cacheUsage)
		os.Exit(1)
	}
}

func (in *Insight) cacheTotal(argv []string) {
	query := `
	SELECT
	  sum(heap_blks_read) as heap_read,
	  sum(heap_blks_hit)  as heap_hit,
	  (sum(heap_blks_hit)*100 / (sum(heap_blks_hit) + sum(heap_blks_read))) as ratio
	FROM
	  pg_statio_user_tables;
	`
	rows, err := in.db.Query(query)
	if err != nil {
		fmt.Println("Couldn't get results\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	results := [][]string{}
	for rows.Next() {
		r := []string{"", "", ""}
		err = rows.Scan(&r[0], &r[1], &r[2])
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
	t.SetHeaders([]string{"Heap Read", "Heap Hit", "Hit Ratio %"})
	t.SetAlign("center")
	fmt.Println()
	fmt.Println(cacheFields)
	fmt.Println(t.Render("grid"))
	fmt.Println()
}

func (in *Insight) cacheTables(argv []string) {
	query := `
	SELECT
      relname,
	  heap_blks_read as heap_read,
	  heap_blks_hit as heap_hit,
	  ( (heap_blks_hit*100) / (heap_blks_hit + heap_blks_read)) as ratio
	FROM
	  pg_statio_user_tables;
	`
	rows, err := in.db.Query(query)
	if err != nil {
		fmt.Println("Couldn't get results\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	results := [][]string{}
	for rows.Next() {
		r := []string{"", "", "", ""}
		err = rows.Scan(&r[0], &r[1], &r[2], &r[3])
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
	t.SetHeaders([]string{"Table", "Heap Read", "Heap Hit", "Hit Ratio %"})
	t.SetAlign("center")
	fmt.Println()
	fmt.Println(cacheFields)
	fmt.Println(t.Render("grid"))
	fmt.Println()
}
