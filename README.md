# Coppe

日本語版READMEは[こちら](./README.jp.md)

Coppe is a data quality monitoring tool for BigQuery. Write your YAML and SQL files for adding monitoring items, and Coppe will take care of deployment and periodic check of your data in BigQuery. Alert messages are thrown to Slack channel that you specified.

The following is the architecture of Coppe. Easily deployable on Google Cloud Platform with prepared Terraform and GitHub Actions config files.

![Coppe Infra Diagram drawio](https://user-images.githubusercontent.com/36804811/138837195-c01eea1f-710e-4112-b3b2-aa3759f5adc2.png)

## Required Settings

1. Github Actions Secrets

- GOOGLE_APPLICATION_CREDENTIALS_JSON - GCP credential key in JSON format

2. Notification Channel ID for Cloud Monitoring

- Create Notification Channel with Slack in Google Cloud Console. Refer to [Managing notification channels | Cloud Monitoring | Google Cloud](https://cloud.google.com/monitoring/support/notification-options)
- Get the notification channel ID retrieval from the following command. Refer to [gcloud alpha monitoring channels | Cloud SDK Documentation](https://cloud.google.com/sdk/gcloud/reference/alpha/monitoring/channels/)

```
$ gcloud alpha monitoring channels list | grep 'name:'
```

- Paste the ID to `locals.slack_channel_id` in coppe.tf

3. Set the following environment variables in **.env.yaml**

- GCP_PROJECT_ID (GCP project ID for deploy)
- SLACK_HOOK_URL (Default Slack Webhook URL for alert. For more info: [Alert Channel](#alert-channel))
- TIMEZONE (For scheduling purpose)


## Get Started

1. Run `./setup.sh` to set up Service Account and Storage Bucket for remote management of tfstate


2. Write YAML files in `yaml` directory in the following format

```
- schedule: "* * * * *"
  sql: "SELECT count(*) as cnt FROM foo"
  expect:
    expression: "cnt == 0"
  description: "Foo table check"

- schedule: "* * * * *"
  sql: "SELECT * as cnt FROM foo"
  expect:
    row_count: 0
  description: "Foo table check"

```

3. Push/Merge to main - Github Actions and Terraform will take care of other infrastructure setup and deployment in GCP.


## Format of YAML file

### Schedule

Coppe uses cron expression parser for scheduling. Cron expression consists of 5 segments:

```
<minute> <hour> <day> <month> <weekday>
```

For more details and examples, please refer to https://github.com/adhocore/gronx#cron-expression.


### SQL in File

Instead of directly writing SQL in `sql:` row, You can put a path to SQL file in `sql_file:`. Detectable are only files in `sql` directory.

```
- ...
  sql_file: sample.foo
  ...
```

In `sql/sample.foo`,

```
SELECT count(*) FROM foo
```

\* Either `sql:` or `sql_file` must represent in a monitoring item.



### SQL Parameters

You can write parameters to set in SQL with an associative aray in `params:` in YAML files.

e.g.
```
- ...
  sql: SELECT count(*) FROM `{{.table}}` limit {{.limit}}
  params:
    table: "zz.foo"
    limit: 100
  ...

```

\* Coppe parses by itself using Go's standard library. For more details of syntax, please refer to the document (https://pkg.go.dev/text/template).


### SQL Matrix

e.g.
```
- ...
  sql: SELECT count(*) FROM `project-{{.env}}.schema.table_name` ...
  matrix:
    env: [stg, prd]
  ...

```

You can write matrix to set in SQL in YAML files. The example above is equivalent to:

```
- ...
  sql: SELECT count(*) FROM `project-stg.schema.table_name` ...
  ...

- ...
  sql: SELECT count(*) FROM `project-prd.schema.table_name` ...
  ...

```

\* Coppe parses by itself using Go's standard library, text/template. For more details of syntax, please refer to the document (https://pkg.go.dev/text/template).


### Expected Row Count / Expression

Coppe accepts either `row_count:` or `expression:` under `expectation:`. 

In `row_count:`, you can specify the expected row size of the query result.

In `expression:`, you can write an expression that should be true using the column names from the query result. For example:

```
- ...
  sql: SELECT table_name, count(*) as error_count ...
  ...
  expected:
    expression: table_name == "foo" && error_count > 10
```

\* You cannot write both row_count and expression together under expectation in one monitoring item.

### Alert Channel

You can specify other Slack channel than the default one for alert with `channel:`. For this, You need to add specified Slack Incoming Webhook URL in .env.yaml.

By default, the URL used for Slack notification is 'SLACK_HOOK_URL' in .env.yaml if nothing specified in `channel:`. If specified, an environment variable for the word in `channel` followed by 'SLACK_HOOK_URL_' is used. 

For example, 

```
channel: CRITICAL
```

Coppe looks for 'SLACK_HOOK_URL_CRITICAL' in .env.yaml, and use it for notification. If not exists, the default URL 'SLACK_HOOK_URL' would be used. `channel:` is case-insensitive, but in .env.yaml, you need to write in capitalized.

### Alert Message

You can write an alert message based on the values of query_result, params, and matrix.

Each of their types is:
- query_result: []map[string]interface{}
- params: map[string]interface{}
- matrix: map[string]interface{}

e.g.
```
- ...
  sql: SELECT table_name, count(*) as cnt, avg(diff) as delay_avg, max(diff) as delay_max ...
  params:
    interval_time: 5
  matrix:
    env: [prd, stg]
  description: |
    ENV: {{ .matrix.env }}
    Detected more than {{ .params.interval_time }} min data transfer delay in the following tables
    {{ range .query_result }}
    {{ .table_name }} : {{ .cnt }} (cnt) : {{ .delay_avg }} (delay_avg) : {{ .delay_max }} (max_delay)
    {{ end }}
```

\* For more details of syntax, please refer to the document for text/template (https://pkg.go.dev/text/template).
