# dtt

Data transformation testing !  (In development now!)

The tool is aiming to simplify SQL transformation testing. Source and expected data are defined in csv, and the test is defined in YAML file.

## How it works

The idea is really simple. The tool just builds a SQL statement that combines source data with `WITH` statement, and runs `WITH` and transformation SQL in real database without creating physical tables. Currently only postgres is supported.

## How to run it

Example:

```bash
export DATABASE_DRIVER=postgres
export DATA_SOURCE_NAME="host=192.168.64.6 user=postgres password=postgres dbname=postgres sslmode=disable"
go run main.go --spec tests/spec1/spec.yaml
go run main.go --spec tests/spec2/spec.yaml
```

## Spec Example

### Spec with schema defined

```yaml
tables:
  - name: orders
    columns:
      - name: id
        data_type: int
      - name: product_id
        data_type: int
      - name: quantity
        data_type: int
      - name: created_date
        data_type: timestamp
  - name: products
    columns:
      - name: id
        data_type: int
      - name: description
        data_type: varchar

tests:
  - name: orders test
    sources:
      - table_name: orders
        csv: fixtures/orders.csv
      - table_name: products
        csv: fixtures/products.csv
    transformation:
      query: | 
        SELECT orders.id, orders.quantity, products.description from orders
        INNER JOIN products on orders.product_id = products.id
    expected_result:
      csv: fixtures/expectations/test1.csv

  - name: orders test 2
    sources:
      - table_name: orders
        csv: fixtures/orders.csv
      - table_name: products
        csv: fixtures/products.csv
    transformation:
      query: | 
        SELECT count(*) from orders
        INNER JOIN products on orders.product_id = products.id
    expected_result:
      csv: fixtures/expectations/test2.csv

  - name: orders test 3
    sources:
      - table_name: orders
        csv: fixtures/orders.csv
      - table_name: products
        csv: fixtures/products.csv
    transformation:
      query: |
        WITH vars as (
            select '2021-01-01'::timestamp as processed_date
        )
        SELECT orders.id, orders.quantity, products.description, vars.processed_date 
        from vars, orders
        INNER JOIN products on orders.product_id = products.id
    expected_result:
      csv: fixtures/expectations/test3.csv

  - name: orders test 4
    sources:
      - table_name: orders
        csv: fixtures/orders.csv
      - table_name: products
        csv: fixtures/products.csv
    transformation: 
      query: |
        SELECT 1 from orders
        INNER JOIN products on orders.product_id = products.id
    expected_result:
      csv: fixtures/expectations/test4.csv
```

### Spec without schema defined (but the data must be defined in certain format)

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

