## Go SQL and Struct Mapping (`sqlm`)

> [!WARNING]
> This package is in an alpha stage. The API is subject to change.

**sqlm** is a high-performance, declarative struct mapping toolkit for Go. It automates common CRUD operations, supports advanced mapping for complex queries.

## Philosophy & Key Features

- **Declarative**: Define your data model and mapping rules with struct tags, not imperative code.
- **Automatic CRUD**: 80% of typical project needs—insert, load, update, delete—are handled automatically and declaratively.
- **Advanced Query Mapping**: For complex queries, combine with `go-sqlb` to map results to nested structs with ease.
- **High Performance**: Uses code generation to avoid reflection at runtime, delivering near-manual performance even for large result sets.
- **Flexible**: Full control when you need it, automation when you don't.

---

## Example: Declarative CRUD (80% Use Case)

```go
func Example_cRUD() {
	// Model represents the base model with common fields.
	type Model struct {
		// ID is the model ID.
		// model tag is required for Load / Update / Delete / Insert operations.
		// pk means primary key.
		// returning means the ID will be returned after insertion.
		ID int64 `sqlb:"model;col:id;pk;returning"`
		// Created is the creation time.
		// readonly means this field will be excluded from INSERT / UPDATE, created usually set by DB default value.
		Created *time.Time `sqlb:"col:created;readonly"`
		// Updated is the last update time.
		// conflict_set means when inserting an existing record, the Updated will be updated.
		Updated *time.Time `sqlb:"col:updated;conflict_set"`
		// soft_delete indicates this column is used for soft deletion.
		// When deleting, the column will be set to current time or true instead of actually deleting the record.
		// When loading, records with non-zero value in this column will be ignored.
		// conflict_set means when inserting an deleted record, undelete it.
		Deleted *time.Time `sqlb:"col:deleted;soft_delete;conflict_set"`
	}

	type User struct {
		// The value defined by 'tables' and 'from' can be inherited by nested fields
		// and by subsequent sibling fields of the current struct.
		Model `sqlb:"table:users"`
		// unique indicates this column has a unique constraint, which can be used to locate records for Load / Delete operation.
		// conflict_on indicates the column(s) to check for conflict during insert.
		Email string `sqlb:"col:email;required;unique;conflict_on"`
		// conflict_set without value means to use excluded column value
		Name string `sqlb:"col:name;conflict_set"`
	}

	ctx := sqlb.NewContext(context.Background(), dialect.PostgreSQL{})

	_, err := sqlm.Load(ctx, nil, &User{Model: Model{ID: 1}}, option.WithDebug())
	if err != nil && !errors.Is(err, sqlm.ErrNilDB) {
		fmt.Println(err)
	}
	// Output:
	// [Load(*sqlm_test.User)] SELECT "created", "updated", "email", "name" FROM "users" WHERE "id" = 1 AND "deleted" IS NULL
}
```

## Example: Complex Query Mapping (20% Use Case)

```go
// Model represents the base model with common fields.
type Model struct {
	ID      int        `sqlb:"model;col:id"`
	Created *time.Time `sqlb:"col:created"`
	Updated *time.Time `sqlb:"col:updated"`
	Deleted *time.Time `sqlb:"col:deleted"`
}

type User struct {
	Model `sqlb:"table:users"`
	OrgID int    `sqlb:"col:org_id"`
	Name  string `sqlb:"col:name"`
}

type Org struct {
	Model `sqlb:"table:orgs"`
	Name  string `sqlb:"col:name"`
}
type userListItem struct {
	User
	Org
}

func Example_complexSelect() {
	Users := sqlb.NewTable("users")
	Orgs := sqlb.NewTable("orgs")
	b := sqlb.NewSelectBuilder().
		From(Users).
		InnerJoin(Orgs, sqlf.F(
			"? = ?",
			// Without code generation, you have to write the column names manually,
			// which is error-prone and not refactor-friendly.
			Users.Column("org_id"),
			Orgs.Column("id"),
		)).
		WhereEquals(Orgs.Column("id"), 1).
		WhereIsNull(Users.Column("deleted"))
		// With code generation, you can use the generated column accessors, 
		// which are type-safe and refactor-friendly.
		// WhereEquals(org.ColumnID(), 1).
		// WhereIsNull(user.ColumnDeleted())
	ctx := sqlb.NewContext(context.Background(), dialect.PostgreSQL{})
	// Select uses generated codes of userListItem if available, 
	// otherwise it will fallback to reflection.
	_, err := sqlm.Select[*userListItem](ctx, nil, b, option.WithDebug())
	if err != nil && !errors.Is(err, sqlm.ErrNilDB) {
		fmt.Println(err)
	}
	// Output:
	// [Select(*sqlm_test.userListItem)] SELECT "users"."id", "users"."created", "users"."updated", "users"."deleted", "users"."org_id", "users"."name", "orgs"."id", "orgs"."created", "orgs"."updated", "orgs"."deleted", "orgs"."name" FROM "users" INNER JOIN "orgs" ON "users"."org_id" = "orgs"."id" WHERE "orgs"."id" = 1 AND "users"."deleted" IS NULL
}
```

For maximum performance, especially when scanning a large number of rows, use the code generation tool (`sqlmgen`). This eliminates all runtime reflection and delivers performance close to hand-written code.

---

## Part of the Go SQL Tools Family

This project is part of a family of Go SQL tools, each designed for a different level of abstraction and automation:

1. **[go-sqlf](https://github.com/qjebbs/go-sqlf)** — Minimalist SQL fragment builder. For simple, manual SQL composition with parameter binding and zero magic.
2. **[go-sqlb](https://github.com/qjebbs/go-sqlb)** — Advanced SQL builder. For programmatically building complex queries (CTE, JOIN, expressions, etc.) with chainable, declarative, and composable APIs.
3. **go-sqlm** (this project) — Struct mapping. Declarative struct mapping, automatic CRUD, batch operations, and high-performance zero-reflection code generation.