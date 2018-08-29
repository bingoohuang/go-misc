package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func updateAcl() {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("Disabling MySQL tests:\n    %v", err)
		return
	}
	defer db.Close()

	rights := []RoleRight{}

	err = db.Select(&rights, `
		select ID, ROLE_ID, COALESCE(SERVPACK_ID,'') SERVPACK_ID 
		from role_funcright`)
	if err != nil {
		fmt.Printf("Select:\n    %v", err)
		return
	}

	for _, v := range rights {
		if v.RoleId != "" && v.ServpackId != "" {
			acl := `@Acl(roles = "` + v.RoleId + `", pkgs = "` + v.ServpackId + `")`
			db.MustExec(`update role_funcright set ACL = ? where id = ?`, acl, v.Id)
		}
	}

}
