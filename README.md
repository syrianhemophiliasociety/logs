<div align="center">
  <a href="https://logs.syrianhemophiliasociety.com" target="_blank"><img src="https://logs.syrianhemophiliasociety.com/assets/web-app-manifest-192x192.png" width="150" /></a>

  <h1>SyrianHemophiliaSocietyLogs</h1>
  <p>
    <strong>A patient care follow up platform for Syrian Hemophilia Society</strong>
  </p>
  <p>
    <a href="https://goreportcard.com/report/github.com/syrianhemophiliasociety/log"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/syrianhemophiliasociety/logs"/></a>
    <a href="https://github.com/syrianhemophiliasociety/logs/actions/workflows/rex-deploy.yml"><img alt="rex-deployment" src="https://github.com/syrianhemophiliasociety/logs/actions/workflows/rex-deploy.yml/badge.svg"/></a>
  </p>
</div>

## About

**SyrianHemophiliaSocietyLogs** is something idk.

## Contributing

IDK, it would be really nice of you to contribute, check the poorly written [CONTRIBUTING.md](/CONTRIBUTING.md) for more info.

## Run locally

1. Clone the repo.

```bash
git clone https://github.com/syrianhemophiliasociety/logs
```

2. Create the docker environment file

```bash
cp .env.example .env.docker
```

3. Run it with docker compose.

```bash
docker compose up -f docker-compose-all.yml
```

3. Visit http://localhost:23103
4. Don't ask why I chose this weird port.

---

Made with 🧉 by [Baraa Al-Masri](https://syrianhemophiliasociety.com)
