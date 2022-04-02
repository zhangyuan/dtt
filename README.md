# dtt
Data transformation testing !

The tool is aiming to simplify SQL transformation testing. Source and expected data are defined in csv, and the test is defined in YAML file.

## How it works

The idea is really simple. The tool just builds a SQL statement that combine source data with `WITH` statement, and run it in real database without creating physical tables. Currently only postgres is supported.

## How to run it

Example:

```bash
export DATABASE_DRIVER=postgres
export DATA_SOURCE_NAME="host=192.168.64.6 user=postgres password=postgres dbname=postgres sslmode=disable"
go run main.go --spec tests/spec1/spec.yaml
go run main.go --spec tests/spec2/spec.yaml
```

## Spec Example

```yaml
tests:
  - name: orders test
    sources:
      - table_name: orders
        csv: fixtures/spec2/orders.csv
      - table_name: products
        csv: fixtures/spec2/products.csv
    transformation:
      query: | 
        SELECT orders.id, orders.quantity, products.description from orders
        INNER JOIN products on orders.product_id = products.id
    expected_result:
      csv: fixtures/spec2/expectations/test1.csv
```

