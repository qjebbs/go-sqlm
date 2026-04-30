package option

import (
	"github.com/qjebbs/go-sqlb"
	"github.com/qjebbs/go-sqlm/internal/util"
)

// SelectOptions defines options for Select().
type SelectOptions struct {
	Tags           util.List[string]
	CoalesceTables util.List[string]
}

// WithSelectTags is an option for Select() which sets the scan field tags to select.
func WithSelectTags(tags ...string) Option {
	return func(o *Options) {
		o.Select.Tags = tags
	}
}

// WithSelectCoalesce is an option for Select() which enables COALESCE for fields of the specified tables.
// The zero value is used as the default value for NULL fields.
// To decide whether to enable COALESCE for a field, it matches the tables.Name
// against the effective `table` key value (e.g. `sqlb:"table:foo"`) of the field.
//
// Example:
//
//	type Foo struct {
//		ID  int64  `sqlb:"col:id;table:foo"`
//		Bar string `sqlb:"col:bar"`
//	}
//	foo := sqlb.NewTable("foo")
//	// SELECT COALESCE("foo"."id", 0), COALESCE("foo"."bar", '') FROM ...
//	sqlm.Select[*Foo](..., option.WithSelectCoalesce(foo))
func WithSelectCoalesce(tables ...sqlb.Table) Option {
	return func(o *Options) {
		o.Select.CoalesceTables = util.Map(tables, func(t sqlb.Table) string {
			return t.Name
		})
	}
}
