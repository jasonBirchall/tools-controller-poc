# Airflow Helm Chart

Run Airflow on Kubernetes.

It looks for DAGs in the airflow-dags repository. Tasks are executed as new k8s pods if you use KubernetesPodOperator.

(gitSync is used to get the airflow-dags repository, polling it regularly.)

## Install/Upgrade the Chart

Run:

```bash
$ helm upgrade --dry-run --debug --install airflow charts/airflow-k8s --namespace airflow -f chart-env-config/ENV/airflow-k8s.yml
```

**NOTE**: double-check output, than once you're happy with the changes remove
          the `--dry-run` flag and re-run.

Airflow will be available at <https://airflow.tools.ENV.mojanalytics.xyz>.

**NOTE**: Change the environment config file to deploy in a different environment
          (the URL will change accordingly)


## Configuration

**IMPORTANT**: Please change credentials from the default/weak ones to proper/strong credentials.


| Parameter           | Description     | Default.    |
| ------------------- | --------------- | ----------- |
| `airflow.admin.username` | Airflow admin username | `""` |
| `airflow.admin.password` | Airflow admin password | `""` |
| `airflow.auth.type`         | Type of authentication to use. Can be `"AUTH_DB"` which will allow you to login using username/password for a user in the DB or `"AUTH_OID"` which will allow you to login using OIDC. | `"AUTH_DB"` |
| `airflow.auth.oauth.*` | OAuth configuration, required only when `airflow.auth.type` is `"AUTH_OID"` | |
| `airflow.config.core.parallelism` | Airflow configuration determining how many tasks it will run at the same time | `32` |
| `airflow.volumes.dags.nfsServer` | Host of the NFS server where DAGs will be stored. | `""` |
| `airflow.volumes.dags.path` | Path on the NFS server where DAGs will be stored. | `"/dags"` |
| `airflow.volumes.logs.nfsServer` | Host of the NFS server where logs will be stored. | `""` |
| `airflow.volumes.logs.path` | Path on the NFS server where DAGs will be stored. | `"/logs"` |
| `gitSync.interval`| Number of seconds between DAGs syncing from Git repository | `180` (3 minutes) |
| `gitSync.repository`| Git repository with the DAGs. If it is not a public repo, specify creds: `https://username:password@github.com/username/repository.git` | `"https://github.com/ministryofjustice/analytics-platform-airflow-example-dags"` |
| `gitSync.branch`| Branch to use when syncing the Git repository with the DAGs | `master` |
| `postgres.host`     | (required)      | `""`        |
| `postgres.port`     |                 | `5432`      |
| `postgres.user`     | (required)      | `"airflow"` |
| `postgres.password` | (required)      | `""`        |
| `postgres.db`       | (required)      | `"airflow"` |
| `redis.password`    | password for redis-server (required)      | `"airflow"` |
