# CEP Weather

HTTP API that receives a Brazilian CEP, resolves the city with ViaCEP, and returns the current temperature in Celsius, Fahrenheit, and Kelvin. The service runs on Google Cloud Run.

## Cloud Run

The live URL is printed by the Deploy workflow after the first successful apply (`terraform output service_url`). Replace this placeholder with that URL:

```text
https://cep-weather-4bggcd6kca-rj.a.run.app```

Request:

```bash
curl -sS "https://cep-weather-4bggcd6kca-rj.a.run.app/weather/01001000"```

## API

`GET /weather/{zipcode}`

`zipcode` must be exactly 8 digits. A hyphenated CEP is rejected.

Success `200`:

```json
{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}
```

| Condition | Status | Body |
| --- | --- | --- |
| CEP is not 8 digits | 422 | `invalid zipcode` |
| CEP is unknown | 404 | `can not find zipcode` |
| Location or weather provider fails | 502 | `weather service unavailable` |

Conversions, rounded to two decimal places:

- Fahrenheit: `F = C × 1.8 + 32`
- Kelvin: `K = C + 273`

## Run the tests

```bash
go test -race ./...
```

## Run locally with Docker

Create a key at [WeatherAPI](https://www.weatherapi.com/signup.aspx), then:

```bash
docker build -t cep-weather .
docker run --rm -p 8080:8080 -e WEATHER_API_KEY="$WEATHER_API_KEY" cep-weather
curl -sS "http://localhost:8080/weather/01001000"
```

The container listens on `PORT` (default `8080`). Optional overrides: `VIACEP_BASE_URL`, `WEATHER_BASE_URL`.

## Deploy

Infrastructure is Terraform. GitHub Actions on `main` runs the tests, pushes the image to Artifact Registry, and applies the Cloud Run stack.

### 1. Bootstrap the GCP project once

```bash
gcloud auth application-default login
terraform -chdir=terraform/bootstrap init
terraform -chdir=terraform/bootstrap apply \
  -var=project_id="$GCP_PROJECT_ID" \
  -var=github_repository="OWNER/REPO"
```

Copy the outputs into GitHub:

| Output | GitHub secret |
| --- | --- |
| `workload_identity_provider` | `GCP_WORKLOAD_IDENTITY_PROVIDER` |
| `deployer_service_account` | `GCP_SERVICE_ACCOUNT` |

Also set secret `GCP_PROJECT_ID` and `WEATHER_API_KEY`.

Optional repository variable `GCP_REGION` (default `southamerica-east1`). The state bucket is `${GCP_PROJECT_ID}-cep-weather-tfstate`.

The WeatherAPI key is stored in Secret Manager. Cloud Run reads it through a secret reference. The runtime service account can access only that secret.

If `allUsers` cannot be granted `roles/run.invoker`, the organization policy `iam.allowedPolicyMemberDomains` is blocking public access. Allow `allUsers` for this project or the service stays private.

### 2. Push to `main`

The Deploy workflow creates the registry, publishes the image, then applies `terraform/service`. The job summary prints the public URL. Put that URL in the Cloud Run section above.
