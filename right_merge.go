package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

func mergeRights() {
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

	merges := 0
	duplicateIds := []string{}
	for i := 0; i+1 < len(rights); {
		l := rights[i]
		r := rights[i+1]

		lid := l.Id
		rid := r.Id
		l.Id = ""
		r.Id = ""

		lsp := l.ServpackId
		lspMerged := lsp
		rsp := r.ServpackId

		l.ServpackId = ""
		r.ServpackId = ""

		hasMerged := false
		if l == r {
			fmt.Printf("Found Mergables: %v\n", lid)
			hasMerged = true
		}

		i++

		for ; l == r; i++ {
			lspMerged += "," + rsp
			merges++
			duplicateIds = append(duplicateIds, rid)
			r.Id = rid
			r.ServpackId = rsp
			fmt.Printf("Mergeable:%v\n", r)
			r = rights[1+i]
			rid = r.Id
			r.Id = ""

			rsp = r.ServpackId
			r.ServpackId = ""
		}

		l.Id = lid
		r.Id = rid

		l.ServpackId = lsp
		r.ServpackId = rsp

		if hasMerged {
			db.MustExec("update tt_d_funcright set SERVPACK_ID = ? where id = ?", lspMerged, lid)
		}
	}

	fmt.Printf("merges:%v\n", merges)

	if merges > 0 {
		db.MustExec(`delete from tt_d_funcright where id in (` + strings.Join(duplicateIds, ",") + `)`)
	}
}
