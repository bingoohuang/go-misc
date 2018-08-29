package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"strings"
)

type Right struct {
	Id              string `db:"ID"`
	RightCode       string `db:"RIGHT_CODE"`
	RightName       string `db:"RIGHT_NAME"`
	RightType       string `db:"RIGHT_TYPE"`
	RightDesc       string `db:"RIGHT_DESC"`
	ServpackId      string `db:"SERVPACK_ID"`
	Url             string `db:"URL"`
	ChainName       string `db:"CHAIN_NAME"`
	ChainDefinition string `db:"CHAIN_DEFINITION"`
}

func init()  {
	dsn = os.Args[1] // "user:pass@tcp(ip:3306)/db?charset=utf8"
}

var dsn  string

func removeDuplicate() {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("Disabling MySQL tests:\n    %v", err)
		return
	}
	defer db.Close()

	rights := []Right{}

	err = db.Select(&rights, `select ID, RIGHT_CODE, RIGHT_NAME, RIGHT_TYPE, RIGHT_DESC, COALESCE(SERVPACK_ID,'') SERVPACK_ID, 
		COALESCE(URL, '') URL, COALESCE(CHAIN_NAME, '') CHAIN_NAME, COALESCE(CHAIN_DEFINITION, '') CHAIN_DEFINITION 
		from tt_d_funcright order by right_code, RIGHT_NAME, servpack_id`)
	if err != nil {
		fmt.Printf("Select:\n    %v", err)
		return
	}

	duplicates := 0
	duplicateIds := []string{}
	for i := 0; i+1 < len(rights); {
		l := rights[i]
		r := rights[i+1]

		lid := l.Id
		rid := r.Id

		l.Id = ""
		r.Id = ""

		if l == r {
			fmt.Printf("Found Duplicats: %v\n", lid)
		}

		i++

		for ; l == r; i++ {
			duplicates++
			r.Id = rid
			duplicateIds = append(duplicateIds, r.Id)
			fmt.Printf("Duplicate:%v\n", r)
			r = rights[1+i]
			rid = r.Id
			r.Id = ""
		}

		l.Id = lid
		r.Id = rid
	}

	fmt.Printf("Duplicates:%v\n", duplicates)

	if duplicates > 0 {
		db.MustExec(`delete from tt_d_funcright where id in (` + strings.Join(duplicateIds, ",") + `)`)
	}
}
