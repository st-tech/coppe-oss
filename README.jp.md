# Coppe

Coppeは、BigQuery用データ品質監視ツールです。定期的なBigQueryへの監視やSlackへの通知を行います。監視項目の追加はYAMLとSQLファイルのみで可能です。

インフラ構成は下図の通りです。TerraformとGitHub Actionsによる自動デプロイが可能です。

![Coppe Infra Diagram drawio](https://user-images.githubusercontent.com/36804811/138837195-c01eea1f-710e-4112-b3b2-aa3759f5adc2.png)

## Required Settings

1. GitHub Actions Secretsの追加

- GOOGLE_APPLICATION_CREDENTIALS_JSON - JSON形式のGCP認証キー

2. Cloud Monitoring用のNotification Channel IDの設定

- Google Cloud Consoleで通知チャンネルにSlackを追加。公式ドキュメントは[こちら](https://cloud.google.com/monitoring/support/notification-options)
- 設定後、以下のgcloudコマンドで通知チャンネルIDを取得し、**coppe.tf**の `locals.slack_channel_id`に貼り付ける。公式ドキュメントは[こちら](https://cloud.google.com/sdk/gcloud/reference/alpha/monitoring/channels/)

```
$ gcloud alpha monitoring channels list | grep 'name:'
```

3. **.env.yaml**に環境変数を設定

- GCP_PROJECT_ID (デプロイ用のGCPプロジェクトID)
- SLACK_HOOK_URL (アラート通知用のデフォルトのSlack Webhook URL。詳しくは[Alert Channel](#alert-channel))
- TIMEZONE (スケジュール用。日本時間の場合、Asia/Tokyo)


## Get Started

1. `./setup.sh`の実行で、Service Accountとtfstate管理用のStorage Bucketを作成

2. `yaml`フォルダーにYAMLファイルを以下のフォーマットで記述

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

3. GitHubにPush/Merge - GitHub ActionsとTerraformによって、インフラの設定やデプロイは行われます


## Format of YAML file

### Schedule

監視項目のスケジュールはcron形式で記述してください。cron形式は以下の方法で記述できます。

```
<minute> <hour> <day> <month> <weekday>
```

例：

毎時0分と30分：0,30 * * * *

毎時10分、11分、12分：10-12 * * * *

5分毎：*/5 * * * *

毎週月曜日の12時：0 12 * * MON

詳しくはこちら（https://github.com/adhocore/gronx#cron-expression）


### SQL in File


`sql:`に直接SQLを書く代わりに、SQLファイル`sql`フォルダーに置いて、`sql`フォルダーからの相対パスにを`sql_file:`に書くこともできます。

例：
```
- ...
  sql_file: sample.foo
  ...
```

In `sql/sample.foo`,

```
SELECT count(*) FROM foo
```

\* 注意： `sql:`か`sql_file`のどちらかは必須です



### SQL Parameters

SQLにパラメーターを設定することもできます。その場合、`params:`に連想配列として記述してください。

例：
```
- ...
  sql: SELECT count(*) FROM `{{.table}}` limit {{.limit}}
  params:
    table: "zz.foo"
    limit: 100
  ...

```

\* Coppeはテキストテンプレートライブラリを使用しています。書き方など詳しくはこちら（https://pkg.go.dev/text/template）


### SQL Matrix

複数のパラメータの組み合わせを使用したい場合、`matrix:`を利用することも可能です。

例：
```
- ...
  sql: SELECT count(*) FROM `project-{{.env}}.schema.table_name` ...
  matrix:
    env: [stg, prd]
  ...

```
上の監視項目は、下の2通りの組み合わせに分解されます。

```
- ...
  sql: SELECT count(*) FROM `project-stg.schema.table_name` ...
  ...

- ...
  sql: SELECT count(*) FROM `project-prd.schema.table_name` ...
  ...

```

\* Coppeはテキストテンプレートライブラリを使用しています。書き方など詳しくはこちら（https://pkg.go.dev/text/template）


### Expected Row Count / Expression

Coppeは期待するクエリ結果として、`row_count:`もしくは`expression:`を利用可能です。

`row_count:`：クエリ結果の列数

`expression:`：正常であれば`true`になるべきクエリ結果を利用した式。例は以下を参照。

```
- ...
  sql: SELECT table_name, count(*) as error_count ...
  ...
  expected:
    expression: table_name == "foo" && error_count > 10
```

\* 注意：row_countとexpressionは、どちらか一方のみ記載可能です

### Alert Channel

監視項目によって、Slackの通知チャンネルを使い分けたい場合、環境変数の追加と`channel:`で通知するURLを指定できます。指定しなかった場合、デフォルトでSLACK_HOOK_URLが使われますが、`channel:`で指定されている場合、SLACK_HOOK_URL_ + `channel:`の値を環境変数から取得し、通知に使用します。

例えば、

```
channel: CRITICAL
```

のように指定する場合、環境変数（.env.yaml）に'SLACK_HOOK_URL_CRITICAL'を取得します。もし、環境変数になかった場合は代わりにデフォルトのURLが使用されます。また、`channel:`の値はケース無視（criticalと書いてもCRITICALとして扱われる）されますが、環境変数は大文字で書くようにしてください。

### Alert Message

クエリ結果が期待される値でなかった場合、Slackへの通知を行います。アラートメッセージはクエリ結果を利用・展開して書くこともできます。また、SQLに使用したparamsやmatrixも使用することが可能です。

それぞれの型は、
- query_result: []map[string]interface{}
- params: map[string]interface{}
- matrix: map[string]interface{}

例：
```
- ...
  sql: SELECT table_name, count(*) as cnt, avg(diff) as delay_avg, max(diff) as delay_max ...
  params:
    interval_time: 5
  matrix:
    env: [prd, stg]
  description: |
    ENV: {{ .matrix.env }}
    Detected more than {{ .interval_time }} min data transfer delay in the following tables
    {{ range . }}
    {{ .table_name }} : {{ .cnt }} (cnt) : {{ .delay_avg }} (delay_avg) : {{ .delay_max }} (max_delay)
    {{ end }}
```

クエリ結果はの型は`[]map[string]interface{}`なので、基本的な書き方としては、

```
{{ . }}
```
で、クエリ結果を[]map[string]interface{}のまま、printfします

```
{{ range . }}
{{ .column_name }}
{{ end }}
```
クエリ結果をループして表示可能です。なお、column_nameには、クエリで取得したカラム名のみ使用可能です。

```
{{ range $i, $row := . }}
{{　$i }} :  {{ $row.column_name }}
{{ end }}
```
また、インデックスはこのように取得できます。


```
{{ range . }}
{{ if ne .column_name "" }}
- foo
{{ else }} 
- bar
{{ end }}
{{ end }}
```
if文はこのように記載可能です。

\* テキストテンプレートライブラリを使用しているので、書き方など詳しくは[公式ドキュメント](https://pkg.go.dev/text/template)を参照してください。

## Contributing

IssueやPull Requestの作成など、コントリビューションは誰でも歓迎です。コントリビューターは[Contributor Covenant](http://contributor-covenant.org/)を遵守することを期待します。

## License

CoppeはMITライセンスを適応してオープンソースとして公開しています。[LICENSE](./LICENSE.txt) を参照してください。
