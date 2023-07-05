# Testing and Databases

Project demonstrates how to create a new database for every test.
If the test fails, a relevant database isn't deleted, and you can check what happened; Otherwise the database is deleted.

```
 === CONT  Test/create_user:_false_negative
        	Error Trace:	/home/runner/work/test-repository-demo/test-repository-demo/app/database_test.go:49
        	Error:      	An error is expected but got nil.
        	Test:       	Test/create_user:_false_negative
    database_test.go:50: database: test_01h4kcddzx7mhgep640zcpq1qh
```
