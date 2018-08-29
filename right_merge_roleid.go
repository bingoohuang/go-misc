package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

type RoleRight struct {
	Id              string `db:"ID"`
	RoleId          string `db:"ROLE_ID"`
	RightCode       string `db:"RIGHT_CODE"`
	RightName       string `db:"RIGHT_NAME"`
	RightType       string `db:"RIGHT_TYPE"`
	RightDesc       string `db:"RIGHT_DESC"`
	ServpackId      string `db:"SERVPACK_ID"`
	Url             string `db:"URL"`
	ChainName       string `db:"CHAIN_NAME"`
	ChainDefinition string `db:"CHAIN_DEFINITION"`
}

func mergeRoleIds() {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("Disabling MySQL tests:\n    %v", err)
		return
	}
	defer db.Close()

	rights := []RoleRight{}

	err = db.Select(&rights, `
select ID, ROLE_ID, RIGHT_CODE, RIGHT_NAME, RIGHT_TYPE, RIGHT_DESC,
COALESCE(SERVPACK_ID,'') SERVPACK_ID, COALESCE(URL, '') URL, COALESCE(CHAIN_NAME, '') CHAIN_NAME, COALESCE(CHAIN_DEFINITION, '') CHAIN_DEFINITION 
from role_funcright
order by right_code, RIGHT_NAME, servpack_id`)
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

		lRoldId := l.RoleId
		mergedRoleIds := lRoldId
		rRoleId := r.RoleId

		l.RoleId = ""
		r.RoleId = ""

		hasMerged := false
		if l == r {
			fmt.Printf("Found Mergables: %v\n", lid)
			hasMerged = true
		}

		i++

		for ; l == r; i++ {
			mergedRoleIds += "," + rRoleId
			merges++
			duplicateIds = append(duplicateIds, rid)
			r.Id = rid
			r.RoleId = rRoleId
			fmt.Printf("Mergeable:%v\n", r)

			if i + 1 >= len(rights) {
				break
			}
			r = rights[1+i]
			rid = r.Id
			r.Id = ""

			rRoleId = r.RoleId
			r.RoleId = ""
		}

		l.Id = lid
		r.Id = rid

		l.RoleId = lRoldId
		r.RoleId = rRoleId

		if hasMerged {
			db.MustExec("update role_funcright set ROLE_ID = ? where id = ?", mergedRoleIds, lid)
		}
	}

	fmt.Printf("merges:%v\n", merges)

	if merges > 0 {
		db.MustExec(`delete from role_funcright where id in (` + strings.Join(duplicateIds, ",") + `)`)
	}
}
