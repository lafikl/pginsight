package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lafikl/gotabulate"
)

func (in *Insight) Disk(args []string) {
	if len(args) == 0 || len(args) > 2 {
		fmt.Println(diskUsage)
		os.Exit(0)
	}
	command := strings.ToLower(args[0])
	switch command {
	case "db":
		in.diskDB(args[1:])
		break
	case "total":
		in.diskTotal(args[1:])
		break
	default:
		fmt.Println("Unknown command ", command)
		fmt.Printf("\n%s", diskUsage)
		os.Exit(1)
	}
}

func (in *Insight) diskDB(argv []string) {
	query := `
    SELECT d.datname AS Name,  pg_catalog.pg_get_userbyid(d.datdba) AS Owner,
        CASE WHEN pg_catalog.has_database_privilege(d.datname, 'CONNECT')
            THEN pg_catalog.pg_size_pretty(pg_catalog.pg_database_size(d.datname))
            ELSE 'No Access'
        END AS SIZE
    FROM pg_catalog.pg_database d
        ORDER BY
        CASE WHEN pg_catalog.has_database_privilege(d.datname, 'CONNECT')
            THEN pg_catalog.pg_database_size(d.datname)
            ELSE NULL
        END DESC -- nulls first
        LIMIT 20;
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
	t.SetHeaders([]string{"DB Name", "Owner", "Size"})
	t.SetAlign("center")
	fmt.Println()
	fmt.Println(t.Render("grid"))
	fmt.Println()
}

func (in *Insight) diskTotal(argv []string) {
	query := `
    SELECT nspname || '.' || relname AS "relation",
        pg_size_pretty(pg_relation_size(C.oid)) AS "size"
      FROM pg_class C
      LEFT JOIN pg_namespace N ON (N.oid = C.relnamespace)
      WHERE nspname NOT IN ('pg_catalog', 'information_schema')
      ORDER BY pg_relation_size(C.oid) DESC LIMIT 20;
	`
	rows, err := in.db.Query(query)
	if err != nil {
		fmt.Println("Couldn't get results\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	results := [][]string{}
	for rows.Next() {
		r := []string{"", ""}
		err = rows.Scan(&r[0], &r[1])
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
	t.SetHeaders([]string{"Relation", "Size"})
	t.SetAlign("center")
	fmt.Println()
	fmt.Println(t.Render("grid"))
	fmt.Println()
}
