package postgres

import (
	"errors"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/orm"
)

// translate converts an ORM query to a PostgreSQL query and arguments.
func translate(q orm.Query) (string, []any, error) {
	c := fmt.Convert() // Replaces strings.Builder
	var args []any
	argIndex := 1

	switch q.Action {
	case orm.ActionCreate:
		c.Write("INSERT INTO ")
		c.Write(q.Table)
		c.Write(" (")
		c.Write(fmt.Convert(q.Columns).Join(", ").String())
		c.Write(") VALUES (")
		for i, v := range q.Values {
			if i > 0 {
				c.Write(", ")
			}
			c.Write(fmt.Sprintf("$%d", argIndex))
			args = append(args, v)
			argIndex++
		}
		c.Write(")")

	case orm.ActionReadOne, orm.ActionReadAll:
		c.Write("SELECT ")
		if len(q.Columns) == 0 {
			c.Write("*")
		} else {
			c.Write(fmt.Convert(q.Columns).Join(", ").String())
		}
		c.Write(" FROM ")
		c.Write(q.Table)
		if err := buildConditions(c, q.Conditions, &args, &argIndex); err != nil {
			return "", nil, err
		}
		if len(q.OrderBy) > 0 {
			c.Write(" ORDER BY ")
			for i, o := range q.OrderBy {
				if i > 0 {
					c.Write(", ")
				}
				c.Write(o.Column())
				c.Write(" ")
				c.Write(o.Dir())
			}
		}
		if q.Limit > 0 {
			c.Write(fmt.Sprintf(" LIMIT %d", q.Limit))
		}
		if q.Offset > 0 {
			c.Write(fmt.Sprintf(" OFFSET %d", q.Offset))
		}

	case orm.ActionUpdate:
		c.Write("UPDATE ")
		c.Write(q.Table)
		c.Write(" SET ")
		for i, col := range q.Columns {
			if i > 0 {
				c.Write(", ")
			}
			c.Write(col)
			c.Write(fmt.Sprintf(" = $%d", argIndex))
			args = append(args, q.Values[i])
			argIndex++
		}
		if err := buildConditions(c, q.Conditions, &args, &argIndex); err != nil {
			return "", nil, err
		}

	case orm.ActionDelete:
		c.Write("DELETE FROM ")
		c.Write(q.Table)
		if err := buildConditions(c, q.Conditions, &args, &argIndex); err != nil {
			return "", nil, err
		}

	default:
		return "", nil, errors.New(fmt.Sprintf("unsupported action: %d", q.Action))
	}

	return c.String(), args, nil
}

func buildConditions(c *fmt.Conv, conditions []orm.Condition, args *[]any, argIndex *int) error {
	if len(conditions) == 0 {
		return nil
	}

	c.Write(" WHERE ")
	for i, cond := range conditions {
		if i > 0 {
			logic := cond.Logic()
			if logic == "" {
				logic = "AND"
			}
			c.Write(fmt.Sprintf(" %s ", logic))
		}

		c.Write(cond.Field())
		c.Write(" ")
		c.Write(cond.Operator())
		c.Write(" ")
		c.Write(fmt.Sprintf("$%d", *argIndex))
		*args = append(*args, cond.Value())
		*argIndex++
	}
	return nil
}
