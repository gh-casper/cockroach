// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package delegate

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// delegateShowDefaultPrivileges implements SHOW DEFAULT PRIVILEGES
// which returns default privileges for a specified role.
func (d *delegator) delegateShowDefaultPrivileges(
	n *tree.ShowDefaultPrivileges,
) (tree.Statement, error) {
	currentDatabase, err := d.getSpecifiedOrCurrentDatabase("")
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(
		"SELECT role, object_type, grantee, privilege_type FROM crdb_internal.default_privileges WHERE database_name = '%s'",
		currentDatabase.Normalize(),
	)

	targetRoles, err := n.Roles.ToSQLUsernames()
	if err != nil {
		return nil, err
	}
	if len(targetRoles) == 0 {
		targetRoles = append(targetRoles, d.evalCtx.SessionData.User())
	}

	query = fmt.Sprintf("%s AND role IN (", query)
	for i, role := range targetRoles {
		if i != 0 {
			query += fmt.Sprintf(", '%s'", role.Normalized())
		} else {
			query += fmt.Sprintf("'%s'", role.Normalized())
		}
	}

	query += ")"

	query += " ORDER BY 1,2,3,4"
	return parse(query)
}
