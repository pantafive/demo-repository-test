# Testing and Databases

Make this text clear and concise. It should be easy to understand.:

We want to be able to run database-related tests in parallel and in any order. Naive approach would be to run tests in
separate transactions, but it makes hard to inspect the database state after the failed test.

The project demonstrates how to generate a unique database for each test. If the test fails, the relevant database is
preserved for inspection; otherwise, it gets deleted.deleted.

```
--- FAIL: Test/create_user:_false_negative (0.02s)
    database_test.go:47:
            Error Trace:    /workspace/app/database_test.go:47
            Error:          An error is expected but got nil.
            Test:           Test/create_user:_false_negative
    database_test.go:48: database: test_create_user_false_negative
```

To make tests fast, this demo uses two tricks:

- Use in-memory database (check tmpfs in docker-compose.yaml)
- Use [database-templates](https://www.postgresql.org/docs/current/manage-ag-templatedbs.html) - database is created
  from a template, so it's much faster than creating a new database from scratch.
