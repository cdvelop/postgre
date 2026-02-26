package postgre

import (
	"fmt"
	"strings"

	"github.com/tinywasm/orm"
)

// translate converts an ORM query to a PostgreSQL query and arguments.
func translate(q orm.Query) (string, []any, error) {
	var sb strings.Builder
	var args []any
	argIndex := 1

	switch q.Action {
	case orm.ActionCreate:
		sb.WriteString("INSERT INTO ")
		sb.WriteString(q.Table)
		sb.WriteString(" (")
		sb.WriteString(strings.Join(q.Columns, ", "))
		sb.WriteString(") VALUES (")
		for i, v := range q.Values {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("$%d", argIndex))
			args = append(args, v)
			argIndex++
		}
		sb.WriteString(")")
		// Append RETURNING id if it's likely expected, although not strictly in generic ORM spec.
		// However, many ORMs rely on LastInsertId which Postgres doesn't support via Result.
		// So we often need `RETURNING id`.
		// Let's assume the user model has an 'id' column for now or handle it via Execute scan.
		// If we don't add it here, `Execute` won't get it back.
		// But generic `translate` shouldn't assume column names unless specified.
		// The `orm` might handle ID assignment via UUIDs generated in Go, in which case RETURNING isn't needed.
		// If the DB generates IDs (SERIAL/IDENTITY), we need it.
		// Let's stick to standard INSERT for now. If tests fail on ID retrieval, we'll revisit.

	case orm.ActionReadOne, orm.ActionReadAll:
		sb.WriteString("SELECT ")
		if len(q.Columns) == 0 {
			sb.WriteString("*")
		} else {
			sb.WriteString(strings.Join(q.Columns, ", "))
		}
		sb.WriteString(" FROM ")
		sb.WriteString(q.Table)
		if err := buildConditions(&sb, q.Conditions, &args, &argIndex); err != nil {
			return "", nil, err
		}
		if len(q.OrderBy) > 0 {
			sb.WriteString(" ORDER BY ")
			for i, o := range q.OrderBy {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(o.Column()) // Changed from o.Field to o.Column()
				sb.WriteString(" ")
				sb.WriteString(o.Dir())    // Changed from checking o.Desc to o.Dir() which returns "ASC" or "DESC"
			}
		}
		if q.Limit > 0 {
			sb.WriteString(fmt.Sprintf(" LIMIT %d", q.Limit))
		}
		if q.Offset > 0 {
			sb.WriteString(fmt.Sprintf(" OFFSET %d", q.Offset))
		}

	case orm.ActionUpdate:
		sb.WriteString("UPDATE ")
		sb.WriteString(q.Table)
		sb.WriteString(" SET ")
		for i, c := range q.Columns {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(c)
			sb.WriteString(fmt.Sprintf(" = $%d", argIndex))
			args = append(args, q.Values[i])
			argIndex++
		}
		if err := buildConditions(&sb, q.Conditions, &args, &argIndex); err != nil {
			return "", nil, err
		}

	case orm.ActionDelete:
		sb.WriteString("DELETE FROM ")
		sb.WriteString(q.Table)
		if err := buildConditions(&sb, q.Conditions, &args, &argIndex); err != nil {
			return "", nil, err
		}

	default:
		return "", nil, fmt.Errorf("unsupported action: %d", q.Action)
	}

	return sb.String(), args, nil
}

func buildConditions(sb *strings.Builder, conditions []orm.Condition, args *[]any, argIndex *int) error {
	if len(conditions) == 0 {
		return nil
	}

	// Check if the first condition has logic, if not, it's implicitly part of the WHERE block start.
	// But `orm.Condition` logic applies to how it connects to the *previous* condition.
	// The first condition's logic is ignored or should be empty/AND (implicitly).

	sb.WriteString(" WHERE ")
	for i, c := range conditions {
		if i > 0 {
			logic := c.Logic()
			if logic == "" {
				logic = "AND"
			}
			sb.WriteString(fmt.Sprintf(" %s ", logic))
		}

		sb.WriteString(c.Field())
		sb.WriteString(" ")
		sb.WriteString(c.Operator())
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf("$%d", *argIndex))
		*args = append(*args, c.Value())
		*argIndex++
	}
	return nil
}
